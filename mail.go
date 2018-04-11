package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

//{"code":200,"error":false,"message":"OK","data":{"sid":"AZvjVfwKeJuNIGoeuDS8","guide":null}}
type MailResp struct {
	Code    int
	Error   bool
	Message string
	Data    map[string]string
}

// var Conf *config.Config = config.InitConfig()
// var CMail chan []string = make(chan []string, 100)

func MailSend() {
	var MailReceipts []string = Conf.Receipts.Mail
	MailAddr := Conf.Api.Mail
	//信道读入
	mailInfo := <-CMail
	//nowTime := mailInfo[0]
	account := mailInfo[0]
	coin := mailInfo[1]
	amount := mailInfo[2]
	suc := mailInfo[3]
	errMsg := mailInfo[4]
	//rechargeAmount := Conf.Charges[coin].RechargeAmount
	var msgbody string
	if suc != "true" { //邮件发送内容
		msgbody = "On " + time.Now().Format("2006-01-02 15:04:05") + ",  asset monitor tool failed to recharged to account " + account + " with " + amount + " " + coin + ",  reason：[" + errMsg + "]。"
	} else {
		msgbody = "On " + time.Now().Format("2006-01-02 15:04:05") + ",  asset monitor tool had successfully recharged to account " + account + " with " + amount + " " + coin + "。"
	}
	for _, v := range MailReceipts {
		//PostForm发送form-data
		resp, err := http.PostForm(MailAddr,
			url.Values{"email": {v}, "codetype": {"notice"}, "vparam": {msgbody}})
		if err != nil {
			log.Println(v, ":Mail send failed")
			continue
		}
		defer resp.Body.Close()
		//得到返回值
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(v, ":Mail send failed")
			continue
		}
		var mailResp MailResp
		err = json.Unmarshal(body, &mailResp)
		if err != nil {
			log.Println(err)
		}
		if mailResp.Error {
			log.Println(v, ":Mail send failed")
			continue
		}
		log.Println(v, ":Mail send succeeded")
		//log.Println(string(body))
	}
}
