package vm

import (
	"fmt"
	"math/rand"
	"os"
)

//VM implemets chip8 vm
type VM struct {
	Memory     [4096]uint8
	Gfx        [32][64]uint8
	Stack      [16]uint16
	V          [16]uint8
	I, PC, SP  uint16
	DelayTimer uint8
	SoundTimer uint8
	OpCode     uint16
	DrawFlag   bool
	Keys       [16]bool
}

//Fonts bytes
var Fonts = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

//New returns new VM
func New() VM {
	vm := VM{}
	vm.Initialize()
	return vm
}

//Initialize the vm
func (vm *VM) Initialize() {
	vm.PC = 0x200
	vm.DrawFlag = true

	//load the fontsets
	for i := 0; i < len(Fonts); i++ {
		vm.Memory[i] = Fonts[i] //loadfontset
	}

}

//Draw called to draw
func (vm *VM) Draw() bool {
	sd := vm.DrawFlag
	vm.DrawFlag = false
	return sd
}

//EmulateCycle fetch->decode->execute->update-timers
func (vm *VM) EmulateCycle() {
	//fetch opcode
	vm.OpCode = uint16(vm.Memory[vm.PC])<<8 | uint16(vm.Memory[vm.PC+1])
	x := (vm.OpCode & 0x0F00) >> 8
	y := (vm.OpCode & 0x00F0) >> 4
	nn := vm.OpCode & 0x00FF
	nnn := vm.OpCode & 0x0FFF
	vm.PC += 2

	//execute
	switch vm.OpCode & 0xF000 {
	case 0x0000:
		switch vm.OpCode & 0x000F {
		case 0x0000: // clear screen
			for i := 0; i < len(vm.Gfx); i++ {
				for j := 0; j < len(vm.Gfx[i]); j++ {
					vm.Gfx[i][j] = 0x0
				}
			}
			vm.DrawFlag = true
		case 0x000E: // return from subroutine
			vm.SP = vm.SP - 1
			vm.PC = vm.Stack[vm.SP]
		default:
			fmt.Println("Unknown opcode")
		}
	case 0x1000: // Jumps to NNN
		vm.PC = nnn
	case 0x2000: // Call subroutine at NNN
		vm.Stack[vm.SP] = vm.PC
		vm.SP++
		vm.PC = 0x0FFF & vm.OpCode
	case 0x3000: // vx equals NN at memory
		if uint16(vm.V[x]) == nn {
			vm.PC += 2
		}
	case 0x4000: // vx not equals NN at memory
		if uint16(vm.V[x]) != nn {
			vm.PC += 2
		}
	case 0x5000: // vx equals vy
		if vm.V[x] == vm.V[y] {
			vm.PC += 2
		}
	case 0x6000: // vx to NN
		vm.V[x] = uint8(nn)
	case 0x7000: // vx += NN
		vm.V[x] += uint8(nn)
	case 0x8000:
		switch vm.OpCode & 0x000F {
		case 0x0000: // vx = vy
			vm.V[x] = vm.V[y]
		case 0x0001: // vx |= vy
			vm.V[x] |= vm.V[y]
		case 0x0002: // vx &= vy
			vm.V[x] &= vm.V[y]
		case 0x0003: // vx ^= vy
			vm.V[x] ^= vm.V[y]
		case 0x0004: // vx += vy
			vm.V[0xF] = 0
			//check if addition > 255 and set flag reg = 1 if true
			if vm.V[x] > (0xFF - vm.V[y]) {
				vm.V[0xF] = 1
			}
			vm.V[x] += vm.V[y]
		case 0x0005: // vx -= vy
			vm.V[0xF] = 0
			//flag = 1 if not borrow i.e vx - vy
			if vm.V[x] > vm.V[y] {
				vm.V[0xF] = 1
			}
			vm.V[x] -= vm.V[y]
		case 0x0006: // vx >> 1
			vm.V[0xF] = vm.V[x] & 0x1
			vm.V[x] >>= 1
		case 0x0007: // vx = vy - vx
			vm.V[0xF] = 0
			//flag = 1 if not borrow i.e vx - vy
			if vm.V[y] > vm.V[x] {
				vm.V[0xF] = 1
			}
			vm.V[x] = vm.V[y] - vm.V[x]
		case 0x000E: // vx << 1
			vm.V[0xF] = vm.V[x] >> 7
			vm.V[x] <<= 1

		default:
			fmt.Println("Unknown opcode")
		}
	case 0x9000: // vx != vy conditional
		if vm.V[x] != vm.V[y] {
			vm.PC += 2
		}
	case 0xA000: // load
		vm.I = nnn
	case 0xB000: // Jump
		vm.PC = uint16(vm.V[0x0]) + nnn
	case 0xC000: // Rand & nn
		vm.V[x] = uint8(nn) & (uint8(rand.Intn(256)))
	case 0xD000: //Draw
		vm.DrawFlag = true
		h := vm.OpCode & 0x000F
		vm.V[0xF] = 0
		var j uint16 = 0
		var i uint16 = 0
		for j = 0; j < h; j++ {
			pixel := vm.Memory[vm.I+j]
			for i = 0; i < 8; i++ {
				if (pixel & (0x80 >> i)) != 0 {
					if vm.Gfx[(vm.V[y] + uint8(j))][vm.V[x]+uint8(i)] == 1 {
						vm.V[0xF] = 1
					}
					vm.Gfx[(vm.V[y] + uint8(j))][vm.V[x]+uint8(i)] ^= 1
				}
			}
		}
		vm.DrawFlag = true
	case 0xE000:
		switch 0x00FF & vm.OpCode {
		case 0x009E:
			if vm.Keys[vm.V[x]] {
				vm.PC += 2
			}
		case 0x00A1:
			if !vm.Keys[vm.V[x]] {
				vm.PC += 2
			}
		default:
			fmt.Println("Unknown opcode")
		}
	case 0xF000:
		switch 0x00FF & vm.OpCode {
		case 0x0007:
			vm.V[x] = vm.DelayTimer
		case 0x000A:
			//wait until key is pressed
			vm.PC -= 2
			for i := uint8(0); i < 16; i++ {
				if vm.Keys[i] {
					vm.V[x] = i
					vm.PC += 2
					break
				}
			}
		case 0x0015:
			vm.DelayTimer = vm.V[x]
		case 0x0018:
			vm.SoundTimer = vm.V[x]
		case 0x001E:
			vm.V[0xF] = 0
			if vm.I > 0xFFF-uint16(vm.V[x]) {
				vm.V[0xF] = 1
			}
			vm.I += uint16(vm.V[x])
		case 0x0029:
			vm.I += uint16(vm.V[x]) * 5
		case 0x0033:
			vm.Memory[vm.I] = vm.V[x] / 100
			vm.Memory[vm.I+1] = (vm.V[x] / 10) % 10
			vm.Memory[vm.I+2] = vm.V[x] % 10
		case 0x0055:
			for i := uint16(0); i < x; i++ {
				vm.Memory[vm.I+i] = vm.V[i]
			}
			vm.I = ((vm.OpCode & 0x0F00) >> 8) + 1
		case 0x0065:
			for i := uint16(0); i < x; i++ {
				vm.V[i] = vm.Memory[vm.I+i]
			}
			vm.I = ((vm.OpCode & 0x0F00) >> 8) + 1
		}
	default:
		fmt.Println("Unknown opcode")
	}

	if vm.DelayTimer > 0 {
		vm.DelayTimer--
	}
	if vm.SoundTimer > 0 {
		if vm.SoundTimer == 1 {
			fmt.Println("BEEP")
		}
		vm.SoundTimer--
	}

}

//LoadProgram in the memory
func (vm *VM) LoadProgram(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	fStat, err := file.Stat()
	if err != nil {
		return err
	}
	if int64(len(vm.Memory)-512) < fStat.Size() { // program is loaded at 0x200
		return fmt.Errorf("Program size bigger than memory")
	}
	buffer := make([]byte, fStat.Size())
	if _, err = file.Read(buffer); err != nil {
		return err
	}

	for i := 0; i < len(buffer); i++ {
		vm.Memory[i+512] = buffer[i]
	}
	fmt.Println("here ")
	return nil
}
