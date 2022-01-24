package captcha

import (
	"encoding/json"
	"fmt"
	"github.com/louismax/HumanCaptcha/toft"
	"github.com/sirupsen/logrus"
	"math"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"
)

type ClickCaptcha struct {
	// 随机字符串集合
	chars *[]string
	// 点选验证码配置
	config *ClickCaptchaConfig
	// 验证图像数据
	captchaDraw *Drawing
}

var clickCaptcha *ClickCaptcha

func NewClickCaptcha(opts ...ClickCaptchaConfigOption) *ClickCaptcha {
	if clickCaptcha == nil {
		clickCaptcha = initClickCaptcha()
		for _, opt := range opts {
			if err := opt.Join(clickCaptcha); err != nil {
				logrus.Warn("自定义配置插入异常!![Custom configuration insert exception!!]")
				return clickCaptcha
			}
		}
	}
	return clickCaptcha
}

func initClickCaptcha() *ClickCaptcha {
	return &ClickCaptcha{
		//chars:       GetCaptchaDefaultChars(),
		config:      GetClickCaptchaDefaultConfig(),
		captchaDraw: &Drawing{},
	}
}

// Generate is a function
/**
 * @Description: 			根据设置的尺寸生成验证码图片
 * @return CaptchaCharDot	位置信息
 * @return string			主图Base64
 * @return string			验证码KEY
 * @return string			缩略图Base64
 * @return error
 */
//func (cc *ClickCaptcha) Generate() (map[int]CharDot, string, string, string, error) {
//	dots, ib64, tb64, key, err := cc.GenerateWithSize(, )
//	return dots, ib64, tb64, key, err
//}

// GenerateWithSize is a function
/**
 * @Description: 			生成验证码图片
 * @param imageSize			主图尺寸
 * @param thumbnailSize		缩略图尺寸
 * @return CaptchaCharDot	位置信息
 * @return string			主图Base64
 * @return string			验证码KEY
 * @return string			缩略图Base64
 * @return error
 */
func (cc *ClickCaptcha) GenerateWithSize() (map[int]CharDot, string, string, string, error) {
	length := toft.RandInt(cc.config.rangTextLen.Min, cc.config.rangTextLen.Max)
	chars := cc.getClickCaptchaChars(length)
	if chars == "" {
		return nil, "", "", "", fmt.Errorf("获取随机字符串失败")
	}

	var allDots, thumbDots, checkDots map[int]CharDot
	var imageBase64, tImageBase64 string
	var checkChars string

	allDots = cc.genDots(cc.config.imageSize, cc.config.rangFontSize, chars, 10)
	// checkChars = "A:B:C"
	checkDots, checkChars = cc.rangeCheckDots(allDots)
	thumbDots = cc.genDots(cc.config.thumbnailSize, cc.config.rangCheckFontSize, checkChars, 0)
	imageBase64, err = cc.genCaptchaImage(cc.config.imageSize, allDots)
	if err != nil {
		return nil, "", "", "", err
	}
	tImageBase64, err = cc.genCaptchaThumbImage(cc.config.thumbnailSize, thumbDots)
	if err != nil {
		return nil, "", "", "", err
	}

	str, _ := json.Marshal(checkDots)
	key, _ := cc.genCaptchaKey(string(str))
	return checkDots, imageBase64, tImageBase64, key, err
}

/**
 * @Description: 生成字符在图片上的点
 * @receiver cc
 * @param imageSize
 * @param fontSize
 * @param chars
 * @param padding
 * @return []*CaptchaCharDot
 */
func (cc *ClickCaptcha) genDots(imageSize Size, fontSize RangeVal, chars string, padding int) map[int]CharDot {
	dots := make(map[int]CharDot) // 各个文字点位置
	width := imageSize.Width
	height := imageSize.Height
	if padding > 0 {
		width -= padding
		height -= padding
	}

	strS := strings.Split(chars, ":")
	for i := 0; i < len(strS); i++ {
		str := strS[i]
		// 随机角度
		randAngle := cc.getRandAngle()
		// 随机颜色
		randColor := cc.getRandColor(cc.config.rangFontColors)
		randColor2 := cc.getRandColor(cc.config.rangThumbFontColors)

		// 随机文字大小
		randFontSize := toft.RandInt(fontSize.Min, fontSize.Max)
		fontHeight := randFontSize
		fontWidth := randFontSize

		if utf8.RuneCountInString(str) > 1 {
			fontWidth = randFontSize * utf8.RuneCountInString(str)

			if randAngle > 0 {
				surplus := fontWidth - randFontSize
				ra := randAngle % 90
				pr := float64(surplus) / 90
				h := math.Max(float64(ra)*pr, 1)
				fontHeight = fontHeight + int(h)
			}
		}

		_w := width / len(strS)
		rd := math.Abs(float64(_w) - float64(fontWidth))
		x := (i * _w) + toft.RandInt(0, int(math.Max(rd, 1)))
		x = int(math.Min(math.Max(float64(x), 10), float64(width-10-(padding*2))))
		y := toft.RandInt(10, height+fontHeight)
		y = int(math.Min(math.Max(float64(y), float64(fontHeight+10)), float64(height+(fontHeight/2)-(padding*2))))
		text := fmt.Sprintf("%s", str)

		dot := CharDot{i, x, y, randFontSize, fontWidth, fontHeight, text, randAngle, randColor, randColor2}
		dots[i] = dot
	}

	return dots
}

/**
 * @Description: 随机检测点
 * @receiver cc
 * @param dots
 * @return map[int]CaptchaCharDot
 * @return string
 */
func (cc *ClickCaptcha) rangeCheckDots(dots map[int]CharDot) (map[int]CharDot, string) {
	rand.Seed(time.Now().UnixNano())
	rs := rand.Perm(len(dots))
	chkDots := make(map[int]CharDot)
	count := toft.RandInt(cc.config.rangCheckTextLen.Min, cc.config.rangCheckTextLen.Max)
	var chars []string
	for i, value := range rs {
		if i >= count {
			continue
		}
		dot := dots[value]
		dot.Index = i
		chkDots[i] = dot
		chars = append(chars, chkDots[i].Text)
	}
	return chkDots, strings.Join(chars, ":")
}

/**
 * @Description: 验证码画图
 * @receiver cc
 * @param size
 * @param dots
 * @return string
 * @return error
 */
func (cc *ClickCaptcha) genCaptchaImage(size Size, dots map[int]CharDot) (base64 string, erro error) {
	var drawDots []DrawDot
	for _, dot := range dots {
		drawDot := DrawDot{
			Dx:      dot.Dx,
			Dy:      dot.Dy,
			FontDPI: cc.config.fontDPI,
			Text:    dot.Text,
			Angle:   dot.Angle,
			Color:   dot.Color,
			Size:    dot.Size,
			Width:   dot.Width,
			Height:  dot.Height,
			Font:    cc.genRandWithString(cc.config.rangFont),
		}

		drawDots = append(drawDots, drawDot)
	}

	img, err := cc.captchaDraw.Draw(DrawCanvas{
		Width:             size.Width,
		Height:            size.Height,
		Background:        cc.genRandWithString(cc.config.rangBackground),
		BackgroundDistort: cc.getRandDistortWithLevel(cc.config.imageFontDistort),
		TextAlpha:         cc.config.imageFontAlpha,
		FontHinting:       cc.config.fontHinting,
		CaptchaDrawDot:    drawDots,

		ShowTextShadow:  cc.config.showTextShadow,
		TextShadowColor: cc.config.textShadowColor,
		TextShadowPoint: cc.config.textShadowPoint,
	})
	if err != nil {
		erro = err
		return
	}

	// 转 base64
	base64 = cc.EncodeB64stringWithJpeg(img)
	return
}

/**
 * @Description: 获取随机角度
 * @receiver cc
 * @return int
 */
func (cc *ClickCaptcha) getRandAngle() int {
	angles := cc.config.rangTexAnglePos
	anglesLen := len(angles)
	index := toft.RandInt(0, anglesLen)
	if index >= anglesLen {
		index = anglesLen - 1
	}

	angle := angles[index]
	res := toft.RandInt(angle.Min, angle.Max)

	return res
}

/**
 * @Description: 随机获取颜色
 * @param colors
 * @return string
 */
func (cc *ClickCaptcha) getRandColor(colors []string) string {
	colorLen := len(colors)
	index := toft.RandInt(0, colorLen)
	if index >= colorLen {
		index = colorLen - 1
	}

	return colors[index]
}

func (_ *ClickCaptcha) getClickCaptchaChars(length int) string {
	var strA []string
	r := make(map[string]interface{})
	for len(strA) < length {
		uChar, char := toft.RandomCreateZHCNUnicode()
		if _, ok := r[uChar]; !ok {
			r[uChar] = char
			strA = append(strA, char)
		}
	}

	return strings.Join(strA, ":")
}

/**
 * @Description: 随机获取值
 * @param strs
 * @return string
 */
func (cc *ClickCaptcha) genRandWithString(strs []string) string {
	strLen := len(strs)
	if strLen == 0 {
		return ""
	}

	index := toft.RandInt(0, strLen)
	if index >= strLen {
		index = strLen - 1
	}

	return strs[index]
}

// CharDot is a type
/**
 * @Description: 图片点数据
 */
type CharDot struct {
	// 顺序索引
	Index int
	// x,y位置
	Dx int
	Dy int
	// 字体大小
	Size int
	// 字体宽
	Width int
	// 字体高
	Height int
	// 字符文本
	Text string
	// 字体角度
	Angle int
	// 颜色
	Color string
	// 颜色2
	Color2 string
}
