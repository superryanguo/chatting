package config

type MysqlConfig interface {
	GetURL() string
	GetName() string
	GetPsw() string
	GetDbName() string
	GetEnabled() bool
	GetMigrate() bool
	GetMaxIdleConnection() int
	GetMaxOpenConnection() int
	GetConnMaxLifetime() int
}

type defaultMysqlConfig struct {
	URL               string `json:"url"`
	Name              string `json:"name"`
	Psw               string `json:"psw"`
	DbName            string `json:"dbname"`
	Enable            bool   `json:"enabled"`
	Migrate           bool   `json:"migrate"`
	MaxIdleConnection int    `json:"maxIdleConnection"`
	MaxOpenConnection int    `json:"maxOpenConnection"`
	ConnMaxLifetime   int    `json:"connMaxLifetime"`
}

func (m defaultMysqlConfig) GetURL() string {
	return m.URL
}

func (m defaultMysqlConfig) GetName() string {
	return m.Name
}
func (m defaultMysqlConfig) GetPsw() string {
	return m.Psw
}

func (m defaultMysqlConfig) GetDbName() string {
	return m.DbName
}

func (m defaultMysqlConfig) GetEnabled() bool {
	return m.Enable
}

func (m defaultMysqlConfig) GetMigrate() bool {
	return m.Migrate
}

func (m defaultMysqlConfig) GetMaxIdleConnection() int {
	return m.MaxIdleConnection
}

func (m defaultMysqlConfig) GetMaxOpenConnection() int {
	return m.MaxOpenConnection
}

func (m defaultMysqlConfig) GetConnMaxLifetime() int {
	return m.ConnMaxLifetime
}
