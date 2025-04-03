package chessman

import "github.com/CXeon/xiangqi/core"

type ChessRUle func(matrix [][]ChessmanInterface, rowGroup map[int]core.ChessmanGroup, source, target core.Coordinate) (bool, error)

type ChessmanInterface interface {
	//获取棋子代号
	GetChessmanCode() core.ChessmanCode
	//获取棋子名称
	GetChessmanName() string
	//获取棋子阵营
	GetChessmanGroup() core.ChessmanGroup
	//获取棋子默认坐标
	GetChessmanDefaultCoordinate() core.Coordinate

	//设置存活状态
	SetIsDead(isDead bool)

	//获取存活状态
	GetIsDead() bool

	// 绑定棋子规则
	BindRule(rule ChessRUle)

	// 校验棋子移动是否符合规则
	CheckMove(matrix [][]ChessmanInterface, rowGroup map[int]core.ChessmanGroup, source, target core.Coordinate) (bool, error)
}
