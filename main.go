package main

import (
	"encoding/json"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"
	"sort"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

// Game is an ebiten Game interface implemetation plus custom struct data
type Game struct {
	Player              Player
	Characters          []Character
	Enemies             []Enemy
	Weapons             []Weapon
	Projectiles         []Projectile
	Doodads             []Doodad
	Tiles               []Tile
	RenderTargets       []*RenderTarget
	Sprites             map[string]Sprite
	Font                font.Face
	Options             *ebiten.DrawImageOptions
	InteractionTarget   InteractionTarget // The target of another game element that the player is having a dialogue interaction with, or nil.
	EnemyCollision      *Enemy
	ProjectileCollision *Projectile
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

// Player represents the player character
type Player struct {
	X         int        // The current X screen offset of the player
	Y         int        // The current Y screen offset of the player
	Animation bool       // Whether or not the player is in a special animation or the normal stand/walk cycle.
	LastDir   ebiten.Key // The last direction the player faced (never -1)
	Sprite    Sprite     // The current sprite for the player
	FrameNum  int        // The current frame of the sprite for the player
	FrameDur  int        // The duration of the current frame of the sprite for the player
	Health    int        // How much health the player has
	Weapon    *Weapon    // The weapon the player has equipped
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
	Behavior  Behavior   // The active behavior of the enemy
}

// Weapon represents a weapon held by something
type Weapon struct {
	Sprite      Sprite       // The current sprite for the weapon
	FrameNum    int          // The current frame of the sprite for the weapon
	Wielder     RenderTarget // The wielder of the weapon. The weapon is drawn relative to the wielder.
	IsAttacking bool         // Whether or not the weapon is attacking and should be drawn.
}

// Projectile represents a projectile
type Projectile struct {
	X        int        // The current X screen offset of the projectile
	Y        int        // The current Y screen offset of the projectile
	Sprite   Sprite     // The current sprite for the projectile
	FrameNum int        // The current frame of the sprite for the projectile
	Speed    int        // The number of pixels the projectile moves per frame
	Dir      ebiten.Key // The direction the projection is travelling
	IsEnemy  bool       // Whether or not the projecile is enemy or friendly
}

// Doodad represents a static environmental item
type Doodad struct {
	X        int    // The current X screen offset of the doodad
	Y        int    // The current Y screen offset of the doodad
	Sprite   Sprite // The current sprite for the doodad
	FrameNum int    // The current frame of the sprite for the doodad
}

// Tile represents a floor texture
type Tile struct {
	X        int    // The current X screen offset of the doodad
	Y        int    // The current Y screen offset of the doodad
	Sprite   Sprite // The current sprite for the doodad
	FrameNum int    // The current frame of the sprite for the doodad
	Collider bool   // Whether or not the tile can be collided with
}

// Sprite represents an image with a number of sub-frames in it to be rendered via rectangles
type Sprite struct {
	FrameWidth  int
	FrameHeight int
	FrameLen    int           // How many frames are in the sprite
	FrameDur    int           // How many render frames to display a single frame of the sprite
	Handles     []image.Point // An array of coordinates inside the sprite for attaching other sprites to
	Image       *ebiten.Image
}

// SpriteJSON represents the json to be read from the sprite json file.
type SpriteJSON struct {
	FrameDur    int           `json:"frameDuration"`
	FrameLen    int           `json:"frameLen"`
	FrameHeight int           `json:"frameHeight"`
	FrameWidth  int           `json:"frameWidth"`
	Handles     []image.Point `json:"handles"`
	Image       string        `json:"image"`
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
// This is used to reset the walk cycle animation for a new direction
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

func (g *Game) Update() error {

	UpdateInteraction(g)
	if g.InteractionTarget != nil {
		return nil
	}

	UpdatePlayer(g)
	UpdateCharacters(g)
	UpdateEnemies(g)
	UpdateProjectiles(g)
	UpdateDamage(g)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	render := []RenderTarget{&g.Player}
	if g.Player.Weapon.IsAttacking {
		render = append(render, g.Player.Weapon)
	}

	for i := range g.Characters {
		render = append(render, &g.Characters[i])
	}

	for i := range g.Enemies {
		render = append(render, &g.Enemies[i])
	}

	for i := range g.Doodads {
		render = append(render, &g.Doodads[i])
	}

	for i := range g.Tiles {
		render = append(render, &g.Tiles[i])
	}

	for i := range g.Projectiles {
		render = append(render, &g.Projectiles[i])
	}

	sort.Slice(render, func(i, j int) bool { return render[i].RenderOrder() < render[j].RenderOrder() })

	for _, t := range render {
		screen.DrawImage(t.RenderImage(), t.RenderOptions())
	}

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
			FrameDur:    v.FrameDur,
			FrameLen:    v.FrameLen,
			FrameHeight: v.FrameHeight,
			FrameWidth:  v.FrameWidth,
			Handles:     v.Handles,
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

	var tiles []Tile

	for x := 8; x < 320; x += 16 {
		for y := 16; y < 320; y += 16 {
			tiles = append(tiles, Tile{
				X:        x,
				Y:        y,
				FrameNum: 0,
				Sprite:   linkSprites["grass"],
			})
		}
	}

	tiles = append(tiles, Tile{
		X:        64,
		Y:        128,
		FrameNum: 0,
		Sprite:   linkSprites["stump"],
		Collider: true,
	})

	tiles = append(tiles, Tile{
		X:        96,
		Y:        160,
		FrameNum: 0,
		Sprite:   linkSprites["stump"],
		Collider: true,
	})

	weapons := []Weapon{
		{
			Sprite: linkSprites["swordEast"],
		},
	}

	game := &Game{
		Player: Player{
			X:         8,
			Y:         21,
			LastDir:   ebiten.KeyDown,
			Animation: false,
			FrameNum:  0,
			Sprite:    sprite,
			Health:    100,
			Weapon:    &weapons[0],
		},
		Characters: []Character{
			{
				X:              32,
				Y:              32,
				FrameNum:       0,
				Sprite:         linkSprites["elderStandSouth"],
				DialogueGraphs: dialogueGraphs,
				DialogueKey:    "elder",
			},
		},
		Enemies: []Enemy{
			{
				X:        256,
				Y:        128,
				FrameNum: 0,
				Sprite:   linkSprites["skeletonWizardStandSouth"],
				Behavior: Behavior{
					Pause: 60,
				},
			},
		},
		Doodads: []Doodad{
			{
				X:        128,
				Y:        128,
				FrameNum: 0,
				Sprite:   linkSprites["tree"],
			},
			{
				X:        256,
				Y:        256,
				FrameNum: 0,
				Sprite:   linkSprites["tree"],
			},
		},
		Tiles:   tiles,
		Weapons: weapons,
		Sprites: linkSprites,
		Font:    face,
		Options: op,
	}

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
