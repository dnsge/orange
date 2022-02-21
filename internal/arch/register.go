package arch

type RegisterValue = uint8

const (
	ZeroRegister          RegisterValue = 0
	SyscallResultRegister RegisterValue = 7
	SyscallErrorRegister  RegisterValue = 8
	SyscallRegister       RegisterValue = 9
	StackRegister         RegisterValue = 14
	ReturnRegister        RegisterValue = 15
)

const (
	ENO_IO                = 1
	ENO_BadFileDescriptor = 2
)
