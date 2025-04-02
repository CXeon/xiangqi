package chessgame

import (
	"fmt"
	"github.com/CXeon/xiangqi/core"
	"github.com/CXeon/xiangqi/core/chessboard"
	"github.com/CXeon/xiangqi/core/chessman"
	"github.com/CXeon/xiangqi/core/player"
)

type ChessGame struct {
	playerDown     player.PlayerInterface         //玩家1号
	playerUp       player.PlayerInterface         //玩家2号
	board          chessboard.ChessboardInterface //棋盘
	nextRoundGroup core.ChessmanGroup             //下一回合应该哪个阵营下棋

	quit chan struct{} //退出通道，用于随时终止棋局
}

// 初始化棋局
func (game *ChessGame) InitialGame(player1, player2 player.PlayerInterface) error {

	//引入玩家
	if player1.GetIsDown() {
		game.playerDown = player1
		game.playerUp = player2
	} else {
		game.playerDown = player2
		game.playerUp = player1
	}

	//创建棋盘
	board := chessboard.NewChessboard()

	//划分棋盘区域
	if player1.GetIsDown() {
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

	if game.playerDown.GetIsFirst() {
		game.nextRoundGroup = game.playerDown.GetGroup()
	} else {
		game.nextRoundGroup = game.playerUp.GetGroup()
	}

	game.quit = make(chan struct{}, 1) //初始化终止通道

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

	//玩家的座位和阵营可能改变，所以重新划分棋盘区域
	if !game.playerDown.GetIsDown() {
		game.playerDown, game.playerUp = game.playerUp, game.playerDown
	}

	game.board.DivideGroup(game.playerDown.GetGroup(), []int{0, 1, 2, 3, 4})
	game.board.DivideGroup(game.playerUp.GetGroup(), []int{5, 6, 7, 8, 9})

	//创建棋子
	chessmenOfPlayerDown, chessmenOfPlayerUp := game.newChessmen(game.playerDown.GetGroup(), game.playerUp.GetGroup())

	allChessmen := append(chessmenOfPlayerDown, chessmenOfPlayerUp...)

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

	//重置回合标记
	if game.playerDown.GetIsFirst() {
		game.nextRoundGroup = game.playerDown.GetGroup()
	} else {
		game.nextRoundGroup = game.playerUp.GetGroup()
	}

	return nil
}

// 关闭棋局
func (game *ChessGame) Close() error {
	game.playerDown.ClearOwnChessman()
	game.playerDown.ClearLostChessman()
	game.playerDown.ClearWonChessman()

	game.playerUp.ClearOwnChessman()
	game.playerUp.ClearLostChessman()
	game.playerUp.ClearWonChessman()

	game.board = nil

	game.quit <- struct{}{} //发送终止信号，使Run方法退出
	close(game.quit)

	return nil
}

// 运行棋局
func (game *ChessGame) Run(downPlayerCh, upPlayerCh chan player.Statement) (msgChan chan GameMsg) {

	msgChan = make(chan GameMsg, 1)
	go func() {
		defer close(msgChan)
		for {
			switch game.nextRoundGroup {
			case core.Group1:
				var pl player.PlayerInterface
				if game.playerDown.GetGroup() == core.Group1 {
					pl = game.playerDown
				} else {
					pl = game.playerUp
				}
				//棋局开始或者后手已经下完棋，先手下棋
				st, err := pl.ReceiveStatement(downPlayerCh, game.quit)
				if err != nil {

					msgChan <- GameMsg{
						Event:           Err,
						WonChessmanCode: "",
						WonGroup:        core.GroupNone,
						Msg:             err.Error(),
					}
					return
				}
				//移动棋子
				wonCode, err := game.board.MoveChessman(st.Group, st.Code, st.Source, st.Target)
				if err != nil {

					msgChan <- GameMsg{
						Event:           Err,
						WonChessmanCode: "",
						WonGroup:        core.GroupNone,
						Msg:             err.Error(),
					}
					continue
				}

				game.nextRoundGroup = core.Group2 //修改下一回合下棋阵营

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
			case core.Group2:
				var pl player.PlayerInterface
				if game.playerDown.GetGroup() == core.Group2 {
					pl = game.playerDown
				} else {
					pl = game.playerUp
				}
				st, err := pl.ReceiveStatement(upPlayerCh, game.quit)
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
					msgChan <- GameMsg{
						Event:           Err,
						WonChessmanCode: "",
						WonGroup:        core.GroupNone,
						Msg:             err.Error(),
					}

					continue
				}

				game.nextRoundGroup = core.Group1 //棋子移动完毕先手已经执行

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

func (game *ChessGame) Show() {
	matrix := game.board.GetMatrix()

	lenRows := len(matrix)

	for i := lenRows - 1; i >= 0; i-- {
		row := matrix[i]
		lenCols := len(row)
		for j := lenCols - 1; j >= 0; j-- {
			c := row[j]
			if c == nil {
				print("  ")
			} else {
				print(c.GetChessmanName())
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func (game *ChessGame) newChessmen(downGroup, upGroup core.ChessmanGroup) (chessmenOfPlayerDown, chessmenOfPlayerUp []chessman.ChessmanInterface) {

	chessmenOfPlayerDown = make([]chessman.ChessmanInterface, 16)

	ju1 := chessman.NewChessman(core.Ju, "车", downGroup, core.Coordinate{
		X: 0,
		Y: 0,
	})
	ju1.BindRule(chessman.RuleJu)
	chessmenOfPlayerDown[0] = ju1

	ju2 := chessman.NewChessman(core.Ju, "车", downGroup, core.Coordinate{
		X: 8,
		Y: 0,
	})
	ju2.BindRule(chessman.RuleJu)
	chessmenOfPlayerDown[1] = ju2

	ma1 := chessman.NewChessman(core.Ma, "马", downGroup, core.Coordinate{
		X: 1,
		Y: 0,
	})
	ma1.BindRule(chessman.RuleMa)
	chessmenOfPlayerDown[2] = ma1

	ma2 := chessman.NewChessman(core.Ma, "马", downGroup, core.Coordinate{
		X: 7,
		Y: 0,
	})
	ma2.BindRule(chessman.RuleMa)
	chessmenOfPlayerDown[3] = ma2

	xiang1 := chessman.NewChessman(core.Xiang, "象", downGroup, core.Coordinate{
		X: 2,
		Y: 0,
	})
	xiang1.BindRule(chessman.RuleXiang)
	chessmenOfPlayerDown[4] = xiang1

	xiang2 := chessman.NewChessman(core.Xiang, "象", downGroup, core.Coordinate{
		X: 6,
		Y: 0,
	})
	xiang2.BindRule(chessman.RuleXiang)
	chessmenOfPlayerDown[5] = xiang2

	shi1 := chessman.NewChessman(core.Shi, "士", downGroup, core.Coordinate{
		X: 3,
		Y: 0,
	})
	shi1.BindRule(chessman.RuleShi)
	chessmenOfPlayerDown[6] = shi1

	shi2 := chessman.NewChessman(core.Shi, "士", downGroup, core.Coordinate{
		X: 5,
		Y: 0,
	})
	shi2.BindRule(chessman.RuleShi)
	chessmenOfPlayerDown[7] = shi2

	jiangShuai1 := chessman.NewChessman(core.JiangShuai, "帅", downGroup, core.Coordinate{
		X: 4,
		Y: 0,
	})
	jiangShuai1.BindRule(chessman.RuleJiangShuai)
	chessmenOfPlayerDown[8] = jiangShuai1

	pao1 := chessman.NewChessman(core.Pao, "炮", downGroup, core.Coordinate{
		X: 1,
		Y: 2,
	})
	pao1.BindRule(chessman.RulePao)
	chessmenOfPlayerDown[9] = pao1

	pao2 := chessman.NewChessman(core.Pao, "炮", downGroup, core.Coordinate{
		X: 7,
		Y: 2,
	})
	pao2.BindRule(chessman.RulePao)
	chessmenOfPlayerDown[10] = pao2

	bingzu1 := chessman.NewChessman(core.BingZu, "兵", downGroup, core.Coordinate{
		X: 0,
		Y: 3,
	})
	bingzu1.BindRule(chessman.RuleBingZu)
	chessmenOfPlayerDown[11] = bingzu1

	bingzu2 := chessman.NewChessman(core.BingZu, "兵", downGroup, core.Coordinate{
		X: 2,
		Y: 3,
	})
	bingzu2.BindRule(chessman.RuleBingZu)
	chessmenOfPlayerDown[12] = bingzu2

	bingzu3 := chessman.NewChessman(core.BingZu, "兵", downGroup, core.Coordinate{
		X: 4,
		Y: 3,
	})
	bingzu3.BindRule(chessman.RuleBingZu)
	chessmenOfPlayerDown[13] = bingzu3

	bingzu4 := chessman.NewChessman(core.BingZu, "兵", downGroup, core.Coordinate{
		X: 6,
		Y: 3,
	})
	bingzu4.BindRule(chessman.RuleBingZu)
	chessmenOfPlayerDown[14] = bingzu4

	bingzu5 := chessman.NewChessman(core.BingZu, "兵", downGroup, core.Coordinate{
		X: 8,
		Y: 3,
	})
	bingzu5.BindRule(chessman.RuleBingZu)
	chessmenOfPlayerDown[15] = bingzu5

	chessmenOfPlayerUp = make([]chessman.ChessmanInterface, 16)

	ju3 := chessman.NewChessman(core.Ju, "车", upGroup, core.Coordinate{
		X: 0,
		Y: 9,
	})
	ju3.BindRule(chessman.RuleJu)
	chessmenOfPlayerUp[0] = ju3

	ju4 := chessman.NewChessman(core.Ju, "车", upGroup, core.Coordinate{
		X: 8,
		Y: 9,
	})
	ju4.BindRule(chessman.RuleJu)
	chessmenOfPlayerUp[1] = ju4

	ma3 := chessman.NewChessman(core.Ma, "马", upGroup, core.Coordinate{
		X: 1,
		Y: 9,
	})
	ma3.BindRule(chessman.RuleMa)
	chessmenOfPlayerUp[2] = ma3

	ma4 := chessman.NewChessman(core.Ma, "马", upGroup, core.Coordinate{
		X: 7,
		Y: 9,
	})
	ma4.BindRule(chessman.RuleMa)
	chessmenOfPlayerUp[3] = ma4

	xiang3 := chessman.NewChessman(core.Xiang, "象", upGroup, core.Coordinate{
		X: 2,
		Y: 9,
	})
	xiang3.BindRule(chessman.RuleXiang)
	chessmenOfPlayerUp[4] = xiang3

	xiang4 := chessman.NewChessman(core.Xiang, "象", upGroup, core.Coordinate{
		X: 6,
		Y: 9,
	})
	xiang4.BindRule(chessman.RuleXiang)
	chessmenOfPlayerUp[5] = xiang4

	shi3 := chessman.NewChessman(core.Shi, "士", upGroup, core.Coordinate{
		X: 3,
		Y: 9,
	})
	shi3.BindRule(chessman.RuleShi)
	chessmenOfPlayerUp[6] = shi3

	shi4 := chessman.NewChessman(core.Shi, "士", upGroup, core.Coordinate{
		X: 5,
		Y: 9,
	})
	shi4.BindRule(chessman.RuleShi)
	chessmenOfPlayerUp[7] = shi4

	jiangShuai2 := chessman.NewChessman(core.JiangShuai, "将", upGroup, core.Coordinate{
		X: 4,
		Y: 9,
	})
	jiangShuai2.BindRule(chessman.RuleJiangShuai)
	chessmenOfPlayerUp[8] = jiangShuai2

	pao3 := chessman.NewChessman(core.Pao, "炮", upGroup, core.Coordinate{
		X: 1,
		Y: 7,
	})
	pao3.BindRule(chessman.RulePao)
	chessmenOfPlayerUp[9] = pao3

	pao4 := chessman.NewChessman(core.Pao, "炮", upGroup, core.Coordinate{
		X: 7,
		Y: 7,
	})
	pao4.BindRule(chessman.RulePao)
	chessmenOfPlayerUp[10] = pao4

	bingzu6 := chessman.NewChessman(core.BingZu, "卒", upGroup, core.Coordinate{
		X: 0,
		Y: 6,
	})
	bingzu6.BindRule(chessman.RuleBingZu)
	chessmenOfPlayerUp[11] = bingzu6

	bingzu7 := chessman.NewChessman(core.BingZu, "卒", upGroup, core.Coordinate{
		X: 2,
		Y: 6,
	})
	bingzu7.BindRule(chessman.RuleBingZu)
	chessmenOfPlayerUp[12] = bingzu7

	bingzu8 := chessman.NewChessman(core.BingZu, "卒", upGroup, core.Coordinate{
		X: 4,
		Y: 6,
	})
	bingzu8.BindRule(chessman.RuleBingZu)
	chessmenOfPlayerUp[13] = bingzu8

	bingzu9 := chessman.NewChessman(core.BingZu, "卒", upGroup, core.Coordinate{
		X: 6,
		Y: 6,
	})
	bingzu9.BindRule(chessman.RuleBingZu)
	chessmenOfPlayerUp[14] = bingzu9

	bingzu10 := chessman.NewChessman(core.BingZu, "卒", upGroup, core.Coordinate{
		X: 8,
		Y: 6,
	})
	bingzu10.BindRule(chessman.RuleBingZu)
	chessmenOfPlayerUp[15] = bingzu10

	return chessmenOfPlayerDown, chessmenOfPlayerUp
}

func (game *ChessGame) JiangShuaiFace2Face() bool {
	//定位先手方“将/帅”当前位置
	matrix := game.board.GetMatrix()
	downJiangShuaiCoordinate := core.Coordinate{}
	for i := 0; i < 3; i++ {
		row := matrix[i]
		for j, r := range row {
			if r == nil {
				continue
			}
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
			if r == nil {
				continue
			}
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
	for i := downJiangShuaiCoordinate.Y + 1; i < upJiangShuaiCoordinate.Y; i++ {
		if matrix[i][downJiangShuaiCoordinate.X] != nil {
			return false
		}
	}

	return true
}
