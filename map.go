package main

import (
	"math"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"github.com/gopxl/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/aquilax/go-perlin"
)

// Width of each texel in pixels
var mapTextureTexelWidth int = 16

// Texel is an id that describes a single block
type Texel int

// Define the texels
const (
	WaterTexel Texel = iota
	RockTexel
	SandTexel
)

// Map is used to store information about the envrionment, primarily the texels
type Map struct {
	texels       [][]Texel
	spriteSheet  pixel.Picture
	sprites      map[Texel]*pixel.Sprite
	texelsCanvas *pixelgl.Canvas
	imd          *imdraw.IMDraw
	dirty        bool
}

// NewGeneratedMap generates a new environment using the given params.
// It also loads up all textures and initialises all other elements of the environment.
func NewGeneratedMap(genParams MapGenerationParams) *Map {
	texels := make([][]Texel, genParams.Length)
	perlinGen := perlin.NewPerlin(2, 2, 5, genParams.Seed)
	for tx := range texels {
		texels[tx] = make([]Texel, genParams.Height)
		height := int(math.Round(perlinGen.Noise1D(float64(tx)/genParams.PerlinWidth) * genParams.PerlinHeight))
		sandWidth := int(perlinGen.Noise1D(float64(tx)/64) * 15)
		for ty := range texels[tx] {
			density := perlinGen.Noise2D(float64(tx)/(genParams.CaveWidth*genParams.CaveAR), float64(ty)/genParams.CaveWidth)
			if tx == 0 || ty == 0 || tx == len(texels)-1 || ty == len(texels[tx])-1 {
				// Border rock
				texels[tx][ty] = RockTexel
			} else if ty < height+genParams.Height/2 && !(density > -genParams.CaveThresh && density < genParams.CaveThresh) {
				// This is terrain
				if ty >= height+genParams.Height/2-sandWidth-1 {
					// Terrain sand
					texels[tx][ty] = SandTexel
				} else {
					// Terrain rock
					texels[tx][ty] = RockTexel
				}
			} else {
				// Terrain air
				texels[tx][ty] = WaterTexel
			}
		}
	}
	spriteSheet := GetSpritePicture("textures")
	spritesMap := make(map[Texel]*pixel.Sprite)
	spritesMap[RockTexel] = spriteFromTileSheet(spriteSheet, 0, 14, mapTextureTexelWidth)
	spritesMap[WaterTexel] = spriteFromTileSheet(spriteSheet, 0, 1, mapTextureTexelWidth)
	spritesMap[SandTexel] = spriteFromTileSheet(spriteSheet, 0, 6, mapTextureTexelWidth)

	return &Map{
		texels:       texels,
		spriteSheet:  spriteSheet,
		sprites:      spritesMap,
		dirty:        true,
		texelsCanvas: pixelgl.NewCanvas(pixel.R(0, 0, float64(mapTextureTexelWidth*genParams.Length), float64(mapTextureTexelWidth*genParams.Height))),
		imd:          imdraw.New(nil),
	}
}

func (m *Map) Render(rd *RenderData) {
	if m.dirty {
		m.texelsCanvas.Clear(pixel.Alpha(0))
		m.imd.Clear()
		for tx := range m.texels {
			for ty := range m.texels[tx] {
				texel := m.texels[tx][ty]
				screenPos := pixel.V(float64(tx), float64(ty)).Scaled(float64(mapTextureTexelWidth))
				worldPos := pixel.V(float64(tx), float64(ty))
				//depth := m.GetDepthAt(worldPos)
				light := m.GetLightAt(worldPos)
				if texel != WaterTexel {
					sprite := m.sprites[texel]
					drawMat := pixel.IM.Moved(screenPos)
					sprite.Draw(m.texelsCanvas, drawMat)
				} else {
					col := pixel.ToRGBA(colornames.Skyblue).Scaled(light).Add(pixel.ToRGBA(pixel.RGB(17.0/255, 42.0/255, 82.0/255)).Scaled(1 - light))
					//col = col.Scaled(light)
					m.imd.Color = col
					squareRad := float64(mapTextureTexelWidth) / 2
					m.imd.Push(screenPos.Sub(pixel.V(squareRad, squareRad)), screenPos.Add(pixel.V(squareRad, squareRad)))
					m.imd.Rectangle(0)
				}
			}
		}
		m.imd.Draw(m.texelsCanvas)
		m.dirty = false
	}
	m.texelsCanvas.Draw(rd.Target, pixel.IM.Moved(m.texelsCanvas.Bounds().Center()).Scaled(pixel.ZV, 1.0/float64(mapTextureTexelWidth)).Moved(rd.CameraWorldPos.Scaled(-1)).Scaled(pixel.ZV, rd.PixelsPerMeter).Moved(rd.TargetRect.Center()))
}

// spriteFromTileSheet extracts a single sprite from a tilesheet of uniform sized square tiles
func spriteFromTileSheet(pic pixel.Picture, coordx, coordy int, tileSize int) *pixel.Sprite {
	sprite := pixel.NewSprite(pic, pixel.R(float64(coordx*tileSize), float64(coordy*tileSize), float64((coordx+1)*tileSize), float64((coordy+1)*tileSize)))
	return sprite
}

// Returns the depth, between 1 and 0, of the provided point
func (m *Map) GetDepthAt(pos pixel.Vec) float64 {
	return 1 - pos.Y/float64(len(m.texels[0]))
}

// Returns the light level between 0 and 1 of the provided point
func (m *Map) GetLightAt(pos pixel.Vec) float64 {
	xPos := int(math.Round(pos.X))
	yPos := int(math.Round(pos.Y))
	covered := false
	// -1 is here so the top border does not cover
	for y := yPos; y < len(m.texels[xPos])-1; y++ {
		if m.texels[xPos][y] != WaterTexel {
			covered = true
			break
		}
	}
	if covered {
		return 0.0
	} else {
		return 1 - m.GetDepthAt(pos)
	}
}
