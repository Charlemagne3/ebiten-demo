package main

import (
	"image"
)

type Collider interface {
	Hitbox(x, y int) image.Rectangle
}

// Hitbox returns a player hitbox rectangle offset by x and y, and simulates perspective
func (p *Player) Hitbox(x, y int) image.Rectangle {
	// ebiten renders from the min vertex (top left)
	// To simulate render from the center of "feet" of sprites, we tranlate up (negative Y) by the sprite height and left (negative X) by half the sprite width
	// To simulate perspective, we also limit the hitbox to the bottom half of the sprite by translating the min point down (positive Y) by half the sprite height
	// This results in a translating up (negative Y by half the sprite height)
	offset := p.Sprite.FrameHeight / 2
	return image.Rect(p.X+x-p.Sprite.FrameWidth/2, p.Y+y-offset, p.X+x+p.Sprite.FrameWidth/2, p.Y+y)
}

// Hitbox returns a character hitbox rectangle offset by x and y, and simulates perspective
func (c *Character) Hitbox(x, y int) image.Rectangle {
	// ebiten renders from the min vertex (top left)
	// To simulate render from the center of "feet" of sprites, we tranlate up (negative Y) by the sprite height and left (negative X) by half the sprite width
	// To simulate perspective, we also limit the hitbox to the bottom half of the sprite by translating the min point down (positive Y) by half the sprite height
	// This results in a translating up (negative Y by half the sprite height)
	offset := c.Sprite.FrameHeight / 2
	return image.Rect(c.X+x-c.Sprite.FrameWidth/2, c.Y+y-offset, c.X+x+c.Sprite.FrameWidth/2, c.Y+y)
}

// Hitbox returns a character hitbox rectangle offset by x and y, and simulates perspective
func (e *Enemy) Hitbox(x, y int) image.Rectangle {
	// ebiten renders from the min vertex (top left)
	// To simulate render from the center of "feet" of sprites, we tranlate up (negative Y) by the sprite height and left (negative X) by half the sprite width
	// To simulate perspective, we also limit the hitbox to the bottom half of the sprite by translating the min point down (positive Y) by half the sprite height
	// This results in a translating up (negative Y by half the sprite height)
	offset := e.Sprite.FrameHeight / 2
	return image.Rect(e.X+x-e.Sprite.FrameWidth/2, e.Y+y-offset, e.X+x+e.Sprite.FrameWidth/2, e.Y+y)
}

// Hitbox returns a doodad hitbox rectangle offset by x and y
func (d *Doodad) Hitbox(x, y int) image.Rectangle {
	offset := d.Sprite.FrameHeight / 2
	return image.Rect(d.X+x-d.Sprite.FrameWidth/2, d.Y+y-offset, d.X+x+d.Sprite.FrameWidth/2, d.Y+y)
}

// Hitbox returns a tile hitbox rectangle offset by x and y
func (t *Tile) Hitbox(x, y int) image.Rectangle {
	if t.Collider {
		offset := t.Sprite.FrameHeight
		return image.Rect(t.X+x-t.Sprite.FrameWidth/2, t.Y+y-offset, t.X+x+t.Sprite.FrameWidth/2, t.Y+y)
	} else {
		return image.Rectangle{}
	}
}

// Hitbox returns a character hitbox rectangle offset by x and y, and simulates perspective
func (p Projectile) Hitbox(x, y int) image.Rectangle {
	// ebiten renders from the min vertex (top left)
	// To simulate render from the center of "feet" of sprites, we tranlate up (negative Y) by the sprite height and left (negative X) by half the sprite width
	// To simulate perspective, we also limit the hitbox to the bottom half of the sprite by translating the min point down (positive Y) by half the sprite height
	// This results in a translating up (negative Y by half the sprite height)
	offset := p.Sprite.FrameHeight
	return image.Rect(p.X+x-p.Sprite.FrameWidth/2, p.Y+y-offset, p.X+x+p.Sprite.FrameWidth/2, p.Y+y)
}
