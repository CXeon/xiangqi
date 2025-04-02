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
		if matrix[source.Y][source.X].GetChessmanGroup() == rowGroup[source.Y] {
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

// "车"的移动规则校验
// “车”走直线
func RuleJu(matrix [][]ChessmanInterface, rowGroup map[int]core.ChessmanGroup, source, target core.Coordinate) (bool, error) {
	//检查象棋的移动是否在棋盘内
	if verifyCoordinateOnTheBoard(source) == false || verifyCoordinateOnTheBoard(target) == false {
		return false, errors.New("invalid move")
	}

	//确定移动方向
	moveY := target.Y - source.Y
	moveX := target.X - source.X

	vh := ""
	if moveX == 0 && moveY > 0 {
		vh = "UP"
	}
	if moveX == 0 && moveY < 0 {
		vh = "DOWN"
	}
	if moveX > 0 && moveY == 0 {
		vh = "LEFT"
	}
	if moveX < 0 && moveY == 0 {
		vh = "RIGHT"
	}

	switch vh {
	case "UP":
		//向上移动，中间不能有其他棋子存在
		for i := source.Y + 1; i < target.Y; i++ {
			if matrix[i][source.X] != nil {
				return false, errors.New("invalid move")
			}
		}
	case "DOWN":
		//向下移动，中间不能有其他棋子存在
		for i := source.Y - 1; i > target.Y; i-- {
			if matrix[i][source.X] != nil {
				return false, errors.New("invalid move")
			}
		}

	case "LEFT":
		//向左移动，中间不能有其他棋子存在
		for i := source.X + 1; i < target.X; i++ {
			if matrix[source.Y][i] != nil {
				return false, errors.New("invalid move")
			}
		}
	case "RIGHT":
		for i := source.X - 1; i > target.X; i-- {
			if matrix[source.Y][i] != nil {
				return false, errors.New("invalid move")
			}
		}
	default:
		return false, errors.New("invalid move")
	}

	if targetCoordinateHaveSameGroupChess(matrix, source, target) {
		return false, errors.New("target has a chessman of the same group")
	}
	return true, nil

}

// "炮"的移动规则
// “炮”翻山
func RulePao(matrix [][]ChessmanInterface, rowGroup map[int]core.ChessmanGroup, source, target core.Coordinate) (bool, error) {
	//检查象棋的移动是否在棋盘内
	if verifyCoordinateOnTheBoard(source) == false || verifyCoordinateOnTheBoard(target) == false {
		return false, errors.New("invalid move")
	}

	//确定移动方向
	moveY := target.Y - source.Y
	moveX := target.X - source.X

	vh := ""
	if moveX == 0 && moveY > 0 {
		vh = "UP"
	}
	if moveX == 0 && moveY < 0 {
		vh = "DOWN"
	}
	if moveX > 0 && moveY == 0 {
		vh = "LEFT"
	}
	if moveX < 0 && moveY == 0 {
		vh = "RIGHT"
	}
	//判定是否发生吃棋
	if targetCoordinateHaveSameGroupChess(matrix, source, target) {
		return false, errors.New("target has a chessman of the same group")
	}
	if chessman := matrix[target.Y][target.X]; chessman != nil {
		//发生吃棋，那目标坐标和起始坐标之间必须要存在1颗棋子
		switch vh {
		case "UP":
			count := 0
			//向上移动
			for i := source.Y + 1; i < target.Y; i++ {
				if matrix[i][source.X] != nil {
					count++
				}
			}
			if count != 1 {
				return false, errors.New("invalid move")
			}
		case "DOWN":
			count := 0
			//向下移动
			for i := source.Y - 1; i > target.Y; i-- {
				if matrix[i][source.X] != nil {
					count++
				}
			}
			if count != 1 {
				return false, errors.New("invalid move")
			}
		case "LEFT":
			count := 0
			//向左移动
			for i := source.X + 1; i < target.X; i++ {
				if matrix[source.Y][i] != nil {
					count++
				}
			}
			if count != 1 {
				return false, errors.New("invalid move")
			}
		case "RIGHT":
			count := 0
			for i := source.X - 1; i > target.X; i-- {
				if matrix[source.Y][i] != nil {
					count++
				}
			}
			if count != 1 {
				return false, errors.New("invalid move")
			}
		default:
			return false, errors.New("invalid move")
		}
	}
	if chessman := matrix[target.Y][target.X]; chessman == nil {
		//没有发生吃棋，那目标坐标和起始坐标之间不能存在棋子
		switch vh {
		case "UP":
			//向上移动，中间不能有其他棋子存在
			for i := source.Y + 1; i < target.Y; i++ {
				if matrix[i][source.X] != nil {
					return false, errors.New("invalid move")
				}
			}
		case "DOWN":
			//向下移动，中间不能有其他棋子存在
			for i := source.Y - 1; i > target.Y; i-- {
				if matrix[i][source.X] != nil {
					return false, errors.New("invalid move")
				}
			}

		case "LEFT":
			//向左移动，中间不能有其他棋子存在
			for i := source.X + 1; i < target.X; i++ {
				if matrix[source.Y][i] != nil {
					return false, errors.New("invalid move")
				}
			}
		case "RIGHT":
			for i := source.X - 1; i > target.X; i-- {
				if matrix[source.Y][i] != nil {
					return false, errors.New("invalid move")
				}
			}
		default:
			return false, errors.New("invalid move")
		}
	}

	return true, nil
}

// “士”的移动规则校验
func RuleShi(matrix [][]ChessmanInterface, rowGroup map[int]core.ChessmanGroup, source, target core.Coordinate) (bool, error) {
	//首先验证棋子是否在规定范围内活动
	if verifyJiangShuaiOrShiInZone(matrix, rowGroup, source, target) == false {
		return false, errors.New("invalid move")
	}

	//验证目标坐标是否存在同阵营棋子
	if targetCoordinateHaveSameGroupChess(matrix, source, target) {
		return false, errors.New("target has a chessman of the same group")
	}

	//确认移动方式
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
		vertical = "LEFT1"
	}
	if moveX == -1 {
		vertical = "RIGHT1"
	}

	vh := vertical + horizontal

	switch vh {
	case "UP1LEFT1", "UP1RIGHT1", "DOWN1LEFT1", "DOWN1RIGHT1":
		return true, nil
	}
	return false, errors.New("invalid move")
}

// “将/帅”的移动规则校验
func RuleJiangShuai(matrix [][]ChessmanInterface, rowGroup map[int]core.ChessmanGroup, source, target core.Coordinate) (bool, error) {
	//首先验证棋子是否在规定范围内活动
	if verifyJiangShuaiOrShiInZone(matrix, rowGroup, source, target) == false {
		return false, errors.New("invalid move")
	}

	//验证目标坐标是否存在同阵营棋子
	if targetCoordinateHaveSameGroupChess(matrix, source, target) {
		return false, errors.New("target has a chessman of the same group")
	}

	//确认移动方式
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
		vertical = "LEFT1"
	}
	if moveX == -1 {
		vertical = "RIGHT1"
	}

	vh := vertical + horizontal

	switch vh {
	case "UP1", "RIGHT1", "DOWN1", "LEFT1":
		return true, nil
	}
	return false, errors.New("invalid move")
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

func verifyJiangShuaiOrShiInZone(matrix [][]ChessmanInterface, rowGroup map[int]core.ChessmanGroup, source, target core.Coordinate) bool {
	//首先确定棋子阵营
	chessman := matrix[source.Y][source.X]
	if chessman.GetChessmanGroup() == rowGroup[0] {
		//起始坐标和目标坐标都应该在棋盘下方中心的田字格
		if source.X < 3 || source.X > 5 {
			return false
		}
		if target.X < 3 || target.X > 5 {
			return false
		}

		if source.Y < 0 || source.Y > 2 {
			return false
		}
		if target.Y < 0 || target.Y > 2 {
			return false
		}
		return true
	} else {
		//起始坐标和目标坐标都应该棋盘上方的田字格
		if source.X < 3 || source.X > 5 {
			return false
		}
		if target.X < 3 || target.X > 5 {
			return false
		}

		if source.Y < 7 || source.Y > 9 {
			return false
		}
		if target.Y < 7 || target.Y > 9 {
			return false
		}
		return true
	}
}
