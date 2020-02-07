package animal
//go:generate stringer -type=MessageStatus
type MessageStatus int

const (
	Sent MessageStatus = iota
	Received
	Rejected
)
