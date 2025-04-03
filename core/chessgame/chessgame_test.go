package chessgame

import (
	"fmt"
	"github.com/CXeon/xiangqi/core"
	"github.com/CXeon/xiangqi/core/player"
	"testing"
)

func TestChessGame(t *testing.T) {
	//玩家就位
	var p1, p2 player.PlayerInterface
	p1 = player.NewPlayer()
	p2 = player.NewPlayer()

	p1.SetGroup(core.Group1)
	p1.SetIsFirst(true)
	p2.SetGroup(core.Group2)

	var chP1 = make(chan player.Statement, 1)
	defer close(chP1)
	var chP2 = make(chan player.Statement, 1)
	defer close(chP2)
	//创建棋局
	chessGame := new(ChessGame)
	//初始化棋局
	err := chessGame.InitialGame(p1, p2)
	if err != nil {
		t.Error(err)
		return
	}

	//打印棋局
	chessGame.Show()

	msgChan := chessGame.Run(chP1, chP2)

	//测试先手“炮二进四”，后手“卒9进1”
	chP1 <- player.Statement{
		Group: core.Group1,
		Code:  core.Pao,
		Source: core.Coordinate{
			X: 1,
			Y: 2,
		},
		Target: core.Coordinate{
			X: 1,
			Y: 6,
		},
	}

	msg := <-msgChan
	fmt.Println(msg)

	chP2 <- player.Statement{
		Group: core.Group2,
		Code:  core.BingZu,
		Source: core.Coordinate{
			X: 0,
			Y: 6,
		},
		Target: core.Coordinate{
			X: 0,
			Y: 5,
		},
	}

	msg = <-msgChan
	fmt.Println(msg)

	//打印棋局
	chessGame.Show()

	//先手再“炮二进三”，吃马
	chP1 <- player.Statement{
		Group: core.Group1,
		Code:  core.Pao,
		Source: core.Coordinate{
			X: 1,
			Y: 6,
		},
		Target: core.Coordinate{
			X: 1,
			Y: 9,
		},
	}
	msg = <-msgChan
	fmt.Println(msg)

	//打印棋局
	chessGame.Show()

	//停止棋局
	err = chessGame.Close()
	if err != nil {
		t.Error(err)
		return
	}

	return
}
