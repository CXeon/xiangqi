package player

import "github.com/CXeon/xiangqi/core"

type Statement struct {
	Group  core.ChessmanGroup //阵营
	Code   core.ChessmanCode  //棋子code
	Source core.Coordinate    //起始坐标
	Target core.Coordinate    //目的坐标
}
