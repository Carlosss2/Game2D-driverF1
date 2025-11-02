package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"game/src/engine"
	"game/src/utils"
)

func main() {
	// Carga assets globales
	utils.InitAssets()

	g := engine.NewGame()

	ebiten.SetWindowSize(960, 1200)
	ebiten.SetWindowTitle("Three-Lane Racer")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
