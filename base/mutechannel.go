package base

const (
	MUTE_CHANNEL = "MuteChannel"
)

type MuteChannel interface {
	Toggle(onOff bool)
}
