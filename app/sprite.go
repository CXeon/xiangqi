package app

import (
	"github.com/CXeon/xiangqi/core"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
)

// Sprite 代表棋子的图片.
type Sprite struct {
	image      *ebiten.Image
	alphaImage *image.Alpha
	x          int                //在棋盘的坐标
	y          int                //在棋盘的坐标
	clicked    bool               //被点击选中
	group      core.ChessmanGroup //棋子所属阵营
	code       core.ChessmanCode  //棋子代号

}

// In returns true if (x, y) is in the sprite, and false otherwise.
func (s *Sprite) In(x, y int) bool {
	return s.alphaImage.At(x-s.x, y-s.y).(color.Alpha).A > 0
}

// MoveTo moves the sprite to the position (x, y).
func (s *Sprite) MoveTo(x, y int) {
	//w, h := s.image.Bounds().Dx(), s.image.Bounds().Dy()

	s.x = x
	s.y = y
	if s.x < 0 {
		s.x = 0
	}
	if s.x > ScreenWidth {
		s.x = ScreenWidth
	}
	if s.y < 0 {
		s.y = 0
	}
	if s.y > ScreenHeight {
		s.y = ScreenHeight
	}
}

// Draw draws the sprite.
func (s *Sprite) Draw(screen *ebiten.Image, alpha float32) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(s.x), float64(s.y))
	op.ColorScale.ScaleAlpha(alpha)
	screen.DrawImage(s.image, op)
}
