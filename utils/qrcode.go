package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/ean"
	"github.com/skip2/go-qrcode"
	"image/png"
	"os"
)

func Generate(content string, compress bool, size int) (qrCode []byte, err error) {
	if compress {
		content, err = GzipEncode(content)
		if err != nil {
			return nil, err
		}
	}

	encode, err := qrcode.Encode(content, qrcode.Low, size)
	if err != nil {
		return nil, err
	}
	return encode, err
}

// GzipEncode 使用gzip算法压缩字符串并输出base64编码的字符串
func GzipEncode(input string) (string, error) {
	var b bytes.Buffer // 创建一个缓冲区用于存储压缩数据

	zw := gzip.NewWriter(&b)
	defer zw.Close()

	// 将字符串写入gzip writer
	if _, err := zw.Write([]byte(input)); err != nil {
		return "", fmt.Errorf("failed to write input to gzip: %w", err)
	}

	if err := zw.Flush(); err != nil {
		return "", fmt.Errorf("failed to flush gzip writer: %w", err)
	}

	// 获取压缩后的字节切片
	compressedBytes := b.Bytes()

	// 对压缩数据进行base64编码
	encodedStr := base64.StdEncoding.EncodeToString(compressedBytes)

	return encodedStr, nil
}

func GenerateBarcode(data string, filename string) error {
	// 要编码的EAN-13号码
	code := "0123456789012"

	// 创建一个EAN-13类型的条形码
	bc, err := ean.Encode(code)
	if err != nil {
		panic(err)
	}

	// 将条形码绘制到一个图像上
	img, err := barcode.Scale(bc, 300, 150)
	if err != nil {
		panic(err)
	}

	// 创建输出文件
	outFile, err := os.Create("barcode.png")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	// 将图像以PNG格式写入文件
	if err := png.Encode(outFile, img); err != nil {
		panic(err)
	}

	fmt.Println("条形码已成功生成并保存为 barcode.png")
	return nil
}
