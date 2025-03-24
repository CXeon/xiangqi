package chessboard

import (
	"github.com/CXeon/xiangqi/core"
	"github.com/CXeon/xiangqi/core/chessman"
)

type ChessboardInterface interface {

	//棋盘划分阵营
	DivideGroup(group core.ChessmanGroup, rowIndex []int)

	//获取某一行属于哪个阵营
	GetGroupInRow(rowIndex int) core.ChessmanGroup

	//放置棋子到棋盘
	PutChessmenOnBoard(chessmen []chessman.ChessmanInterface) error

	//移动棋子
	MoveChessman(group core.ChessmanGroup, code core.ChessmanCode, source, target core.Coordinate) (won core.ChessmanCode, err error)
}
