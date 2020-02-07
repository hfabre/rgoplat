package main

import (
	"encoding/json"
	"fmt"
	r "github.com/lachee/raylib-goplus/raylib"
	"io/ioutil"
)

//
//  Tileset
//

type Tileset struct {
	tiles []r.Texture2D
}

func NewTileset(path string, tileWidth, tileHeight int) Tileset {
	var tiles []r.Texture2D
	tileset := r.LoadImage(path)
	horizontalTileCount := int(tileset.Width) / tileWidth
	verticalTileCount := int(tileset.Height) / tileHeight

	for y := 0; y < verticalTileCount; y++ {
		for x := 0; x < horizontalTileCount; x++ {
			rect := r.Rectangle{X: float32(x * tileWidth), Y: float32(y * tileHeight), Width: float32(tileWidth), Height: float32(tileHeight)}
			tiles = append(tiles, r.LoadTextureFromImage(r.ImageFromImage(tileset, rect)))
		}
	}

	return Tileset{tiles: tiles}
}

func (ts Tileset) Unload() {
	for _, tile := range ts.tiles {
		r.UnloadTexture(tile)
	}
}

func (ts Tileset) Debug() {
	for i, tile := range ts.tiles {
		r.DrawTexture(tile, i*32, 0, r.White)
		r.DrawText(fmt.Sprintf("<%d>", i), i*32, 0, 10, r.Red)
	}
}

//
//  MapConfiguration
//

type MapConfiguration struct {
	Width      int
	Height     int
	TileWidth  int
	TileHeight int
	Board      [][]int
}

func NewMapConfiguration(path string) MapConfiguration {
	var mc MapConfiguration

	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Print(err)
	}

	err = json.Unmarshal(data, &mc)
	if err != nil {
		fmt.Println("error:", err)
	}

	return mc
}

//
//  Map
//

type Map struct {
	mc         MapConfiguration
	ts	       Tileset

	width      int
	height     int
	tileWidth  int
	tileHeight int
	board      [][]int
}

func NewMap(mc MapConfiguration, ts Tileset) Map {
	m := Map{
		mc: mc,
		ts: ts,
		width: mc.Width,
		height: mc.Height,
		tileWidth: mc.TileWidth,
		tileHeight: mc.TileHeight,
		board: mc.Board,
	}

	return m
}

func (m Map) Draw() {
	for y := 0; y < len(m.board); y++ {
		for x := 0; x < len(m.board[y]); x++ {
			r.DrawTexture(m.ts.tiles[m.board[y][x]], x * m.tileWidth, y * m.tileHeight, r.White)
		}
	}
}

//
//  Main
//

func main() {
	r.InitWindow(800, 450, "Raylib Go Plus")
	r.SetTargetFPS(60)

	mc := NewMapConfiguration("./media/map.json")
	tileset := NewTileset("./media/tileset.png", mc.TileWidth, mc.TileHeight)
	mmap := NewMap(mc, tileset)
	background := r.LoadTexture("./media/background.png")
	x := 20
	y := 20
	cam := r.Camera2D{
		Offset: r.Vector2{ X: 400, Y: 225 },
		Target: r.Vector2{X: float32(x + 20), Y: float32(y + 20)},
		Rotation: 0,
		Zoom: 1,
	}

	for !r.WindowShouldClose() {

		if r.IsKeyDown(r.KeyRight) {
			x += 2
		} else if r.IsKeyDown(r.KeyLeft) {
			x -= 2
		} else if r.IsKeyDown(r.KeyUp) {
			y -= 2
		} else if r.IsKeyDown(r.KeyDown) {
			y += 2
		}

		cam.Target = r.Vector2{X: float32(x + 20), Y: float32(y + 20)}

		r.BeginDrawing()
		r.ClearBackground(r.RayWhite)

		// Scrolling background but meh ?
		r.DrawTextureEx(background, r.Vector2{X: float32(-x), Y: float32(0)}, 0, 1, r.White)
		r.DrawTextureEx(background, r.Vector2{X: float32(background.Width * 2 + int32(-x)), Y: float32(0)}, 0, 1, r.White)

		r.BeginMode2D(cam)
		mmap.Draw()
		r.EndMode2D()

		r.EndDrawing()
	}

	r.CloseWindow()
	r.UnloadTexture(background)
	tileset.Unload()
}
