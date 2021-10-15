
function installMy(){
    echo ">> 开始安装必要组件"
    sleep 2s
    echo "<< 必要组件安装完毕"
}
function sshkeygen(){
    echo ">> 开始配置服务器免秘钥"
    sleep 2s
    echo "<< 服务器免秘钥配置完毕"
}
function updateHosts(){
    echo ">> 开始更新各节点hosts"
    sleep 2s
    echo "<< hosts更新完毕"
}
function shutdownFirewall(){
    echo "${master[0]}: 关闭和禁用防火墙"
    sleep 2s
}
function installJDK(){
    echo ">> 准备安装java运行环境"
    sleep 2s
    echo ">> java运行环境安装完毕"
}
function installMysql(){
    echo ">> 准备在${master[0]}安装 mysql"
    sleep 2s
    echo "<< mysql 安装完毕"
}
function uninstallMysql(){
    echo ">> 准备在${master[0]}卸载 mysql"
    sleep 2s
    echo ">> mysql 卸载完毕"
}
function installHttpd(){
    echo ">> 准备在本机安装 httpd"
    sleep 2s
    echo "<< httpd 安装完毕"
}
function mountYum(){
    echo "开始挂载centos镜像"
    sleep 2s
}
function putYumRepo(){
    echo ">> 准备配置离线yum仓库"
    sleep 2s
    echo "<< 离线yum仓库配置完毕"
}
function uploadCDH(){
    echo ">> 开始上传CDH安装包到各节点"
    sleep 2s
    echo "<< CDH安装包上传完毕"
}
function setupMysql(){
    echo ">> 开始配置mysql数据库"
    sleep 2s
    echo "<< mysql 数据库配置完毕"
}
function optimizeServer(){
    echo ">> 开始优化各节点配置"
    sleep 2s
    echo "<< 节点配置优化完毕"
}
function installCM(){
    echo ">> 准备安装 cloudera manager"
    sleep 2s
    echo ">> cloudera manager 安装完毕"
}
function updateCmSetting(){
    echo ">> 准备更新 cloudera manager 配置"
    sleep 2s
    echo ">> cloudera manager 配置更新完毕"
}
function startServer(){
    echo "${master[0]}: 启动 cloudera-scm-server"
    sleep 2s
}
function startAgent(){
    echo "${master[0]}: 启动 cloudera-scm-agent"
    sleep 2s
}
function stopServer(){
    echo "${master[0]}: 关闭 cloudera-scm-server"
    sleep 2s
}
function stopAgent(){
    echo "${master[0]}: 关闭 cloudera-scm-agent"
    sleep 2s
}

$1