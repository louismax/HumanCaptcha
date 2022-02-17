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
	"strconv"
	"strings"
	"time"
)

var capt *captcha.ClickCaptcha

func main() {
	capt = captcha.NewClickCaptcha(
		//captcha.InjectTextRangLenConfig(15, 25),
		captcha.InjectCompleteGB2312CharsConfig(true),
		captcha.InjectFontConfig([]string{"resources/fonts/simhei.ttf"}),
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

func getCaptchaData(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	//chars := "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	//_ = capt.SetRangChars(strings.Split(chars, ""))

	//chars := []string{"HE","CA","WO","NE","HT","IE","PG","GI","CH","CO","DA"}
	//_ = capt.SetRangChars(chars)

	//chars := []string{"你","好","呀","这","是","点","击","验","证","码","哟"}
	//_ = capt.SetRangChars(chars)

	//capt.SetTextRangFontColors([]string{
	//	"#006600",
	//	"#005db9",
	//	"#aa002a",
	//	"#875400",
	//	"#6e3700",
	//	"#333333",
	//	"#660033",
	//})
	//
	//capt.SetFont([]string{
	//	getPWD() + "/resources/fonts/fzshengsksjw_cu.ttf",
	//	getPWD() + "/resources/fonts/fzssksxl.ttf",
	//	getPWD() + "/resources/fonts/hyrunyuan.ttf",
	//})

	// capt.SetBackground([]string{
	// 	getPWD() + "/resources/images/1.jpg",
	// 	getPWD() + "/resources/images/2.jpg",
	// 	getPWD() + "/resources/images/3.jpg",
	// 	getPWD() + "/resources/images/4.jpg",
	// 	getPWD() + "/resources/images/5.jpg",
	// })

	//capt.SetThumbBackground([]string{
	//	getPWD() + "/resources/images/thumb/r1.jpg",
	//	getPWD() + "/resources/images/thumb/r2.jpg",
	//	getPWD() + "/resources/images/thumb/r3.jpg",
	//	getPWD() + "/resources/images/thumb/r4.jpg",
	//	getPWD() + "/resources/images/thumb/r5.jpg",
	//})

	//capt.SetThumbBgCirclesNum(200)
	//capt.SetImageFontAlpha(0.5)

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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	code := 1
	_ = r.ParseForm()
	dots := r.Form.Get("dots")
	key := r.Form.Get("key")
	if dots == "" || key == "" {
		bt, _ := json.Marshal(map[string]interface{}{
			"code":    code,
			"message": "dots or key param is empty",
		})
		_, _ = fmt.Fprintf(w, string(bt))
		return
	}

	cacheData := readCache(key)
	if cacheData == "" {
		bt, _ := json.Marshal(map[string]interface{}{
			"code":    code,
			"message": "illegal key",
		})
		_, _ = fmt.Fprintf(w, string(bt))
		return
	}
	src := strings.Split(dots, ",")

	var dct map[int]captcha.CharDot
	if err := json.Unmarshal([]byte(cacheData), &dct); err != nil {
		bt, _ := json.Marshal(map[string]interface{}{
			"code":    code,
			"message": "illegal key",
		})
		_, _ = fmt.Fprintf(w, string(bt))
		return
	}

	chkRet := false
	if len(src) >= len(dct)*2 {
		chkRet = true
		for i, dot := range dct {
			j := i * 2
			k := i*2 + 1
			sx, _ := strconv.ParseFloat(fmt.Sprintf("%v", src[j]), 64)
			sy, _ := strconv.ParseFloat(fmt.Sprintf("%v", src[k]), 64)
			// 检测点位置
			chkRet = CheckPointDist(int64(sx), int64(sy), int64(dot.Dx), int64(dot.Dy), int64(dot.Width), int64(dot.Height))
			if !chkRet {
				break
			}
		}
	}

	if chkRet && (len(dct)*2) == len(src) {
		code = 0
	}

	bt, _ := json.Marshal(map[string]interface{}{
		"code": code,
	})
	_, _ = fmt.Fprintf(w, string(bt))
	return
}

func CheckPointDist(sx, sy, dx, dy, width, height int64) bool {
	return sx >= dx &&
		sx <= dx+width &&
		sy <= dy &&
		sy >= dy-height
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
	defer logFile.Close()
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
