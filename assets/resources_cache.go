package assets

import "github.com/louismax/HumanCaptcha/assets/default_fonts"

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

/**
 * @Description: 获取默认资源
 * @param path
 * @return []byte
 * @return error
 */
//func findImagesAsset(path string) ([]byte, error) {
//	//return images.FindAsset(path)
//	return default_fonts.FindAsset(path)
//}

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

	//ret, err = findImagesAsset(path)
	//if len(ret) > 0 {
	//	cache = append(cache, &AssetData{
	//		Path:    path,
	//		Content: ret,
	//	})
	//	return
	//}
	return
}
