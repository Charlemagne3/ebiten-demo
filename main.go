package main

import (
	"encoding/json"
	"image"
	_ "image/png"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// Game is an ebiten Game interface implemetation plus custom struct data
type Game struct {
	Player  Player
	Sprites map[string]Sprite
	Options *ebiten.DrawImageOptions
}

// Player represents the player character
type Player struct {
	X         int        // The current X screen offset of the player
	Y         int        // The current Y screen offset of the player
	Animation bool       // Whether or not the player is in a special animation or the normal stand/walk cycle.
	KeyBuf    ebiten.Key // The key pressed on the last frame or -1
	LastDir   ebiten.Key // The last direction the player faced (never -1)
	Sprite    Sprite     // The current sprite for the player
	FrameNum  int        // The current frame of the sprite for the player
}

// Sprite represents an image with a number of sub-frames in it to be rendered via rectangles
type Sprite struct {
	FrameWidth  int
	FrameHeight int
	FrameLen    int // How many frames are in the sprite
	Image       *ebiten.Image
}

// SpriteJSON represents the json to be read from the sprite json file.
type SpriteJSON struct {
	FrameWidth  int    `json:"frameWidth"`
	FrameHeight int    `json:"frameHeight"`
	FrameLen    int    `json:"frameLen"`
	Image       string `json:"image"`
}

type Rect struct {
	X int
	Y int
	W int
	H int
}

// IsCollision checks if each vertex of rectangle is inside another rectangle
func IsCollision(r0, r1 Rect) bool {
	return r0.X > r1.X && r0.X < r1.X+r1.W && r0.Y > r1.Y && r0.Y < r1.Y+r1.H ||
		r0.X+r0.W > r1.X && r0.X+r0.W < r1.X+r1.W && r0.Y > r1.Y && r0.Y < r1.Y+r1.H ||
		r0.X > r1.X && r0.X < r1.X+r1.W && r0.Y+r0.H > r1.Y && r0.Y+r0.H < r1.Y+r1.H ||
		r0.X+r0.W > r1.X && r0.X+r0.W < r1.X+r1.W && r0.Y+r0.H > r1.Y && r0.Y+r0.H < r1.Y+r1.H
}

func (g *Game) Update(screen *ebiten.Image) error {
	playerRect := Rect{
		X: g.Player.X,
		Y: g.Player.Y,
		W: g.Player.Sprite.FrameWidth,
		H: g.Player.Sprite.FrameHeight,
	}

	otherRect := Rect{
		X: 100,
		Y: 100,
		W: 32,
		H: 32,
	}

	if g.Player.Sprite.FrameLen > 1 {
		g.Player.FrameNum++
		// FrameNum is zero indexed and FrameLen is a natural number, so subtract 1 for the mod operation
		g.Player.FrameNum = g.Player.FrameNum % (g.Player.Sprite.FrameLen - 1)
		// End the animation if the last render was the last frame
		if g.Player.Animation && g.Player.FrameNum == 0 {
			g.Player.Animation = false
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if g.Player.KeyBuf != ebiten.KeyLeft {
			g.Player.KeyBuf = ebiten.KeyLeft
			g.Player.LastDir = ebiten.KeyLeft
			g.Player.FrameNum = 0
			g.Player.Sprite = g.Sprites["linkWalkWest"]
		}
		if !IsCollision(playerRect, otherRect) {
			g.Player.X--
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if g.Player.KeyBuf != ebiten.KeyRight {
			g.Player.KeyBuf = ebiten.KeyRight
			g.Player.LastDir = ebiten.KeyRight
			g.Player.FrameNum = 0
			g.Player.Sprite = g.Sprites["linkWalkEast"]
		}
		if !IsCollision(playerRect, otherRect) {
			g.Player.X++
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyUp) {
		if g.Player.KeyBuf != ebiten.KeyUp {
			g.Player.KeyBuf = ebiten.KeyUp
			g.Player.LastDir = ebiten.KeyUp
			g.Player.FrameNum = 0
			g.Player.Sprite = g.Sprites["linkWalkNorth"]
		}
		if !IsCollision(playerRect, otherRect) {
			g.Player.Y--
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		if g.Player.KeyBuf != ebiten.KeyDown {
			g.Player.KeyBuf = ebiten.KeyDown
			g.Player.LastDir = ebiten.KeyDown
			g.Player.FrameNum = 0
			g.Player.Sprite = g.Sprites["linkWalkSouth"]
		}
		if !IsCollision(playerRect, otherRect) {
			g.Player.Y++
		}
	} else {
		if !g.Player.Animation && g.Player.LastDir == ebiten.KeyLeft {
			g.Player.KeyBuf = -1
			g.Player.FrameNum = 0
			g.Player.Sprite = g.Sprites["linkStandWest"]
		} else if !g.Player.Animation && g.Player.LastDir == ebiten.KeyRight {
			g.Player.KeyBuf = -1
			g.Player.FrameNum = 0
			g.Player.Sprite = g.Sprites["linkStandEast"]
		} else if !g.Player.Animation && g.Player.LastDir == ebiten.KeyUp {
			g.Player.KeyBuf = -1
			g.Player.FrameNum = 0
			g.Player.Sprite = g.Sprites["linkStandNorth"]
		} else if !g.Player.Animation && g.Player.LastDir == ebiten.KeyDown {
			g.Player.KeyBuf = -1
			g.Player.FrameNum = 0
			g.Player.Sprite = g.Sprites["linkStandSouth"]
		}
	}

	// If starting an animation
	if !g.Player.Animation && ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.Player.Animation = true
		g.Player.FrameNum = 0
		if g.Player.LastDir == ebiten.KeyLeft {
			g.Player.Sprite = g.Sprites["linkAttackWest"]
		} else if g.Player.LastDir == ebiten.KeyRight {
			g.Player.Sprite = g.Sprites["linkAttackEast"]
		} else if g.Player.LastDir == ebiten.KeyUp {
			g.Player.Sprite = g.Sprites["linkAttackNorth"]
		} else if g.Player.LastDir == ebiten.KeyDown {
			g.Player.Sprite = g.Sprites["linkAttackSouth"]
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Options.GeoM.Reset()
	g.Options.GeoM.Translate(100, 100)
	screen.DrawImage(g.Sprites["stump"].Image, g.Options)

	g.Options.GeoM.Reset()
	g.Options.GeoM.Translate(float64(g.Player.X), float64(g.Player.Y))
	// sub-rect is the width of a frame times the frame number, plus the frame number for the 1-pixel buffer between frames
	screen.DrawImage(g.Player.Sprite.Image.SubImage(image.Rect(g.Player.Sprite.FrameWidth*g.Player.FrameNum+g.Player.FrameNum, 0, g.Player.Sprite.FrameWidth*g.Player.FrameNum+g.Player.FrameNum+g.Player.Sprite.FrameWidth, g.Player.Sprite.FrameHeight)).(*ebiten.Image), g.Options)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("grame")

	sprites, err := os.ReadFile("./sprites/sprites.json")
	var jsonSprites []SpriteJSON
	err = json.Unmarshal(sprites, &jsonSprites)
	if err != nil {
		log.Fatal(err)
	}

	linkSprites := map[string]Sprite{}
	for _, v := range jsonSprites {
		sprite, _, err := ebitenutil.NewImageFromFile("./sprites/"+v.Image, 0)
		if err != nil {
			log.Fatal(err)
		}

		// camelCase the filenames without the extension to make sprite keys
		var snek bool
		var k string
		image := v.Image[:len(v.Image)-4]
		for i := 0; i < len(image); i++ {
			char := string(image[i])
			if char == "_" {
				snek = true
				continue
			}
			if snek == true {
				snek = false
				k += strings.ToUpper(char)
			} else {
				k += char
			}
		}

		linkSprites[k] = Sprite{
			FrameLen:    v.FrameLen,
			FrameHeight: v.FrameHeight,
			FrameWidth:  v.FrameWidth,
			Image:       sprite,
		}
	}

	sprite := linkSprites["linkStandSouth"]
	op := &ebiten.DrawImageOptions{}

	game := &Game{
		Player: Player{
			X:         0,
			Y:         0,
			KeyBuf:    -1,
			LastDir:   ebiten.KeyDown,
			Animation: false,
			FrameNum:  0,
			Sprite:    sprite,
		},
		Sprites: linkSprites,
		Options: op,
	}

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
