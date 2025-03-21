package qrcode

import (
	"encoding/base64"
	"fmt"

	"github.com/skip2/go-qrcode"
)

func CreateQrCode(s string) string {
	rs, err := qrcode.Encode(s, qrcode.Medium, 256)
	if err != nil {
		fmt.Println("error create QR cod", err)
	}
	return base64.StdEncoding.EncodeToString(rs)
}
