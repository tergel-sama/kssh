package models

type HostConfig struct {
	Name     string `yaml:"name"`
	Hostname string `yaml:"hostname"`
	User     string `yaml:"user"`
	Port     int    `yaml:"port"`
	Key      string `yaml:"key"`
}

type SSHConfig struct {
	Hosts []HostConfig `yaml:"hosts"`
}
