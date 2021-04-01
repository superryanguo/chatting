package config

type MisConfig interface {
	GetImageAddr() string
	GetMailUser() string
	GetMailPass() string
}

type defaultMiscConfig struct {
	ImageAddr string `json:"imageaddr"`
	MailUser  string `json:"mailuser"`
	MailPass  string `json:"mailpass"`
}

func (d defaultMiscConfig) GetImageAddr() string {
	return d.ImageAddr
}

func (d defaultMiscConfig) GetMailUser() string {
	return d.MailUser
}

func (d defaultMiscConfig) GetMailPass() string {
	return d.MailPass
}
