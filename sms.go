package main

import (
	// Standard
	"errors"
	"fmt"
	"image"
	_ "image/png" // _ = not directly referenced in code, but is used indirectly
	"io/ioutil"
	"log"
	"math"
	"os"
	"time"

	// Semi-Standard
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"

	// Foreign
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"github.com/gopxl/pixel/text"
)

var (
	WINDOW *pixelgl.Window // Not a constant, but its close enough (constant in action no definition)

	win_size        = pixel.Vec{X: 960, Y: 640}
	win_origin      = pixel.Vec{X: 0, Y: 0}
	win_vsync  bool = true

	text_atlas      *text.Atlas
	text_debug      *text.Text
	text_debug_size float64 = 16

	current_time time.Time = time.Now()
	runtime      float64   = 0.0

	//cam_pos pixel.Vec = pixel.ZV

	glep       Entity
	glep_speed float64 = 250.0
	dvd        Entity
	title      Entity

	smileyball Item
)

func loadPicture(path string) pixel.Picture {
	// Check if file path is valid
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("Image path `%s` does not exist.\n Error: %v", path, err)
	}

	// ========== Standard golang get image from file ==========
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening image file `%s`.\n Error: %v", path, err)
	}
	defer file.Close() // closes when function ends
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Error decoding image file `%s`.\n Error: %v", path, err)
	}

	// === Convert image to Pixel libraries "Picture" class ===
	return pixel.PictureDataFromImage(img)
}

func refreshDeltaTime() float64 {
	dt := time.Since(current_time).Seconds()
	current_time = time.Now()
	return dt
}

func createWindow() {
	cfg := pixelgl.WindowConfig{
		Title:  "Smilemon Sapphire",
		Bounds: pixel.R(win_origin.X, win_origin.Y, win_size.X, win_size.Y),
		VSync:  win_vsync,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	WINDOW = win
	WINDOW.SetSmooth(false) // pixel art

	fontface, err := loadTruetypeFont("eurostile-round.ttf", text_debug_size)
	if err != nil {
		panic(err)
	}

	text_atlas = text.NewAtlas(fontface, text.ASCII)
	text_debug = text.New(pixel.V(text_debug_size, win_size.Y-text_debug_size*2), text_atlas)
	text_debug.Color = colornames.Black
	text_debug.LineHeight = text_atlas.LineHeight() * 1.25

	win.Clear(colornames.Skyblue)
}

func createSprite(image_path string) *pixel.Sprite {
	img := loadPicture(image_path)
	return pixel.NewSprite(img, img.Bounds())
}

func input() {
	vel := pixel.Vec{X: 0.0, Y: 0.0}

	if WINDOW.Pressed(pixelgl.KeyW) {
		vel.Y = 1
	}
	if WINDOW.Pressed(pixelgl.KeyA) {
		vel.X = -1
	}
	if WINDOW.Pressed(pixelgl.KeyS) {
		vel.Y = -1
	}
	if WINDOW.Pressed(pixelgl.KeyD) {
		vel.X = 1
	}

	glep.vel = vel
	glep.normalizeVelocity()
}

func loadTruetypeFont(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}

func gameLoop() {
	createWindow()

	glep = createEntity(
		"images/glep.png",
		WINDOW.Bounds().Center(),
		pixel.Vec{X: 0.0, Y: 0.0},
		pixel.Vec{X: 1.5, Y: 1.5},
		0.0,
	)
	dvd = createEntity(
		"images/dvd.png",
		pixel.Vec{X: win_size.X - 96, Y: 96},
		pixel.Vec{X: 0.0, Y: 0.0},
		pixel.Vec{X: 2.0, Y: 2.0},
		0.0,
	)
	title = createEntity(
		"images/titlepage.png",
		WINDOW.Bounds().Center(),
		pixel.Vec{X: 0.0, Y: 0.0},
		pixel.Vec{X: 1.2, Y: 1.2},
		0.0,
	)

	smileyball = createItem(
		"images/smileyball.png",
		pixel.Vec{X: 500.0, Y: 300.0},
		pixel.Vec{X: 0.7, Y: 0.7},
		0.0,
		"A ball that catches critters. I don't know how. Please don't ask. Please. PLEASE.",
		true,
		true,
	)

	for !WINDOW.Closed() {
		WINDOW.Clear(colornames.Skyblue)

		deltaTime := refreshDeltaTime()
		runtime += deltaTime

		if runtime < 5.0 {
			title.sprite.Draw(WINDOW, pixel.IM.Moved(WINDOW.Bounds().Center()))

			dvd.updateTransform()
			dvd.rot += 7 * deltaTime
			dvd.sprite.Draw(WINDOW, dvd.transf)
			WINDOW.Update()
			continue
		}

		input()

		glep.updatePosition(glep_speed, deltaTime)
		glep.updateTransform()
		glep.sprite.Draw(WINDOW, glep.transf)

		smileyball.sprite.Draw(WINDOW, smileyball.transf)

		text_debug.Clear()
		fmt.Fprintf(text_debug, "FPS: %.0f", math.Round(1/deltaTime))
		text_debug.Draw(WINDOW, pixel.IM)

		WINDOW.Update()
	}
}

func main() {
	pixelgl.Run(gameLoop)
}
