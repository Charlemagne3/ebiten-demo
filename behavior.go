package main

import "github.com/hajimehoshi/ebiten/v2"

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
	if enemyRect.Min.X == playerRect.Min.X {
		if enemyRect.Min.Y < playerRect.Min.Y {
			e.Behavior.Command = "attack_south"
			e.Sprite = g.Sprites["skeletonWizardAttackSouth"]
			g.Projectiles = append(g.Projectiles, Projectile{
				X:        e.X,
				Y:        e.Y,
				Sprite:   g.Sprites["fireballSouth"],
				FrameNum: 0,
				Dir:      ebiten.KeyDown,
			})
		} else if enemyRect.Min.Y > playerRect.Min.Y {
			e.Behavior.Command = "attack_north"
			e.Sprite = g.Sprites["skeletonWizardAttackNorth"]
			g.Projectiles = append(g.Projectiles, Projectile{
				X:        e.X,
				Y:        e.Y,
				Sprite:   g.Sprites["fireballNorth"],
				FrameNum: 0,
				Dir:      ebiten.KeyUp,
			})
		}
	} else if enemyRect.Min.Y == playerRect.Min.Y {
		if enemyRect.Min.X < playerRect.Min.X {
			e.Behavior.Command = "attack_east"
			e.Sprite = g.Sprites["skeletonWizardAttackEast"]
			g.Projectiles = append(g.Projectiles, Projectile{
				X:        e.X,
				Y:        e.Y,
				Sprite:   g.Sprites["fireballEast"],
				FrameNum: 0,
				Dir:      ebiten.KeyRight,
			})
		} else if enemyRect.Min.X > playerRect.Min.X {
			e.Behavior.Command = "attack_west"
			e.Sprite = g.Sprites["skeletonWizardAttackWest"]
			g.Projectiles = append(g.Projectiles, Projectile{
				X:        e.X,
				Y:        e.Y,
				Sprite:   g.Sprites["fireballWest"],
				FrameNum: 0,
				Dir:      ebiten.KeyLeft,
			})
		}
	} else {
		xDiff := AbsDiff(enemyRect.Min.X, playerRect.Min.X)
		yDiff := AbsDiff(enemyRect.Min.Y, playerRect.Min.Y)
		if xDiff < yDiff {
			if enemyRect.Min.X < playerRect.Min.X {
				e.Behavior.Command = "walk_east"
				e.Sprite = g.Sprites["skeletonWizardWalkEast"]
				e.X++
			} else {
				e.Behavior.Command = "walk_west"
				e.Sprite = g.Sprites["skeletonWizardWalkWest"]
				e.X--
			}
		} else {
			if enemyRect.Min.Y < playerRect.Min.Y {
				e.Behavior.Command = "walk_south"
				e.Sprite = g.Sprites["skeletonWizardWalkSouth"]
				e.Y++
			} else {
				e.Behavior.Command = "walk_north"
				e.Sprite = g.Sprites["skeletonWizardWalkNorth"]
				e.Y--
			}
		}
	}
}
