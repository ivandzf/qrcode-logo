package qrcodelogo

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/draw"
	"image/png"

	"github.com/skip2/go-qrcode"
)

// Encoder configuration
type Encoder struct {
	Content         string
	Size            int
	QRCodeImage     image.Image
	Logo            image.Image
	QRRecoveryLevel qrcode.RecoveryLevel
}

var defaultEncoder = Encoder{
	Size:            256,           // default size
	QRRecoveryLevel: qrcode.Medium, // default recovery level
}

// Encode the content to QR Code with logo
func Encode(content string, logo image.Image, size int) (*bytes.Buffer, error) {
	if logo == nil {
		return nil, errors.New("logo cannot nil")
	}

	e := defaultEncoder
	e.Content = content
	e.Logo = logo

	if size < 0 {
		e.Size = size
	}

	// create qrcode first
	qrCodeImage, err := e.createQRCode()
	if err != nil {
		return nil, err
	}

	// set qrcode image
	e.QRCodeImage = qrCodeImage

	// set logo
	result := e.overlayLogo()

	// encode the result to png
	pngEncoder := png.Encoder{CompressionLevel: png.BestCompression}

	var buf bytes.Buffer
	err = pngEncoder.Encode(&buf, result)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

// EncodeToBase64 encode the content to QR Code with logo, but return base64
func EncodeToBase64(content string, logo image.Image, size int) (string, error) {
	buf, err := Encode(content, logo, size)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func (e *Encoder) createQRCode() (image.Image, error) {
	encodedImg, err := qrcode.Encode(e.Content, e.QRRecoveryLevel, e.Size)
	if err != nil {
		return nil, err
	}

	qrImg, _, err := image.Decode(bytes.NewBuffer(encodedImg))
	if err != nil {
		return nil, err
	}

	return qrImg, nil
}

func (e *Encoder) overlayLogo() *image.RGBA {
	var (
		qrCodeImageBounds = e.QRCodeImage.Bounds()
		logoImageBounds   = e.Logo.Bounds()
		newImage          = image.NewRGBA(qrCodeImageBounds)
	)

	// because of qrcode generatorn return image black and white only
	// we need to copy qrcode image to new image with RGBA
	// FIXME: the cost of copy image is expensive
	draw.Draw(newImage, newImage.Bounds(), e.QRCodeImage, qrCodeImageBounds.Min, draw.Src)

	// set middle offset
	offset := qrCodeImageBounds.Max.X/2 - logoImageBounds.Max.X/2

	for x := 0; x < logoImageBounds.Max.X; x++ {
		for y := 0; y < logoImageBounds.Max.Y; y++ {
			// apply logo to the middle of QR Code
			newImage.Set(x+offset, y+offset, e.Logo.At(x, y))
		}
	}

	return newImage
}
