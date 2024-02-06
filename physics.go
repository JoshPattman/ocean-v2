package main

import (
	"math"

	"github.com/gopxl/pixel"
)

// CollideEntityEntity collides two entities together and moves them to new correct positions
func CollideEntityEntity(e1, e2 Entity) {
	delta := e1.Position().Sub(e2.Position())
	dist := delta.Len()
	overlap := (e1.Radius() + e2.Radius()) - dist
	if overlap > 0 {
		correction := delta.Scaled(overlap / dist / 2)
		e1.SlideToPosition(e1.Position().Add(correction))
		e2.SlideToPosition(e2.Position().Sub(correction))
	}
}

// CollideMapEntity moves an entity to a new valid position after colliding it with a map
func CollideMapEntity(m *Map, e Entity) {
	texelRadius := int(math.Round(e.Radius() + 0.5))
	texelPosX := int(math.Round(e.Position().X))
	texelPosY := int(math.Round(e.Position().Y))
	totalCorrectedPos := pixel.ZV
	numDetectedCollisions := 0
	for tx := texelPosX - texelRadius; tx <= texelPosX+texelRadius; tx++ {
		for ty := texelPosY - texelRadius; ty <= texelPosY+texelRadius; ty++ {
			if tx >= 0 && ty >= 0 && tx < len(m.texels) && ty < len(m.texels[tx]) && m.texels[tx][ty] != WaterTexel {
				// Incorrect but good enough for now - approximate all squares to be circles
				delta := e.Position().Sub(pixel.V(float64(tx), float64(ty)))
				dist := delta.Len()
				overlap := (e.Radius() + 0.5) - dist
				if overlap > 0 {
					totalCorrectedPos = totalCorrectedPos.Add(e.Position().Add(delta.Scaled(overlap / dist)))
					numDetectedCollisions++
				}
			}
		}
	}
	if numDetectedCollisions > 0 {
		newPos := totalCorrectedPos.Scaled(1.0 / float64(numDetectedCollisions))
		e.SlideToPosition(newPos)
	}
}

func DragForce(vel pixel.Vec, coeff float64) pixel.Vec {
	return pixel.V(math.Abs(vel.X)*vel.X, math.Abs(vel.Y)*vel.Y).Scaled(-coeff)
}
