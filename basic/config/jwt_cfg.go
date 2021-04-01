package config

type JwtConfig interface {
	GetSecretKey() string
}

type defaultJwtConfig struct {
	SecretKey string `json:"secretKey"`
}

func (m defaultJwtConfig) GetSecretKey() string {
	return m.SecretKey
}
