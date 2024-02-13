package common

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"log"
	"strings"
)

func GenBufferQRPDF(text string) bytes.Buffer {

	base64Image := "data:image/png;base64,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"

	imageData, err := base64.StdEncoding.DecodeString(strings.TrimSpace(base64Image))
	if err != nil {
		log.Fatalf("Error decoding base64 image: %v", err)
	}

	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		log.Fatalf("Error decoding image: %v", err)
	}

	// Convert the image to JPEG format (you may skip this step if the image is already JPEG)
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, nil)
	if err != nil {
		log.Fatalf("Error encoding image to JPEG: %v", err)
	}

	return buf
}
