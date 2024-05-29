package main

import (
	"image"
	"os"
	"time"

	_ "image/png"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var (
	WINDOW *pixelgl.Window // Not a constant, but its close enough (constant in action no definition)

	win_size        = pixel.Vec{1024, 768}
	win_origin      = pixel.Vec{0, 0}
	win_vsync  bool = false

	current_time time.Time = time.Now()
	runtime      float64   = 0.0
)

type Entity struct {
	sprite *pixel.Sprite
	transf pixel.Matrix
	pos    pixel.Vec
	rot    float64
	scale  pixel.Vec
}

func (e *Entity) updateTransform() {
	new_transf := pixel.IM
	new_transf = new_transf.Moved(e.pos)
	new_transf = new_transf.Rotated(e.pos, e.rot)
	new_transf = new_transf.ScaledXY(e.pos, e.scale)

	e.transf = new_transf
}

func (e *Entity) rotateBy(angle float64) {
	e.transf.Rotated(e.pos, angle)
}

func (e *Entity) rotateTo(angle float64) {
	e.transf.Rotated(e.pos, e.rot-angle)
}

func loadPicture(path string) (pixel.Picture, error) {
	// ========== Standard golang get image from file ==========
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close() // closes when function ends
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	// === Convert image to Pixel libraries "Picture" class ===
	return pixel.PictureDataFromImage(img), nil
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

	win.Clear(colornames.Skyblue)
}

func createSprite(image_path string) (*pixel.Sprite, error) {
	img, err := loadPicture(image_path)
	if err != nil {
		return nil, err
	}

	return pixel.NewSprite(img, img.Bounds()), nil
}

func createEntity(image_path string, pos pixel.Vec, rot float64, scale pixel.Vec) Entity {
	var new_entity Entity

	new_entity.sprite, _ = createSprite(image_path)
	new_entity.pos = pos
	new_entity.rot = rot
	new_entity.transf = pixel.IM
	new_entity.transf = new_entity.transf.Moved(pos)
	new_entity.transf = new_entity.transf.Rotated(pos, rot)
	new_entity.scale = scale

	return new_entity
}

func refreshDeltaTime() float64 {
	dt := time.Since(current_time).Seconds()
	current_time = time.Now()
	return dt
}

func gameLoop() {
	createWindow()

	glep := createEntity("images/glep.png", WINDOW.Bounds().Center(), 0.0, pixel.Vec{2.0, 2.0})
	dvd := createEntity("images/dvd.png", pixel.Vec{win_size.X - 96, 96}, 0.0, pixel.Vec{2.0, 2.0})
	title := createEntity("images/titlepage.png", WINDOW.Bounds().Center(), 0.0, pixel.Vec{1.0, 1.0})

	for !WINDOW.Closed() {
		deltaTime := refreshDeltaTime()
		runtime += deltaTime

		if runtime < 5.0 {
			title.sprite.Draw(WINDOW, glep.transf)

			dvd.updateTransform()
			dvd.rot += 7 * deltaTime
			dvd.sprite.Draw(WINDOW, dvd.transf)
			WINDOW.Update()
			continue
		}

		WINDOW.Clear(colornames.Skyblue)

		glep.updateTransform()

		glep.sprite.Draw(WINDOW, glep.transf)

		WINDOW.Update()
	}
}

func main() {
	pixelgl.Run(gameLoop)
}
