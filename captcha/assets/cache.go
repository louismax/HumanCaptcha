package assets

import (
	"github.com/louismax/HumanCaptcha/captcha/assets/fonts"
	"github.com/louismax/HumanCaptcha/captcha/assets/images"
)

type AssetData struct {
	// 路径
	Path string
	// 内容
	Content []byte
}

var cache []*AssetData

/**
 * @Description: 获取默认字体资源
 * @param path
 * @return []byte
 * @return error
 */
func findFontsAsset(path string) ([]byte, error) {
	return fonts.FindAsset(path)
}

/**
 * @Description: 获取默认图片资源
 * @param path
 * @return []byte
 * @return error
 */
func findImagesAsset(path string) ([]byte, error) {
	return images.FindAsset(path)
}

// GetAssetCache is a function
/**
 * @Description: 获取缓存资源
 * @param path
 * @return []byte
 * @return error
 */
func GetAssetCache(path string) (ret []byte, erro error) {
	if len(cache) > 0 {
		for _, asset := range cache {
			if asset.Path == path {
				ret = asset.Content
				return
			}
		}
	}

	ret, erro = findFontsAsset(path)
	if len(ret) > 0 {
		cache = append(cache, &AssetData{
			Path:    path,
			Content: ret,
		})
		return
	}

	ret, erro = findImagesAsset(path)
	if len(ret) > 0 {
		cache = append(cache, &AssetData{
			Path:    path,
			Content: ret,
		})
		return
	}
	return
}
