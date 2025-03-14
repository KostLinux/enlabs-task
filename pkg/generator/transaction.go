package generator

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// TransactionID generates a transaction ID with a 'tx_' prefix and a random component
func TransactionID() string {
	randomBytes := make([]byte, 12)
	_, err := rand.Read(randomBytes)
	if err != nil {
		// Fall back to a timestamp-based ID if random generation fails
		return fmt.Sprintf("tx_time_%d", time.Now().UnixNano())
	}

	return fmt.Sprintf("tx_%s", hex.EncodeToString(randomBytes))
}
