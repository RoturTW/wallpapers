package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"os"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
)

func serveSizedPreview(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(400, gin.H{"error": "Name is required"})
		return
	}

	filePath := "static/" + name
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(404, gin.H{"error": "File not found"})
		return
	}

	f, err := os.Open(filePath)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to open file"})
		return
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to decode image"})
		return
	}

	resizedImg := resizeImage(img, 200)

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, resizedImg, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to encode image"})
		return
	}

	c.Data(200, "image/jpeg", buf.Bytes())
}

func resizeImage(img image.Image, height int) image.Image {
	bounds := img.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()

	aspectRatio := float64(originalWidth) / float64(originalHeight)
	newWidth := int(float64(height) * aspectRatio)

	return imaging.Resize(img, newWidth, height, imaging.Lanczos)
}
