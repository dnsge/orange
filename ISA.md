# orange ISA

* 32-bit word size
* 32-bit (1 word) instruction size
* 16 64-bit registers

Register allocation:
* `r0`: Zero register (returns 0, ignores writes)
* `r15`: Return address (`BL` instruction)
* `r14`: Stack pointer
* `r1-6`: Syscall arguments
* `r7`: Syscall return value
* `r8`: Syscall error (or zero)
* `r9`: Syscall number

ABI:
* Caller saved: `r1-r9`
* Callee saved: `r10-13`
* Function arguments: `r1-r4`
* Return value(s): saved to stack, must be removed

## Opcodes
| Opcode   | Type    | Description                                    |
|----------|---------|------------------------------------------------|
| `ADD`    | A-Type  | Add                                            |
| `SUB`    | A-Type  | Subtract                                       |
| `AND`    | A-Type  | Bitwise AND                                    |
| `OR`     | A-Type  | Bitwise OR                                     |
| `XOR`    | A-Type  | Bitwise XOR                                    |
| `CMP`    | A-Type  | Compare                                        |
| `ADDI`   | AI-Type | Add immediate                                  |
| `SUBI`   | AI-Type | Subtract immediate                             |
| `LSL`    | AI-Type | Left shift by immediate                        |
| `LSR`    | AI-Type | Right shift by immediate                       |
| `CMPI`   | AI-Type | Compare with immediate                         |
| `LDREG`  | M-Type  | Load 8 byte register                           |
| `LDWORD` | M-Type  | Load 4 byte word                               |
| `LDHWRD` | M-Type  | Load 2 byte half-word                          |
| `LDBYTE` | M-Type  | Load 1 byte                                    |
| `STREG`  | M-Type  | Store 8 byte register                          |
| `STWORD` | M-Type  | Store lower 4 byte word                        |
| `STHWRD` | M-Type  | Store lower 2 byte half-word                   |
| `STBYTE` | M-Type  | Store lowest byte                              |
| `ADR`    | M-Type  | Pseudo-instruction to load label address       |
| `MOVZ`   | E-Type  | Zero register and move immediate               |
| `MOVK`   | E-Type  | Move immediate into lower 16 bits              |
| `BREG`   | B-Type  | Branch to address in register                  |
| `BL`     | B-Type  | Branch to address in register with link in r15 |
| `B`      | BI-Type | Branch to relative offset/label                |
| `B.EQ`   | BI-Type | Branch to relative offset/label if ==          |
| `B.NEQ`  | BI-Type | Branch to relative offset/label if !=          |
| `B.LT`   | BI-Type | Branch to relative offset/label if <           |
| `B.LE`   | BI-Type | Branch to relative offset/label if <=          |
| `B.GT`   | BI-Type | Branch to relative offset/label if >           |
| `B.GE`   | BI-Type | Branch to relative offset/label if >=          |
| `PUSH`   | R-Type  | Push register to stack                         |
| `POP`    | R-Type  | Pop register from stack                        |
| `NOOP`   | O-Type  | No-op                                          |
| `HALT`   | O-Type  | Halt VM                                        |

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
