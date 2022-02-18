package default_fonts

//FindFontsAsset 查找字体资源
func FindFontsAsset(path string) ([]byte, error) {
	return Asset(path)
}
