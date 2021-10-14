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
	"github.com/hajimehoshi/ebiten/inpututil"
)

type Collider interface {
	Hitbox(x, y int, flat bool) image.Rectangle
}

// Game is an ebiten Game interface implemetation plus custom struct data
type Game struct {
	Player  Player
	Doodads []Doodad
	Sprites map[string]Sprite
	Options *ebiten.DrawImageOptions
}

// Player represents the player character
type Player struct {
	X         int        // The current X screen offset of the player
	Y         int        // The current Y screen offset of the player
	Animation bool       // Whether or not the player is in a special animation or the normal stand/walk cycle.
	LastDir   ebiten.Key // The last direction the player faced (never -1)
	Sprite    Sprite     // The current sprite for the player
	FrameNum  int        // The current frame of the sprite for the player
}

// Doodad represents a static environmental item
type Doodad struct {
	X        int    // The current X screen offset of the player
	Y        int    // The current Y screen offset of the player
	Sprite   Sprite // The current sprite for the player
	FrameNum int    // The current frame of the sprite for the player
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

// IsOtherDirectionJustReleased checks if one of the three cardinal directions other than the key passed in was just released
// This is used to reset the walk cycel animation for a new direction
func IsOtherDirectionJustReleased(key ebiten.Key) bool {
	switch key {
	case ebiten.KeyLeft:
		return inpututil.IsKeyJustReleased(ebiten.KeyRight) || inpututil.IsKeyJustReleased(ebiten.KeyUp) || inpututil.IsKeyJustReleased(ebiten.KeyDown)
	case ebiten.KeyRight:
		return inpututil.IsKeyJustReleased(ebiten.KeyLeft) || inpututil.IsKeyJustReleased(ebiten.KeyUp) || inpututil.IsKeyJustReleased(ebiten.KeyDown)
	case ebiten.KeyUp:
		return inpututil.IsKeyJustReleased(ebiten.KeyLeft) || inpututil.IsKeyJustReleased(ebiten.KeyRight) || inpututil.IsKeyJustReleased(ebiten.KeyDown)
	case ebiten.KeyDown:
		return inpututil.IsKeyJustReleased(ebiten.KeyLeft) || inpututil.IsKeyJustReleased(ebiten.KeyRight) || inpututil.IsKeyJustReleased(ebiten.KeyUp)
	default:
		return false
	}
}

// IsLeastKeyPressDuration checks if a direction key is the most recently pressed
// This is used to set the correct walk direction if the player is holding down multiple walk buttons at once and releases one of them
func IsLeastKeyPressDuration(key ebiten.Key) bool {
	d := inpututil.KeyPressDuration(key)
	up := inpututil.KeyPressDuration(ebiten.KeyUp)
	down := inpututil.KeyPressDuration(ebiten.KeyDown)
	left := inpututil.KeyPressDuration(ebiten.KeyLeft)
	right := inpututil.KeyPressDuration(ebiten.KeyRight)
	return d > 0 &&
		(key == ebiten.KeyUp || d < up || up == 0) &&
		(key == ebiten.KeyDown || d < down || down == 0) &&
		(key == ebiten.KeyLeft || d < left || left == 0) &&
		(key == ebiten.KeyRight || d < right || right == 0)
}

// Hitbox returns a player hitbox rectanle offset by x and y, and simulates perspective if flat is false
func (p *Player) Hitbox(x, y int, flat bool) image.Rectangle {
	// ebiten renders from the min vertex (top left)
	// To simulate render from the center of "feet" of sprites, we tranlate up (negative Y) by the sprite height and left (negative X) by half the sprite width
	// To simulate perspective, we also limit the hitbox to the bottom half of the sprite by translating the min point down (positive Y) by half the sprite height
	// This results in a translating up (negative Y by half the sprite height)
	offset := p.Sprite.FrameHeight
	if !flat {
		offset /= 2
	}
	return image.Rect(p.X+x-p.Sprite.FrameWidth/2, p.Y+y-offset, p.X+x+p.Sprite.FrameWidth/2, p.Y+y)
}

// Hitbox returns a doodad hitbox rectanle offset by x and y, and simulates perspective if flat is false
func (d *Doodad) Hitbox(x, y int, flat bool) image.Rectangle {
	offset := d.Sprite.FrameHeight
	if !flat {
		offset /= 2
	}
	return image.Rect(d.X+x-d.Sprite.FrameWidth/2, d.Y+y-offset, d.X+x+d.Sprite.FrameWidth/2, d.Y+y)
}

func (g *Game) Update(screen *ebiten.Image) error {

	animEnd := false

	if g.Player.Sprite.FrameLen > 1 {
		g.Player.FrameNum++
		// FrameNum is zero indexed and FrameLen is a natural number, so subtract 1 for the mod operation
		g.Player.FrameNum = g.Player.FrameNum % (g.Player.Sprite.FrameLen - 1)
		// End the animation if the last render was the last frame
		if g.Player.Animation && g.Player.FrameNum == 0 {
			g.Player.Animation = false
			animEnd = true
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		// Start the walk left animation if the player just pressed left or if an animation ended and the player was already moving left
		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || (IsOtherDirectionJustReleased(ebiten.KeyLeft) && IsLeastKeyPressDuration(ebiten.KeyLeft)) || animEnd {
			g.Player.LastDir = ebiten.KeyLeft
			g.Player.FrameNum = 0
			g.Player.Sprite = g.Sprites["linkWalkWest"]
		}
		playerRect := g.Player.Hitbox(-1, 0, false)
		move := true
		for _, v := range g.Doodads {
			isCollision := playerRect.Overlaps(v.Hitbox(0, 0, true))
			if isCollision {
				move = false
				break
			}
		}
		if move {
			g.Player.X--
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		// Start the walk right animation if the player just pressed right or if an animation ended and the player was already moving right
		if inpututil.IsKeyJustPressed(ebiten.KeyRight) || (IsOtherDirectionJustReleased(ebiten.KeyRight) && IsLeastKeyPressDuration(ebiten.KeyRight)) || animEnd {
			g.Player.LastDir = ebiten.KeyRight
			g.Player.FrameNum = 0
			g.Player.Sprite = g.Sprites["linkWalkEast"]
		}
		playerRect := g.Player.Hitbox(1, 0, false)
		move := true
		for _, v := range g.Doodads {
			isCollision := playerRect.Overlaps(v.Hitbox(0, 0, true))
			if isCollision {
				move = false
				break
			}
		}
		if move {
			g.Player.X++
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		// Start the walk up animation if the player just pressed up or if an animation ended and the player was already moving up
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) || (IsOtherDirectionJustReleased(ebiten.KeyUp) && IsLeastKeyPressDuration(ebiten.KeyUp)) || animEnd {
			g.Player.LastDir = ebiten.KeyUp
			g.Player.FrameNum = 0
			g.Player.Sprite = g.Sprites["linkWalkNorth"]
		}
		playerRect := g.Player.Hitbox(0, -1, false)
		move := true
		for _, v := range g.Doodads {
			isCollision := playerRect.Overlaps(v.Hitbox(0, 0, true))
			if isCollision {
				move = false
				break
			}
		}
		if move {
			g.Player.Y--
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		// Start the walk down animation if the player just pressed down or if an animation ended and the player was already moving down
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) || (IsOtherDirectionJustReleased(ebiten.KeyDown) && IsLeastKeyPressDuration(ebiten.KeyDown)) || animEnd {
			g.Player.LastDir = ebiten.KeyDown
			g.Player.FrameNum = 0
			g.Player.Sprite = g.Sprites["linkWalkSouth"]
		}
		playerRect := g.Player.Hitbox(0, 1, false)
		move := true
		for _, v := range g.Doodads {
			isCollision := playerRect.Overlaps(v.Hitbox(0, 0, true))
			if isCollision {
				move = false
				break
			}
		}
		if move {
			g.Player.Y++
		}
	}

	// If no direction is pressed and the player is not in an animation, select a standing sprite based on the last direction the player moved
	if !ebiten.IsKeyPressed(ebiten.KeyLeft) && !ebiten.IsKeyPressed(ebiten.KeyRight) && !ebiten.IsKeyPressed(ebiten.KeyUp) && !ebiten.IsKeyPressed(ebiten.KeyDown) && !g.Player.Animation {
		g.Player.FrameNum = 0
		if g.Player.LastDir == ebiten.KeyLeft {
			g.Player.Sprite = g.Sprites["linkStandWest"]
		} else if g.Player.LastDir == ebiten.KeyRight {
			g.Player.Sprite = g.Sprites["linkStandEast"]
		} else if g.Player.LastDir == ebiten.KeyUp {
			g.Player.Sprite = g.Sprites["linkStandNorth"]
		} else if g.Player.LastDir == ebiten.KeyDown {
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
	for _, d := range g.Doodads {
		g.Options.GeoM.Reset()
		g.Options.GeoM.Translate(float64(d.X-d.Sprite.FrameWidth/2), float64(d.Y-d.Sprite.FrameHeight))
		screen.DrawImage(g.Sprites["stump"].Image, g.Options)
	}

	g.Options.GeoM.Reset()
	g.Options.GeoM.Translate(float64(g.Player.X-g.Player.Sprite.FrameWidth/2), float64(g.Player.Y-g.Player.Sprite.FrameHeight))
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
			X:         8,
			Y:         21,
			LastDir:   ebiten.KeyDown,
			Animation: false,
			FrameNum:  0,
			Sprite:    sprite,
		},
		Doodads: []Doodad{
			{
				X:        100,
				Y:        100,
				FrameNum: 0,
				Sprite:   linkSprites["stump"],
			},
			{
				X:        132,
				Y:        132,
				FrameNum: 0,
				Sprite:   linkSprites["stump"],
			},
		},
		Sprites: linkSprites,
		Options: op,
	}

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
