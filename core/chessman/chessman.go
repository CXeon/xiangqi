package chessman

import "github.com/CXeon/xiangqi/core"

type Chessman struct {
	code              core.ChessmanCode  //棋子code
	name              string             //棋子名称
	group             core.ChessmanGroup //棋子阵营
	isDead            bool               //是否被吃掉
	defaultCoordinate core.Coordinate    //默认坐标

	rule ChessRUle
}

// 创建棋子对象
func NewChessman(code core.ChessmanCode, name string, group core.ChessmanGroup, defaultCoordinate core.Coordinate) *Chessman {
	return &Chessman{
		code:              code,
		name:              name,
		group:             group,
		isDead:            false,
		defaultCoordinate: defaultCoordinate,
	}
}

// 获取棋子代号
func (cm *Chessman) GetChessmanCode() core.ChessmanCode {
	return cm.code
}

// 获取棋子名称
func (cm *Chessman) GetChessmanName() string {
	return cm.name
}

// 获取棋子阵营
func (cm *Chessman) GetChessmanGroup() core.ChessmanGroup {
	return cm.group
}

// 获取棋子默认坐标
func (cm *Chessman) GetChessmanDefaultCoordinate() core.Coordinate {
	return cm.defaultCoordinate
}

// 设置存活状态
func (cm *Chessman) SetIsDead(isDead bool) {
	cm.isDead = isDead
	return
}

// 获取存活状态
func (cm *Chessman) GetIsDead() bool {
	return cm.isDead
}

// 绑定棋子规则
func (cm *Chessman) BindRule(rule ChessRUle) {
	cm.rule = rule
}

// 校验棋子移动是否符合规则
func (cm *Chessman) CheckMove(matrix [][]ChessmanInterface, rowGroup map[int]core.ChessmanGroup, source, target core.Coordinate) (bool, error) {
	if cm.rule == nil {
		return false, nil
	}
	return cm.rule(matrix, rowGroup, source, target)
}
