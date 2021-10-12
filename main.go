package main

import (
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var x int
var y int
var keyBuf ebiten.Key
var lastDir ebiten.Key
var imgFrameNum int
var animation bool
var img Sprite
var linkSprites map[string]Sprite
var op *ebiten.DrawImageOptions

type Sprite struct {
	FrameLen    int
	FrameHeight int
	FrameWidth  int
	Image       *ebiten.Image
}

func init() {
	x = 0
	y = 0
	imgFrameNum = 0
	keyBuf = -1
	lastDir = ebiten.KeyDown
	animation = false
	op = &ebiten.DrawImageOptions{}

	linkSprites = map[string]Sprite{}

	linkStandSouth, _, err := ebitenutil.NewImageFromFile("link_stand_south.png", 0)
	if err != nil {
		log.Fatal(err)
	}
	linkSprites["linkStandSouth"] = Sprite{
		FrameLen:    1,
		FrameHeight: 26,
		FrameWidth:  16,
		Image:       linkStandSouth,
	}

	linkStandNorth, _, err := ebitenutil.NewImageFromFile("link_stand_north.png", 0)
	if err != nil {
		log.Fatal(err)
	}
	linkSprites["linkStandNorth"] = Sprite{
		FrameLen:    1,
		FrameHeight: 26,
		FrameWidth:  16,
		Image:       linkStandNorth,
	}

	linkStandWest, _, err := ebitenutil.NewImageFromFile("link_stand_west.png", 0)
	if err != nil {
		log.Fatal(err)
	}
	linkSprites["linkStandWest"] = Sprite{
		FrameLen:    1,
		FrameHeight: 24,
		FrameWidth:  17,
		Image:       linkStandWest,
	}

	linkStandEast, _, err := ebitenutil.NewImageFromFile("link_stand_east.png", 0)
	if err != nil {
		log.Fatal(err)
	}
	linkSprites["linkStandEast"] = Sprite{
		FrameLen:    1,
		FrameHeight: 24,
		FrameWidth:  17,
		Image:       linkStandEast,
	}

	linkWalkSouth, _, err := ebitenutil.NewImageFromFile("link_walk_south.png", 0)
	if err != nil {
		log.Fatal(err)
	}
	linkSprites["linkWalkSouth"] = Sprite{
		FrameLen:    8,
		FrameHeight: 26,
		FrameWidth:  16,
		Image:       linkWalkSouth,
	}

	linkWalkNorth, _, err := ebitenutil.NewImageFromFile("link_walk_north.png", 0)
	if err != nil {
		log.Fatal(err)
	}
	linkSprites["linkWalkNorth"] = Sprite{
		FrameLen:    8,
		FrameHeight: 26,
		FrameWidth:  16,
		Image:       linkWalkNorth,
	}

	linkWalkWest, _, err := ebitenutil.NewImageFromFile("link_walk_west.png", 0)
	if err != nil {
		log.Fatal(err)
	}
	linkSprites["linkWalkWest"] = Sprite{
		FrameLen:    8,
		FrameHeight: 25,
		FrameWidth:  18,
		Image:       linkWalkWest,
	}

	linkWalkEast, _, err := ebitenutil.NewImageFromFile("link_walk_east.png", 0)
	if err != nil {
		log.Fatal(err)
	}
	linkSprites["linkWalkEast"] = Sprite{
		FrameLen:    8,
		FrameHeight: 25,
		FrameWidth:  18,
		Image:       linkWalkEast,
	}

	linkAttackSouth, _, err := ebitenutil.NewImageFromFile("link_attack_south.png", 0)
	if err != nil {
		log.Fatal(err)
	}
	linkSprites["linkAttackSouth"] = Sprite{
		FrameLen:    8,
		FrameHeight: 24,
		FrameWidth:  16,
		Image:       linkAttackSouth,
	}

	linkAttackNorth, _, err := ebitenutil.NewImageFromFile("link_attack_north.png", 0)
	if err != nil {
		log.Fatal(err)
	}
	linkSprites["linkAttackNorth"] = Sprite{
		FrameLen:    8,
		FrameHeight: 29,
		FrameWidth:  16,
		Image:       linkAttackNorth,
	}

	linkAttackWest, _, err := ebitenutil.NewImageFromFile("link_attack_west.png", 0)
	if err != nil {
		log.Fatal(err)
	}
	linkSprites["linkAttackWest"] = Sprite{
		FrameLen:    8,
		FrameHeight: 24,
		FrameWidth:  20,
		Image:       linkAttackWest,
	}

	linkAttackEast, _, err := ebitenutil.NewImageFromFile("link_attack_east.png", 0)
	if err != nil {
		log.Fatal(err)
	}
	linkSprites["linkAttackEast"] = Sprite{
		FrameLen:    8,
		FrameHeight: 24,
		FrameWidth:  20,
		Image:       linkAttackEast,
	}

	img = linkSprites["linkStandSouth"]
}

type Game struct{}

func (g *Game) Update(screen *ebiten.Image) error {

	if img.FrameLen > 1 {
		imgFrameNum++
		imgFrameNum = imgFrameNum % (img.FrameLen - 1)
		// End the animation if the last render was the last frame
		if animation && imgFrameNum == 0 {
			animation = false
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if keyBuf != ebiten.KeyLeft {
			keyBuf = ebiten.KeyLeft
			lastDir = ebiten.KeyLeft
			imgFrameNum = 0
			img = linkSprites["linkWalkWest"]
		}
		x--
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if keyBuf != ebiten.KeyRight {
			keyBuf = ebiten.KeyRight
			lastDir = ebiten.KeyRight
			imgFrameNum = 0
			img = linkSprites["linkWalkEast"]
		}
		x++
	} else if ebiten.IsKeyPressed(ebiten.KeyUp) {
		if keyBuf != ebiten.KeyUp {
			keyBuf = ebiten.KeyUp
			lastDir = ebiten.KeyUp
			imgFrameNum = 0
			img = linkSprites["linkWalkNorth"]
		}
		y--
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		if keyBuf != ebiten.KeyDown {
			keyBuf = ebiten.KeyDown
			lastDir = ebiten.KeyDown
			imgFrameNum = 0
			img = linkSprites["linkWalkSouth"]
		}
		y++
	} else {
		if !animation && lastDir == ebiten.KeyLeft {
			keyBuf = -1
			imgFrameNum = 0
			img = linkSprites["linkStandWest"]
		} else if !animation && lastDir == ebiten.KeyRight {
			keyBuf = -1
			imgFrameNum = 0
			img = linkSprites["linkStandEast"]
		} else if !animation && lastDir == ebiten.KeyUp {
			keyBuf = -1
			imgFrameNum = 0
			img = linkSprites["linkStandNorth"]
		} else if !animation && lastDir == ebiten.KeyDown {
			keyBuf = -1
			imgFrameNum = 0
			img = linkSprites["linkStandSouth"]
		}
	}

	// If starting an animation
	if !animation && ebiten.IsKeyPressed(ebiten.KeySpace) {
		animation = true
		imgFrameNum = 0
		if lastDir == ebiten.KeyLeft {
			img = linkSprites["linkAttackWest"]
		} else if lastDir == ebiten.KeyRight {
			img = linkSprites["linkAttackEast"]
		} else if lastDir == ebiten.KeyUp {
			img = linkSprites["linkAttackNorth"]
		} else if lastDir == ebiten.KeyDown {
			img = linkSprites["linkAttackSouth"]
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op.GeoM.Reset()
	op.GeoM.Translate(float64(x), float64(y))
	// sub-rect is the width of a frame times the frame number, plus the frame number for the 1-pixel buffer between frames
	screen.DrawImage(img.Image.SubImage(image.Rect(img.FrameWidth*imgFrameNum+imgFrameNum, 0, img.FrameWidth*imgFrameNum+imgFrameNum+img.FrameWidth, img.FrameHeight)).(*ebiten.Image), op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("grame")
	game := &Game{}
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
