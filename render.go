package main

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type RenderTarget interface {
	RenderSprite() *ebiten.Image
	RenderOptions() *ebiten.DrawImageOptions
	RenderOrder() int
}

func (p *Player) RenderSprite() *ebiten.Image {
	// sub-rect is the width of a frame times the frame number, plus the frame number for the 1-pixel buffer between frames
	return p.Sprite.Image.SubImage(image.Rect(p.Sprite.FrameWidth*p.FrameNum+p.FrameNum, 0, p.Sprite.FrameWidth*p.FrameNum+p.FrameNum+p.Sprite.FrameWidth, p.Sprite.FrameHeight)).(*ebiten.Image)
}

func (p *Player) RenderOptions() *ebiten.DrawImageOptions {
	o := ebiten.DrawImageOptions{}
	// ebiten renders from the min vertex (top left). Offset by the frameheight and half the framewidth to emulate rendering from the "feet" of the sprite
	o.GeoM.Translate(float64(p.X-p.Sprite.FrameWidth/2), float64(p.Y-p.Sprite.FrameHeight))
	return &o
}

func (p *Player) RenderOrder() int {
	return p.Y
}

func (c *Character) RenderSprite() *ebiten.Image {
	return c.Sprite.Image.SubImage(image.Rect(c.Sprite.FrameWidth*c.FrameNum+c.FrameNum, 0, c.Sprite.FrameWidth*c.FrameNum+c.FrameNum+c.Sprite.FrameWidth, c.Sprite.FrameHeight)).(*ebiten.Image)
}

func (c *Character) RenderOptions() *ebiten.DrawImageOptions {
	o := ebiten.DrawImageOptions{}
	o.GeoM.Translate(float64(c.X-c.Sprite.FrameWidth/2), float64(c.Y-c.Sprite.FrameHeight))
	return &o
}

func (c *Character) RenderOrder() int {
	return c.Y
}

func (e *Enemy) RenderSprite() *ebiten.Image {
	return e.Sprite.Image
}

func (e *Enemy) RenderOptions() *ebiten.DrawImageOptions {
	o := ebiten.DrawImageOptions{}
	o.GeoM.Translate(float64(e.X-e.Sprite.FrameWidth/2), float64(e.Y-e.Sprite.FrameHeight))
	return &o
}

func (e *Enemy) RenderOrder() int {
	return e.Y
}

func (d *Doodad) RenderSprite() *ebiten.Image {
	return d.Sprite.Image
}

func (d *Doodad) RenderOptions() *ebiten.DrawImageOptions {
	o := ebiten.DrawImageOptions{}
	o.GeoM.Translate(float64(d.X-d.Sprite.FrameWidth/2), float64(d.Y-d.Sprite.FrameHeight))
	return &o
}

func (d *Doodad) RenderOrder() int {
	return d.Y
}

func (t *Tile) RenderSprite() *ebiten.Image {
	return t.Sprite.Image
}

func (t *Tile) RenderOptions() *ebiten.DrawImageOptions {
	o := ebiten.DrawImageOptions{}
	o.GeoM.Translate(float64(t.X-t.Sprite.FrameWidth/2), float64(t.Y-t.Sprite.FrameHeight))
	return &o
}

func (t *Tile) RenderOrder() int {
	if t.Collider {
		return math.MinInt / 2
	} else {
		return math.MinInt
	}
}

func (p *Projectile) RenderSprite() *ebiten.Image {
	return p.Sprite.Image
}

func (p *Projectile) RenderOptions() *ebiten.DrawImageOptions {
	o := ebiten.DrawImageOptions{}
	o.GeoM.Translate(float64(p.X-p.Sprite.FrameWidth/2), float64(p.Y-p.Sprite.FrameHeight))
	return &o
}

func (p *Projectile) RenderOrder() int {
	return p.Y
}
