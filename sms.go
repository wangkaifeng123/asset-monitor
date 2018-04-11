package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func SmsSend() {

	mailInfo := <-CSms
	account := mailInfo[0]
	currency := mailInfo[1]
	amount := mailInfo[2]
	rea := strings.EqualFold(mailInfo[3], "true")
	errMsg := mailInfo[4]
	phone_numbers := Conf.Receipts.Sms
	//循环每个邮箱，进行发邮件操作
	for i := 0; i < len(phone_numbers); i++ {
		smsSend(phone_numbers[i], account, currency, amount, rea, errMsg)
	}

}

func smsSend(phone_number string, account string, currency string, amount string, rea bool, errMsg string) {
	//当前时间的字符串
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	sendStr := ""
	if rea == true {
		sendStr = fmt.Sprintf("充值成功: On %s,  asset monitor tool had successfully "+
			"recharged to account %s with %s %s。", timeStr, account, amount, currency)
	} else {
		sendStr = fmt.Sprintf("充值失败: On %s,  asset monitor tool "+
			"failed to recharged to  account %s with %s %s， reason：[%s]。", timeStr, account, amount, currency, errMsg)
	}
	mailUrl := Conf.Api.Sms
	resp, err := http.PostForm(mailUrl, url.Values{"mobile": {phone_number}, "codetype": {"notice"}, "vparam": {sendStr}})
	if nil != resp {
		defer resp.Body.Close()
	}
	if err != nil {
		fmt.Println(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var tresult interface{}
	err = json.Unmarshal(body, &tresult)
	if err != nil {
		fmt.Println(err)
	}

	result, ok := tresult.(map[string]interface{})
	if !ok {
		fmt.Println("invalid result")
	}
	//检查code的类型
	var icode int
	switch t := result["code"].(type) {
	case nil:
		fmt.Println("sms send failed : req[code] is nil")
	case string:
		fmt.Println("sms send failed : req[code] is string")
	case bool:
		fmt.Println("sms send failed : req[code] is bool")
	case int:
		icode = t
	case float64:
		icode = int(t)
	default:
		fmt.Println("sms send failed : req[code] is unknow type")
	}
	//根据code值判断短信发送是否成功
	if icode == 1032 {
		errstr := ""
		switch t := result["error"].(type) {
		case nil:
			fmt.Println("error is nil")
		case string:
			errstr = t
			fmt.Println("sms send failed : " + errstr)
		case bool:
			fmt.Println("sms send failed : req[error] is bool")
		case int:
			fmt.Println("sms send failed : req[error] is int")
		default:
			fmt.Println("sms send failed : req[error] is unknow type")
		}
	}
}

//"13858075274","15907332188"
//"hxz@disanbo.com","sss@dds.com"
