package action

import (
	"cminstaller/config"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const (
	//上传失败
	UP_FAIL = 0
	//上传成功
	UP_SUCCESS = 1
	//跳过上传
	UP_SKIP = 2
)

//获取ssh客户端
func GetSshClient(hostInfo config.HostInfo) *ssh.Client {
	sshconfig := &ssh.ClientConfig{
		User:            hostInfo.User,
		Auth:            []ssh.AuthMethod{ssh.Password(hostInfo.Pwd)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	addr := hostInfo.Ip + ":" + strconv.Itoa(hostInfo.Port)
	sshClient, errc := ssh.Dial("tcp", addr, sshconfig)
	if nil != errc {
		panic(errc)
	}
	return sshClient
}

//获取sftp客户端
func GetSftpClient(hostInfo config.HostInfo) (*sftp.Client, *ssh.Client) {
	sshClient := GetSshClient(hostInfo)
	sftpClient, errs := sftp.NewClient(sshClient)
	if nil != errs {
		panic(errs)
	}
	return sftpClient, sshClient
}

//校验本地文件和远程文件是否一致
func Sftp_checkSame(sftpClient *sftp.Client, local string, remote string) bool {
	if !Sftp_exists(sftpClient, remote) {
		return false
	}
	fi, _ := os.Lstat(local)
	fi2, _ := sftpClient.Lstat(remote)
	return fi.Size() == fi2.Size()
}

//上传文件
func Sftp_put2(sftpClient *sftp.Client, local string) int {
	remote := local
	return Sftp_put(sftpClient, local, remote)
}

//上传文件
func Sftp_put(sftpClient *sftp.Client, local string, remote string) int {
	//本地文件和远程文件一样则跳过
	if Sftp_checkSame(sftpClient, local, remote) {
		return UP_SKIP
	}
	dir := filepath.Dir(remote)
	if !Sftp_exists(sftpClient, dir) {
		sftpClient.MkdirAll(dir)
	}
	rfile, err := sftpClient.Create(remote)
	if nil != err {
		panic(err)
	}
	defer rfile.Close()
	lfile, err2 := os.Open(local)
	if err2 != nil {
		panic(err2)
	}
	defer lfile.Close()

	_, err3 := io.Copy(rfile, lfile)
	if err3 != nil {
		panic(err3)
	}
	return UP_SUCCESS
}

//检查源文件是否和目标文件一致
func CheckSame(source string, dist string) bool {
	if !Exists(dist) {
		return false
	}
	fi, _ := os.Lstat(source)
	fi2, _ := os.Lstat(dist)
	return fi.Size() == fi2.Size()
}

//将源文件复制到指定路径
func Copy(source string, dist string) int {
	//文件一样则跳过
	if CheckSame(source, dist) {
		return UP_SKIP
	}
	dir := filepath.Dir(dist)
	if !Exists(dir) {
		os.MkdirAll(dir, 0755)
	}
	d, e := os.Create(dist)
	if nil != e {
		panic(e)
	}
	defer d.Close()
	s, e2 := os.Open(source)
	if nil != e2 {
		panic(e2)
	}
	defer s.Close()
	_, e3 := io.Copy(d, s)
	if nil != e3 {
		panic(e3)
	}
	return UP_SUCCESS
}

//列出远程目录下的文件
func Sftp_ls(sftpClient *sftp.Client, dir string) map[string]int {
	nameMap := map[string]int{}
	fis, err := sftpClient.ReadDir(dir)
	if nil != err {
		return nameMap
	}
	for _, fi := range fis {
		nameMap[fi.Name()] = 1
	}
	return nameMap
}

//判断远程文件是否存在
func Sftp_exists(sftpClient *sftp.Client, p string) bool {
	_, err := sftpClient.Lstat(p)
	return nil == err
}

//列出指定目录下的文件
func Ls(dir string) map[string]int {
	nameMap := map[string]int{}
	fis, err := os.ReadDir(dir)
	if nil != err {
		return nameMap
	}
	for _, fi := range fis {
		if fi.IsDir() {
			for k, v := range Ls(dir + "/" + fi.Name()) {
				nameMap[fi.Name()+"/"+k] = v
			}
			continue
		}
		nameMap[fi.Name()] = 1
	}
	return nameMap
}

//列出指定目录下的文件
func Ls2(dir string) []string {
	names := []string{}
	fis, err := os.ReadDir(dir)
	if nil != err {
		return names
	}
	for _, fi := range fis {
		if fi.IsDir() {
			for k := range Ls(dir + "/" + fi.Name()) {
				names = append(names, fi.Name()+"/"+k)
			}
			continue
		}
		names = append(names, fi.Name())
	}
	return names
}

//判断指定文件是否存在
func Exists(p string) bool {
	_, err := os.Lstat(p)
	return nil == err
}
