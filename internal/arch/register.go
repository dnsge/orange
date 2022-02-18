package arch

type RegisterValue = uint8

const (
	ZeroRegister    RegisterValue = 0
	SyscallRegister RegisterValue = 9
	StackRegister   RegisterValue = 14
	ReturnRegister  RegisterValue = 15
)
