# orange

A virtual machine with a custom ISA. Named after my favorite color because I can.

## About

This repository features an assembler and a virtual machine for assembling and then running programs written in assembly. The instruction set architecture is described in [ISA.md](ISA.md).

## Why?

While taking an introductory computer architecture class, I became interested in exploring more advanced assemblers and simulators than the ones that we were learning about and creating in class.

The ISA for orange is inspired by "LEGv8", a subset of ARMv8 (hence the name), which we discuss in EECS 370. I plan to expand on the ISA to add more interesting features like built-in stack management. 

## Usage

To assemble a program, run the package located in `./cmd/orangeasm`. To run a program, run the package located in `./cmd/orangevm`. To link multiple object files, run the package located in `./cmd/orangelinker`.

If you want to assemble a standalone program (e.g. no linking), use `./orangeasm --executable [input file] [output file]`.

## Examples

- [multiplication.orange](./programs/multiplication.orange)
  - Multiplies two 32-bit numbers using grade-school multiplication algorithm
- [string.orange](./programs/string.orange)
  - Traversal of a null-terminated string
- [stack.orange](./programs/stack.orange)
  - Simple demonstration of calling functions and saving values to the stack
- [print.orange](./programs/print.orange)
  - Example usage of write syscall for printing a string
- [sections.orange](./programs/sections.orange)
  - Semantics for defining different sections
- [linker example](./programs/link)
  - Showcases using multiple files to implement a program
  - The main file calls the strLen function in the second file

## Todo

- [x] Better, more modular parsing of assembly
- [x] Stack management
- [x] Proper error management
- [x] Object files with symbol table, relocation table
- [ ] System calls (for console output)
- [ ] Dynamic memory allocation via syscalls
- [ ] Simple language + compiler
