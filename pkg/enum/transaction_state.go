package enum

// TransactionState represents the state of a transaction
type State string

const (
	TransactionStateWin  State = "win"
	TransactionStateLose State = "lose"
)

// IsValid checks if the transaction state is valid
func (state State) IsTransactionValid() bool {
	switch state {
	case TransactionStateWin, TransactionStateLose:
		return true
	}
	return false
}

// ParseTransactionState safely converts a string to a TransactionState
func ParseTransactionState(transaction string) (State, bool) {
	state := State(transaction)
	return state, state.IsTransactionValid()
}
