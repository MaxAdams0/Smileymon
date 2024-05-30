package main

import "github.com/gopxl/pixel"

type Item struct {
	sprite    *pixel.Sprite
	transf    pixel.Matrix
	pos       pixel.Vec
	scale     pixel.Vec
	rot       float64
	desc      string
	onfloor   bool
	canpickup bool
}

func createItem(image_path string, pos, scale pixel.Vec, rot float64, desc string, onfloor, canpickup bool) Item {
	var new_item Item

	new_item.sprite = createSprite(image_path)
	new_item.pos = pos
	new_item.scale = scale
	new_item.rot = rot
	new_item.desc = desc
	new_item.onfloor = onfloor
	new_item.canpickup = canpickup

	new_item.updateTransform()

	return new_item
}

func (i *Item) updateTransform() {
	new_transf := pixel.IM
	new_transf = new_transf.Moved(i.pos)
	new_transf = new_transf.Rotated(i.pos, i.rot)
	new_transf = new_transf.ScaledXY(i.pos, i.scale)

	i.transf = new_transf
}
