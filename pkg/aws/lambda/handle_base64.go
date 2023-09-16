package lambda

import "encoding/base64"

// HandleBase64 tries to decode the provided byte slice.
func HandleBase64(src []byte) []byte {
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	if decLen, err := base64.StdEncoding.Decode(decoded, src); err == nil {
		return decoded[:decLen]
	}
	return src
}
