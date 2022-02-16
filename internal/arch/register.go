package arch

type RegisterValue = uint8

const (
	ZeroRegister    = 0
	SyscallRegister = 9
	StackRegister   = 14
	ReturnRegister  = 15
)
