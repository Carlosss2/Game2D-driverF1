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

	
	
	ebiten.SetWindowSize(480, 800)
	

	ebiten.SetWindowTitle("F1 Mclaren 2D")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}