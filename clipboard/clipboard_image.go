package clipboard

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"os"
	"strings"

	"golang.design/x/clipboard"
)

func Init() error {
	return clipboard.Init()
}

func ReadText() (string, error) {
	data := clipboard.Read(clipboard.FmtText)
	return string(data), nil
}

func ReadImage() (image.Image, error) {
	data := clipboard.Read(clipboard.FmtImage)
	if len(data) == 0 {
		return nil, fmt.Errorf("no image in clipboard")
	}
	
	img, _, err := image.Decode(strings.NewReader(string(data)))
	return img, err
}

func WriteText(text string) error {
	clipboard.Write(clipboard.FmtText, []byte(text))
	return nil
}

func HasImage() bool {
	data := clipboard.Read(clipboard.FmtImage)
	return len(data) > 0
}

func SaveImageToFile(filename string) error {
	img, err := ReadImage()
	if err != nil {
		return err
	}
	
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	return png.Encode(file, img)
}

func GetImageAsBase64() (string, error) {
	data := clipboard.Read(clipboard.FmtImage)
	if len(data) == 0 {
		return "", fmt.Errorf("no image")
	}
	
	return base64.StdEncoding.EncodeToString(data), nil
}
