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
	UP_FAIL    = 0
	UP_SUCCESS = 1
	UP_SKIP    = 2
)

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

func GetSftpClient(hostInfo config.HostInfo) (*sftp.Client, *ssh.Client) {
	sshClient := GetSshClient(hostInfo)
	sftpClient, errs := sftp.NewClient(sshClient)
	if nil != errs {
		panic(errs)
	}
	return sftpClient, sshClient
}

func Sftp_checkSame(sftpClient *sftp.Client, local string, remote string) bool {
	if !Sftp_exists(sftpClient, remote) {
		return false
	}
	fi, _ := os.Lstat(local)
	fi2, _ := sftpClient.Lstat(remote)
	return fi.Size() == fi2.Size()
}

func Sftp_put2(sftpClient *sftp.Client, local string) int {
	remote := local
	return Sftp_put(sftpClient, local, remote)
}

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

func CheckSame(source string, dist string) bool {
	if !Exists(dist) {
		return false
	}
	fi, _ := os.Lstat(source)
	fi2, _ := os.Lstat(dist)
	return fi.Size() == fi2.Size()
}

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

func Sftp_exists(sftpClient *sftp.Client, p string) bool {
	_, err := sftpClient.Lstat(p)
	return nil == err
}

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

func Exists(p string) bool {
	_, err := os.Lstat(p)
	return nil == err
}
