package main

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/gopxl/pixel"
)

type FishEntity struct {
	EntityBase
	anim        *Animator
	nextDir     pixel.Vec
	lastDirTime time.Time
	col         color.Color
}

func NewFish(pos pixel.Vec) *FishEntity {
	pic := GetSpritePicture("entities")
	anim := NewAnimator(pic, 32,
		map[string]pixel.Vec{
			"swimleft.1":  pixel.V(0+6, 2),
			"swimleft.2":  pixel.V(1+6, 2),
			"swimleft.3":  pixel.V(2+6, 2),
			"swimright.1": pixel.V(0+6, 1),
			"swimright.2": pixel.V(1+6, 1),
			"swimright.3": pixel.V(2+6, 1),
		},
		map[string][]string{
			"swimleft":  {"swimleft.1", "swimleft.2", "swimleft.3", "swimleft.2"},
			"swimright": {"swimright.1", "swimright.2", "swimright.3", "swimright.2"},
		},
		map[string]float64{
			"swimleft":  0.5,
			"swimright": 0.5,
		},
	)
	anim.Play("swimleft")
	return &FishEntity{
		*NewEntityBase(pos, 1, 0.5),
		anim,
		pixel.Unit(rand.Float64() * 3.14 * 2),
		time.Now(),
		pixel.RGB(rand.Float64(), rand.Float64(), rand.Float64()),
	}
}

func (e *FishEntity) Render(rd *RenderData) {
	s := e.anim.CurrentSprite()
	tmat := pixel.IM.Scaled(pixel.ZV, e.Radius()*2/s.Frame().W())
	rotAngle := e.Velocity().Angle()
	if e.Velocity().X < 0 {
		rotAngle += math.Pi
	}
	tmat = tmat.Rotated(pixel.ZV, rotAngle)
	tmat = tmat.Moved(e.Position().Sub(rd.CameraWorldPos))
	tmat = tmat.Scaled(pixel.ZV, rd.PixelsPerMeter)
	tmat = tmat.Moved(rd.TargetRect.Bounds().Center())
	s.DrawColorMask(rd.Target, tmat, e.col)
}

// StepLogic adds forces and stuff to be integrated later, in this case adds gravity
func (e *FishEntity) StepLogic() {
	if time.Since(e.lastDirTime).Seconds() > 5 {
		e.lastDirTime = time.Now()
		e.nextDir = pixel.Unit(rand.Float64() * 3.14 * 2)
	}
	if e.Velocity().X < 0 {
		e.anim.PlayIfNot("swimleft")
	} else {
		e.anim.PlayIfNot("swimright")
	}
	e.anim.Step(1.0 / 60)
	// Move towards target
	e.ApplyForce(e.nextDir.Scaled(5))
	// Drag
	e.ApplyForce(DragForce(e.Velocity(), 1))
}
