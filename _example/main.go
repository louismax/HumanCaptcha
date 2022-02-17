package main

import (
	"encoding/json"
	"fmt"
	"github.com/louismax/HumanCaptcha/captcha"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var capt *captcha.ClickCaptcha

func main() {
	capt = captcha.NewClickCaptcha(
		//captcha.InjectTextRangLenConfig(15, 25),
		//captcha.InjectCompleteGB2312CharsConfig(true),
		//captcha.InjectFontConfig([]string{"resources/fonts/simhei.ttf"}),
	)

	// Example: Get captcha data
	http.HandleFunc("/api/go_captcha_data", getCaptchaData)
	// Example: Post check data
	http.HandleFunc("/api/go_captcha_check_data", checkCaptcha)

	log.Println("ListenAndServe 0.0.0.0:9001")
	err := http.ListenAndServe(":9001", nil)
	if err != nil {
		log.Fatal("ListenAndServe err: ", err)
	}
}

func getCaptchaData(w http.ResponseWriter, _ *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	dots, b64, tb64, key, err := capt.GenerateClickCaptcha()
	if err != nil {
		bt, _ := json.Marshal(map[string]interface{}{
			"code":    1,
			"message": "GenCaptcha err",
		})
		_, _ = fmt.Fprintf(w, string(bt))
		return
	}
	writeCache(dots, key)
	bt, _ := json.Marshal(map[string]interface{}{
		"code":         0,
		"image_base64": b64,
		"thumb_base64": tb64,
		"captcha_key":  key,
	})
	_, _ = fmt.Fprintf(w, string(bt))
}

func checkCaptcha(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	code := 0

	con, _ := ioutil.ReadAll(r.Body)

	type req struct {
		Key  string              `json:"key"`
		Dots []captcha.CheckDots `json:"dots"`
	}

	reqData := req{}
	if err := json.Unmarshal(con, &reqData); err != nil {
		bt, _ := json.Marshal(map[string]interface{}{
			"code":    code,
			"message": "参数转换失败",
		})
		_, _ = fmt.Fprintf(w, string(bt))
		return
	}
	fmt.Println(reqData)

	if reqData.Key == "" || len(reqData.Dots) < 1 {
		bt, _ := json.Marshal(map[string]interface{}{
			"code":    code,
			"message": "点参数不允许为空",
		})
		_, _ = fmt.Fprintf(w, string(bt))
		return
	}

	//获取存储的验证信息
	cacheData := readCache(reqData.Key)
	if cacheData == "" {
		bt, _ := json.Marshal(map[string]interface{}{
			"code":    code,
			"message": "非法key",
		})
		_, _ = fmt.Fprintf(w, string(bt))
		return
	}
	var dct map[int]captcha.CharDot
	if err := json.Unmarshal([]byte(cacheData), &dct); err != nil {
		bt, _ := json.Marshal(map[string]interface{}{
			"code":    code,
			"message": "illegal key",
		})
		_, _ = fmt.Fprintf(w, string(bt))
		return
	}

	if len(reqData.Dots) != len(dct) {
		bt, _ := json.Marshal(map[string]interface{}{
			"code":    code,
			"message": "验证参数长度不够",
		})
		_, _ = fmt.Fprintf(w, string(bt))
		return
	}

	fmt.Println(dct)

	if captcha.CheckPointDist(reqData.Dots, dct) {
		bt, _ := json.Marshal(map[string]interface{}{
			"code":    1,
			"message": "ok",
		})
		_, _ = fmt.Fprintf(w, string(bt))
		return
	} else {
		bt, _ := json.Marshal(map[string]interface{}{
			"code":    code,
			"message": "验证失败",
		})
		_, _ = fmt.Fprintf(w, string(bt))
		return
	}
}

func readCache(file string) string {
	month := time.Now().Month().String()
	cacheDir := getCacheDir() + month + "/"
	file = cacheDir + file + ".json"

	if !checkFileIsExist(file) {
		return ""
	}

	bt, err := ioutil.ReadFile(file)
	err = os.Remove(file)
	if err == nil {
		return string(bt)
	}
	return ""
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func writeCache(v interface{}, file string) {
	bt, _ := json.Marshal(v)
	month := time.Now().Month().String()
	cacheDir := getCacheDir() + month + "/"
	_ = os.MkdirAll(cacheDir, 0660)
	file = cacheDir + file + ".json"
	logFile, _ := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	defer func(logFile *os.File) {
		_ = logFile.Close()
	}(logFile)
	// 检查过期文件
	//checkCacheOvertimeFile()
	_, _ = io.WriteString(logFile, string(bt))
}

func getCacheDir() string {
	return getPWD() + "/.cache/"
}

func getPWD() string {
	path, err := os.Getwd()
	if err != nil {
		return ""
	}
	return path
}
