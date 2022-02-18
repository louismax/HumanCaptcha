package default_images

//FindImagesAsset 查找图片资源
func FindImagesAsset(path string) ([]byte, error) {
	return Asset(path)
}
