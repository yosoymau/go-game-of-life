package main

import (
	"errors"
	"fmt"
	"math/rand"

	//"math/rand"

	"github.com/veandco/go-sdl2/sdl"
)

// Constants - Configuration
const SQR_SIZE = int32(10)
const WINDOW_WIDTH = 800
const WINDOW_HEIGHT = 600

func main() {
	// TODO:
	// * Add text to show stats
	// * Allow to draw while simulating maybe?
	// *

	// SDL Setup
	fmt.Println("starting")
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	// Window creation
	fmt.Println("creating window")
	window, err := sdl.CreateWindow("test", 0,0, WINDOW_WIDTH, WINDOW_HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	// Create renderer
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	// Setup game of life.
	// Slices to contain data
	cellsgrid := make([][]bool, WINDOW_WIDTH/SQR_SIZE)

	for i:= range cellsgrid {
		cellsgrid[i] = make([]bool, WINDOW_HEIGHT/SQR_SIZE)
	}

	for i:= range cellsgrid {
		for j := range cellsgrid[i] {
			cellsgrid[i][j] = false
		}
	}
	for i:= range cellsgrid {
		for j := range cellsgrid[i] {
			cellsgrid[i][j] = (rand.Intn(2) == 0)
		}
	}
	//cellsgrid[40][30] = true

	// Vestigial testing
	//cellsgrid := make([][]bool, 3)
	//cellsgrid[0] = []bool{true, false, true}
	//cellsgrid[1] = []bool{false, true, false}
	//cellsgrid[2] = []bool{true, false, true}




	// Main loop
	running := true
	simulating := false
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t:= event.(type) {
			case *sdl.QuitEvent:
				println("quit")
				running = false
				break

			case *sdl.MouseButtonEvent:
				if t.State == sdl.PRESSED{
					xpos := int(t.X / 10)
					ypos := int(t.Y / 10)
					cellsgrid[xpos][ypos] = !cellsgrid[xpos][ypos]
				}
			case *sdl.KeyboardEvent:
				if t.State == sdl.PRESSED{
					switch t.Keysym.Sym {
					case sdl.K_a:
						simulating = !simulating
					case sdl.K_s:
						for i:= range cellsgrid {
							for j := range cellsgrid[i] {
								cellsgrid[i][j] = (rand.Intn(2) == 0)
							}
						}
					case sdl.K_c:
						for i:= range cellsgrid {
							for j := range cellsgrid[i] {
								cellsgrid[i][j] = false
							}
						}
					}

				}
			}
		}
		// Calculate things
		if simulating {
			cellsgrid = calculations(cellsgrid)
		}

		// clear screen
		renderer.SetDrawColor(100, 0, 100, 255)
		renderer.Clear()

		// Create rects
		renderer.SetDrawColor(250, 250, 250, 250)
		rects, err := render_squares(cellsgrid)
		if err == nil {
			renderer.FillRects(rects)
			renderer.DrawRects(rects)
		}

		// Draw
		renderer.Present()
		//sdl.Delay(132)
		sdl.Delay(33)
	}
}

func render_squares(data [][]bool) ([]sdl.Rect, error) {
	var rects []sdl.Rect
	for i:= range data {
		for j := range data[i] {
			if data[i][j] {
				newRect := sdl.Rect{int32(i)*SQR_SIZE, int32(j)*SQR_SIZE, SQR_SIZE, SQR_SIZE}
				rects = append(rects, newRect)
			}
		}
	}
	if len(rects) == 0 {
		return nil, errors.New("No rectangles")
	} else {
		return rects, nil
	}
}

func calculations(data [][]bool) [][]bool {
	new := make([][]bool, WINDOW_WIDTH/SQR_SIZE)

	for i:= range new {
		new[i] = make([]bool, WINDOW_HEIGHT/SQR_SIZE)
	}
	for i:= range data {
		for j := range data[i] {
			sum := 0
			for x:= -1; x<2 ; x++ {
				for y := -1; y<2 ; y++ {
					if i+x < 0 || i+x > len(data)-1 || x == 0 && y == 0 || j+y < 0 || j+y > len(data[i])-1  {
						continue
					}
					if data[i+x][j+y] {
						sum += 1
					}
				}
			}

			if sum == 3  {
				new[i][j] = true
			}
			if sum == 2 && data[i][j]{
				new[i][j] = true
			}
			if sum < 2 {
				new[i][j] = false
			}
			if sum > 3 {
				new[i][j] = false
			}

		}
	}
	return new
}
