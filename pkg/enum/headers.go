package enum

type SourceType string

const (
	SourceTypeGame    SourceType = "game"
	SourceTypeServer  SourceType = "server"
	SourceTypePayment SourceType = "payment"
)

func (source SourceType) IsValid() bool {
	switch source {
	case SourceTypeGame, SourceTypeServer, SourceTypePayment:
		return true
	}
	return false
}
