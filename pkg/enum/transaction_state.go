package enum

// TransactionState represents the state of a transaction
type TransactionState string

const (
	TransactionStateWin  TransactionState = "win"
	TransactionStateLose TransactionState = "lose"
)

// IsValid checks if the transaction state is valid
func (state TransactionState) IsTransactionValid() bool {
	switch state {
	case TransactionStateWin, TransactionStateLose:
		return true
	}

	return false
}

// ParseTransactionState safely converts a string to a TransactionState
func ValidateTransactionState(value string) (TransactionState, bool) {
	state := TransactionState(value)
	return state, state.IsTransactionValid()
}
