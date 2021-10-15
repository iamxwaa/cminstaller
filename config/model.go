package config

//安装记录
type CminstallerRecord struct {
	//安装时间
	StartTime int64 `yaml:"start-time"`
	//安装脚本
	Name string `yaml:"name"`
	//安装进度
	Progress []string `yaml:"progress"`
	//是否继续上次步骤
	Continue bool
}

//安装配置
type CminstallerConfig struct {
	//路径配置
	PathConfig Path `yaml:"path"`
	//数据库配置
	MysqlConfig Mysql `yaml:"mysql"`
	//主机配置
	HostConfig Host `yaml:"host"`
}

//路径
type Path struct {
	//根目录
	Base string `yaml:"base"`
	//软件存放目录
	Packages string `yaml:"package"`
}

//数据库
type Mysql struct {
	//端口
	Port int `yaml:"port"`
	//root账号密码
	RootPwd string `yaml:"root-pwd"`
	//cdh用户名
	CdhUser string `yaml:"cdh-user"`
	//cdh用户密码
	CdhUserPwd string `yaml:"cdh-user-pwd"`
}

//主机列表
type Host struct {
	//主节点
	Master HostInfo `yaml:"master"`
	//从节点
	Slaves []HostInfo `yaml:"slaves"`
}

//主机信息
type HostInfo struct {
	//主机名
	Hostname string `yaml:"hostname"`
	//主机IP
	Ip string `yaml:"ip"`
	//主机ssh端口
	Port int `yaml:"port"`
	//主机ssh用户
	User string `yaml:"user"`
	//主机ssh用户密码
	Pwd string `yaml:"pwd"`
}

//获取安装包目录
func (cmconfig *CminstallerConfig) GetPackageDir(name string) string {
	return cmconfig.PathConfig.Packages + "/" + name
}
