package action

import (
	"bufio"
	"bytes"
	"cminstaller/config"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

//执行命令获取返回
func RunCommand(name string, params []string) string {
	var buffer bytes.Buffer
	e := exec.Command(name, params...)
	e.Stdout = bufio.NewWriter(&buffer)
	err := e.Run()
	if nil != err {
		panic(err)
	}
	return buffer.String()
}

//执行命令
func sh2(name string, params []string) {
	e := exec.Command(name, params...)
	err := e.Run()
	if nil != err {
		StopProcess(S_FAIL)
		panic(err)
	}
}

//执行命令
func sh(name string, out io.Writer, errio io.Writer) {
	e := exec.Command(InstallRecord.Name, name)
	e.Stderr = errio
	e.Stdout = out
	err := e.Run()
	if nil != err {
		StopProcess(S_FAIL)
		panic(err)
	}
}

//执行命令，记录日志
func shIO(name string, logName string) {
	out, errio := openLog(logName)
	defer out.Close()
	defer errio.Close()
	sh(name, out, errio)
}

//创建日志文件
func openLog(name string) (*os.File, *os.File) {
	logDir := "cmi_" + strconv.FormatInt(InstallRecord.StartTime, 10)
	os.Mkdir(logDir, 0755)
	out, err := os.Create(logDir + "/" + name + ".log")
	if nil != err {
		panic(err)
	}
	errio, err2 := os.Create(logDir + "/" + name + ".err.log")
	if nil != err2 {
		panic(err2)
	}
	return out, errio
}

//安装必备组件
func InstallNeeded_1() {
	StartProcess("开始安装必要组件")
	shIO("installMy", "1")
	StopProcess(S_OK)
	recordInstall(1, S_OK)
}

//生成集群免秘钥
func DoSshkeygen_2() {
	StartProcess("开始配置集群免秘钥")
	shIO("sshkeygen", "2")
	StopProcess(S_OK)
	recordInstall(2, S_OK)
}

//更新集群hosts
func UpdateHosts_3() {
	StartProcess("更新集群hosts配置")
	shIO("updateHosts", "3")
	StopProcess(S_OK)
	recordInstall(3, S_OK)
}

//关闭并禁用集群防火墙
func ShutdownFirewall_4() {
	StartProcess("关闭并禁用集群防火墙")
	shIO("shutdownFirewall", "4")
	StopProcess(S_OK)
	recordInstall(4, S_OK)
}

//优化集群配置
func OptimizeServer_5() {
	StartProcess("优化集群服务器配置")
	shIO("optimizeServer", "5")
	StopProcess(S_OK)
	recordInstall(5, S_OK)
}

//安装mysql数据库
func InstallMysql_6() {
	StartProcess("安装Mysql服务")
	shIO("installMysql", "6")
	StopProcess(S_OK)
	recordInstall(6, S_OK)
}

//安装httpd服务
func InstallHttpd_7() {
	StartProcess("安装httpd服务")
	shIO("installHttpd", "7")
	StopProcess(S_OK)
	recordInstall(7, S_OK)
}

//挂载yum源
func MountYum_8() {
	StartProcess("挂载离线yum源")
	shIO("mountYum", "8")
	StopProcess(S_OK)
	recordInstall(8, S_OK)
}

//更新yum源（go实现）
func PutYumRepo(cmconfig config.CminstallerConfig) {
	repoName := "vrv-centos.repo"
	repoPath := "/etc/yum.repos.d"
	repoFilePath := repoPath + "/" + repoName
	StartProcess("创建" + repoFilePath)
	repoFile, err := os.Create(repoFilePath)
	if nil != err {
		panic(err)
	}
	defer repoFile.Close()
	text := `[vrv-centos]
name=vrv-centos
baseurl=http://` + cmconfig.HostConfig.Master.Ip + `/centos
gpgkey=http://` + cmconfig.HostConfig.Master.Ip + `/centos/RPM-GPG-KEY-CentOS-7
gpgcheck=1
enabled=1
priority=1
`
	repoFile.WriteString(text)
	StopProcess(S_FINISH)

	var wg sync.WaitGroup
	logs := []string{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		sh2("yum", []string{"clean", "all"})
		sh2("yum", []string{"repolist"})
	}()
	for _, v := range cmconfig.HostConfig.Slaves {
		sftpClient, sshClient := GetSftpClient(v)
		defer sftpClient.Close()
		defer sshClient.Close()

		StartProcess("上传" + repoName + "到" + v.Ip + ":/" + repoPath)
		Sftp_put2(sftpClient, repoFilePath)
		StopProcess(S_FINISH)
		ip := v.Ip
		wg.Add(1)
		go func() {
			defer wg.Done()
			session, err := sshClient.NewSession()
			if nil != err {
				panic(err)
			}
			defer session.Close()
			out, errio := openLog("9_" + ip)
			defer out.Close()
			defer errio.Close()
			logs = append(logs, errio.Name())
			session.Stderr = errio
			session.Stdout = out
			e1 := session.Run("yum clean all;yum repolist")
			if nil != e1 {
				panic(e1)
			}
		}()
	}
	StartProcess("等待yum更新完毕")
	wg.Wait()
	fail := false
	failLog := ""
	for _, name := range logs {
		content, _ := os.ReadFile(name)
		scontent := string(content)
		if strings.Contains(scontent, "Timeout on") || strings.Contains(scontent, "repolist: 0") {
			fail = true
			failLog = name
			break
		}
	}
	if fail {
		StopProcess(S_FAIL)
		panic("yum配置失败,请查看 " + failLog)
	} else {
		StopProcess(S_OK)
	}
	recordInstall(9, S_OK)
}

//更新yum源
func PutYumRepo_9() {
	StartProcess("更新集群yum源配置")
	shIO("putYumRepo", "9")
	StopProcess(S_OK)
	recordInstall(9, S_OK)
}

//安装jdk（go实现）
func InstallJDK(cmconfig config.CminstallerConfig) {
	sourcePath := cmconfig.GetPackageDir("jdk")
	var sourceFileName string
	for k := range Ls(sourcePath) {
		sourceFileName = k
	}
	//上传安装包到远程服务器
	dist := sourcePath + "/" + sourceFileName
	for _, v := range cmconfig.HostConfig.Slaves {
		StartProcess("上传" + sourceFileName + "到" + v.Ip + ":/" + sourcePath)
		sftpClient, sshClient := GetSftpClient(v)
		defer sftpClient.Close()
		defer sshClient.Close()
		if UP_SKIP == Sftp_put2(sftpClient, dist) {
			StopProcess(S_SKIP)
		} else {
			StopProcess(S_FINISH)
		}
	}
	StartProcess("安装java环境")
	shIO("installJDK2", "10")
	StopProcess(S_OK)
	recordInstall(10, S_OK)
}

//安装jdk
func InstallJDK_10() {
	StartProcess("安装java环境")
	shIO("installJDK", "10")
	StopProcess(S_OK)
	recordInstall(10, S_OK)
}

//上传安装包
func UploadCDH_11() {
	StartProcess("上传CDH安装包")
	shIO("uploadCDH", "11")
	StopProcess(S_OK)
	recordInstall(11, S_OK)
}

//上传安装包（go实现）
func UploadCDH(cmconfig config.CminstallerConfig) {
	distRepoPath := "/opt/cloudera/parcel-repo"
	if !Exists(distRepoPath) {
		fmt.Println(">  创建目录" + distRepoPath)
		os.MkdirAll(distRepoPath, 0755)
	}
	//获取已经存在的文件
	sourceRepoPath := cmconfig.GetPackageDir("parcels")
	sourceRepo := Ls(sourceRepoPath)
	for name := range sourceRepo {
		StartProcess("复制" + name + "到" + distRepoPath)
		if UP_SKIP == Copy(sourceRepoPath+"/"+name, distRepoPath+"/"+name) {
			StopProcess(S_SKIP)
		} else {
			StopProcess(S_FINISH)
		}
	}

	cmpath := cmconfig.GetPackageDir("cm")
	localRpm, err1 := os.ReadDir(cmpath)
	if nil != err1 {
		panic(err1)
	}
	//获取安装包列表
	rpms := []string{}
	for _, r := range localRpm {
		rpms = append(rpms, r.Name())
	}
	//上传安装包到远程服务器
	for _, v := range cmconfig.HostConfig.Slaves {
		fmt.Printf(">  准备上传安装包到%s\n", v.Ip)
		sftpClient, sshClient := GetSftpClient(v)
		defer sftpClient.Close()
		defer sshClient.Close()
		for _, name := range rpms {
			StartProcess("上传" + name + "到" + v.Ip + ":/" + cmpath)
			if UP_SKIP == Sftp_put2(sftpClient, cmpath+"/"+name) {
				StopProcess(S_SKIP)
			} else {
				StopProcess(S_FINISH)
			}
		}
	}
	recordInstall(11, S_OK)
}

//安装cm
func InstallCM_12() {
	StartProcess("安装Cloudera Manager")
	shIO("installCM", "12")
	StopProcess(S_OK)
	recordInstall(12, S_OK)
}

//配置cm数据库
func SetupMysql_13() {
	StartProcess("配置Mysql数据库")
	shIO("setupMysql", "13")
	StopProcess(S_OK)
	recordInstall(13, S_OK)
}

//更新cm配置
func UpdateCmSetting_14() {
	StartProcess("更新Cloudera Manager配置")
	shIO("updateCmSetting", "14")
	StopProcess(S_OK)
	recordInstall(14, S_OK)
}

//启动cm server
func StartServer_15() {
	StartProcess("启动CM-Server")
	shIO("startServer", "15")
	StopProcess(S_OK)
	recordInstall(15, S_OK)
}

//启动cm agent
func StartAgent_16() {
	StartProcess("启动CM-Agent")
	shIO("startAgent", "16")
	StopProcess(S_OK)
	recordInstall(16, S_OK)
}

//关闭cm server
func StopServer_17() {
	StartProcess("关闭CM-Server")
	shIO("stopServer", "17")
	StopProcess(S_OK)
	recordInstall(17, S_OK)
}

//关闭cm agent
func StopAgent_18() {
	StartProcess("关闭CM-Agent")
	shIO("stopAgent", "18")
	StopProcess(S_OK)
	recordInstall(18, S_OK)
}

//启动cm server
func StartServer() {
	StartProcess("启动CM-Server")
	sh("startServer", os.Stdout, os.Stderr)
	StopProcess(S_FINISH)
}

//启动cm agent
func StartAgent() {
	StartProcess("启动CM-Agent")
	sh("startAgent", os.Stdout, os.Stderr)
	StopProcess(S_FINISH)
}

//关闭cm server
func StopServer() {
	StartProcess("关闭CM-Server")
	sh("stopServer", os.Stdout, os.Stderr)
	StopProcess(S_FINISH)
}

//关闭cm agent
func StopAgent() {
	StartProcess("关闭CM-Agent")
	sh("stopAgent", os.Stdout, os.Stderr)
	StopProcess(S_FINISH)
}
