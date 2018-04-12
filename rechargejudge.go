package main

import (
	"fmt"
	"log"
	//"strconv"
	"errors"
	"strings"
	"time"

	//"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func rechargeJudge(uid string, coin string) error {

	seekStr := uid + "_" + coin + "_" //前缀查询数据
	rechargeTimes := 0

	iter := Db.NewIterator(util.BytesPrefix([]byte(seekStr)), nil)
	if !iter.Last() { //迭代器指向容器末尾
		return nil
	}
	iter.Next()
	for iter.Prev() { //从字典序大到小
		key := iter.Key()
		value := iter.Value()
		if fmt.Sprintf("%s", value) == "true" {
			keySlice := strings.Split(fmt.Sprintf("%s", key), "_")

			rechargeTime, err := time.Parse("2006-01-02 15:04:05", keySlice[2]) //获得写入数据的时间
			if err != nil {
				log.Println(err)
				continue
			}

			if time.Since(rechargeTime) >= 24*time.Hour {
				break
			}
			rechargeTimes++
			log.Println(rechargeTimes)
			if rechargeTimes >= Conf.Charges[coin].MaxRechargeTimes {
				//log.Println("over max recharge times")
				return errors.New("over max recharge times")
			}

		}
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func recharge() {
	for {

		var needRecharge map[string][]string = <-CData //信道接收判断阈值传递的map
		//var needRecharge map[string][]string = monitor() //监听

		for uid, coinSlice := range needRecharge {
			for _, coin := range coinSlice {
				RWMutex.Lock()
				err := rechargeJudge(uid, coin)
				if err != nil {
					log.Println(err)
					continue
				}
				RWMutex.Unlock()

				reChargeResult, err := reChargeMsg.RechargeInbank(uid, coin) //充值

				if reChargeResult.Code != "0" {
					log.Println("recharge failed:", err)
					continue
				}
				log.Println("recharge succeeded")
			}
			time.Sleep(time.Minute)
		}
		CRechargeOK <- true
	}
}
