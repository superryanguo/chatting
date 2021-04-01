package basic

import (
	"github.com/superryanguo/chatting/basic/cache"
	"github.com/superryanguo/chatting/basic/config"
	"github.com/superryanguo/chatting/basic/db"
	"github.com/superryanguo/chatting/basic/rediser"
)

func Init() {
	config.Init()
	rediser.Init()
	db.Init()
	cache.Init()
}
