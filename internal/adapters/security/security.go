package security

// func NewAgentID() string {
// 	// ULID/KSUID también sirven; por simplicidad:
// 	buf := make([]byte, 16)
// 	rand.Read(buf)
// 	return "agt_" + base64.RawURLEncoding.EncodeToString(buf) // corto y único
// }

// func NewAgentToken() (plain, prefix string) {
// 	buf := make([]byte, 32)
// 	rand.Read(buf)
// 	plain = "agt_sk_" + base64.RawURLEncoding.EncodeToString(buf)
// 	prefix = plain[:16] // primeros 16 para lookup (ajusta a 12–20)
// 	return
// }

// func HashTokenArgon2id(plain string) string {
// 	salt := make([]byte, 16)
// 	rand.Read(salt)
// 	sum := argon2.IDKey([]byte(plain), salt, 1, 64*1024, 2, 32)
// 	return "argon2id$" +
// 		base64.RawURLEncoding.EncodeToString(salt) + "$" +
// 		base64.RawURLEncoding.EncodeToString(sum)
// }

// func VerifyTokenArgon2id(hash, plain string) bool {
// 	parts := strings.Split(hash, "$")
// 	if len(parts) != 3 || parts[0] != "argon2id" {
// 		return false
// 	}
// 	salt, _ := base64.RawURLEncoding.DecodeString(parts[1])
// 	want, _ := base64.RawURLEncoding.DecodeString(parts[2])
// 	got := argon2.IDKey([]byte(plain), salt, 1, 64*1024, 2, 32)
// 	if len(got) != len(want) {
// 		return false
// 	}
// 	// constant-time compare
// 	var v byte
// 	for i := range got {
// 		v |= got[i] ^ want[i]
// 	}
// 	return v == 0
// }
