package main

import "fmt"

func main() {
	// ARCHITECTURE

	// memory, 4K in total
	mem := [4096]byte{}

	// 15 8-bit general purpose registers V0-VE
	// and VF which is the carry flag
	registers := [16]byte{}

	// address register, 16-bits, used by opcodes that involve memory ops
	var i int16

	// program counter (pseudo-register), 16-bit
	var pc int16

	// stack pointer (pseudo-register), 8-bit
	var sp byte

	// special purpose registers, 8-bit
	var delay_timer, sound_timer byte

	// stack, array of 16 16-bit values
	stack := [16]int16{}

	// current state of the screen, 64x32 pixels
	gfx := [2048]bool{}

	// current state of the keyboard, 16 keys (0-F)
	keyboard := [16]bool{}

	// current opcode, all opcodes are 2-byte long
	var opcode int16


	// INIT
	// initialize the pc
	pc = 0x200

	// fetch the opcode
	opcode = int(mem[pc]) << 8 | int(mem[pc+1])





	fmt.Println(opcode[0])
}

func emulateCycle(*mem) {
	// fetch opcode
	// decode opcode
	// exec opcode
	// update timers
}
