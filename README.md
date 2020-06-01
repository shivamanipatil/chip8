# Chip 8 

CHIP - 8 emulator written in golang

## Installing and running

Download sdl drivers for your operating system.
- clone this repo

```
git clone https://github.com/shivamanipatil/chip8.git
```
- Running

```
go run main.go <modifier> <rom path>
```
Modifier is used to set logical size to pixel.Default chip8 resolution was 64x32 so e.g modifier of 10 will make window size 640x320 

- example running

```
go run main.go 10 ~/roms/pong.c8
```

## Bindings

```
Chip8 keypad         Keyboard mapping
1 | 2 | 3 | C        1 | 2 | 3 | 4
4 | 5 | 6 | D   =>   Q | W | E | R
7 | 8 | 9 | E   =>   A | S | D | F
A | 0 | B | F        Z | X | C | V
```
## Todo

- [ ] Compile to WebAssembly or js 

## Sources

- [skatiyar](https://github.com/skatiyar) [repo](https://github.com/skatiyar/go-chip8) for sdl code inspiration
- [How to write an emulator chip-8 interpreter](http://www.multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/)
- [Cowgod's Chip-8 Technical Reference](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM)
- [Chip-8 opcode table](https://en.wikipedia.org/wiki/CHIP-8)

