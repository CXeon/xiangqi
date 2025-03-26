package chessgame

import (
	"github.com/CXeon/xiangqi/core"
	"github.com/CXeon/xiangqi/core/chessboard"
	"github.com/CXeon/xiangqi/core/chessman"
	"github.com/CXeon/xiangqi/core/player"
)

type ChessGame struct {
	playerDown   player.PlayerInterface         //玩家1号
	playerUp     player.PlayerInterface         //玩家2号
	board        chessboard.ChessboardInterface //棋盘
	downChessmen []chessman.ChessmanInterface   //棋子
	upChessmen   []chessman.ChessmanInterface
	round        int //记录当前回合执行情况 0：都没下棋，1:先手已下棋，-1:后手已下棋
}

// 初始化棋局
// 先手player的阵营会被初始化在棋盘下方
func (game *ChessGame) InitialGame(player1, player2 player.PlayerInterface) error {

	//引入玩家
	if player1.GetIsFirst() {
		game.playerDown = player1
		game.playerUp = player2
	} else {
		game.playerDown = player2
		game.playerUp = player1
	}

	//创建棋盘
	board := chessboard.NewChessboard()

	//划分棋盘区域
	if player1.GetIsFirst() {
		board.DivideGroup(player1.GetGroup(), []int{0, 1, 2, 3, 4})
		board.DivideGroup(player2.GetGroup(), []int{5, 6, 7, 8, 9})
	} else {
		board.DivideGroup(player2.GetGroup(), []int{0, 1, 2, 3, 4})
		board.DivideGroup(player1.GetGroup(), []int{5, 6, 7, 8, 9})
	}

	//创建棋子
	chessmenOfPlayerDown, chessmenOfPlayerUp := game.newChessmen(game.playerDown.GetGroup(), game.playerUp.GetGroup())

	allChessmen := append(chessmenOfPlayerDown, chessmenOfPlayerUp...)
	//放置棋子
	err := board.PutChessmenOnBoard(allChessmen)
	if err != nil {
		return err
	}
	game.board = board

	game.downChessmen = chessmenOfPlayerDown
	game.upChessmen = chessmenOfPlayerUp

	//记录玩家初始拥有的棋子
	codes1 := make([]core.ChessmanCode, len(chessmenOfPlayerDown))
	for i, c := range chessmenOfPlayerDown {
		codes1[i] = c.GetChessmanCode()
	}
	game.playerDown.AddOwnChessmen(codes1)

	codes2 := make([]core.ChessmanCode, len(chessmenOfPlayerUp))
	for i, c := range chessmenOfPlayerUp {
		codes2[i] = c.GetChessmanCode()
	}
	game.playerUp.AddOwnChessmen(codes2)

	return nil
}

// 重置棋局
func (game *ChessGame) ResetGame() error {
	//清空棋子
	game.board.ClearChessmen()
	//重新放置棋子
	allChessmen := append(game.downChessmen, game.upChessmen...)
	err := game.board.PutChessmenOnBoard(allChessmen)
	if err != nil {
		return err
	}

	//清空玩家记录的棋子获得情况
	game.playerDown.ClearOwnChessman()
	game.playerDown.ClearLostChessman()
	game.playerDown.ClearWonChessman()

	game.playerUp.ClearOwnChessman()
	game.playerUp.ClearLostChessman()
	game.playerUp.ClearWonChessman()

	//重新开始记录
	//记录玩家初始拥有的棋子
	codes1 := make([]core.ChessmanCode, len(game.downChessmen))
	for i, c := range game.downChessmen {
		codes1[i] = c.GetChessmanCode()
	}
	game.playerDown.AddOwnChessmen(codes1)

	codes2 := make([]core.ChessmanCode, len(game.upChessmen))
	for i, c := range game.upChessmen {
		codes2[i] = c.GetChessmanCode()
	}
	game.playerUp.AddOwnChessmen(codes2)

	return nil
}

// 销毁棋局
func (game *ChessGame) Destroy() error {
	game.playerDown.ClearOwnChessman()
	game.playerDown.ClearLostChessman()
	game.playerDown.ClearWonChessman()

	game.playerUp.ClearOwnChessman()
	game.playerUp.ClearLostChessman()
	game.playerUp.ClearWonChessman()

	game.board = nil
	game.downChessmen = nil
	game.upChessmen = nil

	return nil
}

// 运行棋局
func (game *ChessGame) Run(downPlayerCh, upPlayerCh chan player.Statement) (msgChan chan GameMsg) {

	msgChan = make(chan GameMsg, 1)
	go func() {
		defer close(msgChan)
		for {
			switch game.round {
			case 0, -1:
				//棋局开始或者后手已经下完棋，先手下棋
				st, err := game.playerDown.ReceiveStatement(downPlayerCh)
				if err != nil {
					gameMsg := GameMsg{
						Event:           Err,
						WonChessmanCode: "",
						WonGroup:        core.GroupNone,
						Msg:             err.Error(),
					}
					msgChan <- gameMsg
					return
				}
				//移动棋子
				wonCode, err := game.board.MoveChessman(st.Group, st.Code, st.Source, st.Target)
				if err != nil {
					gameMsg := GameMsg{
						Event:           Err,
						WonChessmanCode: "",
						WonGroup:        core.GroupNone,
						Msg:             err.Error(),
					}
					msgChan <- gameMsg
				}

				game.round = 1 //棋子移动完毕先手已经执行

				//判定吃的棋子是否将军，是的话就赢了
				if wonCode == core.JiangShuai {
					msgChan <- GameMsg{
						Event:           Fin,
						WonChessmanCode: wonCode,
						WonGroup:        game.playerDown.GetGroup(),
						Msg:             "Win",
					}
					return
				}

				//检查两个阵营的将帅是否见面了，见面了判移动的这一方输
				if game.JiangShuaiFace2Face() {
					msgChan <- GameMsg{
						Event:           Fin,
						WonChessmanCode: wonCode,
						WonGroup:        game.playerUp.GetGroup(),
						Msg:             "Win",
					}
					return
				}

				//还没觉出胜负移动完毕，向外部发送消息
				if len(wonCode) > 0 {
					game.playerDown.AddWonChessman(wonCode)
					game.playerUp.DelOwnChessman(wonCode)
					game.playerUp.AddLostChessman(wonCode)
				}

				msgChan <- GameMsg{
					Event:           Done,
					WonChessmanCode: wonCode,
					WonGroup:        core.GroupNone,
					Msg:             "moveDone",
				}
			case 1:
				//先手已经下完了棋，后手下棋
				st, err := game.playerUp.ReceiveStatement(upPlayerCh)
				if err != nil {
					gameMsg := GameMsg{
						Event:           Err,
						WonChessmanCode: "",
						WonGroup:        core.GroupNone,
						Msg:             err.Error(),
					}
					msgChan <- gameMsg
					return
				}
				//移动棋子
				wonCode, err := game.board.MoveChessman(st.Group, st.Code, st.Source, st.Target)
				if err != nil {
					gameMsg := GameMsg{
						Event:           Err,
						WonChessmanCode: "",
						WonGroup:        core.GroupNone,
						Msg:             err.Error(),
					}
					msgChan <- gameMsg
				}

				game.round = -1 //棋子移动完毕先手已经执行

				//判定吃的棋子是否将军，是的话就赢了
				if wonCode == core.JiangShuai {
					msgChan <- GameMsg{
						Event:           Fin,
						WonChessmanCode: wonCode,
						WonGroup:        game.playerUp.GetGroup(),
						Msg:             "Win",
					}
					return
				}

				//检查两个阵营的将帅是否见面了，见面了判移动的这一方输
				if game.JiangShuaiFace2Face() {
					msgChan <- GameMsg{
						Event:           Fin,
						WonChessmanCode: wonCode,
						WonGroup:        game.playerDown.GetGroup(),
						Msg:             "Win",
					}
					return
				}

				//还没觉出胜负移动完毕，向外部发送消息
				if len(wonCode) > 0 {
					game.playerUp.AddWonChessman(wonCode)
					game.playerDown.DelOwnChessman(wonCode)
					game.playerDown.AddLostChessman(wonCode)
				}

				msgChan <- GameMsg{
					Event:           Done,
					WonChessmanCode: wonCode,
					WonGroup:        core.GroupNone,
					Msg:             "moveDone",
				}
			}

		}

	}()
	return msgChan
}

func (game *ChessGame) newChessmen(downGroup, upGroup core.ChessmanGroup) (chessmenOfPlayerDown, chessmenOfPlayerUp []chessman.ChessmanInterface) {
	cm := make([]chessman.ChessmanInterface, 32)

	//ju1 := chessman.NewChessman(core.Ju, "车",downGroup,core.Coordinate{
	//	X: 0,
	//	Y: 0,
	//})

	ma1 := chessman.NewChessman(core.Ma, "马", downGroup, core.Coordinate{
		X: 1,
		Y: 0,
	})
	ma1.BindRule(chessman.RuleMa)
	cm[0] = ma1

	ma2 := chessman.NewChessman(core.Ma, "马", downGroup, core.Coordinate{
		X: 7,
		Y: 0,
	})
	ma2.BindRule(chessman.RuleMa)
	cm[1] = ma2

	//TODO
	return cm, cm
}

func (game *ChessGame) JiangShuaiFace2Face() bool {
	//定位先手方“将/帅”当前位置
	matrix := game.board.GetMatrix()
	downJiangShuaiCoordinate := core.Coordinate{}
	for i := 0; i < 3; i++ {
		row := matrix[i]
		for j, r := range row {
			if r.GetChessmanCode() == core.JiangShuai {
				downJiangShuaiCoordinate.X = j
				downJiangShuaiCoordinate.Y = i
			}
		}
	}

	//定位后手方“将/帅”当前位置
	upJiangShuaiCoordinate := core.Coordinate{}
	for i := 7; i < 10; i++ {
		row := matrix[i]
		for j, r := range row {
			if r.GetChessmanCode() == core.JiangShuai {
				upJiangShuaiCoordinate.X = j
				upJiangShuaiCoordinate.Y = i
			}
		}
	}

	//判定两个棋子横坐标是否相同
	if downJiangShuaiCoordinate.X != upJiangShuaiCoordinate.X {
		return false
	}
	//横坐标相同判定棋子中间是否已经没有棋子
	for i := downJiangShuaiCoordinate.Y; i <= upJiangShuaiCoordinate.Y; i++ {
		if matrix[i][downJiangShuaiCoordinate.X] != nil {
			return false
		}
	}

	return true
}
