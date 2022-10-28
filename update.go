package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func UpdateInteraction(g *Game) {
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
				g.Player.FrameDur = 0
				g.Player.Sprite = g.Sprites["linkStandEast"]
			} else if AbsDiff(playerRect.Max.Y, characterRect.Min.Y) <= 1 && playerRect.Min.X+(playerRect.Dx()/2) >= characterRect.Min.X && playerRect.Min.X+(playerRect.Dx()/2) <= characterRect.Max.X {
				g.InteractionTarget = &c
				g.Player.FrameNum = 0
				g.Player.FrameDur = 0
				g.Player.Sprite = g.Sprites["linkStandSouth"]
			} else if AbsDiff(playerRect.Min.X, characterRect.Max.X) <= 1 && playerRect.Min.Y+(playerRect.Dy()/2) >= characterRect.Min.Y && playerRect.Min.Y+(playerRect.Dy()/2) <= characterRect.Max.Y {
				g.InteractionTarget = &c
				g.Player.FrameNum = 0
				g.Player.FrameDur = 0
				g.Player.Sprite = g.Sprites["linkStandWest"]
			} else if AbsDiff(playerRect.Min.Y, characterRect.Max.Y) <= 1 && playerRect.Min.X+(playerRect.Dx()/2) >= characterRect.Min.X && playerRect.Min.X+(playerRect.Dx()/2) <= characterRect.Max.X {
				g.InteractionTarget = &c
				g.Player.FrameNum = 0
				g.Player.FrameDur = 0
				g.Player.Sprite = g.Sprites["linkStandNorth"]
			}
		}
	}

}

func UpdatePlayer(g *Game) {
	animEnd := false

	if g.Player.Sprite.FrameLen > 1 {
		g.Player.FrameDur++
		// Use >= because the 0 frame counts as one
		if g.Player.FrameDur >= g.Player.Sprite.FrameDur {
			g.Player.FrameDur = 0
			g.Player.FrameNum++
			// FrameNum is zero indexed and FrameLen is a natural number, so subtract 1 for the mod operation
			g.Player.FrameNum = g.Player.FrameNum % (g.Player.Sprite.FrameLen - 1)
			// End the animation if the last render was the last frame
			if g.Player.Animation && g.Player.FrameNum == 0 {
				g.Player.Weapon.IsAttacking = false
				g.Player.Weapon.FrameNum = 0
				g.Player.Animation = false
				animEnd = true
			}
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		// Start the walk left animation if the player just pressed left or if an animation ended and the player was already moving left
		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || (IsOtherDirectionJustReleased(ebiten.KeyLeft) && IsLeastKeyPressDuration(ebiten.KeyLeft)) || animEnd {
			g.Player.LastDir = ebiten.KeyLeft
			g.Player.FrameNum = 0
			g.Player.FrameDur = 0
			g.Player.Sprite = g.Sprites["linkWalkWest"]
		}
		playerRect := g.Player.Hitbox(-1, 0)
		move := true
		for _, v := range g.Enemies {
			isCollision := playerRect.Overlaps(v.Hitbox(0, 0))
			if isCollision {
				move = false
				g.EnemyCollision = &v
				break
			}
		}
		for _, v := range g.Tiles {
			isCollision := playerRect.Overlaps(v.Hitbox(0, 0))
			if isCollision {
				move = false
				break
			}
		}
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
			g.Player.FrameDur = 0
			g.Player.Sprite = g.Sprites["linkWalkEast"]
		}
		playerRect := g.Player.Hitbox(1, 0)
		move := true
		for _, v := range g.Enemies {
			isCollision := playerRect.Overlaps(v.Hitbox(0, 0))
			if isCollision {
				move = false
				g.EnemyCollision = &v
				break
			}
		}
		for _, v := range g.Tiles {
			isCollision := playerRect.Overlaps(v.Hitbox(0, 0))
			if isCollision {
				move = false
				break
			}
		}
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
			g.Player.FrameDur = 0
			g.Player.Sprite = g.Sprites["linkWalkNorth"]
		}
		playerRect := g.Player.Hitbox(0, -1)
		move := true
		for _, v := range g.Enemies {
			isCollision := playerRect.Overlaps(v.Hitbox(0, 0))
			if isCollision {
				move = false
				g.EnemyCollision = &v
				break
			}
		}
		for _, v := range g.Tiles {
			isCollision := playerRect.Overlaps(v.Hitbox(0, 0))
			if isCollision {
				move = false
				break
			}
		}
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
			g.Player.FrameDur = 0
			g.Player.Sprite = g.Sprites["linkWalkSouth"]
		}
		playerRect := g.Player.Hitbox(0, 1)
		move := true
		for _, v := range g.Enemies {
			isCollision := playerRect.Overlaps(v.Hitbox(0, 0))
			if isCollision {
				move = false
				g.EnemyCollision = &v
				break
			}
		}
		for _, v := range g.Tiles {
			isCollision := playerRect.Overlaps(v.Hitbox(0, 0))
			if isCollision {
				move = false
				break
			}
		}
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
		g.Player.FrameDur = 0
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
		g.Player.FrameDur = 0
		g.Player.Weapon.FrameNum = 0
		g.Player.Weapon.Wielder = &g.Player
		g.Player.Weapon.IsAttacking = true
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
}

func UpdateCharacters(g *Game) {

}

func UpdateEnemies(g *Game) {
	for i := 0; i < len(g.Enemies); i++ {
		AdvanceBehavior(g, &g.Enemies[i])
	}
}

func UpdateProjectiles(g *Game) {
	var remove []int
	for i := 0; i < len(g.Projectiles); i++ {
		if g.Projectiles[i].Dir == ebiten.KeyLeft {
			g.Projectiles[i].X -= g.Projectiles[i].Speed
		} else if g.Projectiles[i].Dir == ebiten.KeyRight {
			g.Projectiles[i].X += g.Projectiles[i].Speed
		} else if g.Projectiles[i].Dir == ebiten.KeyUp {
			g.Projectiles[i].Y -= g.Projectiles[i].Speed
		} else if g.Projectiles[i].Dir == ebiten.KeyDown {
			g.Projectiles[i].Y += g.Projectiles[i].Speed
		}
		if g.Projectiles[i].X < -640 || g.Projectiles[i].X > 1280 || g.Projectiles[i].Y < -480 || g.Projectiles[i].Y > 960 {
			remove = append(remove, i)
			continue
		}
		hitbox := g.Projectiles[i].Hitbox(0, 0)
		if g.Projectiles[i].IsEnemy {
			if hitbox.Overlaps(g.Player.Hitbox(0, 0)) {
				g.ProjectileCollision = &g.Projectiles[i]
				continue
			}

			for _, c := range g.Characters {
				if hitbox.Overlaps(c.Hitbox(0, 0)) {
					remove = append(remove, i)
					continue
				}
			}

			for _, d := range g.Doodads {
				if hitbox.Overlaps(d.Hitbox(0, 0)) {
					remove = append(remove, i)
					continue
				}
			}
		} else {
			for _, e := range g.Enemies {
				if hitbox.Overlaps(e.Hitbox(0, 0)) {
					remove = append(remove, i)
					continue
				}
			}

			for _, d := range g.Doodads {
				if hitbox.Overlaps(d.Hitbox(0, 0)) {
					remove = append(remove, i)
					continue
				}
			}
		}
	}
	for k, v := range remove {
		g.Projectiles = Remove(g.Projectiles, v+k)
	}
}

func UpdateDamage(g *Game) {
	if g.ProjectileCollision != nil {
		g.Player.Health--
		var i int
		for i = 0; i < len(g.Projectiles); i++ {
			if g.ProjectileCollision == &g.Projectiles[i] {
				break
			}
		}
		g.ProjectileCollision = nil
		g.EnemyCollision = nil
		g.Projectiles = Remove(g.Projectiles, i)
		return
	}

	if g.EnemyCollision != nil {
		g.Player.Health--
		g.ProjectileCollision = nil
		g.EnemyCollision = nil
		return
	}
}
