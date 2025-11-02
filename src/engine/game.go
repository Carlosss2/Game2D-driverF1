package engine

import (
	"fmt"
	"game/src/models"
	"game/src/threards"
	"game/src/utils"
	"image/color"
	"log"
	
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil" 
)

type Game struct {
	player  *models.Player
	spawner *threards.Spawner
	bg      *ebiten.Image

	
	gameOver bool
	victory  bool
}

func NewGame() *Game {
	if utils.BackgroundImg == nil || utils.PlayerImg == nil || utils.EnemyImg == nil {
		log.Fatal("assets not initialized - call utils.InitAssets() before NewGame()")
	}
	player := models.NewPlayer(utils.PlayerImg)
	spawner := threards.NewSpawner(utils.EnemyImg)
	return &Game{
		player:  player,
		spawner: spawner,
		bg:      utils.BackgroundImg,
		
		gameOver: false,
		victory:  false,
	}
}

func (g *Game) Update() error {

	dt := 1.0 / float64(ebiten.DefaultTPS)
	



	if inpututil.IsKeyJustPressed(ebiten.KeyA) || inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.player.MoveLeft()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) || inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.player.MoveRight()
	}
	
	g.player.Update(dt)

	// update spawner con concurrency (fan-out / fan-in)
	g.spawner.Update(dt, g.player.Distance)

	// revisar las coliciones de player rect vs any enemy rect
	playerRect := g.player.GetRect()
	for _, e := range g.spawner.AllEnemies() {
		if e.Alive {
			if rectsOverlap(playerRect, e.GetRect()) {
				g.gameOver = true
			}
		}
	}

	// aca pa elegir kilometraje de meta
	km := pixelsToKm(g.player.Distance)
	if km >= 1000 {
		g.victory = true
	}
	if (g.gameOver || g.victory) && ebiten.IsKeyPressed(ebiten.KeyR) {
		g.reset()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(g.bg, op)
	g.player.Draw(screen)
	g.spawner.Draw(screen)
	km := pixelsToKm(g.player.Distance)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("KM: %.2f / 1000.00\nA: izquierda  D: derecha  R: reiniciar", km))
	if g.gameOver {
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


func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Usamos el tamaño lógico base
	return 480, 800
}

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
