package domain

type BotStatus string

const (
	BotUnknown  BotStatus = ""
	BotIsIdle   BotStatus = "idle"
	BotIsActive BotStatus = "active"
)

func (s BotStatus) String() string {
	switch s {
	case BotIsIdle, BotIsActive:
		return string(s)
	default:
		return ""
	}
}

func ToBotStatus(status string) BotStatus {
	switch status {
	case BotIsIdle.String():
		return BotIsIdle
	case BotIsActive.String():
		return BotIsActive
	default:
		return BotUnknown
	}
}
