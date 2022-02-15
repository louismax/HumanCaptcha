package captcha

import "github.com/louismax/HumanCaptcha/captcha/assets"

/**
 * @Description: 获取缓存资源
 * @param path
 * @return []byte
 * @return error
 */
func getAssetCache(path string) (ret []byte, erro error) {
	return assets.GetAssetCache(path)
}
