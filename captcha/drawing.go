package captcha

import (
	"fmt"
	"github.com/golang/freetype"
	"github.com/louismax/HumanCaptcha/assets"
	"github.com/louismax/HumanCaptcha/toft"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"math"
	mRand "math/rand"
)

// DrawDot 绘制点
type DrawDot struct {
	Dx      int
	Dy      int
	FontDPI int
	Text    string
	Size    int
	Width   int
	Height  int
	Angle   int
	Color   string
	Color2  string
	Font    string
}

// DrawCanvas 绘制画布
type DrawCanvas struct {
	// 长、高
	Width  int
	Height int
	// 背景图片
	Background string
	// 缩略图扭曲程度，值为 Distort...,
	BackgroundDistort int
	// 缩略图小圆点数量
	BackgroundCirclesNum int
	// 缩略图线条数量
	BackgroundSlimLineNum int
	// 文本透明度
	TextAlpha float64
	// FontHinting
	FontHinting font.Hinting
	// 点
	CaptchaDrawDot []DrawDot
	// 文本阴影偏移位置
	ShowTextShadow bool
	// 文本阴影颜色
	TextShadowColor string
	// 文本阴影偏移位置
	TextShadowPoint Point
}

// AreaPoint 区域点信息
type AreaPoint struct {
	MinX, MaxX, MinY, MaxY int
}

//Drawing 绘制
type Drawing struct{}

// Draw 绘制
func (cd *Drawing) Draw(params DrawCanvas) (image.Image, error) {
	dots := params.CaptchaDrawDot
	canvas := cd.CreateCanvas(params, true)

	for _, dot := range dots {
		textImg, areaPoint, _ := cd.DrawTextImg(dot, params)

		minX := areaPoint.MinX
		maxX := areaPoint.MaxX
		minY := areaPoint.MinY
		maxY := areaPoint.MaxY
		width := maxX - minX
		height := maxY - minY
		nW := textImg.Bounds().Max.X
		nH := textImg.Bounds().Max.Y
		for x := 0; x < nW; x++ {
			for y := 0; y < nH; y++ {
				co := textImg.At(x, y)
				if _, _, _, a := co.RGBA(); a > 0 {
					if x >= minX && x <= maxX && y >= minY && y <= maxY {
						canvas.Set(dot.Dx+(x-minX), dot.Dy-height+(y-minY), textImg.At(x, y))
					}
				}
			}
		}
		// 重置尺寸
		dot.Height = height
		dot.Width = width
		// 重置XY位置
		dot.Dx = minX
		dot.Dy = maxY
	}

	bgFile := params.Background
	imgBg, iErr := assets.GetClickCaptchaAssetCache(bgFile)
	if iErr != nil {
		return canvas, iErr
	}

	img, dErr := toft.DecodingBinaryToImage(imgBg)
	if dErr != nil {
		return canvas, dErr
	}

	b := canvas.Bounds()
	m := image.NewNRGBA(b)
	curX, curY := cd.rangCutImage(params.Width, params.Height, img)
	draw.Draw(m, b, img, image.Point{X: curX, Y: curY}, draw.Src)
	draw.Draw(m, canvas.Bounds(), canvas, image.Point{}, draw.Over)
	subImg := m.SubImage(image.Rect(0, 0, params.Width, params.Height)).(*image.NRGBA)
	return subImg, nil
}

// DrawWithPalette 使用调色板的画布绘图
func (cd *Drawing) DrawWithPalette(params DrawCanvas, colorA []color.Color, colorB []color.Color) (image.Image, error) {
	dots := params.CaptchaDrawDot
	p := []color.Color{
		color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0x00},
	}
	for ak := range colorA {
		p = append(p, colorA[ak])
	}
	for bk := range colorB {
		p = append(p, colorB[bk])
	}

	canvas := NewPalette(image.Rect(0, 0, params.Width, params.Height), p)

	if params.BackgroundCirclesNum > 0 {
		cd.rangFillWithCircles(canvas, params.BackgroundCirclesNum, 1, 2)
	}
	if params.BackgroundSlimLineNum > 0 {
		cd.drawLine(canvas, params.BackgroundSlimLineNum, colorB)
	}

	for _, dot := range dots {
		// 读字体数据
		fontBytes, err := assets.GetClickCaptchaAssetCache(dot.Font)
		if err != nil {
			return canvas, err
		}
		fontN, err := freetype.ParseFont(fontBytes)
		if err != nil {
			return canvas, err
		}

		dc := freetype.NewContext()
		dc.SetDPI(float64(dot.FontDPI))
		dc.SetFont(fontN)
		dc.SetClip(canvas.Bounds())
		dc.SetDst(canvas)

		// 文字大小
		dc.SetFontSize(float64(dot.Size))

		// 文字颜色
		hexColor, _ := toft.ParseHexColorToRGBA(dot.Color)
		fontColor := image.NewUniform(hexColor)
		dc.SetSrc(fontColor)

		// 画文本
		text := fmt.Sprintf("%s", dot.Text)
		pt := freetype.Pt(dot.Dx, dot.Dy) // 字出现的位置
		_, err = dc.DrawString(text, pt)
		if err != nil {
			return canvas, err
		}
	}

	if params.Background != "" {
		bgFile := params.Background
		imgBg, iErr := assets.GetClickCaptchaAssetCache(bgFile)
		if iErr != nil {
			return canvas, iErr
		}
		img, dErr := toft.DecodingBinaryToImage(imgBg)
		if dErr != nil {
			return canvas, dErr
		}

		b := img.Bounds()
		m := image.NewNRGBA(b)

		curX, curY := cd.rangCutImage(params.Width, params.Height, img)
		draw.Draw(m, b, img, image.Point{X: curX, Y: curY}, draw.Src)
		canvas.distort(float64(toft.RandInt(5, 10)), float64(toft.RandInt(120, 200)))
		draw.Draw(m, canvas.Bounds(), canvas, image.Point{}, draw.Over)
		rc := m.SubImage(image.Rect(0, 0, params.Width, params.Height)).(*image.NRGBA)
		return rc, nil
	}

	if params.BackgroundDistort > 0 {
		canvas.distort(float64(toft.RandInt(5, 10)), float64(params.BackgroundDistort))
	}

	return canvas, nil
}

// CreateCanvas 创建画布
func (cd *Drawing) CreateCanvas(params DrawCanvas, isAlpha bool) (img *image.NRGBA) {
	width := params.Width
	height := params.Height
	img = image.NewNRGBA(image.Rect(0, 0, width, height))
	// 画背景
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if isAlpha {
				img.Set(x, y, color.Alpha{A: uint8(0)})
			} else {
				img.Set(x, y, color.RGBA{R: 255, G: 255, B: 255, A: 255})
			}
		}
	}
	return
}

// DrawTextImg 绘制文本的图片
func (cd *Drawing) DrawTextImg(dot DrawDot, params DrawCanvas) (*Palette, *AreaPoint, error) {
	// 绘制文本
	textColor, _ := toft.ParseHexColorToRGBA(dot.Color)
	var coArr = []color.RGBA{
		textColor,
	}
	textColor.A = cd.formatDiaphanous(params.TextAlpha)
	textImg := cd.DrawStrImg(dot, coArr, textColor)

	// 主画板
	var colorArr = []color.RGBA{
		textColor,
	}

	// 绘制阴影文本颜色
	shadowColorHex := "#101010"
	if params.TextShadowColor != "" {
		shadowColorHex = params.TextShadowColor
	}

	shadowColor, _ := toft.ParseHexColorToRGBA(shadowColorHex)
	if params.ShowTextShadow {
		colorArr = append(colorArr, shadowColor)
	}

	canvas := cd.CreateCanvasWithPalette(DrawCanvas{
		Width:  dot.Width + 10,
		Height: dot.Height + 10,
	}, colorArr)

	if params.ShowTextShadow {
		// 绘制阴影文本
		var shadowColorArr = []color.RGBA{
			shadowColor,
		}
		shadowImg := cd.DrawStrImg(dot, shadowColorArr, shadowColor)
		pointX := params.TextShadowPoint.X
		pointY := params.TextShadowPoint.Y
		draw.Draw(canvas, shadowImg.Bounds(), shadowImg, image.Point{X: pointX, Y: pointY}, draw.Over)
	}
	draw.Draw(canvas, textImg.Bounds(), textImg, image.Point{}, draw.Over)

	// 旋转效果
	canvas.Rotate(dot.Angle)

	// 扭曲效果
	if params.BackgroundDistort > 0 {
		canvas.distort(float64(toft.RandInt(5, 10)), float64(params.BackgroundDistort))
	}

	// 计算裁剪
	ap := cd.calcImageSpace(canvas)

	return canvas, ap, nil
}

// DrawStrImg 绘制字符图片
func (cd *Drawing) DrawStrImg(dot DrawDot, colorArr []color.RGBA, fc color.Color) *Palette {
	canvas := cd.CreateCanvasWithPalette(DrawCanvas{
		Width:  dot.Width + 10,
		Height: dot.Height + 10,
	}, colorArr)

	// 读字体数据
	fontBytes, err := assets.GetClickCaptchaAssetCache(dot.Font)
	if err != nil {
		return canvas
	}
	fontN, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return canvas
	}

	dc := freetype.NewContext()
	dc.SetDPI(float64(dot.FontDPI))
	dc.SetFont(fontN)
	dc.SetClip(canvas.Bounds())
	dc.SetDst(canvas)

	// 文字大小
	dc.SetFontSize(float64(dot.Size))
	dc.SetHinting(font.HintingFull)

	// 文字颜色
	fontColor := image.NewUniform(fc)
	dc.SetSrc(fontColor)

	// 画文本
	text := fmt.Sprintf("%s", dot.Text)

	pt := freetype.Pt(12, dot.Height-5) // 字出现的位置
	//if IsChineseChar(text) {
	//	pt = freetype.Pt(10, dot.Height) // 字出现的位置
	//}

	_, err = dc.DrawString(text, pt)
	if err != nil {
		return nil
	}

	return canvas
}

//formatDiaphanous 透明样式
func (cd *Drawing) formatDiaphanous(val float64) uint8 {
	a := math.Min(val, 1)
	alpha := a * 255
	return uint8(alpha)
}

// CreateCanvasWithPalette 创建带调色板的画布
func (cd *Drawing) CreateCanvasWithPalette(params DrawCanvas, colorArr []color.RGBA) *Palette {
	width := params.Width
	height := params.Height
	p := []color.Color{
		color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0x00},
	}

	for _, co := range colorArr {
		p = append(p, co)
	}

	return NewPalette(image.Rect(0, 0, width, height), p)
}

//calcImageSpace 计算图片空白区域
func (cd *Drawing) calcImageSpace(pa *Palette) *AreaPoint {
	nW := pa.Bounds().Max.X
	nH := pa.Bounds().Max.Y
	// 计算裁剪的最小及最大的坐标
	minX := nW
	maxX := 0
	minY := nH
	maxY := 0
	for x := 0; x < nW; x++ {
		for y := 0; y < nH; y++ {
			co := pa.At(x, y)
			if _, _, _, a := co.RGBA(); a > 0 {
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}

				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	minX = int(math.Max(0, float64(minX-2)))
	maxX = int(math.Min(float64(nW), float64(maxX+2)))
	minY = int(math.Max(0, float64(minY-2)))
	maxY = int(math.Min(float64(nH), float64(maxY+2)))

	return &AreaPoint{
		minX,
		maxX,
		minY,
		maxY,
	}
}

//rangCutImage 随机裁剪图片
func (cd *Drawing) rangCutImage(width int, height int, img image.Image) (int, int) {
	b := img.Bounds()
	iW := b.Max.X
	iH := b.Max.Y
	curX := 0
	curY := 0

	if iW-width > 0 {
		curX = toft.RandInt(0, iW-width)
	}
	if iH-height > 0 {
		curY = toft.RandInt(0, iH-height)
	}

	return curX, curY
}

//rangFillWithCircles 随机填充圆点
func (cd *Drawing) rangFillWithCircles(m *Palette, n, maxRadius int, circleCount int) {
	maxX := m.Bounds().Max.X
	maxY := m.Bounds().Max.Y
	for i := 0; i < n; i++ {
		colorIdx := uint8(toft.RandInt(1, circleCount-1))
		r := toft.RandInt(1, maxRadius)
		m.drawCircle(toft.RandInt(r, maxX-r), toft.RandInt(r, maxY-r), r, colorIdx)
	}
}

//drawLine 绘制线条
func (cd *Drawing) drawLine(m *Palette, num int, colorB []color.Color) {
	first := m.Bounds().Max.X / 10
	end := first * 9
	y := m.Bounds().Max.Y / 3
	for i := 0; i < num; i++ {
		point1 := Point{X: mRand.Intn(first), Y: mRand.Intn(y)}
		point2 := Point{X: mRand.Intn(first) + end, Y: mRand.Intn(y)}

		if i%2 == 0 {
			point1.Y = mRand.Intn(y) + y*2
			point2.Y = mRand.Intn(y)
		} else {
			point1.Y = mRand.Intn(y) + y*(i%2)
			point2.Y = mRand.Intn(y) + y*2
		}

		co := toft.GetRandomColorValueByRGBA(colorB)
		co.A = uint8(0xee)
		m.drawBeeline(point1, point2, co)
	}
}
