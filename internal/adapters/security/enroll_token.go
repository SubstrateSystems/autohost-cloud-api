package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

const (
	enrollTokenPrefix = "autohost-enr_" // lo que quieras
	enrollTokenBytes  = 32              // 32 bytes => ~43 chars base64
)

func GenerateEnrollToken() (plain string, hash string, err error) {
	// 1. bytes aleatorios seguros
	buf := make([]byte, enrollTokenBytes)
	if _, err = rand.Read(buf); err != nil {
		return "", "", fmt.Errorf("rand.Read: %w", err)
	}

	// 2. los codificas en Base64 URL-safe (para que no tenga +, /, =)
	id := base64.RawURLEncoding.EncodeToString(buf)

	// 3. token final con prefijo
	plain = enrollTokenPrefix + id

	// 4. hash SHA-256 del token COMPLETO
	sum := sha256.Sum256([]byte(plain))
	hash = hex.EncodeToString(sum[:])

	return plain, hash, nil
}
