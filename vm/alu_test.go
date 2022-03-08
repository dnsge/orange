package vm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestALU_Equal(t *testing.T) {
	alu := newALU()
	alu.SUB(5, 5)
	assert.True(t, alu.Equal())
	alu.SUB(5, 6)
	assert.False(t, alu.Equal())
	alu.SUB(6, 5)
	assert.False(t, alu.Equal())
}

func TestALU_NotEqual(t *testing.T) {
	alu := newALU()
	alu.SUB(5, 6)
	assert.True(t, alu.NotEqual())
	alu.SUB(5, 5)
	assert.False(t, alu.NotEqual())
	alu.SUB(6, 5)
	assert.True(t, alu.NotEqual())
}

func TestALU_LessThan(t *testing.T) {
	alu := newALU()
	alu.SUB(5, 10)
	assert.True(t, alu.LessThan())
	alu.SUB(5, 5)
	assert.False(t, alu.LessThan())
	alu.SUB(10, 5)
	assert.False(t, alu.LessThan())
	applyWithSigned(alu.SUB, -5, -10)
	assert.False(t, alu.LessThan())
	applyWithSigned(alu.SUB, -10, -5)
	assert.True(t, alu.LessThan())
	applyWithSigned(alu.SUB, -10, -10)
	assert.False(t, alu.LessThan())
}

func TestALU_GreaterThan(t *testing.T) {
	alu := newALU()
	alu.SUB(5, 10)
	assert.False(t, alu.GreaterThan())
	alu.SUB(5, 5)
	assert.False(t, alu.GreaterThan())
	alu.SUB(10, 5)
	assert.True(t, alu.GreaterThan())
	applyWithSigned(alu.SUB, -5, -10)
	assert.True(t, alu.GreaterThan())
	applyWithSigned(alu.SUB, -10, -5)
	assert.False(t, alu.GreaterThan())
	applyWithSigned(alu.SUB, -10, -10)
	assert.False(t, alu.GreaterThan())
}

func TestALU_GreaterThanEqual(t *testing.T) {
	alu := newALU()
	alu.SUB(5, 10)
	assert.False(t, alu.GreaterThanEqual())
	alu.SUB(5, 5)
	assert.True(t, alu.GreaterThanEqual())
	alu.SUB(10, 5)
	assert.True(t, alu.GreaterThanEqual())
	applyWithSigned(alu.SUB, -5, -10)
	assert.True(t, alu.GreaterThanEqual())
	applyWithSigned(alu.SUB, -10, -5)
	assert.False(t, alu.GreaterThanEqual())
	applyWithSigned(alu.SUB, -10, -10)
	assert.True(t, alu.GreaterThanEqual())
}

func TestALU_LessThanEqual(t *testing.T) {
	alu := newALU()
	alu.SUB(5, 10)
	assert.True(t, alu.LessThanEqual())
	alu.SUB(5, 5)
	assert.True(t, alu.LessThanEqual())
	alu.SUB(10, 5)
	assert.False(t, alu.LessThanEqual())
	applyWithSigned(alu.SUB, -5, -10)
	assert.False(t, alu.LessThanEqual())
	applyWithSigned(alu.SUB, -10, -5)
	assert.True(t, alu.LessThanEqual())
	applyWithSigned(alu.SUB, -10, -10)
	assert.True(t, alu.LessThanEqual())
}

func applyWithSigned(f func(uint64, uint64) uint64, a int64, b int64) uint64 {
	return f(uint64(a), uint64(b))
}
