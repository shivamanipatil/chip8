package main

import (
	"os"
	"strconv"

	sdl "github.com/veandco/go-sdl2/sdl"
	vm "github.com/shivamanipatil/chip8/vm"
)

//Width window  
const Width int32 = 64
//Height window
const Height int32 = 32

func main() {
	if len(os.Args) < 3 {
		panic("Please provide modifier and a c8 file")
	}

	fileName := os.Args[2]
	var modifier int32 = 10

	if len(os.Args) == 3 {
		if val, valErr := strconv.ParseInt(os.Args[1], 10, 32); valErr != nil {
			panic(valErr)
		} else {
			if val > 0 {
				modifier = int32(val)
			}
		}
	}

	// Get vm instance
	c8 := vm.New()
	if loadErr := c8.LoadProgram(fileName); loadErr != nil {
		panic(loadErr)
	}

	// Initialize sdl2
	if sdlErr := sdl.Init(sdl.INIT_EVERYTHING); sdlErr != nil {
		panic(sdlErr)
	}
	defer sdl.Quit()


	// Create window, chip8 resolution with given modifier
	window, windowErr := sdl.CreateWindow("Chip 8 - "+fileName, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, Width*modifier, Height*modifier, sdl.WINDOW_SHOWN)
	if windowErr != nil {
		panic(windowErr)
	}
	defer window.Destroy()

	// Create render surface
	canvas, canvasErr := sdl.CreateRenderer(window, -1, 0)
	if canvasErr != nil {
		panic(canvasErr)
	}
	defer canvas.Destroy()

	for {
		// Emulate one cycle
		c8.EmulateCycle()
		// Draw only if draw flag is set
		if c8.Draw() {
			// Clear the screen
			canvas.SetDrawColor(255, 0, 0, 255)
			canvas.Clear()

			// Get the display buffer and render
			vector := c8.Gfx
			for j := 0; j < len(vector); j++ {
				for i := 0; i < len(vector[j]); i++ {
					// Values of pixel are stored in 1D array of size 64 * 32
					if vector[j][i] != 0 {
						canvas.SetDrawColor(255, 255, 255, 255)
					} else {
						canvas.SetDrawColor(0, 0, 0, 255)
					}
					canvas.FillRect(&sdl.Rect{
						Y: int32(j) * modifier,
						X: int32(i) * modifier,
						W: modifier,
						H: modifier,
					})
				}
			}

			canvas.Present()
		}

		// Poll for Quit and Keyboard events
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch et := event.(type) {
			case *sdl.QuitEvent:
				os.Exit(0)
			case *sdl.KeyboardEvent:
				if et.Type == sdl.KEYUP {
					switch et.Keysym.Sym {
					case sdl.K_1:
						c8.Keys[0x1] = false
					case sdl.K_2:
						c8.Keys[0x2] = false
					case sdl.K_3:
						c8.Keys[0x3] = false
					case sdl.K_4:
						c8.Keys[0xC] = false
					case sdl.K_q:
						c8.Keys[0x4] = false
					case sdl.K_w:
						c8.Keys[0x5] = false
					case sdl.K_e:
						c8.Keys[0x6] = false
					case sdl.K_r:
						c8.Keys[0xD] = false
					case sdl.K_a:
						c8.Keys[0x7] = false
					case sdl.K_s:
						c8.Keys[0x8] = false
					case sdl.K_d:
						c8.Keys[0x9] = false
					case sdl.K_f:
						c8.Keys[0xE] = false
					case sdl.K_z:
						c8.Keys[0xA] = false
					case sdl.K_x:
						c8.Keys[0x0] = false
					case sdl.K_c:
						c8.Keys[0xB] = false
					case sdl.K_v:
						c8.Keys[0xF] = false
					}
				} else if et.Type == sdl.KEYDOWN {
					switch et.Keysym.Sym {
					case sdl.K_1:
						c8.Keys[0x1] = true
					case sdl.K_2:
						c8.Keys[0x2] = true
					case sdl.K_3:
						c8.Keys[0x3] = true
					case sdl.K_4:
						c8.Keys[0xC] = true
					case sdl.K_q:
						c8.Keys[0x4] = true
					case sdl.K_w:
						c8.Keys[0x5] = true
					case sdl.K_e:
						c8.Keys[0x6] = true
					case sdl.K_r:
						c8.Keys[0xD] = true
					case sdl.K_a:
						c8.Keys[0x7] = true
					case sdl.K_s:
						c8.Keys[0x8] = true
					case sdl.K_d:
						c8.Keys[0x9] = true
					case sdl.K_f:
						c8.Keys[0xE] = true
					case sdl.K_z:
						c8.Keys[0xA] = true
					case sdl.K_x:
						c8.Keys[0x0] = true
					case sdl.K_c:
						c8.Keys[0xB] = true
					case sdl.K_v:
						c8.Keys[0xF] = true
					}
				}
			}
		}

		// Chip8 cpu clock worked at frequency of 60Hz, so set delay to (1000/60)ms
		sdl.Delay(1000 / 60)
	}
}