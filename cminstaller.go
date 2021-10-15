package main

import (
	"bufio"
	"cminstaller/config"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v3"

	"io/ioutil"

	"cminstaller/action"
)

var (
	configFile  = flag.String("config", "cminstall.yml", "配置文件")
	showInfo    = flag.Bool("info", false, "显示连接信息")
	startServer = flag.Bool("startServer", false, "启动CM-Server")
	startAgent  = flag.Bool("startAgent", false, "启动CM-Agent")
	stopServer  = flag.Bool("stopServer", false, "关闭CM-Server")
	stopAgent   = flag.Bool("stopAgent", false, "关闭CM-Agent")
	myStep      = flag.Bool("myStep", false, "指定运行步骤")
	step        = flag.Int("step", -1, "步骤编号")

	native = flag.Bool("native", true, "启用本地实现")

	nocheck = flag.Bool("nocheck", false, "跳过所有确认项")
)

const (
	DIR_CM      = "cm"
	DIR_HTTPD   = "httpd"
	DIR_ISO     = "iso"
	DIR_JDK     = "jdk"
	DIR_MY      = "my"
	DIR_MYSQL   = "mysql"
	DIR_PARCELS = "parcels"
)

func main() {
	flag.Parse()
	cmconfig, err := loadConfig()
	if !*nocheck && !*myStep {
		if !checkPackage(cmconfig) {
			fmt.Println("安装包检验错误,请联系管理员.")
			return
		}
	}
	//指定步骤时不记录安装步骤
	config.IsRecordInstall = !*myStep
	if nil != err {
		panic(err)
	}
	if *showInfo {
		printInfo(cmconfig)
		return
	}
	action.PrepareInstall(cmconfig, *myStep)
	if *startServer {
		action.StartServer()
		return
	}
	if *startAgent {
		action.StartAgent()
		return
	}
	if *stopServer {
		action.StopServer()
		return
	}
	if *stopAgent {
		action.StopAgent()
		return
	}
	record := action.InstallRecord
	fmt.Println(">  启用脚本" + filepath.Base(record.Name))
	last := 1
	if *myStep {
		last = *step
	}
	if !*myStep && record.Continue && len(record.Progress) > 0 {
		r := record.Progress[len(record.Progress)-1]
		l := strings.Split(r, "#")
		last, _ = strconv.Atoi(l[0])
		last += 1
		if last > 255 {
			last = 255
		}
		fmt.Printf(">  继续执行: %d\n", last)
	}
	switch last {
	case 1:
		action.InstallNeeded_1()
		if *myStep {
			break
		}
		fallthrough
	case 2:
		action.DoSshkeygen_2()
		if *myStep {
			break
		}
		fallthrough
	case 3:
		action.UpdateHosts_3()
		if *myStep {
			break
		}
		fallthrough
	case 4:
		action.ShutdownFirewall_4()
		if *myStep {
			break
		}
		fallthrough
	case 5:
		action.OptimizeServer_5()
		if *myStep {
			break
		}
		fallthrough
	case 6:
		action.InstallMysql_6()
		if *myStep {
			break
		}
		fallthrough
	case 7:
		action.InstallHttpd_7()
		if *myStep {
			break
		}
		fallthrough
	case 8:
		action.MountYum_8()
		if *myStep {
			break
		}
		fallthrough
	case 9:
		if *native {
			action.PutYumRepo(cmconfig)
		} else {
			action.PutYumRepo_9()
		}
		if *myStep {
			break
		}
		fallthrough
	case 10:
		if *native {
			action.InstallJDK(cmconfig)
		} else {
			action.InstallJDK_10()
		}
		if *myStep {
			break
		}
		fallthrough
	case 11:
		if *native {
			action.UploadCDH(cmconfig)
		} else {
			action.UploadCDH_11()
		}
		if *myStep {
			break
		}
		fallthrough
	case 12:
		if *native {
			action.InstallCM(cmconfig)
		} else {
			action.InstallCM_12()
		}
		if *myStep {
			break
		}
		fallthrough
	case 13:
		action.SetupMysql_13()
		if *myStep {
			break
		}
		fallthrough
	case 14:
		action.UpdateCmSetting_14()
		if *myStep {
			break
		}
		fallthrough
	case 15:
		action.StartServer_15()
		if *myStep {
			break
		}
		fallthrough
	case 16:
		action.StartAgent_16()
	case 255:
		fmt.Println(">  已完成安装,请删除记录后重试.")
		return
	default:
		fmt.Printf("未知的操作步骤#%d\n", last)
		return
	}
	if !*myStep {
		action.FinishInstall()
		fmt.Println("")
		printInfo(cmconfig)
	}
}

func loadConfig() (config.CminstallerConfig, error) {
	action.StartProcess("读取配置文件")
	var cminstallerConfig config.CminstallerConfig
	yamlFile, err := ioutil.ReadFile(*configFile)
	if nil != err {
		return cminstallerConfig, err
	}
	yaml.Unmarshal(yamlFile, &cminstallerConfig)
	action.StopProcess("OK")
	return cminstallerConfig, nil
}

func printInfo(cmconfig config.CminstallerConfig) {
	fmt.Println("#########################################################")
	fmt.Printf("# Cloudera Manager访问地址: http://%s:7180\n", cmconfig.HostConfig.Master.Ip)
	fmt.Printf("# Cloudera Manager默认账号密码: %s/%s\n", "admin", "admin")
	fmt.Printf("# Cloudera Manager镜像地址: http://%s/cdh6_parcel\n", cmconfig.HostConfig.Master.Ip)
	fmt.Printf("# Centos镜像地址: http://%s/centos\n", cmconfig.HostConfig.Master.Ip)
	fmt.Printf("# 数据库访问地址: %s\n", cmconfig.HostConfig.Master.Ip)
	fmt.Printf("# 数据库访问端口: %d\n", cmconfig.MysqlConfig.Port)
	fmt.Printf("# 数据库访问账号: %s\n", cmconfig.MysqlConfig.CdhUser)
	fmt.Printf("# 数据库访问密码: %s\n", cmconfig.MysqlConfig.CdhUserPwd)
	fmt.Println("#########################################################")
}

func checkPackage(cmconfig config.CminstallerConfig) bool {
	version := action.RunCommand("cat", []string{"/etc/redhat-release"})
	version = strings.Replace(version, "\n", "", -1)
	iso := action.Ls2(cmconfig.GetPackageDir(DIR_ISO))
	if len(iso) != 1 {
		return false
	}
	isoName := iso[0]
	y := getInput("系统版本: " + version + "\n镜像版本: " + isoName + "\n是否一致? (y/n)")
	return strings.ToLower(y) == "y"
}

func getInput(msg string) string {
	fmt.Print(msg + " ")
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	in := sc.Text()
	return strings.Replace(in, "\n", "", -1)
}
