package main

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type RenderTarget interface {
	RenderSprite() Sprite
	RenderImage() *ebiten.Image
	RenderOptions() *ebiten.DrawImageOptions
	RenderOrder() int
	RenderHandle() image.Point
	RenderX() int
	RenderY() int
}

func (p *Player) RenderSprite() Sprite {
	return p.Sprite
}

func (p *Player) RenderImage() *ebiten.Image {
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

func (p *Player) RenderHandle() image.Point {
	return p.Sprite.Handles[p.FrameNum]
}

func (p *Player) RenderX() int {
	return p.X
}

func (p *Player) RenderY() int {
	return p.Y
}

func (c *Character) RenderSprite() Sprite {
	return c.Sprite
}

func (c *Character) RenderImage() *ebiten.Image {
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

func (c *Character) RenderHandle() image.Point {
	return c.Sprite.Handles[c.FrameNum]
}

func (c *Character) RenderX() int {
	return c.X
}

func (c *Character) RenderY() int {
	return c.Y
}

func (e *Enemy) RenderSprite() Sprite {
	return e.Sprite
}

func (e *Enemy) RenderImage() *ebiten.Image {
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

func (e *Enemy) RenderHandle() image.Point {
	return e.Sprite.Handles[e.FrameNum]
}

func (e *Enemy) RenderX() int {
	return e.X
}

func (e *Enemy) RenderY() int {
	return e.Y
}

func (d *Doodad) RenderSprite() Sprite {
	return d.Sprite
}

func (d *Doodad) RenderImage() *ebiten.Image {
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

func (d *Doodad) RenderX() int {
	return d.X
}

func (d *Doodad) RenderY() int {
	return d.Y
}

func (d *Doodad) RenderHandle() image.Point {
	return d.Sprite.Handles[d.FrameNum]
}

func (t *Tile) RenderSprite() Sprite {
	return t.Sprite
}

func (t *Tile) RenderImage() *ebiten.Image {
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

func (t *Tile) RenderX() int {
	return t.X
}

func (t *Tile) RenderY() int {
	return t.Y
}

func (t *Tile) RenderHandle() image.Point {
	return t.Sprite.Handles[t.FrameNum]
}

func (w *Weapon) RenderSprite() Sprite {
	return w.Sprite
}

func (w *Weapon) RenderImage() *ebiten.Image {
	return w.Sprite.Image
}

func (w *Weapon) RenderOptions() *ebiten.DrawImageOptions {
	o := ebiten.DrawImageOptions{}
	o.GeoM.Translate(float64(w.Wielder.RenderX()-w.Wielder.RenderSprite().FrameWidth/2+w.Wielder.RenderHandle().X), float64(w.Wielder.RenderY()-w.Sprite.FrameHeight))
	return &o
}

func (w *Weapon) RenderOrder() int {
	return w.Wielder.RenderOrder() + 1
}

func (w *Weapon) RenderHandle() image.Point {
	return w.Sprite.Handles[w.FrameNum]
}

func (w *Weapon) RenderX() int {
	return 0
}

func (w *Weapon) RenderY() int {
	return 0
}

func (p *Projectile) RenderSprite() Sprite {
	return p.Sprite
}

func (p *Projectile) RenderImage() *ebiten.Image {
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

func (p *Projectile) RenderHandle() image.Point {
	return p.Sprite.Handles[p.FrameNum]
}

func (p *Projectile) RenderX() int {
	return p.X
}

func (p *Projectile) RenderY() int {
	return p.Y
}
