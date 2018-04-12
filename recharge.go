package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//解码生成标签中的字段名,编码匹配标签中的字段名
type RechargeMessage struct {
	Symbol string `json:"symbol"` //json格式的币种
	Amount string `json:"amount"` //数量
}

type RechargeResult struct {
	Code string `json:"code"` //"0"成功 "1"失败
	Msg  string `json:"msg"`  //""->成功 "invalid symbol"->失败
	Data string `json:"data"` //"" ""
}

type Recharge struct {
}

func NewRecharge() *Recharge {
	r := &Recharge{}
	return r
}

func (r *Recharge) RechargeInbank(uid string, symbol string) (*RechargeResult, error) {
	//将充值数目以字符串格式化
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	amount := fmt.Sprintf("%d", Conf.Charges[symbol].RechargeAmount)
	reqUrl := Conf.Api.Transfer

	tRechargeMessage := &RechargeMessage{}
	tRechargeMessage.Symbol = symbol
	tRechargeMessage.Symbol = amount

	//将数据类型转换为json格式
	tRechargeMessageBytes, err := json.Marshal(tRechargeMessage)
	if nil != err {
		return nil, err
	}

	//以二进制数据流进行POST
	req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(tRechargeMessageBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	//处理返回数据
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
	fmt.Println(body)

	//将收到的[]byte转成json格式
	err = json.Unmarshal(body, tRechargeResult)
	if err != nil {
		return nil, err
	}

	//添加读写锁
	RWMutex.Lock()

	//对leveldb数据库进行操作
	if tRechargeResult.Code == "0" {
		suc := "true"
		//调用传短信和传邮件程序
		CMail <- []string{uid, symbol, amount, suc, tRechargeResult.Msg}
		CSms <- []string{uid, symbol, amount, suc, tRechargeResult.Msg}
		go SmsSend()
		go MailSend()
		//写入数据库 以用户名_币种_时间_充值数量
		dbKey := ""
		dbKey = fmt.Sprintf("%s_%s_%s_%s", uid, symbol, timeStr, amount)

		err = Db.Put([]byte(dbKey), []byte("true" /*value*/), nil)
		if err != nil {
			return nil, err
		} //充值成功

	} else {
		suc := "false"
		//调用传短信和传邮件程序
		CMail <- []string{uid, symbol, amount, suc, tRechargeResult.Msg}
		CSms <- []string{uid, symbol, amount, suc, tRechargeResult.Msg}
		go SmsSend()
		go MailSend()
		//写入数据库 以用户名_币种_时间_充值数量
		dbKey := ""
		dbKey = fmt.Sprintf("%s_%s_%s_%s", uid, symbol, timeStr, amount)

		err = Db.Put([]byte(dbKey), []byte("false" /*value*/), nil)
		if err != nil {
			return nil, err
		}

	} //充值失败
	//解锁
	RWMutex.Unlock()
	//将网页的body以json格式return
	return tRechargeResult, nil
}

