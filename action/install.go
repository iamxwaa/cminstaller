package action

import (
	"cminstaller/config"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"time"

	yaml "gopkg.in/yaml.v3"

	"github.com/briandowns/spinner"
)

const (
	//执行成功
	S_OK = "OK"
	//执行失败
	S_FAIL = "FAIL"
	//执行完毕
	S_FINISH = "FINISH"
	//跳过执行
	S_SKIP = "[SKIP]"
	//安装记录文件
	RECORD_FILE = "cminstaller-record.yml"
)

var (
	//进度条对象
	s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	//脚本保存路径
	ShellPath string
	//安装配置
	InstallRecord = config.CminstallerRecord{Continue: false}
	//安装开始时间
	startTime int64
)

func addX(path string) {
	err := exec.Command("chmod", "755", path).Run()
	if nil != err {
		panic(err)
	}
}

func recordInstall(num int, state string) {
	if !config.IsRecordInstall {
		return
	}
	InstallRecord.Progress = append(InstallRecord.Progress, strconv.Itoa(num)+"#"+state)
	out, _ := yaml.Marshal(InstallRecord)
	ioutil.WriteFile(RECORD_FILE, out, 0755)
}

//结束安装并记录安装结果
func FinishInstall() {
	StartProcess("CDH 安装完毕")
	recordInstall(255, S_OK)
	StopProcess("")
}

//准备安装
func PrepareInstall(cmconfig config.CminstallerConfig, myStep bool) {
	recFile, err := ioutil.ReadFile(RECORD_FILE)
	if nil == err || os.IsExist(err) {
		if !myStep {
			StartProcess("读取上次安装记录")
		}
		yaml.Unmarshal(recFile, &InstallRecord)
		InstallRecord.Continue = true
		_, err2 := os.Open(InstallRecord.Name)
		if nil == err2 || os.IsExist(err2) {
			//直接返回，用上次生成的脚本
			if !myStep {
				StopProcess("(继续执行上次未完成的安装)")
			}
			return
		}
		if !myStep {
			StopProcess("(重新初始化安装配置)")
		}
	}
	StartProcess("初始化安装配置")
	shell := config.GetShell(cmconfig)
	temp, err := ioutil.TempFile("", "cminstaller_*")
	if nil != err {
		panic(err)
	}
	defer temp.Close()
	temp.WriteString(shell)
	ShellPath = temp.Name()
	InstallRecord.Name = ShellPath
	InstallRecord.StartTime = time.Now().UTC().Unix()
	addX(ShellPath)
	StopProcess(S_OK)
	recordInstall(0, S_OK)
}

//启动进度条
func StartProcess(msg string) {
	s.Suffix = "  " + msg
	s.Start()
	startTime = time.Now().UTC().Unix()
	time.Sleep(300 * time.Millisecond)
}

//停止进度条
func StopProcess(msg string) {
	cost := time.Now().UTC().Unix() - startTime
	if msg != "" {
		s.FinalMSG = strconv.FormatInt(cost, 10) + "s " + msg
	}
	s.Stop()
	fmt.Println("")
}
