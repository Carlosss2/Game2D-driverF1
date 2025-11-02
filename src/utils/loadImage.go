package utils

import (
	"image"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	BackgroundImg *ebiten.Image
	PlayerImg     *ebiten.Image
	EnemyImg      *ebiten.Image
)

func LoadImage(path string) (*ebiten.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}

func InitAssets() {
	var err error
	BackgroundImg, err = LoadImage("assets/background.png")
	if err != nil {
		log.Fatal("failed to load background.png:", err)
	}
	PlayerImg, err = LoadImage("assets/player_car.png")
	if err != nil {
		log.Fatal("failed to load player_car.png:", err)
	}
	EnemyImg, err = LoadImage("assets/enemy_car.png")
	if err != nil {
		log.Fatal("failed to load enemy_car.png:", err)
	}
}
