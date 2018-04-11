package main

import (
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"net/http"
	"sync"
)

var Conf = InitConfig()
var Mutex sync.Mutex
var RWMutex sync.RWMutex
var Db *leveldb.DB
var CData = make(chan map[string][]string, 100) //查询数据信道
var CMail = make(chan []string, 100)            //发送邮件信道
var CSms = make(chan []string, 100)             //发送短信信道
var CRechargeOK = make(chan bool)
var Path = "C:\\DB"
var reChargeMsg = NewRecharge()
var reChargeResult RechargeResult

func main() {
	var err error
	Db, err = leveldb.OpenFile(Path, nil)
	if err != nil {

		log.Println(err)
	}
	defer Db.Close()
	go monitor()
	go recharge()

	CRechargeOK <- true
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	http.HandleFunc("/index", index)
	http.HandleFunc("/show", show)
	server.ListenAndServe()
}
