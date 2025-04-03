package app

import (
	"fmt"
	"github.com/CXeon/xiangqi/core"
	"github.com/CXeon/xiangqi/core/chessgame"
	"github.com/CXeon/xiangqi/core/player"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"log"
)

type Game struct {
	//界面相关
	sprites             []*Sprite
	boardLogicZeroPoint coordinate //棋盘逻辑上的0点处于游戏界面的哪个坐标
	gridLength          int        //棋盘格长宽所占像素
	spriteReparation    int        //设置sprite坐标时，x,y坐标的补偿值

	gameMsg *chessgame.GameMsg //如果不为nil，需要显示在屏幕消息区

	winner       string        //游戏胜利者记录
	onceAgainBtn *onceAgainBtn //再来一次按钮

	//业务逻辑相关
	player1 player.PlayerInterface //玩家1 在棋盘下方
	player2 player.PlayerInterface //玩家2 在棋盘上方

	p1Ch chan player.Statement //玩家1的下棋意图
	p2Ch chan player.Statement //玩家2的下棋意图

	clickedSprite  *Sprite            //当前棋盘被选中的棋子
	nextRoundGroup core.ChessmanGroup //下一回合应该下棋的阵营

	gameCore chessgame.ChessGameInterface //游戏内核
	coreCh   chan chessgame.GameMsg       //管道接收内核返回的消息

}

func NewGame() *Game {

	var p1, p2 player.PlayerInterface

	p1 = player.NewPlayer()
	p1.SetIsFirst(true) //默认p1先手
	p1.SetIsDown(true)  //默认p1在棋盘下方
	p1.SetGroup(core.Group1)

	p2 = player.NewPlayer()

	p2.SetGroup(core.Group2)

	p1Ch := make(chan player.Statement, 1)
	p2Ch := make(chan player.Statement, 1)

	//先手执红棋，根据先手创建棋子

	g := &Game{
		sprites:             nil,
		boardLogicZeroPoint: coordinate{x: boardLogicZeroX, y: boardLogicZeroY},
		gridLength:          gridLength,
		spriteReparation:    spriteReparation,
		player1:             p1,
		player2:             p2,
		p1Ch:                p1Ch,
		p2Ch:                p2Ch,
		gameCore:            nil,
		coreCh:              nil,
	}

	g.initSprites(p1, p2)
	if p1.GetIsFirst() {
		g.nextRoundGroup = p1.GetGroup()
	} else {
		g.nextRoundGroup = p2.GetGroup()
	}

	g.winner = ""
	g.onceAgainBtn = &onceAgainBtn{
		image:      ebitenOnceAgainBtnImage,
		alphaImage: ebitenOnceAgainBtnAlphaImage,
		x:          g.boardLogicZeroPoint.x + 3*g.gridLength,
		y:          g.boardLogicZeroPoint.y + 7*g.gridLength,
	}

	//启动内核
	g.gameCore = new(chessgame.ChessGame)
	err := g.gameCore.InitialGame(p1, p2)
	if err != nil {
		log.Fatal(err)
	}
	g.coreCh = g.gameCore.Run(p1Ch, p2Ch)

	return g
}

func (g *Game) Update() error {

	//如果发生鼠标左键点击事件
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.gameMsg = nil

		//如果已经产生胜者，只需要判断是否点击了再来一局按钮
		if len(g.winner) > 0 {
			if g.onceAgainBtn.In(ebiten.CursorPosition()) {
				//重置棋局

				//切换先手
				if g.player1.GetIsFirst() {
					g.player2.SetIsFirst(true)
					g.player1.SetIsFirst(false)
				} else {
					g.player1.SetIsFirst(true)
					g.player2.SetIsFirst(false)
				}

				//重新摆棋
				g.initSprites(g.player1, g.player2)
				if g.player1.GetIsFirst() {
					g.nextRoundGroup = g.player1.GetGroup()
				} else {
					g.nextRoundGroup = g.player2.GetGroup()
				}

				//内核也要重置棋局
				err := g.gameCore.ResetGame()
				if err != nil {
					return err
				}
				g.coreCh = g.gameCore.Run(g.p1Ch, g.p2Ch)

				//重置获胜记录和其他信息
				g.winner = ""
				g.gameMsg = nil
				g.clickedSprite = nil

			}

			return nil
		}

		//是否点击了棋子
		if sp := g.spriteAt(ebiten.CursorPosition()); sp != nil {
			//如果选中的棋子就是当前应该下棋的阵营
			if sp.group == g.nextRoundGroup {
				//将棋子设置为点击状态并记录
				if g.clickedSprite != nil {
					g.clickedSprite.clicked = false
				}

				sp.clicked = true
				g.clickedSprite = sp
			}
			if sp.group != g.nextRoundGroup {
				//如果之前没有棋子被选中，就不记录。如果之前已经有棋子被选中说明是要吃棋
				if g.clickedSprite != nil {
					g.moveSpriteToFront(sp)

					// 将操作传给核心层校验，然后移动被选中的棋子到目标坐标
					sourceCoreX, sourceCoreY := g.transformCoordinate(g.clickedSprite.x-g.spriteReparation, g.clickedSprite.y-g.spriteReparation)
					targetCoreX, targetCoreY := g.transformCoordinate(sp.x-g.spriteReparation, sp.y-g.spriteReparation)
					if g.player1.GetGroup() == g.nextRoundGroup {

						g.p1Ch <- player.Statement{
							Group: g.player1.GetGroup(),
							Code:  g.clickedSprite.code,
							Source: core.Coordinate{
								X: sourceCoreX,
								Y: sourceCoreY,
							},
							Target: core.Coordinate{
								X: targetCoreX,
								Y: targetCoreY,
							},
						}
					} else {
						g.p2Ch <- player.Statement{
							Group: g.player2.GetGroup(),
							Code:  g.clickedSprite.code,
							Source: core.Coordinate{
								X: sourceCoreX,
								Y: sourceCoreY,
							},
							Target: core.Coordinate{
								X: targetCoreX,
								Y: targetCoreY,
							},
						}
					}

					//接收回复
					msg := <-g.coreCh
					//fmt.Println(msg)
					g.gameMsg = &msg

					if msg.Event == chessgame.Done || msg.Event == chessgame.Fin {
						//校验通过，移动游戏界面的棋子
						g.clickedSprite.MoveTo(sp.x, sp.y)
						//修改坐标
						g.clickedSprite.x = sp.x
						g.clickedSprite.y = sp.y

						//删除原来在这个坐标的棋子
						g.deleteSprite(sp)
						//重置棋子选中状态和修改下一回合标记
						g.clickedSprite.clicked = false
						g.clickedSprite = nil
						if g.nextRoundGroup == g.player1.GetGroup() {
							g.nextRoundGroup = g.player2.GetGroup()
						} else {
							g.nextRoundGroup = g.player1.GetGroup()
						}

						//TODO 如果出现赢家
						if msg.Event == chessgame.Fin {
							//获取胜利阵营
							group := msg.WonGroup
							//判断玩家哪个属于这个阵营
							pl := g.getPlayerByGroup(group)
							if pl.GetIsFirst() {
								g.winner = "红方"
							} else {
								g.winner = "黑方"
							}

						}
					}
				}
			}

		} else {
			//没有点击到棋子。判断是否点击的是棋盘格
			if !g.InBoard(ebiten.CursorPosition()) {
				return nil
			}
			//修正坐标
			x, y := g.revisesCoordinate(ebiten.CursorPosition())
			fmt.Printf("revisesCoordinate to [%d,%d]\n", x, y)
			coreX, coreY := g.transformCoordinate(x, y)
			//如果已经有棋子被选中，说明玩家想把棋子移动到选中坐标
			if g.clickedSprite == nil {
				return nil
			}

			if g.clickedSprite.group != g.nextRoundGroup {
				return nil
			}
			g.moveSpriteToFront(g.clickedSprite)
			// 将操作传给核心层校验，然后移动被选中的棋子到目标坐标
			sourceCoreX, sourceCoreY := g.transformCoordinate(g.clickedSprite.x-g.spriteReparation, g.clickedSprite.y-g.spriteReparation)
			if g.player1.GetGroup() == g.nextRoundGroup {
				g.p1Ch <- player.Statement{
					Group: g.player1.GetGroup(),
					Code:  g.clickedSprite.code,
					Source: core.Coordinate{
						X: sourceCoreX,
						Y: sourceCoreY,
					},
					Target: core.Coordinate{
						X: coreX,
						Y: coreY,
					},
				}
			} else {
				g.p2Ch <- player.Statement{
					Group: g.player2.GetGroup(),
					Code:  g.clickedSprite.code,
					Source: core.Coordinate{
						X: sourceCoreX,
						Y: sourceCoreY,
					},
					Target: core.Coordinate{
						X: coreX,
						Y: coreY,
					},
				}
			}

			//接收回复
			msg := <-g.coreCh
			//fmt.Println(msg)
			g.gameMsg = &msg

			if msg.Event == chessgame.Done || msg.Event == chessgame.Fin {
				//校验通过，移动游戏界面的棋子
				g.clickedSprite.MoveTo(x+g.spriteReparation, y+g.spriteReparation)

				//重置棋子选中状态和修改下一回合标记
				g.clickedSprite.clicked = false
				g.clickedSprite = nil
				if g.nextRoundGroup == g.player1.GetGroup() {
					g.nextRoundGroup = g.player2.GetGroup()
				} else {
					g.nextRoundGroup = g.player1.GetGroup()
				}

				if msg.Event == chessgame.Fin {
					//获取胜利阵营
					group := msg.WonGroup
					//判断玩家哪个属于这个阵营
					pl := g.getPlayerByGroup(group)
					if pl.GetIsFirst() {
						g.winner = "红方"
					} else {
						g.winner = "黑方"
					}

				}

			}

		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//绘制棋盘
	screen.DrawImage(ebitenBoardImage, &ebiten.DrawImageOptions{})

	//绘制棋子
	for _, s := range g.sprites {
		if s.clicked {
			s.Draw(screen, 0.7)
		} else {
			s.Draw(screen, 1)
		}
	}

	g.ShowGameMsg(screen)

	if len(g.winner) > 0 {
		g.drawWinner(screen)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Close() {
	close(g.p1Ch)
	close(g.p2Ch)
	g.gameCore.Close()
}

// 初始化棋盘上各个棋子的精灵：先手执红棋，并且根据玩家意愿确定坐在那一方
func (g *Game) initSprites(p1, p2 player.PlayerInterface) {
	g.sprites = make([]*Sprite, 32)

	if p1.GetIsFirst() {
		//确定p1执红棋
		if p1.GetIsDown() {
			//红棋在棋盘下边儿
			g.initSpritesP1IsFirstAndDown(p1, p2)
		} else {
			//红棋在棋盘上边儿
			g.initSpritesP1IsFirstAndUp(p1, p2)
		}

		return
	}

	//p2先手执红棋
	if p2.GetIsDown() {
		//p2在棋盘下边儿
		g.initSpritesP2IsFirstAndDown(p1, p2)
	} else {
		//p2在棋盘上边儿
		g.initSpritesP2IsFirstAndUp(p1, p2)

	}

	return
}

func (g *Game) initSpritesP1IsFirstAndDown(p1, p2 player.PlayerInterface) {

	//p1先手执红棋
	g.sprites[0] = &Sprite{
		image:      ebitenBlackJuImage,
		alphaImage: ebitenBlackJuAlphaImage,
		x:          g.boardLogicZeroPoint.x + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Ju,
	}
	g.sprites[1] = &Sprite{
		image:      ebitenBlackMaImage,
		alphaImage: ebitenBlackMaAlphaImage,
		x:          g.boardLogicZeroPoint.x + g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Ma,
	}
	g.sprites[2] = &Sprite{
		image:      ebitenBlackXiangImage,
		alphaImage: ebitenBlackXiangAlphaImage,
		x:          g.boardLogicZeroPoint.x + 2*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Xiang,
	}
	g.sprites[3] = &Sprite{
		image:      ebitenBlackShiImage,
		alphaImage: ebitenBlackShiAlphaImage,
		x:          g.boardLogicZeroPoint.x + 3*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Shi,
	}
	g.sprites[4] = &Sprite{
		image:      ebitenBlackJiangImage,
		alphaImage: ebitenBlackJiangAlphaImage,
		x:          g.boardLogicZeroPoint.x + 4*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.JiangShuai,
	}
	g.sprites[5] = &Sprite{
		image:      ebitenBlackShiImage,
		alphaImage: ebitenBlackShiAlphaImage,
		x:          g.boardLogicZeroPoint.x + 5*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Shi,
	}
	g.sprites[6] = &Sprite{
		image:      ebitenBlackXiangImage,
		alphaImage: ebitenBlackXiangAlphaImage,
		x:          g.boardLogicZeroPoint.x + 6*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Xiang,
	}
	g.sprites[7] = &Sprite{
		image:      ebitenBlackMaImage,
		alphaImage: ebitenBlackMaAlphaImage,
		x:          g.boardLogicZeroPoint.x + 7*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Ma,
	}
	g.sprites[8] = &Sprite{
		image:      ebitenBlackJuImage,
		alphaImage: ebitenBlackJuAlphaImage,
		x:          g.boardLogicZeroPoint.x + 8*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Ju,
	}
	g.sprites[9] = &Sprite{
		image:      ebitenBlackPaoImage,
		alphaImage: ebitenBlackPaoAlphaImage,
		x:          g.boardLogicZeroPoint.x + g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 2*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Pao,
	}
	g.sprites[10] = &Sprite{
		image:      ebitenBlackPaoImage,
		alphaImage: ebitenBlackPaoAlphaImage,
		x:          g.boardLogicZeroPoint.x + 7*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 2*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Pao,
	}
	g.sprites[11] = &Sprite{
		image:      ebitenBlackZuImage,
		alphaImage: ebitenBlackZuAlphaImage,
		x:          g.boardLogicZeroPoint.x + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 3*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.BingZu,
	}
	g.sprites[12] = &Sprite{
		image:      ebitenBlackZuImage,
		alphaImage: ebitenBlackZuAlphaImage,
		x:          g.boardLogicZeroPoint.x + 2*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 3*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.BingZu,
	}
	g.sprites[13] = &Sprite{
		image:      ebitenBlackZuImage,
		alphaImage: ebitenBlackZuAlphaImage,
		x:          g.boardLogicZeroPoint.x + 4*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 3*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.BingZu,
	}
	g.sprites[14] = &Sprite{
		image:      ebitenBlackZuImage,
		alphaImage: ebitenBlackZuAlphaImage,
		x:          g.boardLogicZeroPoint.x + 6*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 3*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.BingZu,
	}
	g.sprites[15] = &Sprite{
		image:      ebitenBlackZuImage,
		alphaImage: ebitenBlackZuAlphaImage,
		x:          g.boardLogicZeroPoint.x + 8*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 3*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.BingZu,
	}

	g.sprites[16] = &Sprite{
		image:      ebitenRedBingImage,
		alphaImage: ebitenRedBingAlphaImage,
		x:          g.boardLogicZeroPoint.x + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 6*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.BingZu,
	}

	g.sprites[17] = &Sprite{
		image:      ebitenRedBingImage,
		alphaImage: ebitenRedBingAlphaImage,
		x:          g.boardLogicZeroPoint.x + 2*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 6*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.BingZu,
	}

	g.sprites[18] = &Sprite{
		image:      ebitenRedBingImage,
		alphaImage: ebitenRedBingAlphaImage,
		x:          g.boardLogicZeroPoint.x + 4*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 6*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.BingZu,
	}

	g.sprites[19] = &Sprite{
		image:      ebitenRedBingImage,
		alphaImage: ebitenRedBingAlphaImage,
		x:          g.boardLogicZeroPoint.x + 6*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 6*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.BingZu,
	}

	g.sprites[20] = &Sprite{
		image:      ebitenRedBingImage,
		alphaImage: ebitenRedBingAlphaImage,
		x:          g.boardLogicZeroPoint.x + 8*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 6*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.BingZu,
	}

	g.sprites[21] = &Sprite{
		image:      ebitenRedPaoImage,
		alphaImage: ebitenRedPaoAlphaImage,
		x:          g.boardLogicZeroPoint.x + g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 7*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Pao,
	}

	g.sprites[22] = &Sprite{
		image:      ebitenRedPaoImage,
		alphaImage: ebitenRedPaoAlphaImage,
		x:          g.boardLogicZeroPoint.x + 7*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 7*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Pao,
	}

	g.sprites[23] = &Sprite{
		image:      ebitenRedJuImage,
		alphaImage: ebitenRedJuAlphaImage,
		x:          g.boardLogicZeroPoint.x + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Ju,
	}
	g.sprites[24] = &Sprite{
		image:      ebitenRedMaImage,
		alphaImage: ebitenRedMaAlphaImage,
		x:          g.boardLogicZeroPoint.x + 1*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Ma,
	}
	g.sprites[25] = &Sprite{
		image:      ebitenRedXiangImage,
		alphaImage: ebitenRedXiangAlphaImage,
		x:          g.boardLogicZeroPoint.x + 2*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Xiang,
	}
	g.sprites[26] = &Sprite{
		image:      ebitenRedShiImage,
		alphaImage: ebitenRedShiAlphaImage,
		x:          g.boardLogicZeroPoint.x + 3*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Shi,
	}
	g.sprites[27] = &Sprite{
		image:      ebitenRedShuaiImage,
		alphaImage: ebitenRedShuaiAlphaImage,
		x:          g.boardLogicZeroPoint.x + 4*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.JiangShuai,
	}
	g.sprites[28] = &Sprite{
		image:      ebitenRedShiImage,
		alphaImage: ebitenRedShiAlphaImage,
		x:          g.boardLogicZeroPoint.x + 5*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Shi,
	}

	g.sprites[29] = &Sprite{
		image:      ebitenRedXiangImage,
		alphaImage: ebitenRedXiangAlphaImage,
		x:          g.boardLogicZeroPoint.x + 6*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Xiang,
	}
	g.sprites[30] = &Sprite{
		image:      ebitenRedMaImage,
		alphaImage: ebitenRedMaAlphaImage,
		x:          g.boardLogicZeroPoint.x + 7*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Ma,
	}
	g.sprites[31] = &Sprite{
		image:      ebitenRedJuImage,
		alphaImage: ebitenRedJuAlphaImage,
		x:          g.boardLogicZeroPoint.x + 8*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Ju,
	}

}

func (g *Game) initSpritesP1IsFirstAndUp(p1, p2 player.PlayerInterface) {

	g.sprites[0] = &Sprite{
		image:      ebitenRedJuImage,
		alphaImage: ebitenRedJuAlphaImage,
		x:          g.boardLogicZeroPoint.x + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Ju,
	}
	g.sprites[1] = &Sprite{
		image:      ebitenRedMaImage,
		alphaImage: ebitenRedMaAlphaImage,
		x:          g.boardLogicZeroPoint.x + g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Ma,
	}
	g.sprites[2] = &Sprite{
		image:      ebitenRedXiangImage,
		alphaImage: ebitenRedXiangAlphaImage,
		x:          g.boardLogicZeroPoint.x + 2*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Xiang,
	}
	g.sprites[3] = &Sprite{
		image:      ebitenRedShiImage,
		alphaImage: ebitenRedShiAlphaImage,
		x:          g.boardLogicZeroPoint.x + 3*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Shi,
	}
	g.sprites[4] = &Sprite{
		image:      ebitenRedShuaiImage,
		alphaImage: ebitenRedShuaiAlphaImage,
		x:          g.boardLogicZeroPoint.x + 4*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.JiangShuai,
	}
	g.sprites[5] = &Sprite{
		image:      ebitenRedShiImage,
		alphaImage: ebitenRedShiAlphaImage,
		x:          g.boardLogicZeroPoint.x + 5*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Shi,
	}
	g.sprites[6] = &Sprite{
		image:      ebitenRedXiangImage,
		alphaImage: ebitenRedXiangAlphaImage,
		x:          g.boardLogicZeroPoint.x + 6*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Xiang,
	}
	g.sprites[7] = &Sprite{
		image:      ebitenRedMaImage,
		alphaImage: ebitenRedMaAlphaImage,
		x:          g.boardLogicZeroPoint.x + 7*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Ma,
	}
	g.sprites[8] = &Sprite{
		image:      ebitenRedJuImage,
		alphaImage: ebitenRedJuAlphaImage,
		x:          g.boardLogicZeroPoint.x + 8*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Ju,
	}
	g.sprites[9] = &Sprite{
		image:      ebitenRedPaoImage,
		alphaImage: ebitenRedPaoAlphaImage,
		x:          g.boardLogicZeroPoint.x + g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 2*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Pao,
	}
	g.sprites[10] = &Sprite{
		image:      ebitenRedPaoImage,
		alphaImage: ebitenRedPaoAlphaImage,
		x:          g.boardLogicZeroPoint.x + 7*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 2*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Pao,
	}
	g.sprites[11] = &Sprite{
		image:      ebitenRedBingImage,
		alphaImage: ebitenRedBingAlphaImage,
		x:          g.boardLogicZeroPoint.x + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 3*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.BingZu,
	}
	g.sprites[12] = &Sprite{
		image:      ebitenRedBingImage,
		alphaImage: ebitenRedBingAlphaImage,
		x:          g.boardLogicZeroPoint.x + 2*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 3*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.BingZu,
	}
	g.sprites[13] = &Sprite{
		image:      ebitenRedBingImage,
		alphaImage: ebitenRedBingAlphaImage,
		x:          g.boardLogicZeroPoint.x + 4*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 3*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.BingZu,
	}
	g.sprites[14] = &Sprite{
		image:      ebitenRedBingImage,
		alphaImage: ebitenRedBingAlphaImage,
		x:          g.boardLogicZeroPoint.x + 6*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 3*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.BingZu,
	}
	g.sprites[15] = &Sprite{
		image:      ebitenRedBingImage,
		alphaImage: ebitenRedBingAlphaImage,
		x:          g.boardLogicZeroPoint.x + 8*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 3*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.BingZu,
	}

	g.sprites[16] = &Sprite{
		image:      ebitenBlackZuImage,
		alphaImage: ebitenBlackZuAlphaImage,
		x:          g.boardLogicZeroPoint.x + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 6*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.BingZu,
	}

	g.sprites[17] = &Sprite{
		image:      ebitenBlackZuImage,
		alphaImage: ebitenBlackZuAlphaImage,
		x:          g.boardLogicZeroPoint.x + 2*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 6*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.BingZu,
	}

	g.sprites[18] = &Sprite{
		image:      ebitenBlackZuImage,
		alphaImage: ebitenBlackZuAlphaImage,
		x:          g.boardLogicZeroPoint.x + 4*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 6*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.BingZu,
	}

	g.sprites[19] = &Sprite{
		image:      ebitenBlackZuImage,
		alphaImage: ebitenBlackZuAlphaImage,
		x:          g.boardLogicZeroPoint.x + 6*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 6*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.BingZu,
	}

	g.sprites[20] = &Sprite{
		image:      ebitenBlackZuImage,
		alphaImage: ebitenBlackZuAlphaImage,
		x:          g.boardLogicZeroPoint.x + 8*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 6*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.BingZu,
	}

	g.sprites[21] = &Sprite{
		image:      ebitenBlackPaoImage,
		alphaImage: ebitenBlackPaoAlphaImage,
		x:          g.boardLogicZeroPoint.x + g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 7*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Pao,
	}

	g.sprites[22] = &Sprite{
		image:      ebitenBlackPaoImage,
		alphaImage: ebitenBlackPaoAlphaImage,
		x:          g.boardLogicZeroPoint.x + 7*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 7*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Pao,
	}

	g.sprites[23] = &Sprite{
		image:      ebitenBlackJuImage,
		alphaImage: ebitenBlackJuAlphaImage,
		x:          g.boardLogicZeroPoint.x + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Ju,
	}
	g.sprites[24] = &Sprite{
		image:      ebitenBlackMaImage,
		alphaImage: ebitenBlackMaAlphaImage,
		x:          g.boardLogicZeroPoint.x + 1*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Ma,
	}
	g.sprites[25] = &Sprite{
		image:      ebitenBlackXiangImage,
		alphaImage: ebitenBlackXiangAlphaImage,
		x:          g.boardLogicZeroPoint.x + 2*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Xiang,
	}
	g.sprites[26] = &Sprite{
		image:      ebitenBlackShiImage,
		alphaImage: ebitenBlackShiAlphaImage,
		x:          g.boardLogicZeroPoint.x + 3*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Shi,
	}
	g.sprites[27] = &Sprite{
		image:      ebitenBlackJiangImage,
		alphaImage: ebitenBlackJiangAlphaImage,
		x:          g.boardLogicZeroPoint.x + 4*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.JiangShuai,
	}
	g.sprites[28] = &Sprite{
		image:      ebitenBlackShiImage,
		alphaImage: ebitenBlackShiAlphaImage,
		x:          g.boardLogicZeroPoint.x + 5*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Shi,
	}

	g.sprites[29] = &Sprite{
		image:      ebitenBlackXiangImage,
		alphaImage: ebitenBlackXiangAlphaImage,
		x:          g.boardLogicZeroPoint.x + 6*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Xiang,
	}
	g.sprites[30] = &Sprite{
		image:      ebitenBlackMaImage,
		alphaImage: ebitenBlackMaAlphaImage,
		x:          g.boardLogicZeroPoint.x + 7*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Ma,
	}
	g.sprites[31] = &Sprite{
		image:      ebitenBlackJuImage,
		alphaImage: ebitenBlackJuAlphaImage,
		x:          g.boardLogicZeroPoint.x + 8*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Ju,
	}

}

func (g *Game) initSpritesP2IsFirstAndUp(p1, p2 player.PlayerInterface) {

	//p2先手执红棋
	g.sprites[0] = &Sprite{
		image:      ebitenRedJuImage,
		alphaImage: ebitenRedJuAlphaImage,
		x:          g.boardLogicZeroPoint.x + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Ju,
	}
	g.sprites[1] = &Sprite{
		image:      ebitenRedMaImage,
		alphaImage: ebitenRedMaAlphaImage,
		x:          g.boardLogicZeroPoint.x + g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Ma,
	}
	g.sprites[2] = &Sprite{
		image:      ebitenRedXiangImage,
		alphaImage: ebitenRedXiangAlphaImage,
		x:          g.boardLogicZeroPoint.x + 2*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Xiang,
	}
	g.sprites[3] = &Sprite{
		image:      ebitenRedShiImage,
		alphaImage: ebitenRedShiAlphaImage,
		x:          g.boardLogicZeroPoint.x + 3*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Shi,
	}
	g.sprites[4] = &Sprite{
		image:      ebitenRedShuaiImage,
		alphaImage: ebitenRedShuaiAlphaImage,
		x:          g.boardLogicZeroPoint.x + 4*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.JiangShuai,
	}
	g.sprites[5] = &Sprite{
		image:      ebitenRedShiImage,
		alphaImage: ebitenRedShiAlphaImage,
		x:          g.boardLogicZeroPoint.x + 5*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Shi,
	}
	g.sprites[6] = &Sprite{
		image:      ebitenRedXiangImage,
		alphaImage: ebitenRedXiangAlphaImage,
		x:          g.boardLogicZeroPoint.x + 6*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Xiang,
	}
	g.sprites[7] = &Sprite{
		image:      ebitenRedMaImage,
		alphaImage: ebitenRedMaAlphaImage,
		x:          g.boardLogicZeroPoint.x + 7*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Ma,
	}
	g.sprites[8] = &Sprite{
		image:      ebitenRedJuImage,
		alphaImage: ebitenRedJuAlphaImage,
		x:          g.boardLogicZeroPoint.x + 8*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Ju,
	}
	g.sprites[9] = &Sprite{
		image:      ebitenRedPaoImage,
		alphaImage: ebitenRedPaoAlphaImage,
		x:          g.boardLogicZeroPoint.x + g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 2*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Pao,
	}
	g.sprites[10] = &Sprite{
		image:      ebitenRedPaoImage,
		alphaImage: ebitenRedPaoAlphaImage,
		x:          g.boardLogicZeroPoint.x + 7*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 2*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Pao,
	}
	g.sprites[11] = &Sprite{
		image:      ebitenRedBingImage,
		alphaImage: ebitenRedBingAlphaImage,
		x:          g.boardLogicZeroPoint.x + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 3*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.BingZu,
	}
	g.sprites[12] = &Sprite{
		image:      ebitenRedBingImage,
		alphaImage: ebitenRedBingAlphaImage,
		x:          g.boardLogicZeroPoint.x + 2*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 3*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.BingZu,
	}
	g.sprites[13] = &Sprite{
		image:      ebitenRedBingImage,
		alphaImage: ebitenRedBingAlphaImage,
		x:          g.boardLogicZeroPoint.x + 4*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 3*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.BingZu,
	}
	g.sprites[14] = &Sprite{
		image:      ebitenRedBingImage,
		alphaImage: ebitenRedBingAlphaImage,
		x:          g.boardLogicZeroPoint.x + 6*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 3*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.BingZu,
	}
	g.sprites[15] = &Sprite{
		image:      ebitenRedBingImage,
		alphaImage: ebitenRedBingAlphaImage,
		x:          g.boardLogicZeroPoint.x + 8*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 3*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.BingZu,
	}

	g.sprites[16] = &Sprite{
		image:      ebitenBlackZuImage,
		alphaImage: ebitenBlackZuAlphaImage,
		x:          g.boardLogicZeroPoint.x + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 6*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.BingZu,
	}

	g.sprites[17] = &Sprite{
		image:      ebitenBlackZuImage,
		alphaImage: ebitenBlackZuAlphaImage,
		x:          g.boardLogicZeroPoint.x + 2*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 6*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.BingZu,
	}

	g.sprites[18] = &Sprite{
		image:      ebitenBlackZuImage,
		alphaImage: ebitenBlackZuAlphaImage,
		x:          g.boardLogicZeroPoint.x + 4*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 6*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.BingZu,
	}

	g.sprites[19] = &Sprite{
		image:      ebitenBlackZuImage,
		alphaImage: ebitenBlackZuAlphaImage,
		x:          g.boardLogicZeroPoint.x + 6*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 6*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.BingZu,
	}

	g.sprites[20] = &Sprite{
		image:      ebitenBlackZuImage,
		alphaImage: ebitenBlackZuAlphaImage,
		x:          g.boardLogicZeroPoint.x + 8*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 6*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.BingZu,
	}

	g.sprites[21] = &Sprite{
		image:      ebitenBlackPaoImage,
		alphaImage: ebitenBlackPaoAlphaImage,
		x:          g.boardLogicZeroPoint.x + g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 7*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Pao,
	}

	g.sprites[22] = &Sprite{
		image:      ebitenBlackPaoImage,
		alphaImage: ebitenBlackPaoAlphaImage,
		x:          g.boardLogicZeroPoint.x + 7*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 7*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Pao,
	}

	g.sprites[23] = &Sprite{
		image:      ebitenBlackJuImage,
		alphaImage: ebitenBlackJuAlphaImage,
		x:          g.boardLogicZeroPoint.x + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Ju,
	}
	g.sprites[24] = &Sprite{
		image:      ebitenBlackMaImage,
		alphaImage: ebitenBlackMaAlphaImage,
		x:          g.boardLogicZeroPoint.x + 1*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Ma,
	}
	g.sprites[25] = &Sprite{
		image:      ebitenBlackXiangImage,
		alphaImage: ebitenBlackXiangAlphaImage,
		x:          g.boardLogicZeroPoint.x + 2*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Xiang,
	}
	g.sprites[26] = &Sprite{
		image:      ebitenBlackShiImage,
		alphaImage: ebitenBlackShiAlphaImage,
		x:          g.boardLogicZeroPoint.x + 3*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Shi,
	}
	g.sprites[27] = &Sprite{
		image:      ebitenBlackJiangImage,
		alphaImage: ebitenBlackJiangAlphaImage,
		x:          g.boardLogicZeroPoint.x + 4*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.JiangShuai,
	}
	g.sprites[28] = &Sprite{
		image:      ebitenBlackShiImage,
		alphaImage: ebitenBlackShiAlphaImage,
		x:          g.boardLogicZeroPoint.x + 5*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Shi,
	}

	g.sprites[29] = &Sprite{
		image:      ebitenBlackXiangImage,
		alphaImage: ebitenBlackXiangAlphaImage,
		x:          g.boardLogicZeroPoint.x + 6*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Xiang,
	}
	g.sprites[30] = &Sprite{
		image:      ebitenBlackMaImage,
		alphaImage: ebitenBlackMaAlphaImage,
		x:          g.boardLogicZeroPoint.x + 7*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Ma,
	}
	g.sprites[31] = &Sprite{
		image:      ebitenBlackJuImage,
		alphaImage: ebitenBlackJuAlphaImage,
		x:          g.boardLogicZeroPoint.x + 8*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Ju,
	}

}

func (g *Game) initSpritesP2IsFirstAndDown(p1, p2 player.PlayerInterface) {

	g.sprites[0] = &Sprite{
		image:      ebitenBlackJuImage,
		alphaImage: ebitenBlackJuAlphaImage,
		x:          g.boardLogicZeroPoint.x + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Ju,
	}
	g.sprites[1] = &Sprite{
		image:      ebitenBlackMaImage,
		alphaImage: ebitenBlackMaAlphaImage,
		x:          g.boardLogicZeroPoint.x + g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Ma,
	}
	g.sprites[2] = &Sprite{
		image:      ebitenBlackXiangImage,
		alphaImage: ebitenBlackXiangAlphaImage,
		x:          g.boardLogicZeroPoint.x + 2*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Xiang,
	}
	g.sprites[3] = &Sprite{
		image:      ebitenBlackShiImage,
		alphaImage: ebitenBlackShiAlphaImage,
		x:          g.boardLogicZeroPoint.x + 3*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Shi,
	}
	g.sprites[4] = &Sprite{
		image:      ebitenBlackJiangImage,
		alphaImage: ebitenBlackJiangAlphaImage,
		x:          g.boardLogicZeroPoint.x + 4*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.JiangShuai,
	}
	g.sprites[5] = &Sprite{
		image:      ebitenBlackShiImage,
		alphaImage: ebitenBlackShiAlphaImage,
		x:          g.boardLogicZeroPoint.x + 5*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Shi,
	}
	g.sprites[6] = &Sprite{
		image:      ebitenBlackXiangImage,
		alphaImage: ebitenBlackXiangAlphaImage,
		x:          g.boardLogicZeroPoint.x + 6*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Xiang,
	}
	g.sprites[7] = &Sprite{
		image:      ebitenBlackMaImage,
		alphaImage: ebitenBlackMaAlphaImage,
		x:          g.boardLogicZeroPoint.x + 7*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Ma,
	}
	g.sprites[8] = &Sprite{
		image:      ebitenBlackJuImage,
		alphaImage: ebitenBlackJuAlphaImage,
		x:          g.boardLogicZeroPoint.x + 8*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Ju,
	}
	g.sprites[9] = &Sprite{
		image:      ebitenBlackPaoImage,
		alphaImage: ebitenBlackPaoAlphaImage,
		x:          g.boardLogicZeroPoint.x + g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 2*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Pao,
	}
	g.sprites[10] = &Sprite{
		image:      ebitenBlackPaoImage,
		alphaImage: ebitenBlackPaoAlphaImage,
		x:          g.boardLogicZeroPoint.x + 7*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 2*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.Pao,
	}
	g.sprites[11] = &Sprite{
		image:      ebitenBlackZuImage,
		alphaImage: ebitenBlackZuAlphaImage,
		x:          g.boardLogicZeroPoint.x + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 3*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.BingZu,
	}
	g.sprites[12] = &Sprite{
		image:      ebitenBlackZuImage,
		alphaImage: ebitenBlackZuAlphaImage,
		x:          g.boardLogicZeroPoint.x + 2*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 3*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.BingZu,
	}
	g.sprites[13] = &Sprite{
		image:      ebitenBlackZuImage,
		alphaImage: ebitenBlackZuAlphaImage,
		x:          g.boardLogicZeroPoint.x + 4*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 3*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.BingZu,
	}
	g.sprites[14] = &Sprite{
		image:      ebitenBlackZuImage,
		alphaImage: ebitenBlackZuAlphaImage,
		x:          g.boardLogicZeroPoint.x + 6*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 3*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.BingZu,
	}
	g.sprites[15] = &Sprite{
		image:      ebitenBlackZuImage,
		alphaImage: ebitenBlackZuAlphaImage,
		x:          g.boardLogicZeroPoint.x + 8*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 3*g.gridLength + g.spriteReparation,
		group:      p1.GetGroup(),
		code:       core.BingZu,
	}

	g.sprites[16] = &Sprite{
		image:      ebitenRedBingImage,
		alphaImage: ebitenRedBingAlphaImage,
		x:          g.boardLogicZeroPoint.x + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 6*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.BingZu,
	}

	g.sprites[17] = &Sprite{
		image:      ebitenRedBingImage,
		alphaImage: ebitenRedBingAlphaImage,
		x:          g.boardLogicZeroPoint.x + 2*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 6*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.BingZu,
	}

	g.sprites[18] = &Sprite{
		image:      ebitenRedBingImage,
		alphaImage: ebitenRedBingAlphaImage,
		x:          g.boardLogicZeroPoint.x + 4*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 6*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.BingZu,
	}

	g.sprites[19] = &Sprite{
		image:      ebitenRedBingImage,
		alphaImage: ebitenRedBingAlphaImage,
		x:          g.boardLogicZeroPoint.x + 6*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 6*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.BingZu,
	}

	g.sprites[20] = &Sprite{
		image:      ebitenRedBingImage,
		alphaImage: ebitenRedBingAlphaImage,
		x:          g.boardLogicZeroPoint.x + 8*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 6*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.BingZu,
	}

	g.sprites[21] = &Sprite{
		image:      ebitenRedPaoImage,
		alphaImage: ebitenRedPaoAlphaImage,
		x:          g.boardLogicZeroPoint.x + g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 7*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Pao,
	}

	g.sprites[22] = &Sprite{
		image:      ebitenRedPaoImage,
		alphaImage: ebitenRedPaoAlphaImage,
		x:          g.boardLogicZeroPoint.x + 7*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 7*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Pao,
	}

	g.sprites[23] = &Sprite{
		image:      ebitenRedJuImage,
		alphaImage: ebitenRedJuAlphaImage,
		x:          g.boardLogicZeroPoint.x + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Ju,
	}
	g.sprites[24] = &Sprite{
		image:      ebitenRedMaImage,
		alphaImage: ebitenRedMaAlphaImage,
		x:          g.boardLogicZeroPoint.x + 1*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Ma,
	}
	g.sprites[25] = &Sprite{
		image:      ebitenRedXiangImage,
		alphaImage: ebitenRedXiangAlphaImage,
		x:          g.boardLogicZeroPoint.x + 2*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Xiang,
	}
	g.sprites[26] = &Sprite{
		image:      ebitenRedShiImage,
		alphaImage: ebitenRedShiAlphaImage,
		x:          g.boardLogicZeroPoint.x + 3*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Shi,
	}
	g.sprites[27] = &Sprite{
		image:      ebitenRedShuaiImage,
		alphaImage: ebitenRedShuaiAlphaImage,
		x:          g.boardLogicZeroPoint.x + 4*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.JiangShuai,
	}
	g.sprites[28] = &Sprite{
		image:      ebitenRedShiImage,
		alphaImage: ebitenRedShiAlphaImage,
		x:          g.boardLogicZeroPoint.x + 5*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Shi,
	}

	g.sprites[29] = &Sprite{
		image:      ebitenRedXiangImage,
		alphaImage: ebitenRedXiangAlphaImage,
		x:          g.boardLogicZeroPoint.x + 6*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Xiang,
	}
	g.sprites[30] = &Sprite{
		image:      ebitenRedMaImage,
		alphaImage: ebitenRedMaAlphaImage,
		x:          g.boardLogicZeroPoint.x + 7*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Ma,
	}
	g.sprites[31] = &Sprite{
		image:      ebitenRedJuImage,
		alphaImage: ebitenRedJuAlphaImage,
		x:          g.boardLogicZeroPoint.x + 8*g.gridLength + g.spriteReparation,
		y:          g.boardLogicZeroPoint.y + 9*g.gridLength + g.spriteReparation,
		group:      p2.GetGroup(),
		code:       core.Ju,
	}

}

// 判断坐标是否在sprite范围内
func (g *Game) spriteAt(x, y int) *Sprite {
	for i := len(g.sprites) - 1; i >= 0; i-- {
		s := g.sprites[i]
		if s.In(x, y) {
			return s
		}
	}
	return nil
}

func (g *Game) moveSpriteToFront(sprite *Sprite) {
	index := -1
	for i, ss := range g.sprites {
		if ss == sprite {
			index = i
			break
		}
	}
	g.sprites = append(g.sprites[:index], g.sprites[index+1:]...)
	g.sprites = append(g.sprites, sprite)
}

func (g *Game) deleteSprite(sprite *Sprite) {
	index := -1
	for i, ss := range g.sprites {
		if ss == sprite {
			index = i
			break
		}
	}
	g.sprites = append(g.sprites[:index], g.sprites[index+1:]...)

}

// InBoard 判断坐标是否位于棋盘内部
func (g *Game) InBoard(x, y int) bool {
	if x < boardLogicZeroX || x > boardLogicZeroX+8*gridLength {
		return false
	}
	if y < boardLogicZeroY || y > boardLogicZeroY+9*gridLength {
		return false
	}
	return true
}

// 修正点击棋盘的坐标精确落到棋盘格
func (g *Game) revisesCoordinate(x, y int) (x2, y2 int) {

	x = x - boardLogicZeroX
	y = y - boardLogicZeroY

	x2 = x / boardReparation * boardReparation
	if x%boardReparation > boardReparation/2 {
		x2 += boardReparation
	}
	x2 += boardLogicZeroX

	y2 = y / boardReparation * boardReparation
	if y%boardReparation > boardReparation/2 {
		y2 += boardReparation
	}
	y2 += boardLogicZeroY

	return x2, y2
}

// 将游戏界面坐标转换成核心层棋盘逻辑坐标
func (g *Game) transformCoordinate(x, y int) (coreX, coreY int) {

	//计算棋盘终点
	dx := boardLogicZeroX + 8*gridLength
	dy := boardLogicZeroY + 9*gridLength

	coreX = (dx - x) / gridLength
	coreY = (dy - y) / gridLength

	return coreX, coreY
}

func (g *Game) ShowGameMsg(screen *ebiten.Image) {
	if g.gameMsg == nil {
		return
	}
	str := ""
	str += fmt.Sprintf("结果：%s. ", g.gameMsg.Event)

	if len(g.gameMsg.WonChessmanCode) > 0 {
		str += fmt.Sprintf("吃棋：%s. ", g.gameMsg.WonChessmanCode)
	}

	if len(g.gameMsg.Msg) > 0 {
		str += fmt.Sprintf("消息：%s. ", g.gameMsg.Msg)
	}

	if g.gameMsg.WonGroup > 0 {
		str += fmt.Sprintf("对局结束，胜：%d. ", g.gameMsg.WonGroup)
	}

	f := &text.GoTextFace{
		Source:    hanziFaceSource,
		Direction: text.DirectionLeftToRight,
		Size:      24,
		//Language:  language.Chinese,
	}
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(g.boardLogicZeroPoint.x), float64(g.boardLogicZeroPoint.y+9*gridLength+30))
	text.Draw(screen, str, f, op)
	//ebitenutil.DebugPrintAt(screen, str, g.boardLogicZeroPoint.x, g.boardLogicZeroPoint.y+9*gridLength+30)
}

// 根据阵营获取玩家信息
func (g *Game) getPlayerByGroup(group core.ChessmanGroup) player.PlayerInterface {
	if g.player1.GetGroup() == group {
		return g.player1
	}
	if g.player2.GetGroup() == group {
		return g.player2
	}
	return nil
}

func (g *Game) drawWinner(screen *ebiten.Image) {
	//绘制胜利图标
	op := &ebiten.DrawImageOptions{}
	winLogoX := g.boardLogicZeroPoint.x + 2*g.gridLength
	winLogoY := g.boardLogicZeroPoint.y + 3*g.gridLength
	op.GeoM.Translate(float64(winLogoX), float64(winLogoY))
	screen.DrawImage(ebitenWinImage, op)

	//ebitenutil.DebugPrintAt(screen, g.winner, winLogoX+2*g.gridLength, winLogoY+3*g.gridLength)

	f := &text.GoTextFace{
		Source:    hanziFaceSource,
		Direction: text.DirectionLeftToRight,
		Size:      32,
		//Language:  language.Chinese,
	}
	op2 := &text.DrawOptions{}
	op2.GeoM.Translate(float64(winLogoX+3*g.gridLength/2), float64(winLogoY+3*g.gridLength))
	text.Draw(screen, g.winner, f, op2)

	//绘制再来一次按钮
	onceAgainBtn := &onceAgainBtn{
		image:      ebitenOnceAgainBtnImage,
		alphaImage: ebitenOnceAgainBtnAlphaImage,
		x:          g.boardLogicZeroPoint.x + 3*g.gridLength,
		y:          g.boardLogicZeroPoint.y + 7*g.gridLength,
	}
	onceAgainBtn.Draw(screen, 1)

}
