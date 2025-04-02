package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
)

type onceAgainBtn struct {
	image      *ebiten.Image
	alphaImage *image.Alpha
	x          int //在棋盘的坐标
	y          int //在棋盘的坐标
}

func (oab *onceAgainBtn) In(x, y int) bool {
	return oab.alphaImage.At(x-oab.x, y-oab.y).(color.Alpha).A > 0
}

func (oab *onceAgainBtn) Draw(screen *ebiten.Image, alpha float32) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(oab.x), float64(oab.y))
	op.ColorScale.ScaleAlpha(alpha)
	screen.DrawImage(oab.image, op)
}
