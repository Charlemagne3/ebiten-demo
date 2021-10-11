package main

import (
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var op *ebiten.DrawImageOptions
var x int
var y int
var keyBuf ebiten.Key
var imgFrameNum int
var attack bool
var img MetaImage
var linkSprites map[string]MetaImage

type MetaImage struct {
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
	op = &ebiten.DrawImageOptions{}

	linkSprites = map[string]MetaImage{}

	linkStandSouth, _, err := ebitenutil.NewImageFromFile("link_stand_south.png", 0)
	if err != nil {
		log.Fatal(err)
	}
	linkSprites["linkStandSouth"] = MetaImage{
		FrameLen:    1,
		FrameHeight: 26,
		FrameWidth:  16,
		Image:       linkStandSouth,
	}

	linkStandNorth, _, err := ebitenutil.NewImageFromFile("link_stand_north.png", 0)
	if err != nil {
		log.Fatal(err)
	}
	linkSprites["linkStandNorth"] = MetaImage{
		FrameLen:    1,
		FrameHeight: 26,
		FrameWidth:  16,
		Image:       linkStandNorth,
	}

	linkStandWest, _, err := ebitenutil.NewImageFromFile("link_stand_west.png", 0)
	if err != nil {
		log.Fatal(err)
	}
	linkSprites["linkStandWest"] = MetaImage{
		FrameLen:    1,
		FrameHeight: 24,
		FrameWidth:  17,
		Image:       linkStandWest,
	}

	linkStandEast, _, err := ebitenutil.NewImageFromFile("link_stand_east.png", 0)
	if err != nil {
		log.Fatal(err)
	}
	linkSprites["linkStandEast"] = MetaImage{
		FrameLen:    1,
		FrameHeight: 24,
		FrameWidth:  17,
		Image:       linkStandEast,
	}

	linkWalkSouth, _, err := ebitenutil.NewImageFromFile("link_walk_south.png", 0)
	if err != nil {
		log.Fatal(err)
	}
	linkSprites["linkWalkSouth"] = MetaImage{
		FrameLen:    8,
		FrameHeight: 26,
		FrameWidth:  16,
		Image:       linkWalkSouth,
	}

	linkWalkNorth, _, err := ebitenutil.NewImageFromFile("link_walk_north.png", 0)
	if err != nil {
		log.Fatal(err)
	}
	linkSprites["linkWalkNorth"] = MetaImage{
		FrameLen:    8,
		FrameHeight: 26,
		FrameWidth:  16,
		Image:       linkWalkNorth,
	}

	linkWalkWest, _, err := ebitenutil.NewImageFromFile("link_walk_west.png", 0)
	if err != nil {
		log.Fatal(err)
	}
	linkSprites["linkWalkWest"] = MetaImage{
		FrameLen:    8,
		FrameHeight: 25,
		FrameWidth:  17,
		Image:       linkWalkWest,
	}

	linkWalkEast, _, err := ebitenutil.NewImageFromFile("link_walk_east.png", 0)
	if err != nil {
		log.Fatal(err)
	}
	linkSprites["linkWalkEast"] = MetaImage{
		FrameLen:    8,
		FrameHeight: 25,
		FrameWidth:  17,
		Image:       linkWalkEast,
	}

	linkAttackEast, _, err := ebitenutil.NewImageFromFile("link_attack_east.png", 0)
	if err != nil {
		log.Fatal(err)
	}
	linkSprites["linkAttackEast"] = MetaImage{
		FrameLen:    8,
		FrameHeight: 24,
		FrameWidth:  20,
		Image:       linkAttackEast,
	}

	img = linkSprites["linkStandSouth"]
}

type Game struct{}

func (g *Game) Update(screen *ebiten.Image) error {

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if keyBuf != ebiten.KeyLeft {
			keyBuf = ebiten.KeyLeft
			imgFrameNum = 0
			img = linkSprites["linkWalkWest"]
		} else {
			imgFrameNum++
			imgFrameNum = imgFrameNum % (img.FrameLen - 1)
		}
		x--
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if keyBuf != ebiten.KeyRight {
			keyBuf = ebiten.KeyRight
			imgFrameNum = 0
			img = linkSprites["linkWalkEast"]
		} else {
			imgFrameNum++
			imgFrameNum = imgFrameNum % (img.FrameLen - 1)
		}
		x++
	} else if ebiten.IsKeyPressed(ebiten.KeyUp) {
		if keyBuf != ebiten.KeyUp {
			keyBuf = ebiten.KeyUp
			imgFrameNum = 0
			img = linkSprites["linkWalkNorth"]
		} else {
			imgFrameNum++
			imgFrameNum = imgFrameNum % (img.FrameLen - 1)
		}
		y--
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		if keyBuf != ebiten.KeyDown {
			keyBuf = ebiten.KeyDown
			imgFrameNum = 0
			img = linkSprites["linkWalkSouth"]
		} else {
			imgFrameNum++
			imgFrameNum = imgFrameNum % (img.FrameLen - 1)
		}
		y++
	} else {
		if keyBuf == ebiten.KeyLeft {
			img = linkSprites["linkStandWest"]
		}
		if keyBuf == ebiten.KeyRight {
			img = linkSprites["linkStandEast"]
		}
		if keyBuf == ebiten.KeyUp {
			img = linkSprites["linkStandNorth"]
		}
		if keyBuf == ebiten.KeyDown {
			img = linkSprites["linkStandSouth"]
		}
		keyBuf = -1
		if !attack {
			imgFrameNum = 0
		}
	}

	if attack {
		imgFrameNum++ // passing up the limit and going blank here
		if imgFrameNum >= img.FrameLen {
			attack = false
		}
	}

	if !attack && ebiten.IsKeyPressed(ebiten.KeySpace) {
		attack = true
		img = linkSprites["linkAttackEast"]
		imgFrameNum = 0
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op.GeoM.Reset()
	op.GeoM.Translate(float64(x), float64(y))
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
