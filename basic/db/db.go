package db

import (
	"database/sql"
	"fmt"
	"sync"

	log "github.com/micro/go-micro/v2/logger"
	"github.com/superryanguo/chatting/basic/config"
)

var (
	inited  bool
	mysqlDB *sql.DB
	m       sync.RWMutex
)

func Init() {
	m.Lock()
	defer m.Unlock()

	var err error

	if inited {
		err = fmt.Errorf("[Init] db already init")
		log.Error(err)
		return
	}

	if config.GetMysqlConfig().GetEnabled() {
		initMysql()
	}

	inited = true
}

func GetDB() *sql.DB {
	return mysqlDB
}
