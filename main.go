package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func main() {
	// Ensure graphics are on main thread
	pixelgl.Run(run)
}

func run() {
	// Create a window
	cfg := pixelgl.WindowConfig{
		Title:  "Boids Terrain",
		Bounds: pixel.R(0, 0, 800, 800),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	var settings = DefSimSettings
	if settings.MapGenerationParams.Seed == -1 {
		settings.MapGenerationParams.Seed = time.Now().Unix()
	}
	runSimulation(win, settings, DefUserSettings)
}

func runSimulation(win *pixelgl.Window, simSettings SimulationSettings, userSettings UserSettings) {
	// Setup the camera
	cameraWorldPos := pixel.V(20, 80)
	currentPixelsPerMeter := 50.0

	// Generate a new map
	currentMap := NewGeneratedMap(simSettings.MapGenerationParams)

	// Create the entities
	entities := make([]Entity, 0)
	for i := 0; i < 500; i++ {
		entities = append(entities, NewFish(pixel.V(float64(i)*1+2, 250)))
	}

	// Create the batch so we can draw all entities at once
	entitiesBatch := pixel.NewBatch(&pixel.TrianglesData{}, GetSpritePicture("entities"))

	// Update loop
	for !win.Closed() {
		// Read keypresses and wipe the window
		win.Update()
		win.Clear(colornames.Black)

		// Process player input to move the camera around
		spd := userSettings.CameraSettings.MoveSpeed / currentPixelsPerMeter
		if win.Pressed(pixelgl.KeyW) {
			cameraWorldPos = cameraWorldPos.Add(pixel.V(0, spd/60.0))
		} else if win.Pressed(pixelgl.KeyS) {
			cameraWorldPos = cameraWorldPos.Add(pixel.V(0, -spd/60.0))
		}
		if win.Pressed(pixelgl.KeyA) {
			cameraWorldPos = cameraWorldPos.Add(pixel.V(-spd/60.0, 0))
		} else if win.Pressed(pixelgl.KeyD) {
			cameraWorldPos = cameraWorldPos.Add(pixel.V(spd/60.0, 0))
		}
		scaleSpd := 1.0
		if win.Pressed(pixelgl.KeyQ) {
			scaleSpd = 1.0 / math.Pow(userSettings.CameraSettings.ZoomSpeed, 1.0/60)
		} else if win.Pressed(pixelgl.KeyE) {
			scaleSpd = math.Pow(userSettings.CameraSettings.ZoomSpeed, 1.0/60)
		}
		currentPixelsPerMeter *= scaleSpd

		if win.JustPressed(pixelgl.KeyB) {
			for _, e := range entities {
				if !e.IsKinematic() {
					e.ApplyImpulse(pixel.V(50, 0).Rotated(rand.Float64() * 3.14159 * 2))
				}
			}
		}

		// Update logic for entities
		for _, e := range entities {
			e.StepLogic()
		}

		// Update forces and integrate kinematics
		for _, e := range entities {
			if !e.IsKinematic() {
				e.StepPhysics()
			}
		}

		// Process collisions and ensure the solver ends in a valid state
		for i, e := range entities {
			CollideMapEntity(currentMap, e)
			for j, e2 := range entities {
				if j <= i {
					continue
				}
				CollideEntityEntity(e, e2)
			}
		}

		// Create the render data to draw map with
		renderData := &RenderData{
			Target:         win,
			TargetRect:     win.Bounds(),
			CameraWorldPos: cameraWorldPos,
			PixelsPerMeter: currentPixelsPerMeter,
		}

		// Render the current map
		currentMap.Render(renderData)

		// Create the render data to draw entities with
		renderDataEntities := &RenderData{
			Target:         entitiesBatch,
			TargetRect:     win.Bounds(),
			CameraWorldPos: cameraWorldPos,
			PixelsPerMeter: currentPixelsPerMeter,
		}

		// Render all of the entities
		entitiesBatch.Clear()
		for _, e := range entities {
			e.Render(renderDataEntities)
		}
		entitiesBatch.Draw(win)
	}
}
