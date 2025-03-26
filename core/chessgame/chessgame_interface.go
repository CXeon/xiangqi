package chessgame

import "github.com/CXeon/xiangqi/core/player"

type ChessGameInterface interface {

	//初始化棋局
	//先手player的阵营会被初始化在棋盘下方
	InitialGame(player1, player2 player.PlayerInterface) error

	//重置棋局
	ResetGame() error

	//销毁棋局
	Destroy() error

	//运行棋局
	Run(downPlayerCh, upPlayerCh chan player.Statement) (msgChan <-chan GameMsg)
}
