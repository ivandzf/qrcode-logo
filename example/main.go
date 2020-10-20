package main

import (
	"fmt"
	"image"
	"os"

	qrcodelogo "github.com/ivandzf/qrcode"
)

func main() {
	f, err := os.Open("example/me.png")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	s, err := qrcodelogo.EncodeToBase64("https://github.com/ivandzf/qrcode-logo", img, 500)
	if err != nil {
		panic(err)
	}

	fmt.Println("result: ", s)
}
