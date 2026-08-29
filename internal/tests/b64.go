package tests

import "encoding/base64"

func base64Decode(dst, src []byte) (int, error) {
	out, err := base64.RawStdEncoding.DecodeString(string(trimPad(src)))
	if err != nil {
		return 0, err
	}
	return copy(dst, out), nil
}

func trimPad(b []byte) []byte {
	for len(b) > 0 && b[len(b)-1] == '=' {
		b = b[:len(b)-1]
	}
	return b
}
