package WxService

import (
	"log"
	"testing"
)

func TestClient_UnifiedOrder(t *testing.T) {
	c := NewClient(&ClientConfig{
		AppID:     "wx6eed9af8e39fc073",
		AppSecret: "5ddca8cb278fdbf6ca35baec2495557e",
		MchID:     "1597844801",
		ApiKey:    "8c42c5af8a59e751a625f8e5c5687119",
		IsSandBox: true,
	})

	p := NewParams()

	ip, err := GetIp()
	if err != nil{
		log.Println(err)
		return
	}

	p.SetString("appid", c.Client().AppID).
		SetString("mch_id", c.Client().MchID).
		SetString("nonce_str", "CZGZDFRCWZWMDDXMNPOXAWUFJTYIXUWJ").
		SetString("body", "百汇福-测试商品").
		SetString("out_trade_no", "1000000001").
		SetInt("total_fee", 101).
		SetString("spbill_create_ip", ip).
		SetString("notify_url", `http://47.106.125.164:8888/`).
		SetString("trade_type", "JSAPI")

	prepayId := c.UnifiedOrder(p)
	log.Println("prepay_id: ", prepayId)
	log.Println(c.Err())
}

func TestClient_GetSandboxSignKey(t *testing.T) {
	c := NewClient(&ClientConfig{
		AppID:     "wx6eed9af8e39fc073",
		AppSecret: "5ddca8cb278fdbf6ca35baec2495557e",
		MchID:     "1597844801",
		ApiKey:    "8c42c5af8a59e751a625f8e5c5687119",
		IsSandBox: true,
	})

	p := NewParams()

	p.SetString("mch_id", c.Client().MchID).
		SetString("nonce_str", GeneNonceStr(32))

	t.Log(c.GetSandboxSignKey(p))
}