package usecase

import (
	"fmt"
	"github.com/skip2/go-qrcode"
	"os"
	"path/filepath"
)

func createQrCode(content string, fileName string, folderPath string) (string, error) {

	fullFileName := fileName + ".png"
	fullFileNamePath := folderPath + "/" + fullFileName
	ex, errrr := os.Executable()
	if errrr != nil {
		panic(errrr)
	}
	exPath := filepath.Dir(ex)
	fmt.Println(exPath)

	err := qrcode.WriteFile(content, qrcode.Medium, 256, fullFileNamePath)
	if err != nil {
		return "", err
	}

	return fullFileName, nil
}
