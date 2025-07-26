package czdb

import "encoding/base64"

type Decrypted struct {
	keyBytes []byte
}

func NewDecrypted(key string) *Decrypted {
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil
	}
	return &Decrypted{keyBytes: keyBytes}
}

func (d Decrypted) decrypt(data []byte) []byte {
	result := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		result[i] = data[i] ^ d.keyBytes[i%len(d.keyBytes)]
	}
	return result
}
