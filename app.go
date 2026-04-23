package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx          context.Context
	stitchResult image.Image
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// UploadImages handles image uploads
func (a *App) UploadImages(images []string) map[string]interface{} {
	// Simply return the images as they are already base64 encoded
	return map[string]interface{}{
		"success": true,
		"images":  images,
	}
}

// StitchImages stitches images in different modes
func (a *App) StitchImages(imageData []string, mode string) map[string]interface{} {
	// Decode images
	var images []image.Image
	for _, data := range imageData {
		// Extract base64 data
		parts := strings.Split(data, ",")
		if len(parts) != 2 {
			return map[string]interface{}{
				"success": false,
				"error":   "无效的图片数据",
			}
		}

		// Decode base64
		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"error":   "解码图片数据失败",
			}
		}

		// Decode image
		img, _, err := image.Decode(bytes.NewReader(decoded))
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"error":   "解码图片失败",
			}
		}

		images = append(images, img)
	}

	// Calculate dimensions and draw images based on mode
	var width, height int
	var newImage *image.RGBA

	switch mode {
	case "horizontal":
		if len(images) < 2 {
			return map[string]interface{}{
				"success": false,
				"error":   "至少需要2张图片",
			}
		}
		// Calculate dimensions
		height = images[0].Bounds().Dy()
		for _, img := range images {
			width += img.Bounds().Dx()
			if img.Bounds().Dy() > height {
				height = img.Bounds().Dy()
			}
		}
		// Create new image
		newImage = image.NewRGBA(image.Rect(0, 0, width, height))
		// Draw images
		currentX, currentY := 0, 0
		for _, img := range images {
			draw.Draw(newImage, image.Rect(currentX, currentY, currentX+img.Bounds().Dx(), currentY+img.Bounds().Dy()), img, image.Point{}, draw.Src)
			currentX += img.Bounds().Dx()
		}

	case "vertical":
		if len(images) < 2 {
			return map[string]interface{}{
				"success": false,
				"error":   "至少需要2张图片",
			}
		}
		// Calculate dimensions
		width = images[0].Bounds().Dx()
		for _, img := range images {
			height += img.Bounds().Dy()
			if img.Bounds().Dx() > width {
				width = img.Bounds().Dx()
			}
		}
		// Create new image
		newImage = image.NewRGBA(image.Rect(0, 0, width, height))
		// Draw images
		currentX, currentY := 0, 0
		for _, img := range images {
			draw.Draw(newImage, image.Rect(currentX, currentY, currentX+img.Bounds().Dx(), currentY+img.Bounds().Dy()), img, image.Point{}, draw.Src)
			currentY += img.Bounds().Dy()
		}

	case "grid-2x2":
		if len(images) != 4 {
			return map[string]interface{}{
				"success": false,
				"error":   "2x2网格模式需要4张图片",
			}
		}
		// Calculate dimensions
		maxWidth := 0
		maxHeight := 0
		for _, img := range images {
			if img.Bounds().Dx() > maxWidth {
				maxWidth = img.Bounds().Dx()
			}
			if img.Bounds().Dy() > maxHeight {
				maxHeight = img.Bounds().Dy()
			}
		}
		width = maxWidth * 2
		height = maxHeight * 2
		// Create new image
		newImage = image.NewRGBA(image.Rect(0, 0, width, height))
		// Draw images in 2x2 grid
		draw.Draw(newImage, image.Rect(0, 0, maxWidth, maxHeight), images[0], image.Point{}, draw.Src)
		draw.Draw(newImage, image.Rect(maxWidth, 0, maxWidth*2, maxHeight), images[1], image.Point{}, draw.Src)
		draw.Draw(newImage, image.Rect(0, maxHeight, maxWidth, maxHeight*2), images[2], image.Point{}, draw.Src)
		draw.Draw(newImage, image.Rect(maxWidth, maxHeight, maxWidth*2, maxHeight*2), images[3], image.Point{}, draw.Src)

	case "grid-3x1":
		if len(images) != 3 {
			return map[string]interface{}{
				"success": false,
				"error":   "3x1网格模式需要3张图片",
			}
		}
		// Calculate dimensions
		width = images[0].Bounds().Dx()
		for _, img := range images {
			height += img.Bounds().Dy()
			if img.Bounds().Dx() > width {
				width = img.Bounds().Dx()
			}
		}
		// Create new image
		newImage = image.NewRGBA(image.Rect(0, 0, width, height))
		// Draw images in 3x1 grid
		currentY := 0
		for _, img := range images {
			draw.Draw(newImage, image.Rect(0, currentY, width, currentY+img.Bounds().Dy()), img, image.Point{}, draw.Src)
			currentY += img.Bounds().Dy()
		}

	case "grid-1x3":
		if len(images) != 3 {
			return map[string]interface{}{
				"success": false,
				"error":   "1x3网格模式需要3张图片",
			}
		}
		// Calculate dimensions
		height = images[0].Bounds().Dy()
		for _, img := range images {
			width += img.Bounds().Dx()
			if img.Bounds().Dy() > height {
				height = img.Bounds().Dy()
			}
		}
		// Create new image
		newImage = image.NewRGBA(image.Rect(0, 0, width, height))
		// Draw images in 1x3 grid
		currentX := 0
		for _, img := range images {
			draw.Draw(newImage, image.Rect(currentX, 0, currentX+img.Bounds().Dx(), height), img, image.Point{}, draw.Src)
			currentX += img.Bounds().Dx()
		}

	default:
		return map[string]interface{}{
			"success": false,
			"error":   "不支持的拼接模式",
		}
	}

	// Save result
	a.stitchResult = newImage

	// Convert to base64
	var buf bytes.Buffer
	err := png.Encode(&buf, newImage)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   "编码结果图片失败",
		}
	}

	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())

	return map[string]interface{}{
		"success": true,
		"image":   fmt.Sprintf("data:image/png;base64,%s", base64Str),
	}
}

// DownloadResult downloads the stitched image
func (a *App) DownloadResult() map[string]interface{} {
	if a.stitchResult == nil {
		return map[string]interface{}{
			"success": false,
			"error":   "没有可用的拼接图片",
		}
	}

	// Create temporary file
	tempFile, err := os.CreateTemp("", "stitched-*.png")
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   "创建临时文件失败",
		}
	}
	defer tempFile.Close()

	// Encode image to file
	err = png.Encode(tempFile, a.stitchResult)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   "编码图片失败",
		}
	}

	// Show save dialog
	savePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultFilename: "stitched.png",
		Filters: []runtime.FileFilter{
			{
				Pattern:     "*.png",
				DisplayName: "PNG Images",
			},
		},
	})

	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   "保存对话框已取消",
		}
	}

	// Copy file
	src, err := os.Open(tempFile.Name())
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   "打开临时文件失败",
		}
	}
	defer src.Close()

	dst, err := os.Create(savePath)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   "创建输出文件失败",
		}
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   "保存文件失败",
		}
	}

	return map[string]interface{}{
		"success": true,
	}
}
