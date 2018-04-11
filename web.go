package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type Asset struct {
	Account           string
	InformationAsset  map[string]float64
	InformationRecord []Record
}

type Record struct {
	CoinName string
	Time     string
	Number   string
}

func index(w http.ResponseWriter, r *http.Request) { // /index处理器
	t := template.New("index.html")
	t, _ = t.ParseFiles("index.html")
	t.ExecuteTemplate(w, "index", "")
}

func PraseRecord() []Record {

	informationrecord := make([]Record, 0)
	iter := Db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		all := strings.Split(string(key), "_")
		if string(value) == "true" {
			// do something for the key
			record := Record{CoinName: all[1], Time: all[2], Number: all[3]}
			informationrecord = append(informationrecord, record)
		}
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		log.Println(err)
	}
	return informationrecord
}

func show(w http.ResponseWriter, r *http.Request) { // /show页面
	// 将r传过来的值与数据库进行比对，若有该账号，才跳转，否则失败

	if r.FormValue("account") == "dd" {

		//读取账户信息
		body := Post(`{"uid":"9c885ce7664ed3c3675df8e49590c141665d94add5bf72a816a8d8a78bd8fbe5"}`)
		var informationasset = make(map[string]float64)
		informationasset = ParsingWeb(body)

		//排序
		RWMutex.Lock()
		var informationrecord = PraseRecord()
		RWMutex.Unlock()

		data := Asset{Account: "9c885ce7664ed3c3675df8e49590c141665d94add5bf72a816a8d8a78bd8fbe5", InformationAsset: informationasset, InformationRecord: informationrecord}

		t := template.New("show.html")
		t, _ = t.ParseFiles("show.html")
		t.ExecuteTemplate(w, "show", data)
	} else {
		t := template.New("show_.html")
		t, _ = t.ParseFiles("show_.html")
		t.ExecuteTemplate(w, "show_", "")
	}
	fmt.Println("Done")
}
