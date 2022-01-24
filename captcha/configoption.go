package captcha

import "fmt"

type CustomTextRangLen struct {
	rangeVal RangeVal
}

func (w CustomTextRangLen) Join(o *ClickCaptcha) error {
	o.config.rangTextLen = w.rangeVal
	return nil
}

type CustomRangCheckTextLen struct {
	rangeVal RangeVal
}

func (w CustomRangCheckTextLen) Join(o *ClickCaptcha) error {
	if w.rangeVal.Max > o.config.rangTextLen.Min {
		return fmt.Errorf("RangCheckTextLen.max必须小于或等于RangTextLen.min")
	}
	o.config.rangCheckTextLen = w.rangeVal
	return nil
}

type CustomImageSize struct {
	size Size
}

func (w CustomImageSize) Join(o *ClickCaptcha) error {
	o.config.imageSize = w.size
	return nil
}
