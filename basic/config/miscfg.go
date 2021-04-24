package config

type MisConfig interface {
	GetImageAddr() string
	GetMailUser() string
	GetMailPass() string
	GetDialogPath() string
	GetDialogPrefix() string
}

type defaultMiscConfig struct {
	ImageAddr    string `json:"imageaddr"`
	MailUser     string `json:"mailuser"`
	MailPass     string `json:"mailpass"`
	DialogPath   string `json:"dialogpath"`
	DialogPrefix string `json:"dialogprefix"`
}

func (d defaultMiscConfig) GetImageAddr() string {
	return d.ImageAddr
}

func (d defaultMiscConfig) GetDialogPrefix() string {
	return d.DialogPrefix
}

func (d defaultMiscConfig) GetDialogPath() string {
	return d.DialogPath
}

func (d defaultMiscConfig) GetMailUser() string {
	return d.MailUser
}

func (d defaultMiscConfig) GetMailPass() string {
	return d.MailPass
}
