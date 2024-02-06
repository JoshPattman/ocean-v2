package main

var DefSimSettings = SimulationSettings{
	MapGenerationParams: MapGenerationParams{
		Length:       512,
		Height:       256,
		PerlinWidth:  256,
		PerlinHeight: 256,
		CaveWidth:    30,
		CaveAR:       2,
		CaveThresh:   0.1,
		Seed:         -1,
	},
}

var DefUserSettings = UserSettings{
	CameraSettings: CameraSettings{
		ZoomSpeed: 4.0,
		MoveSpeed: 500,
	},
}

// MapGenerationParams are the parameters used to generate a new environment
type MapGenerationParams struct {
	Length       int     `json:"length"`
	Height       int     `json:"height"`
	PerlinWidth  float64 `json:"perlin-width"`
	PerlinHeight float64 `json:"perlin-height"`
	CaveWidth    float64 `json:"cave-width"`
	CaveAR       float64 `json:"cave-height"`
	CaveThresh   float64 `json:"cave-thresh"`
	Seed         int64   `json:"seed"`
}

type SimulationSettings struct {
	MapGenerationParams MapGenerationParams `json:"map-gen"`
}

type CameraSettings struct {
	ZoomSpeed float64 `json:"zoom-speed"`
	MoveSpeed float64 `json:"move-speed"`
}

type UserSettings struct {
	CameraSettings CameraSettings `json:"camera"`
}
