package captcha

import (
	"encoding/json"
	"fmt"
	"github.com/louismax/HumanCaptcha/toft"
	"github.com/sirupsen/logrus"
	"image/color"
	"math"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"
)

type ClickCaptcha struct {
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
		config:      GetClickCaptchaDefaultConfig(),
		captchaDraw: &Drawing{},
	}
}

//GenerateClickCaptcha 生成点选验证码
func (cc *ClickCaptcha) GenerateClickCaptcha() (map[int]CharDot, string, string, string, error) {
	length := toft.RandInt(cc.config.rangTextLen.Min, cc.config.rangTextLen.Max)

	//获取随机字符串
	chars := cc.getClickCaptchaChars(length)
	if chars == "" {
		return nil, "", "", "", fmt.Errorf("获取随机字符串失败")
	}

	fmt.Println(chars)

	var err error
	var allDots, thumbDots, checkDots map[int]CharDot
	var imageBase64, tImageBase64 string
	var checkChars string
	//生成字符在图片上的点
	allDots = cc.genDots(cc.config.imageSize, cc.config.rangFontSize, chars, 10)

	//随机检测点
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
	key, _ := toft.GenCaptchaKey(string(str))
	return checkDots, imageBase64, tImageBase64, key, nil
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
		randColor := toft.GetRandomStringValue(cc.config.rangFontColors)
		randColor2 := toft.GetRandomStringValue(cc.config.rangThumbFontColors)

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
func (cc *ClickCaptcha) genCaptchaImage(size Size, dots map[int]CharDot) (base64 string, err error) {
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
			Font:    toft.GetRandomStringValue(cc.config.rangFont),
		}

		drawDots = append(drawDots, drawDot)
	}

	img, err := cc.captchaDraw.Draw(DrawCanvas{
		Width:             size.Width,
		Height:            size.Height,
		Background:        toft.GetRandomStringValue(cc.config.rangBackground),
		BackgroundDistort: cc.getRandDistortWithLevel(cc.config.imageFontDistort),
		TextAlpha:         cc.config.imageFontAlpha,
		FontHinting:       cc.config.fontHinting,
		CaptchaDrawDot:    drawDots,

		ShowTextShadow:  cc.config.showTextShadow,
		TextShadowColor: cc.config.textShadowColor,
		TextShadowPoint: cc.config.textShadowPoint,
	})
	if err != nil {
		return
	}

	// 转 base64
	base64 = toft.EncodingImageToBase64StrForJpeg(img, cc.config.imageQuality)
	return
}

/**
 * @Description: 验证码缩略画图
 * @receiver cc
 * @param size
 * @param dots
 * @return string
 * @return error
 */
func (cc *ClickCaptcha) genCaptchaThumbImage(size Size, dots map[int]CharDot) (string, error) {
	var drawDots []DrawDot

	fontWidth := size.Width / len(dots)
	for i, dot := range dots {
		Dx := int(math.Max(float64(fontWidth*i+fontWidth/dot.Width), 8))
		Dy := size.Height/2 + dot.Size/2 - rand.Intn(size.Height/16*len(dot.Text))

		drawDot := DrawDot{
			Dx:      Dx,
			Dy:      Dy,
			FontDPI: cc.config.fontDPI,
			Text:    dot.Text,
			Angle:   dot.Angle,
			Color:   dot.Color2,
			Size:    dot.Size,
			Width:   dot.Width,
			Height:  dot.Height,
			Font:    toft.GetRandomStringValue(cc.config.rangFont),
		}
		drawDots = append(drawDots, drawDot)
	}

	params := DrawCanvas{
		Width:                 size.Width,
		Height:                size.Height,
		CaptchaDrawDot:        drawDots,
		BackgroundDistort:     cc.getRandDistortWithLevel(cc.config.thumbFontDistort),
		BackgroundCirclesNum:  cc.config.thumbBgCirclesNum,
		BackgroundSlimLineNum: cc.config.thumbBgSlimLineNum,
	}

	if len(cc.config.rangThumbBackground) > 0 {
		params.Background = toft.GetRandomStringValue(cc.config.rangThumbBackground)
	}

	var colorA []color.Color
	for _, cStr := range cc.config.rangThumbFontColors {
		co, _ := toft.ParseHexColorToRGBA(cStr)
		colorA = append(colorA, co)
	}

	var colorB []color.Color
	for _, co := range cc.config.rangThumbBgColors {
		rc, _ := toft.ParseHexColorToRGBA(co)
		colorB = append(colorB, rc)
	}

	img, err := cc.captchaDraw.DrawWithPalette(params, colorA, colorB)
	if err != nil {
		return "", err
	}

	// 转 base64
	dist := toft.EncodingImageToBase64StrForPng(img)
	return dist, err
}

//// EncodeB64string is a function
///**
// * @Description: base64编码
// * @receiver cc
// * @param img
// * @return string
// */
//func (cc *ClickCaptcha) EncodeB64stringWithPng(img image.Image) string {
//	return EncodeB64stringWithPng(img)
//}

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

//getClickCaptchaChars 获取点选验证码随机字符
func (cc *ClickCaptcha) getClickCaptchaChars(length int) string {
	var strA []string
	r := make(map[string]interface{})
	if cc.config.HasCompleteGB2312Chars {
		for len(strA) < length {
			uChar, char := toft.RandomCreateZHCNUnicode()
			if _, ok := r[uChar]; !ok {
				r[uChar] = char
				strA = append(strA, char)
			}
		}
	} else {
		for len(strA) < length {
			uChar, char := toft.RandomCreateSimplifyZHCNUnicode()
			if _, ok := r[uChar]; !ok {
				r[uChar] = char
				strA = append(strA, char)
			}
		}
	}
	return strings.Join(strA, ":")
}

/**
 * @Description: 根据级别获取扭曲程序
 * @receiver cc
 * @param level
 * @return int
 */
func (cc *ClickCaptcha) getRandDistortWithLevel(level int) int {
	if level == 1 {
		return toft.RandInt(240, 320)
	} else if level == 2 {
		return toft.RandInt(180, 240)
	} else if level == 3 {
		return toft.RandInt(120, 180)
	} else if level == 4 {
		return toft.RandInt(100, 160)
	} else if level == 5 {
		return toft.RandInt(80, 140)
	}
	return 0
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
