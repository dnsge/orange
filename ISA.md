# orange ISA

* 32-bit instruction size
* 32-bit word size
* 16 64-bit registers

`0000 0000 0000 0000 0000 0000 0000 0000`
`00oo oooo dddd aaaa bbbb 0000 0000 0000`
`00oo oooo dddd aaaa iiii iiii iiii iiii`

## Opcodes
A-Type (Arithmetic/Logical):
- `00oo oooo dddd aaaa bbbb 0000 0000 0000`
  - ex: `ADD r1, r2, r3 ;; r1 = r2 + r3`
- `00oo oooo dddd aaaa iiii iiii iiii iiii`
  - ex: `ADDI r1, r2, #10 ;; r1 = r2 + 10`

* `ADD` Add
* `ADDI` Add immediate
* `SUB` Subtract
* `SUBI` Subtract immediate
* `AND` Logical AND
* `OR` Logical OR
* `XOR` Logical XOR
* `LSL` Logical shift left
* `LSR` Logical shift right
* `CMP` Equiv to `SUB r0, r1, r2`
* `CMPI` Equiv to `SUB r0, r1, imm`

M-Type (Memory Operations):
- `00oo oooo aaaa bbbb ssss ssss ssss ssss`
  - ex: `LDREG r1, r2, #0`
  - ex: `STREG r1, r2, #1`

* `LDREG` Load 64-bit register
* `LDWORD` Load 32-bit word
* `LDHWRD` Load 16-bit half-word
* `LDBYTE` Load byte
* `STREG` Store 64-bit register
* `STWORD` Store lower 32-bit word
* `STHWRD` Store lower 16-bit half-word
* `STBYTE` Store lower byte

E-Type (Exchange):
- `00oo oooo dddd 0000 iiii iiii iiii iiii`
  - ex: `MOVZ r1, #100`

* `MOVZ` Move immediate into zeroed register
* `MOVK` Move immediate into register

B-Type (Branches):
- `00oo oooo 0000 0000 ssss ssss ssss ssss`'
  - Relative branches
- `00oo oooo aaaa 0000 0000 0000 0000 0000`
  - `BREG`

* `B` Branch by PC offset
* `BL` Branch with link in R15
* `BREG` Branch to register memory address
* `B.EQ`
* `B.NEQ`
* `B.LT`
* `B.LE`
* `B.GT`
* `B.GE`

O-Type (Operations):

* `HALT`
* `NOOP`

| test | test |
| ---- | ---- |
| asdf | asdf |