package main

import (
	"encoding/json"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func Contains(a []string, s string) bool {
	for _, v := range a {
		if s == v {
			return true
		}
	}
	return false
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func AbsDiff(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}

func CamelCase(s string) string {
	var snek bool
	var camel string
	for i := 0; i < len(s); i++ {
		char := string(s[i])
		if char == "_" {
			snek = true
			continue
		}
		if snek {
			snek = false
			camel += strings.ToUpper(char)
		} else {
			camel += char
		}
	}
	return camel
}

type Collider interface {
	Hitbox(x, y int) image.Rectangle
}

type InteractionTarget interface {
	Dialogue() string    // The current dialogue to render
	Options() [][]string // The options for the current dialogue, or empty
	SelectOption(int)    // Selects a next or previous option
	SelectedOption() int // Returns the selected option
	AdvanceRune()        // Advances to the next rune
	AdvancePhrase()      // Advances to the next phrase
	IsExhausted() bool   // Returns true if the current dialogue tree is complete
}

// Game is an ebiten Game interface implemetation plus custom struct data
type Game struct {
	Player            Player
	Characters        []Character
	Enemies           []Enemy
	Doodads           []Doodad
	Sprites           map[string]Sprite
	Font              font.Face
	Options           *ebiten.DrawImageOptions
	InteractionTarget InteractionTarget // The target of another game element that the player is having a dialogue interaction with, or nil.
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

type DialogueGraph struct {
	Nodes   map[string]*DialogueNode
	Edges   map[string][]string
	NodeKey string // The current node of dialogue the player is on
	RootKey string // The root node of the current dialogue tree
}

type DialogueNode struct {
	Phrase    string     // The phrase of dialogue
	Options   [][]string // The options on the node, or empty
	RuneNum   int
	OptionNum int  // Which option is selected
	End       bool // Whether or not to end the interaction after this node is completed.
}

// Character represents an npc character
type Character struct {
	X              int                       // The current X screen offset of the character
	Y              int                       // The current Y screen offset of the character
	Animation      bool                      // Whether or not the character is in a special animation or the normal stand/walk cycle.
	LastDir        ebiten.Key                // The last direction the character faced (never -1)
	Sprite         Sprite                    // The current sprite for the character
	FrameNum       int                       // The current frame of the sprite for the character
	DialogueGraphs map[string]*DialogueGraph // The dialogue graphs the character has
	DialogueKey    string                    // The current dialogue graph the character has loaded
}

func (c *Character) Dialogue() string {
	graph := c.DialogueGraphs[c.DialogueKey]
	node := graph.Nodes[graph.NodeKey]
	return node.Phrase[:node.RuneNum]
}

func (c *Character) Options() [][]string {
	graph := c.DialogueGraphs[c.DialogueKey]
	node := graph.Nodes[graph.NodeKey]
	return node.Options
}

func (c *Character) SelectOption(dir int) {
	graph := c.DialogueGraphs[c.DialogueKey]
	node := graph.Nodes[graph.NodeKey]
	node.OptionNum = Min(Max(node.OptionNum+dir, 0), len(node.Options)-1)
}

func (c *Character) SelectedOption() int {
	graph := c.DialogueGraphs[c.DialogueKey]
	node := graph.Nodes[graph.NodeKey]
	return node.OptionNum
}

func (c *Character) AdvanceRune() {
	graph := c.DialogueGraphs[c.DialogueKey]
	node := graph.Nodes[graph.NodeKey]
	node.RuneNum = Min(node.RuneNum+1, len(node.Phrase))
}

func (c *Character) AdvancePhrase() {
	graph := c.DialogueGraphs[c.DialogueKey]
	connections := graph.Edges[graph.NodeKey]
	node := graph.Nodes[graph.NodeKey]
	// If the node has no options then there is only a single node to advance to
	if len(node.Options) == 0 && len(connections) > 0 {
		graph.Nodes[graph.NodeKey].RuneNum = 0
		graph.Nodes[graph.NodeKey].OptionNum = 0
		graph.NodeKey = connections[0]
	} else if len(node.Options) > 0 {
		options := node.Options[node.OptionNum]
		if Contains(connections, options[1]) {
			graph.Nodes[graph.NodeKey].RuneNum = 0
			graph.Nodes[graph.NodeKey].OptionNum = 0
			graph.NodeKey = options[1]
		}
	}
}

func (c *Character) IsExhausted() bool {
	graph := c.DialogueGraphs[c.DialogueKey]
	return graph.Nodes[graph.NodeKey].End
}

// Enemy represents an enemy
type Enemy struct {
	X         int        // The current X screen offset of the enemy
	Y         int        // The current Y screen offset of the enemy
	Animation bool       // Whether or not the enemy is in a special animation or the normal stand/walk cycle.
	LastDir   ebiten.Key // The last direction the enemy faced (never -1)
	Sprite    Sprite     // The current sprite for the enemy
	FrameNum  int        // The current frame of the sprite for the enemy
}

// Doodad represents a static environmental item
type Doodad struct {
	X        int    // The current X screen offset of the doodad
	Y        int    // The current Y screen offset of the doodad
	Sprite   Sprite // The current sprite for the doodad
	FrameNum int    // The current frame of the sprite for the doodad
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

// DialogueJSON represents the json to be read from the dialogue json file.
type DialogueJSON struct {
	ID          string     `json:"id"`
	Phrase      string     `json:"phrase"`
	Options     [][]string `json:"options"`
	Connections []string   `json:"connections"`
	End         bool       `json:"end"`
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

// Hitbox returns a doodad hitbox rectangle offset by x and y
func (d *Doodad) Hitbox(x, y int) image.Rectangle {
	offset := d.Sprite.FrameHeight
	return image.Rect(d.X+x-d.Sprite.FrameWidth/2, d.Y+y-offset, d.X+x+d.Sprite.FrameWidth/2, d.Y+y)
}

func (g *Game) Update() error {

	if g.InteractionTarget != nil {
		// Render the next rune to scroll the text
		g.InteractionTarget.AdvanceRune()
		if inpututil.IsKeyJustReleased(ebiten.KeyLeft) {
			g.InteractionTarget.SelectOption(-1)
		} else if inpututil.IsKeyJustReleased(ebiten.KeyRight) {
			g.InteractionTarget.SelectOption(1)
		} else if inpututil.IsKeyJustReleased(ebiten.KeyEnter) {
			// If out of dialogue, end the interaction
			if g.InteractionTarget.IsExhausted() {
				g.InteractionTarget.AdvancePhrase()
				g.InteractionTarget = nil
			} else {
				g.InteractionTarget.AdvancePhrase()
			}
		}
	} else if !g.Player.Animation && inpututil.IsKeyJustReleased(ebiten.KeyEnter) {
		for _, c := range g.Characters {
			playerRect := g.Player.Hitbox(0, 0)
			characterRect := c.Hitbox(0, 0)
			// Check if a side of the player rect is touching the character rect and the midpoint of that side is touching the character rect
			if AbsDiff(playerRect.Max.X, characterRect.Min.X) <= 1 && playerRect.Min.Y+(playerRect.Dy()/2) >= characterRect.Min.Y && playerRect.Min.Y+(playerRect.Dy()/2) <= characterRect.Max.Y {
				g.InteractionTarget = &c
				g.Player.FrameNum = 0
				g.Player.Sprite = g.Sprites["linkStandEast"]
			} else if AbsDiff(playerRect.Max.Y, characterRect.Min.Y) <= 1 && playerRect.Min.X+(playerRect.Dx()/2) >= characterRect.Min.X && playerRect.Min.X+(playerRect.Dx()/2) <= characterRect.Max.X {
				g.InteractionTarget = &c
				g.Player.FrameNum = 0
				g.Player.Sprite = g.Sprites["linkStandSouth"]
			} else if AbsDiff(playerRect.Min.X, characterRect.Max.X) <= 1 && playerRect.Min.Y+(playerRect.Dy()/2) >= characterRect.Min.Y && playerRect.Min.Y+(playerRect.Dy()/2) <= characterRect.Max.Y {
				g.InteractionTarget = &c
				g.Player.FrameNum = 0
				g.Player.Sprite = g.Sprites["linkStandWest"]
			} else if AbsDiff(playerRect.Min.Y, characterRect.Max.Y) <= 1 && playerRect.Min.X+(playerRect.Dx()/2) >= characterRect.Min.X && playerRect.Min.X+(playerRect.Dx()/2) <= characterRect.Max.X {
				g.InteractionTarget = &c
				g.Player.FrameNum = 0
				g.Player.Sprite = g.Sprites["linkStandNorth"]
			}
		}
	}

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
		playerRect := g.Player.Hitbox(-1, 0)
		move := true
		for _, v := range g.Doodads {
			isCollision := playerRect.Overlaps(v.Hitbox(0, 0))
			if isCollision {
				move = false
				break
			}
		}
		for _, v := range g.Characters {
			isCollision := playerRect.Overlaps(v.Hitbox(0, 0))
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
		playerRect := g.Player.Hitbox(1, 0)
		move := true
		for _, v := range g.Doodads {
			isCollision := playerRect.Overlaps(v.Hitbox(0, 0))
			if isCollision {
				move = false
				break
			}
		}
		for _, v := range g.Characters {
			isCollision := playerRect.Overlaps(v.Hitbox(0, 0))
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
		playerRect := g.Player.Hitbox(0, -1)
		move := true
		for _, v := range g.Doodads {
			isCollision := playerRect.Overlaps(v.Hitbox(0, 0))
			if isCollision {
				move = false
				break
			}
		}
		for _, v := range g.Characters {
			isCollision := playerRect.Overlaps(v.Hitbox(0, 0))
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
		playerRect := g.Player.Hitbox(0, 1)
		move := true
		for _, v := range g.Doodads {
			isCollision := playerRect.Overlaps(v.Hitbox(0, 0))
			if isCollision {
				move = false
				break
			}
		}
		for _, v := range g.Characters {
			isCollision := playerRect.Overlaps(v.Hitbox(0, 0))
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
		screen.DrawImage(d.Sprite.Image, g.Options)
	}

	for _, c := range g.Characters {
		g.Options.GeoM.Reset()
		// ebiten renders from the min vertex (top left). Offset by the frameheight and half the framewidth to emulate rendering from the "feet" of the sprite
		g.Options.GeoM.Translate(float64(c.X-c.Sprite.FrameWidth/2), float64(c.Y-c.Sprite.FrameHeight))
		// sub-rect is the width of a frame times the frame number, plus the frame number for the 1-pixel buffer between frames
		screen.DrawImage(c.Sprite.Image.SubImage(image.Rect(c.Sprite.FrameWidth*c.FrameNum+c.FrameNum, 0, c.Sprite.FrameWidth*c.FrameNum+c.FrameNum+c.Sprite.FrameWidth, c.Sprite.FrameHeight)).(*ebiten.Image), g.Options)
	}

	g.Options.GeoM.Reset()
	// ebiten renders from the min vertex (top left). Offset by the frameheight and half the framewidth to emulate rendering from the "feet" of the sprite
	g.Options.GeoM.Translate(float64(g.Player.X-g.Player.Sprite.FrameWidth/2), float64(g.Player.Y-g.Player.Sprite.FrameHeight))
	// sub-rect is the width of a frame times the frame number, plus the frame number for the 1-pixel buffer between frames
	screen.DrawImage(g.Player.Sprite.Image.SubImage(image.Rect(g.Player.Sprite.FrameWidth*g.Player.FrameNum+g.Player.FrameNum, 0, g.Player.Sprite.FrameWidth*g.Player.FrameNum+g.Player.FrameNum+g.Player.Sprite.FrameWidth, g.Player.Sprite.FrameHeight)).(*ebiten.Image), g.Options)

	// If in a text interaction, draw the text box last over eveything else.
	if g.InteractionTarget != nil {
		leftWidth := g.Sprites["dialogueFrameLeft"].Image.Bounds().Dx()
		rightWidth := g.Sprites["dialogueFrameRight"].Image.Bounds().Dx()

		g.Options.GeoM.Reset()
		g.Options.GeoM.Scale(39, 1)
		g.Options.GeoM.Translate(float64(leftWidth), 0)
		screen.DrawImage(g.Sprites["dialogueFrameCenter"].Image, g.Options)

		g.Options.GeoM.Reset()
		screen.DrawImage(g.Sprites["dialogueFrameLeft"].Image, g.Options)

		g.Options.GeoM.Translate(float64(320-rightWidth), 0)
		screen.DrawImage(g.Sprites["dialogueFrameRight"].Image, g.Options)

		dialogue := strings.Split(g.InteractionTarget.Dialogue(), " ")
		line := 1
		for i := 0; i < len(dialogue); line++ {
			var render []string
			// Loop over words until they surpass the screen length, then back up by one word
			for w := 0; w < 320-leftWidth-rightWidth && i < len(dialogue); i++ {
				render = append(render, dialogue[i])
				join := strings.Join(render, " ")
				// Calculate the rect size of the string
				bound, _ := font.BoundString(g.Font, join)
				w = (bound.Max.X - bound.Min.X).Ceil()
				// If the rect overflows the screen, go back by one word
				if w >= 320-leftWidth-rightWidth {
					i--
				}
			}
			if i < len(dialogue) {
				render = render[:len(render)-1]
			}
			text.Draw(screen, strings.Join(render, " "), g.Font, 8, line*18, color.White)
		}
		options := g.InteractionTarget.Options()
		if len(options) > 0 {
			var o []string
			for i := 0; i < len(options); i++ {
				o = append(o, options[i][0])
			}
			text.Draw(screen, strings.Join(o, " "), g.Font, 8, line*18, color.White)
			g.Options.GeoM.Reset()
			g.Options.GeoM.Translate(float64(8+18*g.InteractionTarget.SelectedOption()), float64(line*12))
			screen.DrawImage(g.Sprites["selectBox"].Image, g.Options)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("grame")

	sprites, err := os.ReadFile("./sprites/sprites.json")
	if err != nil {
		log.Fatal(err)
	}

	var jsonSprites []SpriteJSON
	err = json.Unmarshal(sprites, &jsonSprites)
	if err != nil {
		log.Fatal(err)
	}

	linkSprites := map[string]Sprite{}
	for _, v := range jsonSprites {
		sprite, _, err := ebitenutil.NewImageFromFile("./sprites/" + v.Image)
		if err != nil {
			log.Fatal(err)
		}

		// camelCase the filenames without the extension to make sprite keys
		k := CamelCase(v.Image[:len(v.Image)-4])

		linkSprites[k] = Sprite{
			FrameLen:    v.FrameLen,
			FrameHeight: v.FrameHeight,
			FrameWidth:  v.FrameWidth,
			Image:       sprite,
		}
	}

	sprite := linkSprites["linkStandSouth"]

	dialogues, err := os.ReadFile("./dialogue/dialogue.json")
	if err != nil {
		log.Fatal(err)
	}

	var jsonDialogues map[string][]DialogueJSON
	err = json.Unmarshal(dialogues, &jsonDialogues)
	if err != nil {
		log.Fatal(err)
	}

	dialogueGraphs := map[string]*DialogueGraph{}
	for k, v := range jsonDialogues {
		graph := DialogueGraph{
			Nodes: map[string]*DialogueNode{},
			Edges: map[string][]string{},
		}
		for i := 0; i < len(v); i++ {
			node := DialogueNode{
				Phrase:    v[i].Phrase,
				Options:   v[i].Options,
				RuneNum:   0,
				OptionNum: 0,
				End:       v[i].End,
			}
			graph.Nodes[v[i].ID] = &node
			graph.Edges[v[i].ID] = v[i].Connections
			if i == 0 {
				graph.RootKey = v[i].ID
				graph.NodeKey = v[i].ID
			}
		}

		dialogueGraphs[k] = &graph
	}

	f, err := opentype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    12,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

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
		Characters: []Character{
			{
				X:              100,
				Y:              50,
				FrameNum:       0,
				Sprite:         linkSprites["elderStandSouth"],
				DialogueGraphs: dialogueGraphs,
				DialogueKey:    "elder",
			},
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
		Font:    face,
		Options: op,
	}

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
