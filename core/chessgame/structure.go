package chessgame

import "github.com/CXeon/xiangqi/core"

type GameMsg struct {
	Event           MoveEvent          //事件类型
	WonChessmanCode core.ChessmanCode  //发生吃棋后赢得的棋子code
	WonGroup        core.ChessmanGroup //分出胜负后，赢家的阵营
	Msg             string             //消息
}

type MoveEvent string

const (
	Done MoveEvent = "DONE" //表示移动有效，棋子正常移动
	Err  MoveEvent = "ERR"  //表示移动报错，棋子不能移动
	Fin  MoveEvent = "FIN"  //表示胜负已分，本局对局结束
)
