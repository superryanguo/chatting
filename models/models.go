package models

import (
	"fmt"
	"sync"

	log "github.com/micro/go-micro/v2/logger"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/superryanguo/chatting/basic/config"
)

var (
	inited bool
	mygdb  *gorm.DB
	m      sync.RWMutex
)

const (
	Userid = 11111111
)

type User struct {
	ID           int
	Name         string `gorm:"size:50"  json:"name"`
	PasswordHash string `gorm:"size:256" json:"password"`
	Email        string `gorm:"size:50;unique"  json:"email"`
	RealName     string `gorm:"size:32" json:"real_name"`
	IdCard       string `gorm:"size:20" json:"id_card"`
	AvatarUrl    string `gorm:"size:256" json:"avatar_url"`
	Dialogs      []Dialog
}

type Dialog struct {
	//ID        int    `json:"dialog_id"`
	SessionId string `gorm:"size:256" json:"session_id"`
	UserID    int    `json:"user_id"`
	Path      string `json:"path"`
}

func Init() {
	m.Lock()
	defer m.Unlock()

	var err error

	if inited {
		err = fmt.Errorf("[Init] models already init")
		log.Error(err)
		return
	}

	log.Info("Initing the models.........")
	orm_config := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.GetMysqlConfig().GetName(),
		config.GetMysqlConfig().GetPsw(),
		config.GetMysqlConfig().GetURL(),
		config.GetMysqlConfig().GetDbName())
	log.Debug("connect config=", orm_config)
	db, err := gorm.Open("mysql", orm_config)

	if err != nil {
		panic("failed to connect database")
	}
	//TODO: add a handling to close the db ,such as receive the exit singnal, then call db.close
	//defer db.Close()

	if config.GetMysqlConfig().GetMigrate() {
		DataTableInit(db)
	}

	mygdb = db
	inited = true
	log.Info("Database tables init done")
}
func DataTableInit(db *gorm.DB) {
	log.Debug("gorm automigrate database and init the areas data")

	db.AutoMigrate(&User{}, &Dialog{})

}

func GetGorm() *gorm.DB {
	return mygdb
}
