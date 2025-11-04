F1 Mclaren 2D: Juego Concurrente en Go

Este es un videojuego de carreras 2D, "F1 Mclaren 2D", desarrollado en el lenguaje de programación Go y utilizando la biblioteca gráfica EbitenEngine.
Más allá de ser un juego, este proyecto es una demostración práctica de la aplicación de principios de concurrencia avanzada para gestionar tareas de alto rendimiento en tiempo real. Se utiliza el patrón Fan-out / Fan-in para distribuir eficientemente el cálculo de la lógica y física de los autos enemigos a través de múltiples goroutines, asegurando un rendimiento fluido y escalable.



Características

Juego de carreras 2D de 2 carriles, de esquivar infinito.

Controles simples de un solo toque (Izquierda/Derecha).

Sistema de puntuación basado en la distancia recorrida (KM).

Generación procedural de enemigos que asegura siempre un carril libre.

Núcleo de actualización concurrente: La lógica de los enemigos se procesa en paralelo.

Seguridad de concurrencia: Verificado con go run -race para garantizar la ausencia de condiciones de carrera.

Arquitectura y Patrón de Concurrencia

El desafío técnico principal de este proyecto es la gestión de estado de múltiples entidades (enemigos) sin bloquear el hilo principal del juego.

El Problema

En una implementación secuencial, si hay 100 enemigos en pantalla, el hilo principal (engine.Update) tendría que iterar sobre los 100 en un solo bucle for. Este cálculo puede tomar más de 16.6ms (el tiempo de un frame a 60 FPS), causando caídas de rendimiento (lag) y una mala experiencia.

La Solución: Patrón Fan-out / Fan-in

Para resolver esto, el paquete threards/spawner.go implementa el patrón Fan-out / Fan-in para la actualización de los enemigos.

El flujo de datos en cada frame es el siguiente:

División de Tareas (en spawner.Update): La lista de enemigos se divide en dos grupos de trabajo (group1Jobs y group2Jobs).

Seguridad (Snapshot): Al crear cada "job" (tarea), se crea una copia (snapshot) del enemigo (snap := *en). Esto es VITAL. La goroutine (hilo secundario) solo modificará esta copia, no el dato original.

Fan-out (Distribuir): Se llama a concurrency.FanOut para cada grupo. Esta función lanza una goroutine por cada enemigo, que ejecuta snap.Update(dt) en paralelo. Los resultados se envían a dos canales (resCh1, resCh2).

Fan-in (Combinar): Se llama a concurrency.FanIn para tomar los dos canales de resultados y combinarlos en un solo canal maestro (resCh).

Recolección Segura (en spawner.Update): El hilo principal (aún en spawner.Update) itera sobre el canal resCh. A medida que recibe cada resultado, actualiza el estado del enemigo original (e.Y = r.NewY).

Esta arquitectura garantiza que:

El cálculo pesado (la física de N enemigos) se realiza en paralelo.

El hilo principal (engine.Update) queda libre y nunca se bloquea.

No hay condiciones de carrera, ya que las goroutines solo leen de una copia y solo el hilo principal escribe en la lista original.

Prerrequisitos

Go (versión 1.20 o superior).

Las librerías de EbitenEngine (se instalan automáticamente con go mod tidy).

Instalación y Ejecución

Clonar el repositorio:

git clone https://github.com/Carlosss2/Game2D-driverF1
cd F1-Mclaren-2D


Instalar dependencias:

go mod tidy


Ejecutar el juego:

go run .


Verificar la seguridad de concurrencia (Requerido para la evaluación):
Este comando ejecuta el juego con el detector de condiciones de carrera de Go activado.

go run -race .


Cómo Jugar

Flecha Izquierda / Tecla A: Mover al carril izquierdo.

Flecha Derecha / Tecla D: Mover al carril derecho.

Tecla R: Reiniciar el juego (después de 'Game Over' o al ganar).

Objetivo: Sobrevivir el mayor tiempo posible y alcanzar la meta de 100 KM.

Autor

Carlos Gael Castro Trujillo
