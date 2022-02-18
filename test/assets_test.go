package test

import (
	"github.com/louismax/HumanCaptcha/assets/default_fonts"
	"testing"
)

func TestDefaults(t *testing.T) {
	default_fonts.MustAsset("")
	default_fonts.AssetNames()
	_ = default_fonts.RestoreAssets("", "")
}
