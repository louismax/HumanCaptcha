package toft

import (
	randC "crypto/rand"
	"math"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unsafe"
)

//RandomCreateZHCNUnicode 随机生成中文Unicode编码及转义字符
func RandomCreateZHCNUnicode() (string, string) {
	uStr := "\\u"
	for i := 0; i < 4; i++ {
		if i == 0 {
			uStr += randStr(1, "456789")
		} else if i == 1 && uStr == "\\u4" {
			uStr += randStr(1, "ef")
		} else if i == 2 && uStr == "\\u9f" {
			uStr += randStr(1, "0123456789a")
		} else if i == 3 && uStr == "\\u9fa" {
			uStr += randStr(1, "012345")
		} else {
			uStr += randStr(1, "0123456789abcdef")
		}
	}
	str, _ := strconv.Unquote(strings.Replace(strconv.Quote(uStr), `\\u`, `\u`, -1))
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) {
			return uStr, str
		} else {
			return RandomCreateZHCNUnicode()
		}
	}
	return "", ""
}

const letters = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var src = rand.NewSource(time.Now().UnixNano())

const (
	// 6 bits to represent a letter index
	letterIdBits = 6
	// All 1-bits as many as letterIdBits
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

func randStr(n int, letter ...string) string {
	letterX := letters
	if len(letter) > 0 {
		letterX = letter[0]
	}
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letterX) {
			b[i] = letterX[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

// RandInt 生成区间[-m, n]的安全随机数
func RandInt(min, max int) int {
	if min > max {
		return max
	}

	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int(f64Min)
		result, _ := randC.Int(randC.Reader, big.NewInt(int64(max+1+i64Min)))

		return int(result.Int64() - int64(i64Min))
	}
	result, _ := randC.Int(randC.Reader, big.NewInt(int64(max-min+1)))
	return int(int64(min) + result.Int64())
}
