package app

const (
	boardLogicZeroX  int = 32
	boardLogicZeroY  int = 32
	gridLength       int = 64
	spriteReparation int = -26
	boardReparation  int = 64 //棋盘坐标补偿单位为64个像素
)

type coordinate struct {
	x, y int //定义坐标x，y
}
