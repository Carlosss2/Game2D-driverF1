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

// spawnIfNeeded genera enemigos en los carriles conforme la distancia aumenta.
// playerDistance es la distancia del jugador (en píxels acumulados).
func (s *Spawner) spawnIfNeeded(playerDistance float64) {
	// spawn cada 400 píxels de avance (ajustable)
	threshold := 400.0
	if playerDistance - s.LastSpawnDistance < threshold {
		return
	}
	// generamos 1-2 nuevos autos en carriles aleatorios delante del jugador
	n := 1 + rand.Intn(2)
	startY := -400.0 - float64(rand.Intn(200)) // mucho más arriba // fuera de pantalla por arriba
	for i := 0; i < n; i++ {
		lane := rand.Intn(3) // 0..2
		e := models.NewEnemy(s.NextID, s.EnemyImg, lane, startY - float64(rand.Intn(200)))
		s.NextID++
		s.Enemies = append(s.Enemies, e)
	}
	s.LastSpawnDistance = playerDistance
}

// Update usa Fan-Out/Fan-In para procesar cada enemigo en paralelo.
// dt = delta seconds, playerDistance used to spawn
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
			// --- INICIO DE LA MODIFICACIÓN ---
			// En lugar de recalcular la lógica aquí,
			// llamamos al método Update EN LA SNAPSHOT.
			// Esto es seguro porque 'snap' es una copia.
			//
			// models.Enemy.Update() moverá 'snap.Y' Y
			// pondrá 'snap.Alive = false' si 'snap.Y > 900'.
			snap.Update(dt) 

			// Ahora la snapshot 'snap' tiene los valores actualizados
			res := EnemyResult{
				ID:    snap.ID,
				NewY:  snap.Y,     // <-- Valor de la copia actualizada
				Alive: snap.Alive, // <-- Valor de la copia actualizada
			}
			// --- FIN DE LA MODIFICACIÓN ---
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
