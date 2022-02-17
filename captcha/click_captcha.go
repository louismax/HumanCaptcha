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
	//生成字符在图片上的位置
	allDots = cc.generateCharacterPosition(cc.config.imageSize, cc.config.rangFontSize, chars, 10)

	//随机得到有效的位置
	checkDots, checkChars = cc.rangeValidPosition(allDots)

	thumbDots = cc.generateCharacterPosition(cc.config.thumbnailSize, cc.config.rangCheckFontSize, checkChars, 0)

	imageBase64, err = cc.drawClickCaptchaImage(cc.config.imageSize, allDots)
	if err != nil {
		return nil, "", "", "", err
	}
	tImageBase64, err = cc.drawClickCaptchaThumbImage(cc.config.thumbnailSize, thumbDots)
	if err != nil {
		return nil, "", "", "", err
	}

	str, _ := json.Marshal(checkDots)
	key, _ := toft.GenCaptchaKey(string(str))
	return checkDots, imageBase64, tImageBase64, key, nil
}

//generateCharacterPosition 生成字符在图片上的位置
func (cc *ClickCaptcha) generateCharacterPosition(imageSize Size, fontSize RangeVal, chars string, padding int) map[int]CharDot {
	dots := make(map[int]CharDot) // 位置集合
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
		randAngle := cc.getRandomAngles()
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

//rangeValidPosition 随机得到有效的位置
func (cc *ClickCaptcha) rangeValidPosition(dots map[int]CharDot) (map[int]CharDot, string) {
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

//drawClickCaptchaImage 绘制点选验证码图片
func (cc *ClickCaptcha) drawClickCaptchaImage(size Size, dots map[int]CharDot) (base64 string, err error) {
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
		BackgroundDistort: cc.getContortionsByLevel(cc.config.imageFontDistort),
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

//drawClickCaptchaThumbImage 绘制点选验证码缩略图
func (cc *ClickCaptcha) drawClickCaptchaThumbImage(size Size, dots map[int]CharDot) (string, error) {
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
		BackgroundDistort:     cc.getContortionsByLevel(cc.config.thumbFontDistort),
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

//getRandomAngles 获取随机角度
func (cc *ClickCaptcha) getRandomAngles() int {
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

//getContortionsByLevel 获取扭曲程度
func (cc *ClickCaptcha) getContortionsByLevel(level int) int {
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

func CheckPointDist(cds []CheckDots, dots map[int]CharDot, paddings ...int64) bool {
	chkRet := false
	if len(paddings) > 0 {
		for i, dot := range dots {
			chkRet = checkPointDistWithPadding(int64(cds[i].X), int64(cds[i].Y), int64(dot.Dx), int64(dot.Dy), int64(dot.Width), int64(dot.Height), int64(paddings[0]))
			if !chkRet {
				break
			}
		}
	} else {
		for i, dot := range dots {
			chkRet = checkPointDistSimple(int64(cds[i].X), int64(cds[i].Y), int64(dot.Dx), int64(dot.Dy), int64(dot.Width), int64(dot.Height))
			if !chkRet {
				break
			}
		}
	}
	return chkRet
}

func checkPointDistSimple(sx, sy, dx, dy, width, height int64) bool {
	return sx >= dx &&
		sx <= dx+width &&
		sy <= dy &&
		sy >= dy-height
}

func checkPointDistWithPadding(sx, sy, dx, dy, width, height, padding int64) bool {
	newWidth := width + (padding * 2)
	newHeight := height + (padding * 2)
	newDx := int64(math.Max(float64(dx), float64(dx-padding)))
	newDy := dy + padding

	return sx >= newDx &&
		sx <= newDx+newWidth &&
		sy <= newDy &&
		sy >= newDy-newHeight
}
