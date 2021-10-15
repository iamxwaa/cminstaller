function installMy(){
    echo ">> 开始安装必要组件"
    chmod +x $basePackagePath/my/*.rpm
    rpm -ih $basePackagePath/my/*.rpm
    echo "<< 必要组件安装完毕"
}
function sshkeygen(){
    echo ">> 开始配置服务器免秘钥"
    rm -rf ~/.ssh
    ssh-keygen -t rsa -P "" -f ~/.ssh/id_rsa
    cat ~/.ssh/id_rsa.pub >> ~/.ssh/authorized_keys
    echo "StrictHostKeyChecking no" > ~/.ssh/config
    for (( i=0 ; i < ${#slaves[@]} ; i++ )) do
        slave=${slaves[i]}
        hostInfo=(${slave// / })
        sshpass -p ${hostInfo[3]} ssh -p ${hostInfo[1]} ${hostInfo[2]}@${hostInfo[0]} "rm -rf ~/.ssh;ssh-keygen -t rsa -P \"\" -f ~/.ssh/id_rsa;cat ~/.ssh/id_rsa.pub >> ~/.ssh/authorized_keys"
        sshpass -p ${hostInfo[3]} ssh -p ${hostInfo[1]} ${hostInfo[2]}@${hostInfo[0]} "cat ~/.ssh/authorized_keys" >> ~/.ssh/authorized_keys
    done
    for (( i=0 ; i < ${#slaves[@]} ; i++ )) do
        slave=${slaves[i]}
        hostInfo=(${slave// / })
        sshpass -p ${hostInfo[3]} scp -P ${hostInfo[1]} ~/.ssh/authorized_keys ${hostInfo[2]}@${hostInfo[0]}://root/.ssh/
    done
    echo "<< 服务器免秘钥配置完毕"
}
function updateHosts(){
    echo ">> 开始更新各节点hosts"
    echo "${master[0]}: 更新/etc/hosts"
    echo "${master[0]} ${master[4]}" >> /etc/hosts
    for (( i=0 ; i < ${#slaves[@]} ; i++ )) do
        slave=${slaves[i]}
        hostInfo=(${slave// / })
        echo "${hostInfo[0]} ${hostInfo[4]}" >> /etc/hosts
    done
    for (( i=0 ; i < ${#slaves[@]} ; i++ )) do
        slave=${slaves[i]}
        hostInfo=(${slave// / })
        echo "${hostInfo[0]}: 更新/etc/hosts"
        ssh -p ${hostInfo[1]} ${hostInfo[2]}@${hostInfo[0]} "echo "${master[0]} ${master[4]}" >> /etc/hosts"
        for (( j=0 ; j < ${#slaves[@]} ; j++ )) do
            slave2=${slaves[j]}
            hostInfo2=(${slave2// / })
            ssh -p ${hostInfo[1]} ${hostInfo[2]}@${hostInfo[0]} "echo \"${hostInfo2[0]} ${hostInfo2[4]}\" >> /etc/hosts"
        done
    done
    echo "<< hosts更新完毕"
}
function shutdownFirewall(){
    echo "${master[0]}: 关闭和禁用防火墙"
    systemctl stop firewalld
    systemctl disable firewalld
    for (( i=0 ; i < ${#slaves[@]} ; i++ )) do
        slave=${slaves[i]}
        hostInfo=(${slave// / })
        echo "${hostInfo[0]}: 关闭和禁用防火墙"
        ssh -p ${hostInfo[1]} ${hostInfo[2]}@${hostInfo[0]} "systemctl stop firewalld;systemctl disable firewalld"
    done
}
function installJDK(){
    echo ">> 准备安装java运行环境"
    jdkPackage=`ls $basePackagePath/jdk`
    cat > /tmp/jdkinstall.sh << EOF
#!/bin/bash
mkdir -p /opt/vrv/java
dirName=\`ls /opt/vrv/java\`
if [ "" == "\$dirName" ];then 
    tar -zxvf $basePackagePath/jdk/$jdkPackage -C /opt/vrv/java
    dirName=\`ls /opt/vrv/java\`
    javaHome="/opt/vrv/java/\$dirName"
    echo "export JAVA_HOME=\$javaHome" >> /etc/profile
    echo "export PATH=\\\$PATH:\\\$JAVA_HOME/bin:/usr/share/java" >> /etc/profile
    echo "export CLASSPATH=.:\\\$JAVA_HOME/lib/tools.jar:\\\$JAVA_HOME/lib/dt.jar" >> /etc/profile
    source /etc/profile
    echo "JAVA_HOME=\$javaHome"
    chmod +x \$javaHome/bin/java
    if [ ! -f /usr/local/bin/java ];then
        ln -s \$javaHome/bin/java /usr/local/bin
    fi
    java -version
    mkdir -p /usr/java
    ln -s \$javaHome  /usr/java/default
fi
EOF
    chmod +x /tmp/jdkinstall.sh
    echo "${master[0]}: 开始安装 $jdkPackage"
    /tmp/jdkinstall.sh
    for (( i=0 ; i < ${#slaves[@]} ; i++ )) do
    (
        slave=${slaves[i]}
        hostInfo=(${slave// / })
        echo "${hostInfo[0]}: 上传安装包 $jdkPackage"
        ssh -p ${hostInfo[1]} ${hostInfo[2]}@${hostInfo[0]} "mkdir -p $basePackagePath/jdk"
        scp -P ${hostInfo[1]} $basePackagePath/jdk/$jdkPackage ${hostInfo[2]}@${hostInfo[0]}:/$basePackagePath/jdk
        echo "${hostInfo[0]}: 开始安装 $jdkPackage"
        scp -P ${hostInfo[1]} /tmp/jdkinstall.sh ${hostInfo[2]}@${hostInfo[0]}://tmp
        ssh -p ${hostInfo[1]} ${hostInfo[2]}@${hostInfo[0]} "chmod +x /tmp/jdkinstall.sh;/tmp/jdkinstall.sh"
    ) &
    done
    wait
    echo ">> java运行环境安装完毕"
}
function installJDK2(){
    echo ">> 准备安装java运行环境"
    jdkPackage=`ls $basePackagePath/jdk`
    cat > /tmp/jdkinstall.sh << EOF
#!/bin/bash
mkdir -p /opt/vrv/java
dirName=\`ls /opt/vrv/java\`
if [ "" == "\$dirName" ];then 
    tar -zxvf $basePackagePath/jdk/$jdkPackage -C /opt/vrv/java
    dirName=\`ls /opt/vrv/java\`
    javaHome="/opt/vrv/java/\$dirName"
    echo "export JAVA_HOME=\$javaHome" >> /etc/profile
    echo "export PATH=\\\$PATH:\\\$JAVA_HOME/bin:/usr/share/java" >> /etc/profile
    echo "export CLASSPATH=.:\\\$JAVA_HOME/lib/tools.jar:\\\$JAVA_HOME/lib/dt.jar" >> /etc/profile
    source /etc/profile
    echo "JAVA_HOME=\$javaHome"
    chmod +x \$javaHome/bin/java
    if [ ! -f /usr/local/bin/java ];then
        ln -s \$javaHome/bin/java /usr/local/bin
    fi
    java -version
    mkdir -p /usr/java
    ln -s \$javaHome  /usr/java/default
fi
EOF
    chmod +x /tmp/jdkinstall.sh
    echo "${master[0]}: 开始安装 $jdkPackage"
    /tmp/jdkinstall.sh
    for (( i=0 ; i < ${#slaves[@]} ; i++ )) do
    (
        slave=${slaves[i]}
        hostInfo=(${slave// / })
        echo "${hostInfo[0]}: 开始安装 $jdkPackage"
        scp -P ${hostInfo[1]} /tmp/jdkinstall.sh ${hostInfo[2]}@${hostInfo[0]}://tmp
        ssh -p ${hostInfo[1]} ${hostInfo[2]}@${hostInfo[0]} "chmod +x /tmp/jdkinstall.sh;/tmp/jdkinstall.sh"
    ) &
    done
    wait
    echo ">> java运行环境安装完毕"
}
function installMysql(){
    echo ">> 准备在${master[0]}安装 mysql"
    dirName=`ls $basePackagePath/mysql`
    myPath=$basePackagePath/mysql/$dirName
    
    if [ "$mysqlport" != "3306" ];then
        setenforce 0
    fi

    maria=`rpm -qa | grep mariadb`
    if [ "" != "$maria" ];then
        echo "卸载 mariadb"
        rpm -e --nodeps $maria
    fi
    mysqinstall=`rpm -qa | grep mysql`
    if [ "" != "$mysqinstall" ];then
        echo -e "已安装: \n$mysqinstall\n请卸载后重试"
        return
    fi
    chmod +x $myPath/*.rpm
    yum localinstall -y $myPath/*.rpm
    echo "设置mysql端口: $mysqlport"
    if [ -f "/etc/my.cnf" ];then
        echo "设置 skip-grant-tables"
        echo "port=$mysqlport" >> /etc/my.cnf
        echo "skip-grant-tables" >> /etc/my.cnf
    else
        echo "创建 /etc/my.cnf"
        touch /etc/my.cnf
        echo "# For advice on how to change settings please see
# http://dev.mysql.com/doc/refman/5.7/en/server-configuration-defaults.html

[mysqld]
port=$mysqlport

datadir=/var/lib/mysql
socket=/var/lib/mysql/mysql.sock

symbolic-links=0

log-error=/var/log/mysqld.log
pid-file=/var/run/mysqld/mysqld.pid
skip-grant-tables

" > /etc/my.cnf
    fi
    echo "启动 mysql"
    systemctl start mysqld
    mysql -uroot << EOF
use mysql;
update user set authentication_string=password("$mysqlpwd"),password_expired="N" where user='root' and host='localhost';
flush privileges;
EOF
    echo "移除 skip-grant-tables"
    sed -i 's/^skip-grant-tables//' /etc/my.cnf
    echo "重启 mysql"
    systemctl restart mysqld
    systemctl enable mysqld
    sleep 5s
    if [ ! -d "/usr/share/java" ];then
        mkdir /usr/share/java
    fi
    echo "复制mysql-connector-java.jar到/usr/share/java/mysql-connector-java.jar"
    cp $myPath/mysql-connector-java-* /usr/share/java/mysql-connector-java.jar
    for (( i=0 ; i < ${#slaves[@]} ; i++ )) do
        slave=${slaves[i]}
        hostInfo=(${slave// / })
        echo "复制mysql-connector-java.jar到${hostInfo[0]}://usr/share/java"
        ssh -p ${hostInfo[1]} ${hostInfo[2]}@${hostInfo[0]} "mkdir -p /usr/share/java"
        scp -P ${hostInfo[1]} /usr/share/java/mysql-connector-java.jar ${hostInfo[2]}@${hostInfo[0]}://usr/share/java
    done
    echo "<< mysql 安装完毕"
}
function uninstallMysql(){
    echo ">> 准备在${master[0]}卸载 mysql"
    systemctl stop mysqld

    mysql_package=`rpm -qa | grep mysql`
    if [ "" != "$mysql_package" ];then
        rpm -e --nodeps $mysql_package
    fi

    rm -fv /etc/my.cnf.rpmsave
    rm -rfv /var/lib/mysql*
    echo ">> mysql 卸载完毕"
}
function installHttpd(){
    echo ">> 准备在本机安装 httpd"
    yum localinstall -y $basePackagePath/httpd/*.rpm
    echo "启动 httpd"
    systemctl start httpd
    echo "访问地址: http://${master[0]}"
    echo "<< httpd 安装完毕"
}
function mountYum(){
    echo "开始挂载centos镜像"
    mkdir -p /var/www/html/centos
    umount /var/www/html/centos
    mount $basePackagePath/iso/* /var/www/html/centos
    rm -rf /var/www/html/cdh6_parcel
    mkdir -p /var/www/html/cdh6_parcel
    ln -s $basePackagePath/parcels/* /var/www/html/cdh6_parcel
    echo "###############################"
    echo "访问地址: http://${master[0]}/centos"
    echo "访问地址: http://${master[0]}/cdh6_parcel"
    echo "###############################"
}
function putYumRepo(){
    echo ">> 准备配置离线yum仓库"
    cat > /etc/yum.repos.d/vrv-centos.repo << EOF
[vrv-centos]
name=vrv-centos
baseurl=http://${master[0]}/centos
gpgkey=http://${master[0]}/centos/RPM-GPG-KEY-CentOS-7
gpgcheck=1
enabled=1
priority=1
EOF
    (
    yum clean all
    yum makecache fast
    yum repolist
    ) &
    for (( i=0 ; i < ${#slaves[@]} ; i++ )) do
    (
        slave=${slaves[i]}
        hostInfo=(${slave// / })
        echo "复制vrv-centos.repo到${hostInfo[0]}://etc/yum.repos.d"
        scp -P ${hostInfo[1]} /etc/yum.repos.d/vrv-centos.repo ${hostInfo[2]}@${hostInfo[0]}://etc/yum.repos.d
        ssh -p ${hostInfo[1]} ${hostInfo[2]}@${hostInfo[0]} "yum clean all;yum makecache fast;yum repolist"
    ) &
    done
    wait
    echo "<< 离线yum仓库配置完毕"
}
function uploadCDH(){
    echo ">> 开始上传CDH安装包到各节点"
    mkdir -p /opt/cloudera/parcel-repo
    cp -v $basePackagePath/parcels/* /opt/cloudera/parcel-repo
    for (( i=0 ; i < ${#slaves[@]} ; i++ )) do
    (   slave=${slaves[i]}
        hostInfo=(${slave// / })
        echo "上传cloudera-manager到${hostInfo[0]}:/$basePackagePath/cm"
        ssh -p ${hostInfo[1]} ${hostInfo[2]}@${hostInfo[0]} "mkdir -p $basePackagePath/cm"
        scp -P ${hostInfo[1]} $basePackagePath/cm/* ${hostInfo[2]}@${hostInfo[0]}:/$basePackagePath/cm
    ) &
    done
    wait
    echo "<< CDH安装包上传完毕"
}
function setupMysql(){
    echo ">> 开始配置mysql数据库"
    export MYSQL_PWD=$mysqlpwd

    echo 创建数据库用户: $mysqlcdhuser , 密码: $mysqlcdhpwd
    mysql -uroot -N << EOF
set global validate_password_policy=0;

CREATE USER '$mysqlcdhuser'@'%' IDENTIFIED BY '$mysqlcdhpwd';
CREATE USER '$mysqlcdhuser'@'localhost' IDENTIFIED BY '$mysqlcdhpwd';
CREATE USER '$mysqlcdhuser'@'${master[0]}' IDENTIFIED BY '$mysqlcdhpwd';

GRANT ALL PRIVILEGES ON *.* TO '$mysqlcdhuser'@'%';
GRANT ALL PRIVILEGES ON *.* TO '$mysqlcdhuser'@'localhost';
GRANT ALL PRIVILEGES ON *.* TO '$mysqlcdhuser'@'${master[0]}';

FLUSH PRIVILEGES;
EOF
    echo 创建数据库: amon,cm,hive,hue,oozie
    mysql -uroot -N << EOF
CREATE DATABASE IF NOT EXISTS amon;
CREATE DATABASE IF NOT EXISTS cm;
CREATE DATABASE IF NOT EXISTS hive;
CREATE DATABASE IF NOT EXISTS hue;
CREATE DATABASE IF NOT EXISTS oozie;
EOF

    echo "<< mysql 数据库配置完毕"
}
function optimizeServer(){
    echo ">> 开始优化各节点配置"
    echo 'SELINUX=disabled'
    sed -i "s/SELINUX=enforcing/SELINUX=disabled/g" /etc/selinux/config

    echo 'vm.swappiness=10'
    sysctl vm.swappiness=10
    echo 'vm.swappiness=10' >> /etc/sysctl.conf

    echo 'never > /sys/kernel/mm/transparent_hugepage/defrag'
    echo never > /sys/kernel/mm/transparent_hugepage/defrag
    echo 'echo never > /sys/kernel/mm/transparent_hugepage/defrag' >>  /etc/rc.local

    echo 'never > /sys/kernel/mm/transparent_hugepage/enabled'
    echo never > /sys/kernel/mm/transparent_hugepage/enabled
    echo 'echo never > /sys/kernel/mm/transparent_hugepage/enabled' >>  /etc/rc.local
    
    echo 'ulimit -n 65535 '
    ulimit -n 65535
    echo '* soft nofile 65535'
    echo '* soft nofile 65535' >> /etc/security/limits.conf
    echo '* hard nofile 65535'
    echo '* hard nofile 65535' >> /etc/security/limits.conf

    for (( i=0 ; i < ${#slaves[@]} ; i++ )) do
        slave=${slaves[i]}
        hostInfo=(${slave// / })
        echo "${hostInfo[0]}: SELINUX=disabled"
        ssh -p ${hostInfo[1]} ${hostInfo[2]}@${hostInfo[0]} "sed -i \"s/SELINUX=enforcing/SELINUX=disabled/g\" /etc/selinux/config"
        echo "${hostInfo[0]}: sysctl vm.swappiness=10"
        ssh -p ${hostInfo[1]} ${hostInfo[2]}@${hostInfo[0]} "sysctl vm.swappiness=10"
        ssh -p ${hostInfo[1]} ${hostInfo[2]}@${hostInfo[0]} "echo 'vm.swappiness=10' >> /etc/sysctl.conf"
        echo "${hostInfo[0]}: never > /sys/kernel/mm/transparent_hugepage/defrag"
        ssh -p ${hostInfo[1]} ${hostInfo[2]}@${hostInfo[0]} "echo never > /sys/kernel/mm/transparent_hugepage/defrag"
        ssh -p ${hostInfo[1]} ${hostInfo[2]}@${hostInfo[0]} "echo 'echo never > /sys/kernel/mm/transparent_hugepage/defrag' >>  /etc/rc.local"
        echo "${hostInfo[0]}: never > /sys/kernel/mm/transparent_hugepage/enabled"
        ssh -p ${hostInfo[1]} ${hostInfo[2]}@${hostInfo[0]} "echo never > /sys/kernel/mm/transparent_hugepage/enabled"
        ssh -p ${hostInfo[1]} ${hostInfo[2]}@${hostInfo[0]} "echo 'echo never > /sys/kernel/mm/transparent_hugepage/enabled' >>  /etc/rc.local"
        echo "${hostInfo[0]}: ulimit -n 65535"
        ssh -p ${hostInfo[1]} ${hostInfo[2]}@${hostInfo[0]} "ulimit -n 65535"
        echo "${hostInfo[0]}: * soft nofile 65535"
        ssh -p ${hostInfo[1]} ${hostInfo[2]}@${hostInfo[0]} "echo '* soft nofile 65535' >> /etc/security/limits.conf"
        echo "${hostInfo[0]}: * hard nofile 65535"
        ssh -p ${hostInfo[1]} ${hostInfo[2]}@${hostInfo[0]} "echo '* hard nofile 65535' >> /etc/security/limits.conf"
    done

    echo "<< 节点配置优化完毕"
}
function installCM(){
    echo ">> 准备安装 cloudera manager"
    source /etc/profile
    for (( i=0 ; i < ${#slaves[@]} ; i++ )) do
    (
        slave=${slaves[i]}
        hostInfo=(${slave// / })
        echo "${hostInfo[0]}: 安装 cloudera-manager-daemons"
        ssh -p ${hostInfo[1]} ${hostInfo[2]}@${hostInfo[0]} "source /etc/profile;yum localinstall -y $basePackagePath/cm/cloudera-manager-daemons*"
        echo "${hostInfo[0]}: 安装 cloudera-manager-agent"
        ssh -p ${hostInfo[1]} ${hostInfo[2]}@${hostInfo[0]} "source /etc/profile;yum localinstall -y $basePackagePath/cm/cloudera-manager-agent*"
    ) &
    done
    (
        source /etc/profile
        echo "${master[0]}: 安装 cloudera-manager-daemons"
        yum localinstall -y $basePackagePath/cm/cloudera-manager-daemons*
        echo "${master[0]}: 安装 cloudera-manager-agent"
        yum localinstall -y $basePackagePath/cm/cloudera-manager-agent*
        echo "${master[0]}: 安装 cloudera-manager-server"
        yum localinstall -y $basePackagePath/cm/cloudera-manager-server*
    ) &
    wait
    echo ">> cloudera manager 安装完毕"
}
function updateCmSetting(){
    echo ">> 准备更新 cloudera manager 配置"
    ts=`date '+%s'`
    echo "${master[0]}: 备份 /etc/cloudera-scm-agent/config.ini"
    cp /etc/cloudera-scm-agent/config.ini /etc/cloudera-scm-agent/config.ini.bak.$ts
    echo "${master[0]}: 备份 /etc/cloudera-scm-server/db.properties"
    cp /etc/cloudera-scm-server/db.properties /etc/cloudera-scm-server/db.properties.bak.$ts
    echo "${master[0]}: 更新 /etc/cloudera-scm-agent/config.ini"
    sed -i "s/server_host=localhost/server_host=${master[0]}/g" /etc/cloudera-scm-agent/config.ini
    echo "${master[0]}: 更新 /etc/cloudera-scm-server/db.properties"
    sed -i "s/#com.cloudera.cmf.db.host=localhost/com.cloudera.cmf.db.host=${master[0]}:$mysqlport/g" /etc/cloudera-scm-server/db.properties
    sed -i "s/#com.cloudera.cmf.db.name=cmf/com.cloudera.cmf.db.name=cm/g" /etc/cloudera-scm-server/db.properties
    sed -i "s/#com.cloudera.cmf.db.user=cmf/com.cloudera.cmf.db.user=$mysqlcdhuser/g" /etc/cloudera-scm-server/db.properties
    sed -i "s/#com.cloudera.cmf.db.password=/com.cloudera.cmf.db.password=$mysqlcdhpwd/g" /etc/cloudera-scm-server/db.properties
    sed -i "s/com.cloudera.cmf.db.setupType=INIT/com.cloudera.cmf.db.setupType=EXTERNAL/g" /etc/cloudera-scm-server/db.properties
    for (( i=0 ; i < ${#slaves[@]} ; i++ )) do
        slave=${slaves[i]}
        hostInfo=(${slave// / })
        echo "${hostInfo[0]}: 更新 /etc/cloudera-scm-agent/config.ini"
        ssh -p ${hostInfo[1]} ${hostInfo[2]}@${hostInfo[0]} "sed -i \"s/server_host=localhost/server_host=${master[0]}/g\" /etc/cloudera-scm-agent/config.ini"
    done
    echo ">> cloudera manager 配置更新完毕"
}
function startServer(){
    echo "${master[0]}: 启动 cloudera-scm-server"
    systemctl start cloudera-scm-server
    echo "###############################"
    echo "首次启动时间较长,请查看日志观察启动情况: tail -f /var/log/cloudera-scm-server/cloudera-scm-server.log"
    echo "待 cloudera-scm-server 启动成功后再启动 cloudera-scm-agent"
    echo "访问地址: http://${master[0]}:7180"
    echo "默认账号: admin , 密码: admin"
    echo "###############################"
}
function startAgent(){
    echo "${master[0]}: 启动 cloudera-scm-agent"
    systemctl start cloudera-scm-agent
    for (( i=0 ; i < ${#slaves[@]} ; i++ )) do
        slave=${slaves[i]}
        hostInfo=(${slave// / })
        echo "${hostInfo[0]}: 启动 cloudera-scm-agent"
        ssh -p ${hostInfo[1]} ${hostInfo[2]}@${hostInfo[0]} "systemctl start cloudera-scm-agent"
    done
}
function stopServer(){
    echo "${master[0]}: 关闭 cloudera-scm-server"
    systemctl stop cloudera-scm-server
}
function stopAgent(){
    echo "${master[0]}: 关闭 cloudera-scm-agent"
    systemctl stop cloudera-scm-agent
    for (( i=0 ; i < ${#slaves[@]} ; i++ )) do
        slave=${slaves[i]}
        hostInfo=(${slave// / })
        echo "${hostInfo[0]}: 关闭 cloudera-scm-agent"
        ssh -p ${hostInfo[1]} ${hostInfo[2]}@${hostInfo[0]} "systemctl stop cloudera-scm-agent"
    done
}

$1