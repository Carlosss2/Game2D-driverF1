package models

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// --- INICIO DE LA MODIFICACIÓN (2 Carriles) ---
// Ahora hay dos carriles: izquierda (0), derecha (1)
// Centrados en una pantalla de 480px de ancho
var LaneX = []float64{160, 320}
// --- FIN DE LA MODIFICACIÓN ---


type Player struct {
// ... (struct no cambia) ...
	Lane     int
	Y        float64
	Image    *ebiten.Image
	Alive    bool
	Distance float64
	Speed    float64
}

func NewPlayer(img *ebiten.Image) *Player {
	// ... (escalado de imagen no cambia) ...
	scaled := ebiten.NewImage(img.Bounds().Dx()/2, img.Bounds().Dy()/2)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.5, 0.5)
	scaled.DrawImage(img, op)

	return &Player{
		// --- INICIO DE LA MODIFICACIÓN ---
		Lane:     0, // Empieza en el carril izquierdo
		// --- FIN DE LA MODIFICACIÓN ---
		Y:        400,
		Image:    scaled,
		Alive:    true,
		Distance: 0,
		Speed:    180,
	}
}

// --- INICIO DE LA MODIFICACIÓN (Lógica de Movimiento) ---
func (p *Player) MoveLeft() {
	p.Lane = 0 // Ir al carril izquierdo
}

func (p *Player) MoveRight() {
	p.Lane = 1 // Ir al carril derecho
}
// --- FIN DE LA MODIFICACIÓN ---

func (p *Player) Update(dt float64) {
// ... (no cambia) ...
	p.Distance += p.Speed * dt
}

func (p *Player) Draw(screen *ebiten.Image) {
// ... (no cambia) ...
	op := &ebiten.DrawImageOptions{}
	x := LaneX[p.Lane] - float64(p.Image.Bounds().Dx())/2
	op.GeoM.Translate(x, p.Y)
	screen.DrawImage(p.Image, op)
}

func (p *Player) GetRect() image.Rectangle {
	w, h := p.Image.Size()
	x := LaneX[p.Lane] - float64(w)/2
	y := p.Y

	// --- INICIO DE LA MODIFICACIÓN (Hitbox) ---
	// Encogemos el hitbox del jugador en un 30% (15% de cada lado).
	// El auto del jugador (McLaren) parece mejor recortado que el enemigo.
	paddingX := float64(w) * 0.20
	paddingY := float64(h) * 0.20

	x1 := int(x + paddingX)
	y1 := int(y + paddingY)
	x2 := int(x + float64(w) - paddingX)
	y2 := int(y + float64(h) - paddingY)

	return image.Rect(x1, y1, x2, y2)
	// --- FIN DE LA MODIFICACIÓN ---
}