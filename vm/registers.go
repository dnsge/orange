package vm

type registerFile [16]uint64

func initRegisterFile() registerFile {
	return [16]uint64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
}

var zeroPtr uint64 = 0

func (r *registerFile) Get(regNum uint8) uint64 {
	if regNum == 0 {
		// r0 always returns zero
		return 0
	} else {
		return r[regNum]
	}
}

func (r *registerFile) Set(regNum uint8, val uint64) {
	if regNum == 0 {
		// Silently ignore writes to r0
		return
	} else {
		r[regNum] = val
	}
}

func (r *registerFile) Ref(regNum uint8) *uint64 {
	if regNum == 0 {
		// protect the zero register
		// not thread-safe
		zeroPtr = 0
		return &zeroPtr
	}
	return &r[regNum]
}
