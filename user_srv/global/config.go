package global

type AppConfig struct {
	Mysql  Mysql  `mapstructure:"mysql" yaml:"mysql"`
	Grpc   Grpc   `mapstructure:"grpc" yaml:"grpc"`
	Consul Consul `mapstructure:"consul" yaml:"consul"`
}
type Nacos struct {
	Host      string `json:"host" mapstructure:"host" yaml:"host"`
	Port      uint64 `json:"port" mapstructure:"port" yaml:"port"`
	NameSpace string `json:"nameSpace" mapstructure:"nameSpace" yaml:"nameSpace"`
	User      string `json:"user" mapstructure:"user" yaml:"user"`
	Password  string `json:"password" mapstructure:"password" yaml:"password"`
	DataID    string `json:"dataID" mapstructure:"dataID" yaml:"dataID"`
	Group     string `json:"group" mapstructure:"group" yaml:"group"`
}
type Consul struct {
	IP   string `json:"ip" mapstructure:"ip" yaml:"ip"`
	Port int    `json:"port" mapstructure:"port" yaml:"port"`
}
type Grpc struct {
	IP   string   `json:"ip" mapstructure:"ip" yaml:"ip"`
	Port int      `json:"port" mapstructure:"port" yaml:"port"`
	Name string   `json:"name" mapstructure:"name" yaml:"name"`
	Tags []string `json:"tags" mapstructure:"tags" yaml:"tags"`
}
type Mysql struct {
	IP       string `json:"ip" mapstructure:"ip" yaml:"ip"`
	Port     int    `json:"port" mapstructure:"port" yaml:"port"`
	DbName   string `json:"dbName" mapstructure:"dbName" yaml:"dbName"`
	User     string `json:"user" mapstructure:"user" yaml:"user"`
	Password string `json:"password" mapstructure:"password" yaml:"password"`
}
