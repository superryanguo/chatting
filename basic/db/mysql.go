package db

import (
	"database/sql"
	"fmt"
	"time"

	log "github.com/micro/go-micro/v2/logger"
	"github.com/superryanguo/chatting/basic/config"
)

func initMysql() {
	var err error

	connect := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		config.GetMysqlConfig().GetName(),
		config.GetMysqlConfig().GetPsw(),
		config.GetMysqlConfig().GetURL(),
		config.GetMysqlConfig().GetDbName())
	log.Debug("connect sql=", connect)
	mysqlDB, err = sql.Open("mysql", connect)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	mysqlDB.SetMaxOpenConns(config.GetMysqlConfig().GetMaxOpenConnection())

	mysqlDB.SetMaxIdleConns(config.GetMysqlConfig().GetMaxIdleConnection())

	mysqlDB.SetConnMaxLifetime(time.Second * time.Duration(config.GetMysqlConfig().GetConnMaxLifetime()))
	log.Info("Connecting the mysql database, PING...")
	if err = mysqlDB.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Info("PONG... mysql connected successfully...")
}
