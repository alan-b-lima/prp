package hash

import (
	"golang.org/x/crypto/bcrypt"
)

func Hash(data []byte) [60]byte {
	ingest := To72Bytes(data)

	digest, err := bcrypt.GenerateFromPassword(ingest, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return [60]byte(digest)
}

func Compare(hash, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, To72Bytes(password))
	return err == nil
}

func To72Bytes(data []byte) []byte {
	const size = 72
	result := make([]byte, size)

	if len(data) < size {
		for i := 0; i < size; i += len(data) {
			copy(result[i:min(i+len(data), len(result))], data)
		}
	} else {
		for i, d := range data {
			result[i%size] ^= d
		}
	}

	return result
}
