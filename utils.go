package WxService

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"log"
	"math/rand"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"
)

// 生成随机字符串
func GeneNonceStr(len int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}

	return string(bytes)
}

// 签名
func GeneSign(args map[string]interface{}, key string) string {
	var (
		stringA        string
		stringSignTemp string
		sign           string
		m              = md5.New()
	)
	// 1. 字典序排序
	stringA = ParseMap(args)
	log.Println(stringA)

	// 2. 与key拼接得到stringSignTemp
	if key != "" {
		stringSignTemp = stringA + "&key=" + key
	}

	// 3. 加密
	m.Write([]byte(stringSignTemp))

	sign = strings.ToUpper(hex.EncodeToString(m.Sum(nil)))

	return sign
}

// 校验签名
func CheckSign(key string, params map[string]interface{}) (bool, error) {
	var (
		ok      bool
		sign    string
		newSign string
	)

	sign, ok = params["sign"].(string)
	if !ok {
		return false, errors.New("sign类型错误！")
	}

	delete(params, "sign")

	newSign = GeneSign(params, key)

	return newSign == sign, nil
}

func ParseMap(args map[string]interface{}) string {
	var (
		stringA string
		keys    = make([]string, 0)
	)

	for k, _ := range args {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, key := range keys {
		switch v := args[key].(type) {
		case string:
			stringA += "&" + key + "=" + v

		case int:
			stringA += "&" + key + "=" + strconv.Itoa(v)

		case int32:
			stringA += "&" + key + "=" + strconv.Itoa(int(v))

		case int64:
			stringA += "&" + key + "=" + strconv.Itoa(int(v))

		default:
			return ""
		}
	}

	return stringA[1:]
}

// 获取本机IP
func GetIp() (string, error) {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range addr {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", nil
}

// 获取一个月的起始和最后一天的字符串
func GetBeginAndEndByMonth(year, month int) (string, string) {
	var (
		monthStr string
		lastDay  string
		begin    string
		end      string
		monthMap = map[int]int{
			1: 31, 2: 28, 3: 31, 4: 30, 5: 31, 6: 30, 7: 31, 8: 31, 9: 30, 10: 31, 11: 30, 12: 31,
		}
	)

	// 能被4整除
	// 整百时能被400整除
	if year%4 == 0 {
		if year%100 == 0 {
			if year%400 == 0 {
				monthMap[2] = 29
			}
		} else {
			monthMap[2] = 29
		}
	}

	lastDay = strconv.Itoa(monthMap[month])

	if month < 10 {
		monthStr = "0" + strconv.Itoa(month)
	} else {
		monthStr = strconv.Itoa(month)
	}

	begin = strconv.Itoa(year) + monthStr + "01"
	end = strconv.Itoa(year) + monthStr + lastDay

	return begin, end
}

// 获取上周的起始以及最后一天的字符串
func GetBeginAndEndByWeek() (string, string) {
	var (
		begin   string
		end     string
		weekMap = map[string]int{
			"Monday": 1, "Tuesday": 2, "Wednesday": 3, "Thursday": 4,
			"Friday": 5, "Saturday": 6, "Sunday": 7,
		}
		weekdayNow = time.Now().Weekday().String()
	)
	begin = time.Now().AddDate(0, 0, -weekMap[weekdayNow]-6).Format("20060102")
	end = time.Now().AddDate(0, 0, -weekMap[weekdayNow]).Format("20060102")

	return begin, end
}
