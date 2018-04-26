package wordcloud

import (
	"os"
	"log"
	"math"
	"flag"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"math/rand"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	//external library
	"github.com/golang/freetype/truetype"
	"github.com/elissalim/wordcloudgo/textprocessing"
)

var (
	width               = 800
	height              = 800
	boundsOfDrawnLabels = map[fixed.Rectangle26_6]bool{}
	dpi                 = flag.Float64("dpi", 300, "screen resolution in Dots Per Inch")
	fontFile            = flag.String("fontFile", "Roboto-Medium.ttf", "filename of the ttf font")
)

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func loadFont() *truetype.Font {
	fontBytes, err := ioutil.ReadFile(*fontFile)
	if err != nil {
		log.Println(err)
	}
	fontStyle, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Println(err)
	}
	return fontStyle
}

func getColor(size float64) color.Color {
	r := uint8(randInt(0, 40))
	g := uint8(randInt(50, 120))
	b := uint8(randInt(50, 150))
	alpha := 180 + size
	if alpha > 255 {
		alpha = 255
	}
	a := uint8(alpha)
	return color.RGBA{r, g, b, a}
}

func shortenedList(originalList textprocessing.PairList, cutoff int) textprocessing.PairList {
	shortenedList := originalList[:0]
	for _, v := range originalList {
		if v.Value >= cutoff {
			shortenedList = append(shortenedList, v)
		}
	}
	return shortenedList
}

func dotIsValid(x, y int) bool {
	return x >= 0 && x <= width && y >= 0 && y <= height
}

func canFitIn(newBound fixed.Rectangle26_6) bool {
	if boundOutsideImage(newBound) {
		return false
	}
	for bound := range boundsOfDrawnLabels {
		if colliding(bound, newBound) {
			return false
		}
	}
	return true
}

func colliding(a, b fixed.Rectangle26_6) bool {
	return a.Min.X <= b.Max.X && a.Max.X >= b.Min.X && a.Min.Y <= b.Max.Y && a.Max.Y >= b.Min.Y
}

func boundOutsideImage(bound fixed.Rectangle26_6) bool {
	return bound.Max.X.Round() >= width || bound.Max.Y.Round() >= bound.Min.Y.Round() * 2 || bound.Min.Y.Round() <= 0
}

func calibrateBound(bound fixed.Rectangle26_6, x, y fixed.Int26_6) fixed.Rectangle26_6 {
	bound.Min.X += x
	bound.Max.X += x
	bound.Min.Y += y
	bound.Max.Y += y
	return bound
}

func pickADot(i int) (x, y int) {
	index := float64(i)
	radiusIncrement := 0.15
	thetaIncrement := 0.1
	radius := 1.0 + index * radiusIncrement
	theta := 0.0 + index * thetaIncrement
	x = int(radius * math.Cos(theta))
	y = int(radius * math.Sin(theta))
	return x, y
}

func prepareAndDrawLabel(img *image.RGBA, f *truetype.Font, label string, size float64) {
	face := truetype.NewFace(f, &truetype.Options{
		Size: size,
		DPI:  *dpi,
	})
	x, y, index := 0, 0, 1
	for dotIsValid( x + width/2, y + height/2) {
		bound, _ := font.BoundString(face, label)
		calibrationX := fixed.Int26_6((x + width/2) * 64)
		calibrationY := fixed.Int26_6((y + height/2) * 64)
		bound = calibrateBound(bound, calibrationX, calibrationY)
		if canFitIn(bound) {
			boundsOfDrawnLabels[bound] = true
			d := &font.Drawer{
				Dst:  img,
				Src:  image.NewUniform(getColor(size)),
				Face: face,
				Dot:  fixed.Point26_6{calibrationX, calibrationY},
			}
			d.DrawString(label)
			break
		}
		index++
		x, y = pickADot(index)
	}
}

func WordCloud(processedText textprocessing.PairList) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	fontStyle := loadFont()
	cutoff := 9
	list := shortenedList(processedText, cutoff)
	shuffleIndex := rand.Perm(len(list))
	for _, v := range shuffleIndex {
		ratio := float64(list[v].Value) / float64(cutoff)
		size := ratio * 10
		prepareAndDrawLabel(img, fontStyle, list[v].Key, size)
	}
	file, err := os.Create("Word Cloud.png")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		log.Println(err)
	}
}
