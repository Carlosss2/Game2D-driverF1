package models

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

var LaneX = []float64{160, 320}



type Player struct {

	Lane     int
	Y        float64
	Image    *ebiten.Image
	Alive    bool
	Distance float64
	Speed    float64
}

func NewPlayer(img *ebiten.Image) *Player {
	
	scaled := ebiten.NewImage(img.Bounds().Dx()/2, img.Bounds().Dy()/2)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.5, 0.5)
	scaled.DrawImage(img, op)

	return &Player{
		
		Lane:     0, 
		
		Y:        400,
		Image:    scaled,
		Alive:    true,
		Distance: 0,
		Speed:    180,
	}
}


func (p *Player) MoveLeft() {
	p.Lane = 0 
}

func (p *Player) MoveRight() {
	p.Lane = 1 
}


func (p *Player) Update(dt float64) {

	p.Distance += p.Speed * dt
}

func (p *Player) Draw(screen *ebiten.Image) {

	op := &ebiten.DrawImageOptions{}
	x := LaneX[p.Lane] - float64(p.Image.Bounds().Dx())/2
	op.GeoM.Translate(x, p.Y)
	screen.DrawImage(p.Image, op)
}

func (p *Player) GetRect() image.Rectangle {
	w, h := p.Image.Size()
	x := LaneX[p.Lane] - float64(w)/2
	y := p.Y

	
	paddingX := float64(w) * 0.20
	paddingY := float64(h) * 0.20

	x1 := int(x + paddingX)
	y1 := int(y + paddingY)
	x2 := int(x + float64(w) - paddingX)
	y2 := int(y + float64(h) - paddingY)

	return image.Rect(x1, y1, x2, y2)
	
}