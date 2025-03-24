package chessman

import (
	"errors"
	"github.com/CXeon/xiangqi/core"
	"math"
)

// “马”的移动规则校验
// 马走日
func RuleMa(matrix [][]ChessmanInterface, rowGroup map[int]core.ChessmanGroup, source, target core.Coordinate) (bool, error) {

	if verifyCoordinateOnTheBoard(source) == false || verifyCoordinateOnTheBoard(target) == false {
		return false, errors.New("invalid move")
	}
	//“马”是否走了日
	if math.Abs(float64(target.X-source.X))*math.Abs(float64(target.Y-source.Y)) != 2.0 {
		return false, errors.New("invalid move")
	}

	//确认“马”走了个什么日
	moveY := target.Y - source.Y
	moveX := target.X - source.X

	//记录垂直和水平的移动方向，确认是否是撇脚马
	vertical := ""
	horizontal := ""

	if moveY == 2 {
		vertical = "UP2"
	}
	if moveY == 1 {
		vertical = "UP1"
	}
	if moveY == -2 {
		vertical = "DOWN2"
	}
	if moveY == -1 {
		vertical = "DOWN1"
	}
	if moveX == 1 {
		horizontal = "LEFT1"
	}
	if moveX == 2 {
		horizontal = "LEFT2"
	}
	if moveX == -1 {
		horizontal = "RIGHT1"
	}
	if moveX == -2 {
		horizontal = "RIGHT2"
	}

	vh := vertical + horizontal

	switch vh {
	case "UP2LEFT1", "UP2RIGHT1":
		//确认是否是“撇脚马”
		if matrix[source.Y+1][source.X] != nil {
			return false, errors.New("invalid move")
		}
	case "DOWN2LEFT1", "DOWN2RIGHT1":
		if matrix[source.Y-1][source.X] != nil {
			return false, errors.New("invalid move")
		}
	case "UP1LEFT2", "DOWN1LEFT2":
		if matrix[source.Y][source.X+1] != nil {
			return false, errors.New("invalid move")
		}
	case "UP1RIGHT2", "DOWN1RIGHT2":
		if matrix[source.Y][source.X-1] != nil {
			return false, errors.New("invalid move")
		}
	default:
		return false, errors.New("invalid move")
	}

	//目的坐标是否已经被同阵营的棋子占据
	if targetCoordinateHaveSameGroupChess(matrix, source, target) {
		return false, errors.New("target has a chessman of the same group")
	}

	return true, nil
}

// "象"的移动规则校验
// 象飞田
func RuleXiang(matrix [][]ChessmanInterface, rowGroup map[int]core.ChessmanGroup, source, target core.Coordinate) (bool, error) {
	//检查象棋的移动是否在棋盘内
	if verifyCoordinateOnTheBoard(source) == false || verifyCoordinateOnTheBoard(target) == false {
		return false, errors.New("invalid move")
	}

	//限制象只能在一个阵营区域移动
	if rowGroup[source.Y] != rowGroup[target.Y] {
		return false, errors.New("invalid move")
	}

	//象是否飞了田
	if math.Abs(float64(target.X-source.X)) != 2.0 || math.Abs(float64(target.Y-source.Y)) != 2.0 {
		return false, errors.New("invalid move")
	}

	//象飞了个什么田
	moveY := target.Y - source.Y
	moveX := target.X - source.X

	vertical := ""
	horizontal := ""

	if moveY == 2 {
		vertical = "UP2"
	}
	if moveY == -2 {
		vertical = "DOWN2"
	}
	if moveX == 2 {
		horizontal = "LEFT2"
	}
	if moveX == -2 {
		horizontal = "RIGHT2"
	}

	vh := vertical + horizontal

	switch vh {
	case "UP2LEFT2":
		if matrix[source.Y+1][source.X+1] != nil {
			return false, errors.New("invalid move")
		}
	case "DOWN2LEFT2":
		if matrix[source.Y-1][source.X+1] != nil {
			return false, errors.New("invalid move")
		}
	case "UP2RIGHT2":
		if matrix[source.Y+1][source.X-1] != nil {
			return false, errors.New("invalid move")
		}
	case "DOWN2RIGHT2":
		if matrix[source.Y-1][source.X-1] != nil {
			return false, errors.New("invalid move")
		}
	default:
		return false, errors.New("invalid move")
	}

	//目标位置是否存在己方棋子
	if targetCoordinateHaveSameGroupChess(matrix, source, target) {
		return false, errors.New("target has a chessman of the same group")
	}

	return true, nil
}

// “兵（卒）的移动规则”
func RuleBingZu(matrix [][]ChessmanInterface, rowGroup map[int]core.ChessmanGroup, source, target core.Coordinate) (bool, error) {
	//检查象棋的移动是否在棋盘内
	if verifyCoordinateOnTheBoard(source) == false || verifyCoordinateOnTheBoard(target) == false {
		return false, errors.New("invalid move")
	}

	//检查象棋进行了哪种移动
	moveY := target.Y - source.Y
	moveX := target.X - source.X

	vertical := ""
	horizontal := ""

	if moveY == 1 {
		vertical = "UP1"
	}
	if moveY == -1 {
		vertical = "DOWN1"
	}
	if moveX == 1 {
		horizontal = "LEFT1"
	}
	if moveX == -1 {
		horizontal = "RIGHT1"
	}
	vh := vertical + horizontal

	switch vh {
	case "UP1":
		if matrix[source.Y][source.X].GetChessmanGroup() != rowGroup[4] {
			return false, errors.New("invalid move")
		}
	case "DOWN1":
		if matrix[source.Y][source.X].GetChessmanGroup() != rowGroup[5] {
			return false, errors.New("invalid move")
		}
	case "LEFT1", "RIGHT1":
		if matrix[source.Y][source.X].GetChessmanGroup() == rowGroup[4] {
			return false, errors.New("invalid move")
		}
	default:
		return false, errors.New("invalid move")
	}

	//目标位置是否存在己方棋子
	if targetCoordinateHaveSameGroupChess(matrix, source, target) {
		return false, errors.New("target has a chessman of the same group")
	}

	return true, nil
}

// 目的坐标是否已经被同阵营的棋子占据
func targetCoordinateHaveSameGroupChess(matrix [][]ChessmanInterface, source, target core.Coordinate) bool {

	if matrix[target.Y][target.X] == nil {
		return false
	}

	if matrix[source.Y][source.X].GetChessmanGroup() == matrix[target.Y][target.X].GetChessmanGroup() {
		return true
	}
	return false
}

// 确认坐标仍在在棋盘范围内
func verifyCoordinateOnTheBoard(co core.Coordinate) bool {
	x := co.X
	y := co.Y
	if x < 0 || x > 8 {
		return false
	}
	if y < 0 || y > 9 {
		return false
	}
	return true
}
