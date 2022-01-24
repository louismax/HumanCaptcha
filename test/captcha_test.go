package test

import (
	"github.com/louismax/HumanCaptcha/captcha"
	"testing"
)

func TestName(t *testing.T) {
	capt := captcha.NewClickCaptcha(
		captcha.InjectTextRangLenConfig(10, 12),
	)

	t.Log(capt)
}
