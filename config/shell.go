package config

import (
	_ "embed"
	"strconv"
)

//go:embed cdhinstall.sh
var shell string

var IsRecordInstall = true

func GetShell(cmconfig CminstallerConfig) string {
	s := "#!/bin/bash"
	s += "\nbasePath=" + cmconfig.PathConfig.Base
	s += "\nbasePackagePath=" + cmconfig.PathConfig.Packages
	s += "\nmysqlpwd=" + cmconfig.MysqlConfig.RootPwd
	s += "\nmysqlcdhuser=" + cmconfig.MysqlConfig.CdhUser
	s += "\nmysqlcdhpwd=" + cmconfig.MysqlConfig.CdhUserPwd
	s += "\nmysqlport=" + strconv.Itoa(cmconfig.MysqlConfig.Port)
	s += "\nmaster=(" + cmconfig.HostConfig.Master.Ip + " " + strconv.Itoa(cmconfig.HostConfig.Master.Port) + " " + cmconfig.HostConfig.Master.User + " " + cmconfig.HostConfig.Master.Pwd + " " + cmconfig.HostConfig.Master.Hostname + ")"
	for index, slave := range cmconfig.HostConfig.Slaves {
		s += "\nslaves[" + strconv.Itoa(index) + "]=\"" + slave.Ip + " " + strconv.Itoa(slave.Port) + " " + slave.User + " " + slave.Pwd + " " + slave.Hostname + "\""
	}
	s += "\n\n"
	s += shell
	return s
}
