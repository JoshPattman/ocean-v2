package main

import (
	"image/color"
	"math/rand"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
)

// DummyEntity is a simple entity for testing that uses gravity to test collisions
type DummyEntity struct {
	EntityBase
	imd *imdraw.IMDraw
	col color.Color
}

// NewDummyEntity creates a new dummy entity with position and radius
func NewDummyEntity(pos pixel.Vec, radius float64) *DummyEntity {
	return &DummyEntity{
		*NewEntityBase(pos, 1, radius),
		imdraw.New(nil),
		pixel.RGB(rand.Float64(), rand.Float64(), rand.Float64()),
	}
}

// Render draws the entity to the screen in the correct place
func (e *DummyEntity) Render(rd *RenderData) {
	e.imd.Clear()
	e.imd.Color = e.col
	e.imd.Push(e.Position().Sub(rd.CameraWorldPos).Scaled(rd.PixelsPerMeter).Add(rd.TargetRect.Center()))
	e.imd.Circle(e.radius*rd.PixelsPerMeter, 0)
	e.imd.Draw(rd.Target)
}

// StepLogic adds forces and stuff to be integrated later, in this case adds gravity
func (e *DummyEntity) StepLogic() {
	e.ApplyForce(pixel.V(0, -9.81*e.Mass()))
}

func (e *DummyEntity) Tags() []string { return []string{} }
