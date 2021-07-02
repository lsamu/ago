package pay

import (
    "fmt"
    "github.com/iGoogle-ink/gopay"
    "github.com/iGoogle-ink/gopay/pkg/xlog"
    "github.com/iGoogle-ink/gopay/wechat/v3"
    "github.com/iGoogle-ink/gopay/alipay"
)

func InitWechatPay()  {
    xlog.Debug("GoPay Version: ", gopay.Version)
    client,_ :=wechat.NewClientV3("","","","","")
    bm:=gopay.BodyMap{}
    bm.Set("out_order_no","202104021339585117785701")
    bm.Set("out_request_no","20210402133958511778570101")
    bm.Set("remark","测试取消")
    rsp,err:=client.V3TransactionJsapi(bm)
    if err!=nil{
        fmt.Println(err)
    }
    fmt.Println(rsp.SignInfo)
}

func InitAliPay()  {
    xlog.Debug("GoPay Version: ", gopay.Version)
    client:=alipay.NewClient("","",false)
    bm:=gopay.BodyMap{}
    bm.Set("out_order_no","202104021339585117785701")
    bm.Set("out_request_no","20210402133958511778570101")
    bm.Set("remark","测试取消")
    rsp,err:=client.TradePrecreate(bm)
    if err!=nil{
        fmt.Println(err)
    }
    fmt.Println(rsp.Response)
}

func InitPayPal()  {
    
}