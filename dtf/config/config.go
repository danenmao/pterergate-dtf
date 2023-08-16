package config

// Redis地址
type RedisAddress struct {
	Name     string `mapstructure:"name" json:"name"`
	Type     string `mapstructure:"type" json:"type"`
	Address  string `mapstructure:"address" json:"address"`
	DB       string `mapstructure:"db" json:"db"`
	Password string `mapstructure:"password" json:"password"`
}

// MySQL地址
type MySQLAddress struct {
	Name     string `mapstructure:"name" json:"name"`
	Type     string `mapstructure:"type" json:"type"`
	Protocol string `mapstructure:"protocol" json:"protocol"`
	Address  string `mapstructure:"address" json:"address"`
	DB       string `mapstructure:"db" json:"db"`
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
}

// Mongo地址
type MongoAddress struct {
	Username   string `mapstructure:"username" json:"username"`
	Password   string `mapstructure:"password" json:"password"`
	Address    string `mapstructure:"address" json:"address"`
	Database   string `mapstructure:"database" json:"database"`
	ReplicaSet string `mapstructure:"replica-set" json:"replica-set"`
}

// RPC service address
type RPCServiceAddress struct {
	Hostname string `json:"hostname"`
	Port     uint16 `json:"port"`
}
