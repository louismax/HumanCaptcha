package captcha

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
