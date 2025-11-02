package engine

import (
	"fmt"
	"game/src/models"
	"game/src/threards"
	"game/src/utils"
	"image/color"
	"log"
	"time"
"image"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	player  *models.Player
	spawner *threards.Spawner
	bg      *ebiten.Image

	lastTime time.Time
	gameOver bool
	victory  bool
}

func NewGame() *Game {
	if utils.BackgroundImg == nil || utils.PlayerImg == nil || utils.EnemyImg == nil {
		// ensure assets initialized; InitAssets called in main
		log.Fatal("assets not initialized - call utils.InitAssets() before NewGame()")
	}
	player := models.NewPlayer(utils.PlayerImg)
	spawner := threards.NewSpawner(utils.EnemyImg)
	return &Game{
		player: player,
		spawner: spawner,
		bg: utils.BackgroundImg,
		lastTime: time.Now(),
		gameOver: false,
		victory: false,
	}
}

func (g *Game) Update() error {
	// delta time
	now := time.Now()
	dt := now.Sub(g.lastTime).Seconds()
	if dt > 0.1 {
		dt = 0.1 // clamp
	}
	g.lastTime = now

	// input: move left/right (A/D or arrow keys)
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.MoveLeft()
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.MoveRight()
	}

	// update player distance
	g.player.Update(dt)

	// update spawner with concurrency (fan-out / fan-in)
	g.spawner.Update(dt, g.player.Distance)

	// check collisions: player rect vs any enemy rect
	playerRect := g.player.GetRect()
	for _, e := range g.spawner.AllEnemies() {
		if e.Alive {
			if rectsOverlap(playerRect, e.GetRect()) {
				g.gameOver = true
			}
		}
	}

	// check victory: player.Distance in pixels -> convert to km. We'll define 100 px = 1 km (configurable)
	// To reach 1000 km: need 1000 * 100 px = 100000 px (you can adjust conversion). We'll compute km dynamically.
	km := pixelsToKm(g.player.Distance)
	if km >= 1000 {
		g.victory = true
	}

	// restart with R
	if (g.gameOver || g.victory) && ebiten.IsKeyPressed(ebiten.KeyR) {
		g.reset()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// draw background (stretch to screen)
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(g.bg, op)

	// draw player
	g.player.Draw(screen)

	// draw enemies
	g.spawner.Draw(screen)

	// hud
	km := pixelsToKm(g.player.Distance)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("KM: %.2f / 1000.00\nA: izquierda  D: derecha  R: reiniciar", km))

	if g.gameOver {
		// overlay red
		overlay := ebiten.NewImage(480, 800)
		overlay.Fill(color.RGBA{150, 0, 0, 120})
		screen.DrawImage(overlay, nil)
		ebitenutil.DebugPrintAt(screen, "GAME OVER - Presiona R", 140, 380)
	}
	if g.victory {
		overlay := ebiten.NewImage(480, 800)
		overlay.Fill(color.RGBA{0, 150, 0, 120})
		screen.DrawImage(overlay, nil)
		ebitenutil.DebugPrintAt(screen, "¡GANASTE! 1000 km alcanzados - Presiona R", 60, 380)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) { return 720, 1200 }

func rectsOverlap(a, b image.Rectangle) bool {
	return a.Overlaps(b)
}

func resetGameState(p *models.Player, s *threards.Spawner) {
	p.Lane = 1
	p.Distance = 0
	p.Alive = true
	s.Enemies = []*models.Enemy{}
	s.NextID = 1
	s.LastSpawnDistance = 0
}

func (g *Game) reset() {
	resetGameState(g.player, g.spawner)
	g.gameOver = false
	g.victory = false
}

func pixelsToKm(pixels float64) float64 {
	// conversión arbitraria: 100 pixels = 1 km
	return pixels / 100.0
}
