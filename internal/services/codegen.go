// internal/services/codegen.go
package services

import (
	"crypto/rand"
	"encoding/binary"
)

// generate4DigitCode генерирует криптографически стойкий код в диапазоне 1000..9999.
func generate4DigitCode() (int, error) {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return 0, err
	}
	// Берём 64-битное число и приводим к диапазону 0..8999, затем +1000 -> 1000..9999
	n := int(binary.BigEndian.Uint64(b[:]) % 9000)
	return 1000 + n, nil
}
