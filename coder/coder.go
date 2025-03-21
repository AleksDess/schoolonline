package coder

import (
	"encoding/base64"
	"schoolonline/config"
	"strings"
)

func cod3byte(a []byte, key string) []byte {

	nn := int(a[0])*256*256 + int(a[1])*256 + int(a[2])
	n1 := nn & 63
	n2 := nn >> 6 & 63
	n3 := nn >> 12 & 63
	n4 := nn >> 18 & 63
	res := make([]byte, 4)
	res[0] = key[n4]
	res[1] = key[n3]
	res[2] = key[n2]
	res[3] = key[n1]

	return res
}

func cod2byte(a []byte, key string) []byte {

	nn := int(a[0])*256*256 + int(a[1])*256
	n2 := nn >> 6 & 63
	n3 := nn >> 12 & 63
	n4 := nn >> 18 & 63
	res := make([]byte, 4)
	res[0] = key[n4]
	res[1] = key[n3]
	res[2] = key[n2]
	res[3] = 60

	return res
}

func cod1byte(a []byte, key string) []byte {

	nn := int(a[0]) * 256 * 256
	n3 := nn >> 12 & 63
	n4 := nn >> 18 & 63
	res := make([]byte, 4)
	res[0] = key[n4]
	res[1] = key[n3]
	res[2] = 62
	res[3] = 60

	return res
}

func Code(b []byte) []byte {

	key := config.C.CoderKey

	res := make([]byte, 0)
	n := len(b) / 3
	os := len(b) % 3
	for i := 0; i < n; i++ {
		bn := b[i*3 : i*3+3]
		res = append(res, cod3byte(bn, key)...)
	}
	if os == 2 {
		res = append(res, cod2byte(b[n*3:], key)...)

	}
	if os == 1 {
		res = append(res, cod1byte(b[n*3:], key)...)
	}

	return res
}

func uncod(a byte, key string) int {
	for i, j := range key {
		if a == byte(j) {
			return i
		}
	}
	return 0
}

func Decod(b []byte) []byte {

	key := config.C.CoderKey

	n := (len(b) / 4) - 1
	res := make([]byte, 0)
	if len(b) < 4 {
		return res
	}
	for i := 0; i < n; i++ {
		bn := b[i*4 : i*4+4]
		nn := uncod(bn[0], key)*64*64*64 + uncod(bn[1], key)*64*64 + uncod(bn[2], key)*64 + uncod(bn[3], key)
		n1 := byte(nn >> 16 & 255)
		n2 := byte(nn >> 8 & 255)
		n3 := byte(nn & 255)
		res = append(res, n1, n2, n3)
	}
	bn := b[n*4:]
	if bn[3] != 60 {
		nn := uncod(bn[0], key)*64*64*64 + uncod(bn[1], key)*64*64 + uncod(bn[2], key)*64 + uncod(bn[3], key)
		n1 := byte(nn >> 16 & 255)
		n2 := byte(nn >> 8 & 255)
		n3 := byte(nn & 255)
		res = append(res, n1, n2, n3)
	}
	if bn[3] == 60 && bn[2] != 62 {
		nn := uncod(bn[0], key)*64*64*64 + uncod(bn[1], key)*64*64 + uncod(bn[2], key)*64 + uncod(bn[3], key)
		n1 := byte(nn >> 16 & 255)
		n2 := byte(nn >> 8 & 255)
		res = append(res, n1, n2)
	}
	if bn[3] == 60 && bn[2] == 62 {
		nn := uncod(bn[0], key)*64*64*64 + uncod(bn[1], key)*64*64 + uncod(bn[2], key)*64 + uncod(bn[3], key)
		n1 := byte(nn >> 16 & 255)
		res = append(res, n1)
	}

	return res
}

func CreateCodeName(city, driver string) string {
	f := city + "|" + driver
	f = base64.URLEncoding.EncodeToString([]byte(f))
	return f
}

func UnCodeName(s string) (city, driver string) {
	s = strings.ReplaceAll(s, ".HTML", "")
	g, _ := base64.URLEncoding.DecodeString(s)
	f := string(g)
	r := strings.Split(f, "|")
	if len(r) != 2 {
		return
	}
	city = r[0]
	driver = r[1]
	return
}
