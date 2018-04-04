# 监控公司银行账号资产余额

## 背景

有 N 个账号，每个账号名下有 M 个币种，为了防止某个币种余额过低导致出现用户充值失败的情况，有必要对这些币种的余额信息进行监控，使得当余额低于某个阀值时我们能够及时发现并充值。

## 功能

1. 当某个账号的某个币种余额低于某个阀值时，往该账号充币，同时发送短信+邮件通知
2. 提供 rpc接口+web页面 查询历史充值记录（指因为余额过低而充值的情况）
3. 限制每天单个账号单个币种的充值次数 

## 配置

见 `etc/config.toml`

### 短息/邮件内容

```
成功：On 2018-04-03 11:21:09,  asset monitor tool had successfully recharged to account 3f093fed33d7 with 1 BTC。

失败：On 2018-04-03 11:21:09,  asset monitor tool failed to recharged to  account 3f093fed33d7 with 1 BTC， reason：[ErrNoMargin]。
```
