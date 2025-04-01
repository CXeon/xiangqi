package main

import (
	"github.com/CXeon/xiangqi/app"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	ebiten.SetWindowSize(app.ScreenWidth, app.ScreenHeight)
	ebiten.SetWindowTitle("XiangQi Demo")
	game := app.NewGame()
	defer game.Close()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
