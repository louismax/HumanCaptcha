package captcha

import (
	"github.com/louismax/HumanCaptcha/toft"
	"golang.org/x/image/font"
)

type ClickCaptchaConfigOption interface {
	Join(settings *ClickCaptcha) error
}

// RangeVal is a type
/**
 * @Description: 范围值
 * @Example: {min: 0, max: 45} 从0-45中取任意值
 */
type RangeVal struct {
	Min, Max int
}

// Size is a type
/**
 * @Description: 尺寸
 * @Example: {width: 0, height: 45} 从0-45中取任意值
 */
type Size struct {
	Width, Height int
}

// Point is a type
/**
 * @Description: 点
 */
type Point struct {
	X, Y int
}

/**
 * @Description: 扭曲程度
 */
const (
	DistortNone   = iota // 无扭曲
	DistortLevel1        // 扭曲程度 1级别
	DistortLevel2        // 扭曲程度 2级别
	DistortLevel3        // 扭曲程度 3级别
	DistortLevel4        // 扭曲程度 4级别
	DistortLevel5        // 扭曲程度 5级别
)

/**
 * @Description: 质量压缩程度
 */
const (
	QualityCompressNone = iota // 无压缩质量,原图

	QualityCompressLevel1 = 100 // 质量压缩程度 1-5 级别，压缩级别越低图像越清晰
	QualityCompressLevel2 = 80
	QualityCompressLevel3 = 60
	QualityCompressLevel4 = 40
	QualityCompressLevel5 = 20
)

var textColors = []string{
	"#fde98e",
	"#60c1ff",
	"#fcb08e",
	"#fb88ff",
	"#b4fed4",
	"#cbfaa9",
}

var thumbTextColors = []string{
	"#006600",
	"#005db9",
	"#aa002a",
	"#875400",
	"#6e3700",
	"#660033",
}

var textShadowColor = "#101010"

type ClickCaptchaConfig struct {
	// 随机字符串长度范围
	rangTextLen RangeVal
	// 随机验证字符串长度范围, 注意：RangCheckTextLen < RangTextLen
	rangCheckTextLen RangeVal
	// 随机文本角度范围集合
	rangTexAnglePos []RangeVal
	// 随机文本尺寸范围集合
	rangFontSize RangeVal
	// 随机缩略文本尺寸范围集合
	rangCheckFontSize RangeVal
	// 随机文本颜色	格式："#541245"
	rangFontColors []string
	// 文本阴影偏移位置
	showTextShadow bool
	// 文本阴影颜色
	textShadowColor string
	// 文本阴影偏移位置
	textShadowPoint Point
	// 缩略图随机文本颜色	格式："#541245"
	rangThumbFontColors []string
	// 随机字体	格式：字体绝对路径字符串, /home/..../xxx.ttf
	rangFont []string
	// 屏幕每英寸的分辨率
	fontDPI int
	// 随机验证码背景图		格式：图片绝对路径字符串, /home/..../xxx.png
	rangBackground []string
	// 验证码尺寸, 注意：高度 > RangFontSize.max , 长度 > RangFontSize.max * RangFontSize.max
	imageSize Size
	// 图片清晰度 1-101
	imageQuality int
	// 验证码文本扭曲程度
	imageFontDistort int
	// 验证码文本透明度 0-1
	imageFontAlpha float64
	// 缩略图尺寸, 注意：高度 > RangCheckFontSize.max , 长度 > RangCheckFontSize.max * RangFontSize.max
	thumbnailSize Size
	// 字体Hinting
	fontHinting font.Hinting
	// 随机缩略背景图		格式：图片绝对路径字符串, /home/..../xxx.png
	rangThumbBackground []string
	// 缩略图背景随机色	格式："#541245"
	rangThumbBgColors []string
	// 缩略图扭曲程度，值为 Distort...,
	thumbBgDistort int
	// 缩略图文字扭曲程度，值为 Distort...,
	thumbFontDistort int
	// 缩略图小圆点数量
	thumbBgCirclesNum int
	// 缩略图线条数量
	thumbBgSlimLineNum int
}

func GetClickCaptchaDefaultConfig() *ClickCaptchaConfig {
	return &ClickCaptchaConfig{
		rangTextLen:      RangeVal{6, 7},
		rangCheckTextLen: RangeVal{2, 4},
		rangTexAnglePos: []RangeVal{
			{20, 35},
			{35, 45},
			{45, 60},
			{60, 75},

			{285, 305},
			{305, 325},
			{325, 345},
			{345, 365},
		},
		rangFontSize:        RangeVal{28, 38},
		fontDPI:             72,
		rangCheckFontSize:   RangeVal{24, 32},
		imageFontDistort:    DistortNone,
		imageFontAlpha:      1,
		rangFontColors:      getDefaultTextColors(),
		showTextShadow:      true,
		textShadowColor:     getDefaultTextShadowColor(),
		textShadowPoint:     Point{-1, -1},
		rangThumbFontColors: getDefaultThumbTextColors(),
		fontHinting:         font.HintingNone,
		imageSize:           Size{300, 240},
		imageQuality:        QualityCompressLevel1,
		thumbnailSize:       Size{150, 40},
		rangThumbBgColors:   getDefaultThumbTextColors(),
		thumbFontDistort:    DistortLevel3,
		thumbBgDistort:      DistortLevel4,
		thumbBgCirclesNum:   24,
		thumbBgSlimLineNum:  2,

		rangFont:       toft.DefaultBinFontList(),
		rangBackground: toft.DefaultBinImageList(),
	}
}

/**
 * @Description: 获取默认文本颜色
 * @return []string
 */
func getDefaultTextColors() []string {
	return textColors
}

/**
 * @Description: 获取默认阴影文本颜色
 * @return string
 */
func getDefaultTextShadowColor() string {
	return textShadowColor
}

/**
 * @Description: 获取默认缩略图文本颜色
 * @return []string
 */
func getDefaultThumbTextColors() []string {
	return thumbTextColors
}
