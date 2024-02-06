package main

import (
	"image/png"
	"os"
	"path"
	"strings"

	"github.com/gopxl/pixel"
)

var globalResPics map[string]pixel.Picture

func init() {
	// Load up all the GlobalResPics
	globalResPics = make(map[string]pixel.Picture)
	sp := path.Join(".", "data", "sprites")
	entries, err := os.ReadDir(sp)
	if err != nil {
		panic("err")
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		f, err := os.Open(path.Join(sp, e.Name()))
		if err != nil {
			panic(err)
		}
		img, err := png.Decode(f)
		if err != nil {
			panic(err)
		}
		globalResPics[strings.TrimSuffix(e.Name(), ".png")] = pixel.PictureDataFromImage(img)
	}
}

func GetSpritePicture(name string) pixel.Picture {
	if pic, ok := globalResPics[name]; ok {
		return pic
	}
	panic("sprite did not exist")
}
