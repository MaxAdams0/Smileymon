package main

import (
	"math"

	"github.com/gopxl/pixel"
)

type Entity struct {
	sprite *pixel.Sprite
	transf pixel.Matrix
	pos    pixel.Vec
	vel    pixel.Vec
	scale  pixel.Vec
	rot    float64
}

func createEntity(image_path string, pos, vel, scale pixel.Vec, rot float64) Entity {
	var new_entity Entity

	new_entity.sprite = createSprite(image_path)
	new_entity.pos = pos
	new_entity.vel = vel
	new_entity.scale = scale
	new_entity.rot = rot

	new_entity.updateTransform()

	return new_entity
}

func (e *Entity) updateTransform() {
	new_transf := pixel.IM
	new_transf = new_transf.Moved(e.pos)
	new_transf = new_transf.Rotated(e.pos, e.rot)
	new_transf = new_transf.ScaledXY(e.pos, e.scale)

	e.transf = new_transf
}

func (e *Entity) updatePosition(speed float64, deltaTime float64) {
	e.pos = e.pos.Add(e.vel.Scaled(speed).Scaled(deltaTime))
}

func (e *Entity) rotateBy(angle float64) {
	e.transf.Rotated(e.pos, angle)
}

func (e *Entity) rotateTo(angle float64) {
	e.transf.Rotated(e.pos, e.rot-angle)
}

func (e *Entity) normalizeVelocity() {
	if e.vel.X != 0 && e.vel.Y != 0 {
		e.vel = e.vel.Scaled(1 / math.Sqrt2)
	}
}
