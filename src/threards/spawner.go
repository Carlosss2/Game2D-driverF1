package threards

import (
	"game/src/concurrency"
	"game/src/models"
	"math/rand"
	"github.com/hajimehoshi/ebiten/v2"
)

// EnemyResult contiene info devuelta por cada job
type EnemyResult struct {
	ID    int
	NewY  float64
	Alive bool
}

type Spawner struct {
	Enemies  []*models.Enemy
	EnemyImg *ebiten.Image
	NextID   int
	LastSpawnDistance float64 // distance at which last spawn occurred
}

func NewSpawner(img *ebiten.Image) *Spawner {
	return &Spawner{
		Enemies: []*models.Enemy{},
		EnemyImg: img,
		NextID: 1,
		LastSpawnDistance: 0,
	}
}

func (s *Spawner) spawnIfNeeded(playerDistance float64) {
	// spawn para que aparezcam los carritos pixels de avance 
	threshold := 1000.0
	if playerDistance-s.LastSpawnDistance < threshold {
		return
	}
	lane := rand.Intn(2) 
	startY := -400.0 - float64(rand.Intn(200)) 
	e := models.NewEnemy(s.NextID, s.EnemyImg, lane, startY) 
	s.NextID++
	s.Enemies = append(s.Enemies, e)
	s.LastSpawnDistance = playerDistance
}

// Update usa Fan-Out/Fan-In para procesar cada enemigo en paralelo.
// 
func (s *Spawner) Update(dt float64, playerDistance float64) {
	s.spawnIfNeeded(playerDistance)

	if len(s.Enemies) == 0 {
		return
	}

	// separar los enemigos en grupos
	var group1Jobs, group2Jobs []concurrency.Job
	for i, en := range s.Enemies {
		snap := *en
		job := func() interface{} {
			snap.Update(dt)
			return EnemyResult{
				ID:    snap.ID,
				NewY:  snap.Y,
				Alive: snap.Alive,
			}
		}
		if i%2 == 0 {
			group1Jobs = append(group1Jobs, job)
		} else {
			group2Jobs = append(group2Jobs, job)
		}
	}

	// fan-out por grupo
	resCh1 := concurrency.FanOut(group1Jobs)
	resCh2 := concurrency.FanOut(group2Jobs)

	// fan-in combina ambos canales
	resCh := concurrency.FanIn(resCh1, resCh2)

	// fan-in: colectar resultados
	newEnemies := []*models.Enemy{}
	for v := range resCh {
		r, ok := v.(EnemyResult)
		if !ok {
			continue
		}
		for _, e := range s.Enemies {
			if e.ID == r.ID {
				e.Y = r.NewY
				e.Alive = r.Alive
				if e.Alive {
					newEnemies = append(newEnemies, e)
				}
				
			}
		}
	}
	s.Enemies = newEnemies
}

// Draw dibuja todos los enemigos
func (s *Spawner) Draw(screen *ebiten.Image) {
	for _, e := range s.Enemies {
		e.Draw(screen)
	}
}

// AllEnemies retorna slice (para colisiones)
func (s *Spawner) AllEnemies() []*models.Enemy {
	return s.Enemies
}
