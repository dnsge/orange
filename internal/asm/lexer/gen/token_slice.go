package main

var (
	NoSlice = TokenSlice{0, 0}
)

type TokenSlice struct {
	Start int
	End   int
}

func Slice(start, end int) TokenSlice {
	return TokenSlice{start, end}
}
