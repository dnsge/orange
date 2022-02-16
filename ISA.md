# orange ISA

* 32-bit instruction size
* 32-bit word size
* 16 64-bit registers

Register allocation:
* `r0`: Zero register (returns 0, ignores writes)
* `r15`: Return address (`BL` instruction)
* `r14`: Stack pointer
* `r1-8`: Syscall arguments
* `r9`: Syscall number
* Caller saved: `r1-r6`
* Callee saved: `r7-12`

## Opcodes
| Opcode   | Type    |
|----------|---------|
| `ADD`    | A-Type  |
| `SUB`    | A-Type  |
| `AND`    | A-Type  |
| `OR`     | A-Type  |
| `XOR`    | A-Type  |
| `CMP`    | A-Type  |
| `ADDI`   | AI-Type |
| `SUBI`   | AI-Type |
| `LSL`    | AI-Type |
| `LSR`    | AI-Type |
| `CMPI`   | AI-Type |
| `LDREG`  | M-Type  |
| `LDWORD` | M-Type  |
| `LDHWRD` | M-Type  |
| `LDBYTE` | M-Type  |
| `STREG`  | M-Type  |
| `STWORD` | M-Type  |
| `STHWRD` | M-Type  |
| `STBYTE` | M-Type  |
| `ADR`    | M-Type  |
| `MOVZ`   | E-Type  |
| `MOVK`   | E-Type  |
| `BREG`   | B-Type  |
| `BL`     | B-Type  |
| `B`      | BI-Type |
| `B.EQ`   | BI-Type |
| `B.NEQ`  | BI-Type |
| `B.LT`   | BI-Type |
| `B.LE`   | BI-Type |
| `B.GT`   | BI-Type |
| `B.GE`   | BI-Type |
| `PUSH`   | R-Type  |
| `POP`    | R-Type  |
| `NOOP`   | O-Type  |
| `HALT`   | O-Type  |

## Instruction Formats
| Type    | Layout                                   | Description                            |
|---------|------------------------------------------|:---------------------------------------|
| A-Type  | `00oooooo dddd aaaa bbbb 0000 0000 0000` | Operation with 3 registers             |
| AI-Type | `00oooooo dddd aaaa iiii iiii iiii iiii` | Operation with 2 registers + immediate |
| M-Type  | `00oooooo aaaa bbbb ssss ssss ssss ssss` | Memory operation                       |
| E-Type  | `00oooooo dddd 0000 iiii iiii iiii iiii` | Exchange with immediate operation      |
| B-Type  | `00oooooo aaaa 0000 0000 0000 0000 0000` | Branch operation                       |
| BI-Type | `00oooooo 0000 0000 ssss ssss ssss ssss` | Branch with immediate operation        |
| R-Type  | `00oooooo aaaa 0000 0000 0000 0000 0000` | Operation with register                |
| O-Type  | `00oooooo 0000 0000 0000 0000 0000 0000` | Generic, no argument operation         |
