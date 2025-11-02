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

	// --- INICIO DE LA MODIFICACIÓN ---
	// Vamos a usar el tamaño lógico del juego para la ventana.
	ebiten.SetWindowSize(480, 800)
	// --- FIN DE LA MODIFICACIÓN ---

	ebiten.SetWindowTitle("Three-Lane Racer")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}