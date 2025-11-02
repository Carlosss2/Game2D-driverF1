package models

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Ajustamos los carriles para que estén bien centrados dentro de 480 px de ancho
// Ahora hay tres carriles: izquierda, centro, derecha
var LaneX = []float64{200, 360, 520} // posiciones horizontales de cada carril

type Player struct {
	Lane     int
	Y        float64
	Image    *ebiten.Image
	Alive    bool
	Distance float64
	Speed    float64
}

func NewPlayer(img *ebiten.Image) *Player {
	// Redimensionamos el auto para hacerlo más pequeño
	scaled := ebiten.NewImage(img.Bounds().Dx()/2, img.Bounds().Dy()/2)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.5, 0.5)
	scaled.DrawImage(img, op)

	return &Player{
		Lane:     1,        // empieza en el carril central
		Y:        650,      // más arriba, para que se vea completo
		Image:    scaled,
		Alive:    true,
		Distance: 0,
		Speed:    180,      // velocidad base
	}
}

func (p *Player) MoveLeft() {
	if p.Lane > 0 {
		p.Lane--
	}
}

func (p *Player) MoveRight() {
	if p.Lane < 2 {
		p.Lane++
	}
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

	// --- INICIO DE LA MODIFICACIÓN (Más Agresiva) ---

	// Encogemos el hitbox del jugador en un 20% (10% de cada lado).
	// Ajusta este valor (0.10) si es necesario.
	paddingX := float64(w) * 0.20
	paddingY := float64(h) * 0.20

	x1 := int(x + paddingX)
	y1 := int(y + paddingY)
	x2 := int(x + float64(w) - paddingX)
	y2 := int(y + float64(h) - paddingY)

	return image.Rect(x1, y1, x2, y2)
	// --- FIN DE LA MODIFICACIÓN ---
}