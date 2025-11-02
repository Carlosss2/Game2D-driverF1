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
	// spawn para que aparezcam los carritos píxels de avance 
	threshold := 900.0
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
	// possibly spawn
	s.spawnIfNeeded(playerDistance)

	if len(s.Enemies) == 0 {
		return
	}

	jobs := []concurrency.Job{}
	// create a job per enemy using snapshot
	for _, en := range s.Enemies {
		snap := *en
		job := func() interface{} {
		
			snap.Update(dt) 

			res := EnemyResult{
				ID:    snap.ID,
				NewY:  snap.Y,     // Valor de la copia actualizada
				Alive: snap.Alive, //Valor de la copia actualizada
			}
	
			return res
		}
		jobs = append(jobs, job)
	}

	// fan-out
	resCh := concurrency.FanOut(jobs)

	// fan-in: collect results and apply on main thread
	newEnemies := []*models.Enemy{}
	for v := range resCh {
		r, ok := v.(EnemyResult)
		if !ok {
			continue
		}
		// find original
		var orig *models.Enemy
		for _, e := range s.Enemies {
			if e.ID == r.ID {
				orig = e
				break
			}
		}
		if orig == nil {
			continue
		}
		orig.Y = r.NewY
		orig.Alive = r.Alive
		if orig.Alive {
			newEnemies = append(newEnemies, orig)
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
