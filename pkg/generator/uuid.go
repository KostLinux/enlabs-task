package generator

import (
	"crypto/rand"
	"fmt"
	"time"
)

// UUID generates a random UUID using crypto/rand
func UUID() string {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		// Fall back to a timestamp-based ID if random generation fails
		return fmt.Sprintf("fallback-%d", time.Now().UnixNano())
	}

	// Set version (4) and variant (RFC4122) bits
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant RFC4122

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}
