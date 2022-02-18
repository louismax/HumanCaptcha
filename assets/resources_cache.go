package assets

import (
	"github.com/louismax/HumanCaptcha/assets/default_fonts"
	"github.com/louismax/HumanCaptcha/assets/default_images"
)

//AssetData 资源
type AssetData struct {
	// 路径
	Path string
	// 内容
	Content []byte
}

var cache []*AssetData

//findDefaultClickCaptchaFontsAsset 查找默认点选验证码字体资源
func findDefaultClickCaptchaFontsAsset(path string) ([]byte, error) {
	return default_fonts.FindFontsAsset(path)
}

//findDefaultClickCaptchaImagesAsset 查找默认点选验证码图片资源
func findDefaultClickCaptchaImagesAsset(path string) ([]byte, error) {
	return default_images.FindImagesAsset(path)
}

//GetClickCaptchaAssetCache 获取点选验证码资源缓存
func GetClickCaptchaAssetCache(path string) (ret []byte, err error) {
	if len(cache) > 0 {
		for _, asset := range cache {
			if asset.Path == path {
				ret = asset.Content
				return
			}
		}
	}

	ret, err = findDefaultClickCaptchaFontsAsset(path)
	if len(ret) > 0 {
		cache = append(cache, &AssetData{
			Path:    path,
			Content: ret,
		})
		return
	}

	ret, err = findDefaultClickCaptchaImagesAsset(path)
	if len(ret) > 0 {
		cache = append(cache, &AssetData{
			Path:    path,
			Content: ret,
		})
		return
	}
	return
}

//HasAssetCache 资源是否缓存
func HasAssetCache(path string) bool {
	if len(cache) > 0 {
		for _, asset := range cache {
			if asset.Path == path {
				return true
			}
		}
	}
	return false
}

//SetAssetCache 设置资源缓存
func SetAssetCache(path string, content []byte, force bool) bool {
	if len(cache) > 0 {
		for _, asset := range cache {
			if asset.Path == path && !force {
				return true
			}
		}
	}

	cache = append(cache, &AssetData{
		Path:    path,
		Content: content,
	})
	return true
}
