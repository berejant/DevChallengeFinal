package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"net/http"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.POST("/api/image-input", processInputImage)
	r.GET("/healthcheck", func(c *gin.Context) {
		c.String(http.StatusOK, "health")
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}

type ImageInputRequest struct {
	MinLevel     int    `json:"min_level" binding:"min=0,max=100"`
	Based64Image string `json:"image" binding:"required"`
}

type CellMineResult struct {
	X     int `json:"x" binding:"required"`
	Y     int `json:"y" binding:"required"`
	Level int `json:"level" binding:"required"`
}

func processInputImage(c *gin.Context) {
	var imageInputRequest ImageInputRequest
	if err := c.ShouldBindJSON(&imageInputRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	image, err := readBased64InputImage(imageInputRequest.Based64Image)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var response struct {
		Mines []CellMineResult `json:"mines"`
	}
	response.Mines = findMinesInImage(image, imageInputRequest.MinLevel)
	c.JSON(http.StatusOK, response)
}

func readBased64InputImage(dataUri string) (image.Image, error) {
	if dataUri[0:5] != "data:" {
		return nil, errors.New("expected for base64 DataURI string")
	}

	if dataUri[0:22] != "data:image/png;base64," {
		return nil, errors.New(
			"Expected for Data URI with based64 encoded PNG, receive: " + dataUri[11:14],
		)
	}

	imageBytes, err := base64.StdEncoding.DecodeString(dataUri[22:])
	if err != nil {
		return nil, err
	}

	inputImage, err := png.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, err
	}

	if inputImage.ColorModel() == color.GrayModel {
		return inputImage, nil
	}

	// convert to grey
	grayImage := image.NewGray(inputImage.Bounds())
	for y := inputImage.Bounds().Min.Y; y < inputImage.Bounds().Max.Y; y++ {
		for x := inputImage.Bounds().Min.X; x < inputImage.Bounds().Max.X; x++ {
			grayImage.Set(x, y, inputImage.At(x, y))
		}
	}

	return grayImage, nil
}

func findMinesInImage(img image.Image, minLevel int) []CellMineResult {
	horizontalLines := findHorizontalWhileLines(img)
	verticalLines := findVerticalWhileLines(img)

	var maxY int
	var maxX int
	var rect image.Rectangle

	cellMineResults := make([]CellMineResult, 0)

	for cellIndexY, minY := range horizontalLines {
		if cellIndexY+1 == len(horizontalLines) { // skip last grid line
			continue
		}

		maxY = horizontalLines[cellIndexY+1]
		for cellIndexX, minX := range verticalLines {
			if cellIndexX+1 == len(verticalLines) { // skip last grid line
				continue
			}

			maxX = verticalLines[cellIndexX+1]
			rect = image.Rect(minX+1, minY+1, maxX, maxY)

			darkLevel := checkCellForMines(img, rect, minLevel)

			if minLevel <= darkLevel {
				cellMineResults = append(cellMineResults, CellMineResult{
					X:     cellIndexX,
					Y:     cellIndexY,
					Level: darkLevel,
				})
			}
		}
	}

	return cellMineResults
}

func checkCellForMines(img image.Image, rect image.Rectangle, minLevel int) int {
	colorSum := 0

	var colorCode uint32
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			colorCode, _, _, _ = img.At(x, y).RGBA()
			colorCode = colorCode >> 8

			colorSum += int(colorCode)
		}
	}

	pixelCount := rect.Size().X * rect.Size().Y
	averageColorCode := float64(colorSum) / float64(pixelCount)

	// 255 - absolute white
	// 0 - absolute black (dark)
	darkLevel := 100 - int(math.Round(100.0*averageColorCode/255.0))

	return darkLevel
}

func findHorizontalWhileLines(img image.Image) []int {
	var horizontalLines []int
	var prevWhiteLineAtY int

	if isHorizontalWhiteLine(img, 0) {
		horizontalLines = append(horizontalLines, 0)
		prevWhiteLineAtY = 0
	} else {
		horizontalLines = append(horizontalLines, -1)
		prevWhiteLineAtY = -1
	}

	for y := prevWhiteLineAtY + 1; y < img.Bounds().Dy(); y++ {
		if isHorizontalWhiteLine(img, y) {
			if prevWhiteLineAtY+1 != y {
				horizontalLines = append(horizontalLines, y)
			}
			prevWhiteLineAtY = y
		}
	}

	if prevWhiteLineAtY+1 != img.Bounds().Dx() {
		horizontalLines = append(horizontalLines, img.Bounds().Dx())
	}

	return horizontalLines
}

func findVerticalWhileLines(img image.Image) []int {
	var verticalLines []int
	var prevWhiteLineAtX int

	if isVerticalWhiteLine(img, 0) {
		verticalLines = append(verticalLines, 0)
		prevWhiteLineAtX = 0
	} else {
		verticalLines = append(verticalLines, -1)
		prevWhiteLineAtX = -1
	}

	for x := 1; x < img.Bounds().Dx(); x++ {
		if isVerticalWhiteLine(img, x) {
			if prevWhiteLineAtX+1 != x {
				verticalLines = append(verticalLines, x)
			}
			prevWhiteLineAtX = x
		}
	}
	if prevWhiteLineAtX+1 != img.Bounds().Dx() {
		verticalLines = append(verticalLines, img.Bounds().Dx())
	}

	return verticalLines
}

func isHorizontalWhiteLine(img image.Image, y int) bool {
	maxX := 1000
	if img.Bounds().Max.X < maxX {
		maxX = img.Bounds().Max.X
	}

	var colorCode uint32
	for x := img.Bounds().Min.X; x < maxX; x++ {
		colorCode, _, _, _ = img.At(x, y).RGBA()
		// 255 in 32-bit is 65535
		if colorCode < 65280 {
			return false
		}
	}

	return true
}

func isVerticalWhiteLine(img image.Image, x int) bool {
	maxY := 1000
	if img.Bounds().Max.Y < maxY {
		maxY = img.Bounds().Max.Y
	}

	var colorCode uint32
	for y := img.Bounds().Min.Y; y < maxY; y++ {
		colorCode, _, _, _ = img.At(x, y).RGBA()
		// 255 in 32-bit is 65535
		if colorCode != 65535 {
			return false
		}
	}

	return true
}
