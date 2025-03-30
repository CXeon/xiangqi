package chessboard

import (
	"errors"
	"fmt"
	"github.com/CXeon/xiangqi/core"
	"github.com/CXeon/xiangqi/core/chessman"
)

const (
	rows = 10 //棋盘行数
	cols = 9  //棋盘列数
)

type Chessboard struct {
	//棋盘上所有的位置，必须被初始化为 10行9列。
	//索引代表坐标 对应棋盘的方向为从右往左，从下往上
	matrix [][]chessman.ChessmanInterface

	//保存每一行对应哪个阵营，key是行在数组的索引，value是阵营名称
	rowGroup map[int]core.ChessmanGroup
}

// NewChessboard 新建棋盘
func NewChessboard() *Chessboard {
	return &Chessboard{
		matrix:   initMatrix(rows, cols),
		rowGroup: make(map[int]core.ChessmanGroup),
	}
}

// 棋盘划分阵营
func (board *Chessboard) DivideGroup(group core.ChessmanGroup, rowIndex []int) {
	for _, ri := range rowIndex {
		board.rowGroup[ri] = group
	}
}

// 获取某一行属于哪个阵营
func (board *Chessboard) GetGroupInRow(rowIndex int) core.ChessmanGroup {
	return board.rowGroup[rowIndex]
}

// 初始化棋盘格矩阵
func initMatrix(row, column int) [][]chessman.ChessmanInterface {
	matrix := make([][]chessman.ChessmanInterface, row)
	for i := 0; i < rows; i++ {
		matrix[i] = make([]chessman.ChessmanInterface, column)
	}
	return matrix
}

// 放置棋子到棋盘
func (board *Chessboard) PutChessmenOnBoard(chessmen []chessman.ChessmanInterface) error {
	for _, chess := range chessmen {
		defaultX := chess.GetChessmanDefaultCoordinate().X
		defaultY := chess.GetChessmanDefaultCoordinate().Y
		if defaultX > cols-1 {
			msg := fmt.Sprintf("the default x is too big for chess %s", chess.GetChessmanCode())
			//遇到报错重置矩阵
			board.matrix = initMatrix(rows, cols)
			return errors.New(msg)
		}
		if defaultY > rows-1 {
			msg := fmt.Sprintf("the default y is too big for chess %s", chess.GetChessmanCode())
			//遇到报错重置棋盘
			board.matrix = initMatrix(rows, cols)
			return errors.New(msg)
		}
		board.matrix[defaultY][defaultX] = chess
	}
	return nil
}

// 移动棋子
func (board *Chessboard) MoveChessman(group core.ChessmanGroup, code core.ChessmanCode, source, target core.Coordinate) (won core.ChessmanCode, err error) {
	//找到想要移动的棋子
	cm, err := board.getChessman(group, code, source)
	if err != nil {
		return "", err
	}

	//验证启动棋子是否符合规则
	_, err = cm.CheckMove(board.matrix, board.rowGroup, source, target)
	if err != nil {
		return "", err
	}

	//移动棋子，如果target坐标上有棋子，那target坐标的棋子会被吃掉
	if board.matrix[target.Y][target.X] != nil {
		board.matrix[target.Y][target.X].SetIsDead(true)
		won = board.matrix[target.Y][target.X].GetChessmanCode()
	}
	board.matrix[source.Y][source.X] = nil
	board.matrix[target.Y][target.X] = cm

	return won, nil
}

// 在棋盘上查找棋子
func (board *Chessboard) getChessman(group core.ChessmanGroup, code core.ChessmanCode, location core.Coordinate) (chess chessman.ChessmanInterface, err error) {
	x, y := location.X, location.Y
	cm := board.matrix[y][x]
	if cm == nil {
		return nil, errors.New(fmt.Sprintf("the chessman %s is not exist", code))
	}
	if cm.GetChessmanCode() != code || cm.GetChessmanGroup() != group {
		return nil, errors.New(fmt.Sprintf("the location [%d,%d] is other chess group %d, code %s", location.X, location.Y, cm.GetChessmanGroup(), cm.GetChessmanCode()))
	}
	return cm, nil
}

// 返回棋盘的棋子以及位置
func (board *Chessboard) GetMatrix() [][]chessman.ChessmanInterface {
	return board.matrix
}

// 返回棋盘上每一行对应的阵营
func (board *Chessboard) GetRowsGroup() map[int]core.ChessmanGroup {
	return board.rowGroup
}

func (board *Chessboard) ClearChessmen() {
	board.matrix = initMatrix(rows, cols)
}
