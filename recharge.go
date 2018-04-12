package main

import (
	"bytes"
	"encoding/json"
	//"github.com/syndtr/goleveldb/leveldb"
	"io/ioutil"
	"net/http"
	// "sync"
	"fmt"
	"time"
)

type RechargeMessage struct {
	Symbol string `json:"symbol"`
	Amount string `json:"amount"`
}

type RechargeResult struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

type Recharge struct {
}

func NewRecharge() *Recharge {
	r := &Recharge{}
	return r
}

func (r *Recharge) RechargeInbank(uid string, symbol string) (*RechargeResult, error) {
	amount := fmt.Sprintf("%d", Conf.Charges[symbol].RechargeAmount)
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	reqUrl := Conf.Api.Transfer

	tRechargeMessage := &RechargeMessage{}
	tRechargeMessage.Symbol = symbol
	//POST请求
	tRechargeMessageBytes, err := json.Marshal(tRechargeMessage)
	if nil != err {
		return nil, err
	}

	req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(tRechargeMessageBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if nil != resp {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, err
	}
	//读取信息
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	tRechargeResult := &RechargeResult{}
	//将收到的[]byte转成json格式
	err = json.Unmarshal(body, tRechargeResult)
	if err != nil {
		return nil, err
	}
	//对leveldb数据库进行操作
	//添加读写锁
	RWMutex.Lock()
	if tRechargeResult.Code == "0" {
		suc := "true"
		CMail <- []string{uid, symbol, amount, suc, tRechargeResult.Msg}
		CSms <- []string{uid, symbol, amount, suc, tRechargeResult.Msg}
		go SmsSend()
		go MailSend()
		dbKey := ""
		dbKey = fmt.Sprintf("%s_%s_%s_%s", uid, symbol, timeStr, amount)

		err = Db.Put([]byte(dbKey), []byte("true" /*value*/), nil)
		if err != nil {
			return nil, err
		}

		// fmt.Println("充值成功")
	} else {
		suc := "false"
		CMail <- []string{uid, symbol, amount, suc, tRechargeResult.Msg}
		CSms <- []string{uid, symbol, amount, suc, tRechargeResult.Msg}
		go SmsSend()
		go MailSend()
		dbKey := ""
		dbKey = fmt.Sprintf("%s_%s_%s_%s", uid, symbol, timeStr, amount)

		err = Db.Put([]byte(dbKey), []byte("false" /*value*/), nil)
		if err != nil {
			return nil, err
		}

		// fmt.Println("充值失败")
	} //充值失败
	RWMutex.Unlock() //解锁
	return tRechargeResult, nil
}
