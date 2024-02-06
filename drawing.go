package main

import "github.com/gopxl/pixel"

type RenderData struct {
	Target         pixel.Target
	TargetRect     pixel.Rect
	CameraWorldPos pixel.Vec
	PixelsPerMeter float64
}

type Renderable interface {
	Render(*RenderData)
}
