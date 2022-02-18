package captcha

import "fmt"

//CustomTextRangLen is CustomTextRangLen
type CustomTextRangLen struct {
	rangeVal RangeVal
}

// Join CustomTextRangLen.Join
func (w CustomTextRangLen) Join(o *ClickCaptcha) error {
	o.config.rangTextLen = w.rangeVal
	return nil
}

//CustomRangCheckTextLen is CustomRangCheckTextLen
type CustomRangCheckTextLen struct {
	rangeVal RangeVal
}

// Join CustomRangCheckTextLen.Join
func (w CustomRangCheckTextLen) Join(o *ClickCaptcha) error {
	if w.rangeVal.Max > o.config.rangTextLen.Min {
		return fmt.Errorf("RangCheckTextLen.max必须小于或等于RangTextLen.min")
	}
	o.config.rangCheckTextLen = w.rangeVal
	return nil
}

//CustomImageSize is CustomImageSize
type CustomImageSize struct {
	size Size
}

// Join CustomImageSize.Join
func (w CustomImageSize) Join(o *ClickCaptcha) error {
	o.config.imageSize = w.size
	return nil
}

//CustomRangFont is CustomRangFont
type CustomRangFont struct {
	fonts []string
}

// Join CustomRangFont.Join
func (w CustomRangFont) Join(o *ClickCaptcha) error {
	o.config.rangFont = w.fonts
	return nil
}

//CustomCompleteGB2312Chars is CustomCompleteGB2312Chars
type CustomCompleteGB2312Chars struct {
	val bool
}

// Join CustomCompleteGB2312Chars.Join
func (w CustomCompleteGB2312Chars) Join(o *ClickCaptcha) error {
	o.config.HasCompleteGB2312Chars = w.val
	return nil
}
