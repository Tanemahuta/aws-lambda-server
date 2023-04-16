package aws

import "encoding/base64"

// HandleBase64 tries to decode the provided byte slice.
func HandleBase64(src []byte) []byte {
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	if decLen, err := base64.StdEncoding.Decode(decoded, src); err == nil {
		return decoded[:decLen]
	}
	return src
}

// HandleBase64String tries to decode the provided string behind the ptr.
func HandleBase64String(src *string) *string {
	if src != nil {
		str := (string)(HandleBase64([]byte(*src)))
		src = &str
	}
	return src
}
