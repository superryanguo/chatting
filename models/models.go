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

type User struct {
	Uid           string `gorm:"primary_key;size:256" json:"user_id"`
	Name          string `gorm:"size:50"  json:"name"`
	Password_hash string `gorm:"size:256" json:"password"`
	Email         string `gorm:"size:50;unique"  json:"email"`
	Real_name     string `gorm:"size:32" json:"real_name"`
	Id_card       string `gorm:"size:20" json:"id_card"`
	Avatar_url    string `gorm:"size:256" json:"avatar_url"`
}

type Talklog struct {
	//gorm.Model
	ID      int    `json:"talklog_id"`
	UserID  uint   `json:"user_id"`
	Title   string `gorm:"size:64" json:"title"`
	Datalog string `json:"datalog"`
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

	db.AutoMigrate(&User{}, &Talklog{})

}

func GetGorm() *gorm.DB {
	return mygdb
}
