package goemu

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type emuContext struct {
	paused  bool
	running bool
	ticks   u64
}

// Global emu context
var ctxEmu emuContext

// Get emu context
func emuGetContext() *emuContext {
	return &ctxEmu
}

func delay(ms u32) {
	sdl.Delay(uint32(ms))
}

func emuRun(args []string) int {
	if len(args) < 2 {
		fmt.Println("Usage: emu <rom_file>")
		return -1
	}

	if !cartLoad(args[1]) {
		fmt.Printf("Failed to load ROM file: %s\n", args[1])
		return -2
	}

	fmt.Println("Cart loaded..")

	if err := sdl.Init(0x00000001); err != nil {
		fmt.Printf("SDL INIT failed: %s\n", err)
		return -3
	}
	defer sdl.Quit()
	fmt.Println("SDL INIT")

	if err := ttf.Init(); err != nil {
		fmt.Printf("TTF INIT failed: %s\n", err)
		return -4
	}
	defer ttf.Quit()
	fmt.Println("TTF INIT")

	cpuInit()

	ctxEmu.running = true
	ctxEmu.paused = false
	ctxEmu.ticks = 0

	for ctxEmu.running {
		if ctxEmu.paused {
			delay(10)
			continue
		}

		if !cpuStep() {
			fmt.Println("CPU Stopped")
			return -3
		}

		ctxEmu.ticks++
	}

	return 0
}
