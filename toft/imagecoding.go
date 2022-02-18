package toft

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
)

//DecodingBinaryToImage 解码二进制图片
func DecodingBinaryToImage(b []byte) (img image.Image, err error) {
	var buf bytes.Buffer
	buf.Write(b)
	img, err = jpeg.Decode(&buf)
	buf.Reset()
	return
}

//EncodingImageToBinaryWithJpeg 图片编码为二进制(image格式)
func EncodingImageToBinaryWithJpeg(img image.Image, quality int) (ret []byte) {
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality}); err != nil {
		panic(err.Error())
	}
	ret = buf.Bytes()
	buf.Reset()
	return
}

//EncodingImageToBinaryWithPng 图片编码为二进制(png格式)
func EncodingImageToBinaryWithPng(img image.Image) (ret []byte) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		panic(err.Error())
	}
	ret = buf.Bytes()
	buf.Reset()
	return
}

//EncodingImageToBase64StrForJpeg 将普通(JPG,JPEG)图片转为base64字符串
func EncodingImageToBase64StrForJpeg(img image.Image, quality int) string {
	return fmt.Sprintf("data:%s;base64,%s", "image/jpeg", base64.StdEncoding.EncodeToString(EncodingImageToBinaryWithJpeg(img, quality)))
}

//EncodingImageToBase64StrForPng 将png格式图片转为base64字符串
func EncodingImageToBase64StrForPng(img image.Image) string {
	return fmt.Sprintf("data:%s;base64,%s", "image/png", base64.StdEncoding.EncodeToString(EncodingImageToBinaryWithPng(img)))
}

//ParseHexColorToRGBA 十六进制颜色转为color.RGBA
func ParseHexColorToRGBA(s string) (c color.RGBA, err error) {
	c.A = 0xff
	if s[0] != '#' {
		return c, errors.New("十六进制颜色必须以'#'开头")
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = errors.New("十六进制颜色无效")
		return 0
	}

	switch len(s) {
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])

	case 4:
		c.R = hexToByte(s[1]) * 17
		c.G = hexToByte(s[2]) * 17
		c.B = hexToByte(s[3]) * 17
	default:
		err = errors.New("十六进制颜色无效")
	}
	return
}
