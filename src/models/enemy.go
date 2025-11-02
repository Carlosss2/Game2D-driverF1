package models

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Enemy struct {
	ID    int
	Lane  int
	X, Y  float64
	Speed float64
	Image *ebiten.Image
	Alive bool
}

func NewEnemy(id int, img *ebiten.Image, lane int, startY float64) *Enemy {
	// Escalamos los autos enemigos también a la mitad del tamaño
	scaled := ebiten.NewImage(img.Bounds().Dx()/2, img.Bounds().Dy()/2)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.5, 0.5)
	scaled.DrawImage(img, op)

	baseSpeed := 120.0 + rand.Float64()*80.0
	return &Enemy{
		ID:    id,
		Lane:  lane,
		X:     LaneX[lane],
		Y:     startY,
		Speed: baseSpeed,
		Image: scaled,
		Alive: true,
	}
}

func (e *Enemy) Update(dt float64) {
	e.Y += e.Speed * dt
	if e.Y > 900 {
		e.Alive = false
	}
}

func (e *Enemy) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	x := e.X - float64(e.Image.Bounds().Dx())/2
	op.GeoM.Translate(x, e.Y)
	screen.DrawImage(e.Image, op)
}


func (e *Enemy) GetRect() image.Rectangle {
	w, h := e.Image.Size()
	x := e.X - float64(w)/2
	y := e.Y

	
	paddingX := float64(w) * 0.42
	paddingY := float64(h) * 0.42

	x1 := int(x + paddingX)
	y1 := int(y + paddingY)
	x2 := int(x + float64(w) - paddingX)
	y2 := int(y + float64(h) - paddingY)

	return image.Rect(x1, y1, x2, y2)
	
}