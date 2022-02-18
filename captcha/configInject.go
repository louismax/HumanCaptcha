package captcha

import (
	"fmt"
	"github.com/louismax/HumanCaptcha/assets"
	"github.com/louismax/HumanCaptcha/toft"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

// InjectTextRangLenConfig 设置字符随机长度范围
func InjectTextRangLenConfig(min, max int) ClickCaptchaConfigOption {
	return CustomTextRangLen{
		rangeVal: RangeVal{
			Min: min,
			Max: max,
		},
	}
}

// InjectRangCheckTextLenConfig 设置随机验证字符串长度范围
func InjectRangCheckTextLenConfig(min, max int) ClickCaptchaConfigOption {
	return CustomRangCheckTextLen{
		rangeVal: RangeVal{
			Min: min,
			Max: max,
		},
	}
}

// InjectImageSizeConfig 设置验证码图片尺寸
func InjectImageSizeConfig(width, height int) ClickCaptchaConfigOption {
	return CustomImageSize{
		size: Size{
			Width:  width,
			Height: height,
		},
	}
}

// InjectFontConfig 设置字体配置
func InjectFontConfig(fonts []string, args ...bool) ClickCaptchaConfigOption {
	for _, path := range fonts {
		if has, err := toft.PathExists(path); !has || err != nil {
			panic(fmt.Errorf("CaptchaConfig Error: The [%s] file does not exist", path))
		}
		hasCache := assets.HasAssetCache(path)
		if !hasCache || (hasCache && len(args) > 0 && args[0]) {
			bytes, err := ioutil.ReadFile(path)
			if err != nil {
				panic(err)
			}

			assets.SetAssetCache(path, bytes, len(args) > 0 && args[0])
		}
	}
	return CustomRangFont{
		fonts: fonts,
	}
}

//InjectCompleteGB2312CharsConfig 设置是否使用GB2312字符集
func InjectCompleteGB2312CharsConfig(v bool) ClickCaptchaConfigOption {
	if v {
		logrus.Warn("Custom configuration using GB2312 character set, there may be font files cannot render some characters!! You are advised to replace the font file that fully supports GB2312！")
	}
	return CustomCompleteGB2312Chars{
		val: v,
	}
}
