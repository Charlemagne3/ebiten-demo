package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

var AttackCommands = []string{"attack_west", "attack_east", "attack_north", "attack_south"}

type Behavior struct {
	Command string     // The name of the command for the enemy to execute
	Key     ebiten.Key // Any key associated with the command, like a direction
	Pause   int        // The wait time between new commands
	Paused  int        // How long the behavior has been paused since the last action
}

func AdvanceBehavior(g *Game, e *Enemy) {
	if Contains(AttackCommands, e.Behavior.Command) && e.Behavior.Paused < e.Behavior.Pause {
		e.Behavior.Paused++
		return
	} else {
		e.Behavior.Paused = 0
	}
	playerRect := g.Player.Hitbox(0, 0)
	enemyRect := e.Hitbox(0, 0)
	playerX := playerRect.Min.X + playerRect.Dx()/2
	enemyX := enemyRect.Min.X + enemyRect.Dx()/2
	if enemyX == playerX {
		if enemyRect.Max.Y < playerRect.Max.Y {
			e.Behavior.Command = "attack_south"
			e.Sprite = g.Sprites["skeletonWizardAttackSouth"]
			g.Projectiles = append(g.Projectiles, Projectile{
				X:        enemyX,
				Y:        enemyRect.Max.Y,
				Sprite:   g.Sprites["fireballSouth"],
				FrameNum: 0,
				Speed:    3,
				Dir:      ebiten.KeyDown,
			})
		} else if enemyRect.Max.Y > playerRect.Max.Y {
			e.Behavior.Command = "attack_north"
			e.Sprite = g.Sprites["skeletonWizardAttackNorth"]
			g.Projectiles = append(g.Projectiles, Projectile{
				X:        enemyX,
				Y:        enemyRect.Min.Y,
				Sprite:   g.Sprites["fireballNorth"],
				FrameNum: 0,
				Speed:    3,
				Dir:      ebiten.KeyUp,
			})
		}
	} else if enemyRect.Max.Y == playerRect.Max.Y {
		if enemyRect.Max.X < playerRect.Max.X {
			e.Behavior.Command = "attack_east"
			e.Sprite = g.Sprites["skeletonWizardAttackEast"]
			g.Projectiles = append(g.Projectiles, Projectile{
				X:        enemyRect.Max.X,
				Y:        enemyRect.Max.Y,
				Sprite:   g.Sprites["fireballEast"],
				FrameNum: 0,
				Speed:    3,
				Dir:      ebiten.KeyRight,
			})
		} else if enemyRect.Max.X > playerRect.Max.X {
			e.Behavior.Command = "attack_west"
			e.Sprite = g.Sprites["skeletonWizardAttackWest"]
			g.Projectiles = append(g.Projectiles, Projectile{
				X:        enemyRect.Min.X,
				Y:        enemyRect.Max.Y,
				Sprite:   g.Sprites["fireballWest"],
				FrameNum: 0,
				Speed:    3,
				Dir:      ebiten.KeyLeft,
			})
		}
	} else {
		xDiff := AbsDiff(enemyX, playerX)
		yDiff := AbsDiff(enemyRect.Max.Y, playerRect.Max.Y)
		if xDiff < yDiff {
			if enemyX < playerX {
				move := true
				enemyRect := e.Hitbox(1, 0)
				for _, v := range g.Enemies {
					isCollision := *e != v && enemyRect.Overlaps(v.Hitbox(0, 0))
					if isCollision {
						move = false
						break
					}
				}
				if move {
					for _, v := range g.Doodads {
						isCollision := enemyRect.Overlaps(v.Hitbox(0, 0))
						if isCollision {
							move = false
							break
						}
					}
				}
				if move {
					e.Behavior.Command = "walk_east"
					e.Sprite = g.Sprites["skeletonWizardWalkEast"]
					e.X++
				}
			} else {
				move := true
				enemyRect := e.Hitbox(-1, 0)
				for _, v := range g.Enemies {
					isCollision := *e != v && enemyRect.Overlaps(v.Hitbox(0, 0))
					if isCollision {
						move = false
						break
					}
				}
				if move {
					for _, v := range g.Doodads {
						isCollision := enemyRect.Overlaps(v.Hitbox(0, 0))
						if isCollision {
							move = false
							break
						}
					}
				}
				if move {
					e.Behavior.Command = "walk_west"
					e.Sprite = g.Sprites["skeletonWizardWalkWest"]
					e.X--
				}
			}
		} else {
			if enemyRect.Max.Y < playerRect.Max.Y {
				move := true
				enemyRect := e.Hitbox(0, 1)
				for _, v := range g.Enemies {
					isCollision := *e != v && enemyRect.Overlaps(v.Hitbox(0, 0))
					if isCollision {
						move = false
						break
					}
				}
				if move {
					for _, v := range g.Doodads {
						isCollision := enemyRect.Overlaps(v.Hitbox(0, 0))
						if isCollision {
							move = false
							break
						}
					}
				}
				if move {
					e.Behavior.Command = "walk_south"
					e.Sprite = g.Sprites["skeletonWizardWalkSouth"]
					e.Y++
				}
			} else {
				move := true
				enemyRect := e.Hitbox(0, -1)
				for _, v := range g.Enemies {
					isCollision := *e != v && enemyRect.Overlaps(v.Hitbox(0, 0))
					if isCollision {
						move = false
						break
					}
				}
				if move {
					for _, v := range g.Doodads {
						isCollision := enemyRect.Overlaps(v.Hitbox(0, 0))
						if isCollision {
							move = false
							break
						}
					}
				}
				if move {
					e.Behavior.Command = "walk_north"
					e.Sprite = g.Sprites["skeletonWizardWalkNorth"]
					e.Y--
				}
			}
		}
	}
}
