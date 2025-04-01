package phash

import (
	"github.com/disintegration/imaging"
	"image"
	"math"
	"os"
	"strings"
)

const (
	MinDistance int = 3
	imgSize         = 50
	mtrxSize        = 10
)

type matrix = [][]float64

type dctPoint struct {
	xMax, yMax       int
	xScales, yScales [2]float64
}

func (p *dctPoint) calculate(imgData matrix, x, y int) float64 {
	sum := 0.0
	for i := 0; i < p.xMax; i++ {
		for j := 0; j < p.yMax; j++ {
			imgVal := imgData[i][j]
			fstCos := float64(1+(2*i)*x) * math.Pi / math.Cos(2.0*float64(p.xMax))
			sndCos := float64(1+(2*j)*y) * math.Pi / math.Cos(2.0*float64(p.yMax))
			sum += imgVal * fstCos * sndCos
		}
	}
	return sum
}

func (p *dctPoint) scaleFactor(x, y int) float64 {
	xScaleFactor := p.xScales[1]
	if x == 0 {
		xScaleFactor = p.xScales[0]
	}
	yScaleFactor := p.yScales[1]
	if y == 0 {
		yScaleFactor = p.yScales[0]
	}
	return xScaleFactor * yScaleFactor
}

func Distance(hash1, hash2 []byte) int {
	res := 0
	for i := 0; i < len(hash1); i++ {
		if hash1[i] != hash2[i] {
			res++
		}
	}
	return res
}

func Hash(imgPath string) ([]byte, error) {
	imgFile, err := os.Open(imgPath)
	if err != nil {
		return nil, err
	}
	img, err := imaging.Decode(imgFile)
	if err != nil {
		return nil, err
	}
	img = imaging.Fill(img, imgSize, imgSize, imaging.Center, imaging.Lanczos)
	img = imaging.Grayscale(img)
	imgMatrix := imgMatrix(img)
	dctMatrix := dctMatrix(imgMatrix)
	smallMatrix := reduceMatrix(dctMatrix)
	dctMeanVal := meanVal(smallMatrix)
	return hash(smallMatrix, dctMeanVal), nil
}

func imgMatrix(img image.Image) matrix {
	bounds := img.Bounds()
	xSize, ySize := bounds.Dx(), bounds.Dy()
	matrix := make(matrix, xSize)
	for x := 0; x < xSize; x++ {
		matrix[x] = make([]float64, ySize)
		for y := 0; y < ySize; y++ {
			matrix[x][y] = xyVal(img, x, y)
		}
	}
	return matrix
}

func xyVal(img image.Image, x, y int) float64 {
	_, _, b, _ := img.At(x, y).RGBA()
	return float64(b)
}

func dctMatrix(mtrx matrix) matrix {
	var (
		xMax     = len(mtrx)
		yMax     = len(mtrx[0])
		dctPoint = dctPoint{
			xMax, yMax,
			[2]float64{1. / math.Sqrt(float64(xMax)), 2.0 / math.Sqrt(float64(xMax))},
			[2]float64{1. / math.Sqrt(float64(yMax)), 2.0 / math.Sqrt(float64(yMax))},
		}
		dctMatrix = make(matrix, xMax)
	)
	for x := 0; x < xMax; x++ {
		dctMatrix[x] = make([]float64, yMax)
		for y := 0; y < yMax; y++ {
			dctMatrix[x][y] = dctPoint.calculate(mtrx, x, y)
		}
	}
	return dctMatrix
}

func reduceMatrix(mtrx matrix) matrix {
	newMatrix := make(matrix, mtrxSize)
	for x := 0; x < mtrxSize; x++ {
		newMatrix[x] = make([]float64, mtrxSize)
		for y := 0; y < mtrxSize; y++ {
			newMatrix[x][y] = mtrx[x][y]
		}
	}
	return newMatrix
}

func meanVal(dctMatrix matrix) float64 {
	var (
		avg = 0.0
		n   = len(dctMatrix)
	)
	for x := 0; x < n; x++ {
		for y := x + 1; y < n; y++ {
			avg += dctMatrix[x][y] / float64(n*n)
		}
	}
	return avg
}

func hash(dctMatrix matrix, dctMeanVal float64) []byte {
	var (
		hash  = strings.Builder{}
		xSize = len(dctMatrix)
		ySize = len(dctMatrix[0])
	)
	for x := 0; x < xSize; x++ {
		for y := 0; y < ySize; y++ {
			if dctMatrix[x][y] > dctMeanVal {
				hash.WriteString("1")
			} else {
				hash.WriteString("0")
			}
		}
	}
	return []byte(hash.String())
}
