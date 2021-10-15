package config

type CminstallerRecord struct {
	StartTime int64    `yaml:"start-time"`
	Name      string   `yaml:"name"`
	Progress  []string `yaml:"progress"`
	Continue  bool
}

type CminstallerConfig struct {
	PathConfig  Path  `yaml:"path"`
	MysqlConfig Mysql `yaml:"mysql"`
	HostConfig  Host  `yaml:"host"`
}

type Path struct {
	Base     string `yaml:"base"`
	Packages string `yaml:"package"`
}

type Mysql struct {
	Port       int    `yaml:"port"`
	RootPwd    string `yaml:"root-pwd"`
	CdhUser    string `yaml:"cdh-user"`
	CdhUserPwd string `yaml:"cdh-user-pwd"`
}

type Host struct {
	Master HostInfo   `yaml:"master"`
	Slaves []HostInfo `yaml:"slaves"`
}

type HostInfo struct {
	Hostname string `yaml:"hostname"`
	Ip       string `yaml:"ip"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Pwd      string `yaml:"pwd"`
}

func (cmconfig *CminstallerConfig) GetPackageDir(name string) string {
	return cmconfig.PathConfig.Packages + "/" + name
}
