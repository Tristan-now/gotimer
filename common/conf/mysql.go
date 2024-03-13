package conf

type MysqlConfig struct {
	DSN string `yaml:"DSN"`
	// 最大连接数
	MaxOpenConns int `yaml:"MaxOpenConns"`
	// 最大空闲连接数(连接池中已经建立但并没有使用的连接)
	MaxIdleCoons int `yaml:"MaxIdleCoons"`
}

type MysqlConfProvider struct {
	conf *MysqlConfig
}

func NewMysqlConfigProvider(conf *MysqlConfig) *MysqlConfProvider {
	return &MysqlConfProvider{
		conf: conf,
	}
}

func (m *MysqlConfProvider) Get() *MysqlConfig {
	return m.conf
}

var defaultMysqlConfProvider *MysqlConfProvider

func DefaultMysqlConfigProvider() *MysqlConfProvider {
	return defaultMysqlConfProvider
}
