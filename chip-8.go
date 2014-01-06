package main

import (
	"github.com/DeedleFake/sdl"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

type Cpu struct {
	Mem           [4096]uint8
	V             [16]uint8
	Stack         [16]uint16
	Gfx           [2048]uint8
	Keys          [16]bool
	Draw          bool
	SP, DT, ST    uint8
	I, PC, Opcode uint16
}

const (
	TimerFreq = 240
	GfxScale = 10
	screenWidth = 64
	screenHeight = 32
)

var keyMap = map[int]byte{
	30: 0x1, // 1
	31: 0x2, // 2
	32: 0x3, // 3
	33: 0xc, // 4
	34: 0x4, // Q
	26: 0x5, // W
	8:  0x6, // E
	21: 0xd, // R
	4:  0x7, // A
	22: 0x8, // S
	7:  0x9, // D
	9:  0xe, // F
	29: 0xa, // Z
	27: 0x0, // X
	6:  0xb, // C
	25: 0xf, // V
}

func main() {
	c := Cpu{}
	c.init()
	c.loadRom("games/pong")

	err := sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()

	win, err := sdl.CreateWindow(
		"Test",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		screenWidth*GfxScale,
		screenHeight*GfxScale,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		panic(err)
	}
	defer win.Destroy()

	ren, err := win.CreateRenderer(-1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer ren.Destroy()

	ren.SetDrawColor(255, 255, 255, sdl.ALPHA_OPAQUE)

	for {
		c.keySet()
		c.emulateCycle()
		if c.Draw {
			c.updateScreen(ren)
		}
		c.updateTimers()
		time.Sleep(time.Second / TimerFreq)
	}
}

func (c *Cpu) init() {
	c.PC = 0x200

	font := [80]uint8{
		0xf0, 0x90, 0x90, 0x90, 0xf0, 0x20, 0x60, 0x20, 0x20, 0x70,
		0xf0, 0x10, 0xf0, 0x80, 0xf0, 0xf0, 0x10, 0xf0, 0x10, 0xf0,
		0x90, 0x90, 0xf0, 0x10, 0x10, 0xf0, 0x80, 0xf0, 0x10, 0xf0,
		0xf0, 0x80, 0xf0, 0x90, 0xf0, 0xf0, 0x10, 0x20, 0x40, 0x40,
		0xf0, 0x90, 0xf0, 0x90, 0xf0, 0xf0, 0x90, 0xf0, 0x10, 0xf0,
		0xf0, 0x90, 0xf0, 0x90, 0x90, 0xe0, 0x90, 0xe0, 0x90, 0xe0,
		0xf0, 0x80, 0x80, 0x80, 0xf0, 0xe0, 0x90, 0x90, 0x90, 0xe0,
		0xf0, 0x80, 0xf0, 0x80, 0xf0, 0xf0, 0x80, 0xf0, 0x80, 0x80,
	}

	for i, v := range font {
		c.Mem[i] = v
	}
}

func (c *Cpu) loadRom(rom string) {
	f, err := ioutil.ReadFile(rom)
	if err != nil {
		log.Fatal(err)
	}

	for i, content := range f {
		c.Mem[0x200+i] = content
	}
}

func (c *Cpu) emulateCycle() {
	c.Opcode = uint16(c.Mem[c.PC])<<8 | uint16(c.Mem[c.PC+1])

	switch c.Opcode & 0xf000 {
	case 0x0000:
		switch c.Opcode & 0x0fff {
		case 0x00e0:
			c.Gfx = [2048]uint8{}
			c.Draw = true
			c.PC += 2
		case 0x00ee:
			c.SP--
			c.PC = c.Stack[c.SP]
			c.PC += 2
		}
	case 0x1000:
		c.PC = c.Opcode & 0x0fff
	case 0x2000:
		c.Stack[c.SP] = c.PC
		c.SP++
		c.PC = c.Opcode & 0x0fff
	case 0x3000:
		if uint16(c.V[c.Opcode&0x0f00>>8]) == c.Opcode&0x00ff {
			c.PC += 4
		} else {
			c.PC += 2
		}
	case 0x4000:
		if uint16(c.V[c.Opcode&0x0f00>>8]) != c.Opcode&0x00ff {
			c.PC += 4
		} else {
			c.PC += 2
		}
	case 0x5000:
		if c.V[c.Opcode&0x0f00>>8] == c.V[c.Opcode&0x00f0>>4] {
			c.PC += 4
		} else {
			c.PC += 2
		}
	case 0x6000:
		c.V[c.Opcode&0x0f00>>8] = uint8(c.Opcode & 0x00ff)
		c.PC += 2
	case 0x7000:
		c.V[c.Opcode&0x0f00>>8] += uint8(c.Opcode & 0x00ff)
		c.PC += 2
	case 0x8000:
		switch c.Opcode & 0x000f {
		case 0x0000:
			c.V[c.Opcode&0x0f00>>8] = c.V[c.Opcode&0x00f0>>4]
			c.PC += 2
		case 0x0001:
			c.V[c.Opcode&0x0f00>>8] |= c.V[c.Opcode&0x00f0>>4]
			c.PC += 2
		case 0x0002:
			c.V[c.Opcode&0x0f00>>8] &= c.V[c.Opcode&0x00f0>>4]
			c.PC += 2
		case 0x0003:
			c.V[c.Opcode&0x0f00>>8] ^= c.V[c.Opcode&0x00f0>>4]
			c.PC += 2
		case 0x0004:
			if c.V[c.Opcode&0x00f0>>4] > (0xff-c.V[c.Opcode&0x0f00>>8]) {
				c.V[0xf] = 1
			} else {
				c.V[0xf] = 0
			}
			c.V[c.Opcode&0x0f00>>8] += c.V[c.Opcode&0x00f0>>4]
			c.PC += 2
		case 0x0005:
			if c.V[c.Opcode&0x00f0>>4] > c.V[c.Opcode&0x0f00>>8] {
				c.V[0xf] = 0
			} else {
				c.V[0xf] = 1
			}
			c.V[c.Opcode&0x0f00>>8] -= c.V[c.Opcode&0x00f0>>4]
			c.PC += 2
		case 0x0006:
			c.V[0xf] = c.V[c.Opcode&0x0f00>>8] & 0x1
			c.V[c.Opcode&0x0f00>>8] >>= 1
			c.PC += 2
		case 0x0007:
			if c.V[c.Opcode&0x00f0>>4] < c.V[c.Opcode&0x0f00>>8] {
				c.V[0xf] = 0
			} else {
				c.V[0xf] = 1
			}
			c.V[c.Opcode&0x0f00>>8] = c.V[c.Opcode&0x00f0>>4] - c.V[c.Opcode&0x0f00>>8]
			c.PC += 2
		case 0x000e:
			c.V[0xf] = c.V[c.Opcode&0x0f00>>8] >> 7 //& 0x80
			c.V[c.Opcode&0x0f00>>8] <<= 1
			c.PC += 2
		}
	case 0x9000:
		if c.V[c.Opcode&0x0f00>>8] != c.V[c.Opcode&0x00f0>>4] {
			c.PC += 4
		} else {
			c.PC += 2
		}
	case 0xa000:
		c.I = c.Opcode & 0x0fff
		c.PC += 2
	case 0xb000:
		c.PC = c.Opcode&0x0fff + uint16(c.V[0])
	case 0xc000:
		c.V[c.Opcode&0x0f00>>8] = uint8(uint16(rand.Int() % 0xff) & (c.Opcode & 0x00ff))
		c.PC += 2
	case 0xd000:
		height := uint8(c.Opcode & 0x000f)
		x := c.V[c.Opcode&0x0f00>>8]
		y := c.V[c.Opcode&0x00f0>>4]
		c.V[0xf] = 0

		for yline := uint16(0); yline < uint16(height); yline++ {
			pixel := c.Mem[c.I+uint16(yline)]
			for xline := uint16(0); xline < 8; xline++ {
				if (pixel & (0x80 >> xline)) != 0 {
					if c.Gfx[uint16(x)+xline+((uint16(y)+yline)*screenWidth)] == 1 {
						c.V[0xf] = 1
					}
					c.Gfx[uint16(x)+xline+((uint16(y)+yline)*screenWidth)] ^= 1
				}
			}
		}

		c.Draw = true
		c.PC += 2
	case 0xe000:
		switch c.Opcode & 0x00ff {
		case 0x009e:
			if c.Keys[c.V[c.Opcode&0x0f00>>8]] {
				c.PC += 4
			} else {
				c.PC += 2
			}
		case 0x00a1:
			if c.Keys[c.Opcode&0x0f00>>8] {
				c.PC += 2
			} else {
				c.PC += 4
			}
		}
	case 0xf000:
		switch c.Opcode & 0x00ff {
		case 0x0007:
			c.V[c.Opcode&0x0f00>>8] = c.DT
			c.PC += 2
		case 0x000a:
			panic("Not implemented")
		case 0x0015:
			c.DT = c.V[c.Opcode&0x0f00>>8]
			c.PC += 2
		case 0x0018:
			c.ST = c.V[c.Opcode&0x0f00>>8]
			c.PC += 2
		case 0x001e:
			if (c.I + uint16(c.V[c.Opcode&0x0f00>>4])) > 0xfff {
				c.V[0xf] = 1
			} else {
				c.V[0xf] = 0
			}
			c.I += uint16(c.V[c.Opcode&0x0f00>>8])
			c.PC += 2
		case 0x0029:
			c.I = 0x05 * (c.Opcode & 0x0f00 >> 8)
			c.PC += 2
		case 0x0033:
			c.Mem[c.I] = c.V[c.Opcode&0x0f00>>8] / 100
			c.Mem[c.I+1] = (c.V[c.Opcode&0x0f00>>8] / 10) % 10
			c.Mem[c.I+2] = (c.V[c.Opcode&0x0f00>>8] % 100) % 10
			c.PC += 2
		case 0x0055:
			for i := uint16(0); i <= (c.Opcode&0x0f00 >> 8); i++ {
				c.Mem[i + i] = c.V[i]
			}
			c.I += (c.Opcode & 0x0f00 >> 8) + 1
			c.PC += 2
		case 0x0065:
			for i := uint16(0); i <= (c.Opcode&0x0f00 >> 8); i++ {
				c.V[i] = c.Mem[c.I + i]
			}
			c.I += (c.Opcode&0x0f00 >> 8) + 1
			c.PC += 2
		}
	}
}

func (c *Cpu) updateScreen(ren *sdl.Renderer) {
	for k := 0; k < screenHeight; k++ {
		for l := 0; l < screenWidth; l++ {
			if c.Gfx[(l+(k*screenWidth))] == 1 {
				ren.FillRect(&sdl.Rect{int32(l) * GfxScale, int32(k) * GfxScale, GfxScale, GfxScale})
			}
		}
	}
	ren.Present()
	c.Draw = false
}

func (c *Cpu) updateTimers() {
	if c.ST > 0 {
		c.ST--
	}

	if c.ST == 1 {
		log.Printf("BEEP!")
	}

	if c.DT > 0 {
		c.DT--
	}
}

func (c *Cpu) keySet() {
	var ev sdl.Event
	for sdl.PollEvent(&ev) {
		switch ev := ev.(type) {
		case *sdl.KeyboardEvent:
			if ev.Type == sdl.KEYDOWN {
				c.Keys[keyMap[int(ev.Keysym.Scancode)]] = true
			} else if ev.Type == sdl.KEYUP {
				c.Keys[keyMap[int(ev.Keysym.Scancode)]] = false
			}
		}
	}
}
