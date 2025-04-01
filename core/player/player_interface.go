package player

import (
	"github.com/CXeon/xiangqi/core"
)

type PlayerInterface interface {
	ReceiveStatement(ch chan Statement, quit chan struct{}) (Statement, error) //player接收意图,quit用来终止读取防止阻塞
	SetGroup(group core.ChessmanGroup)                                         //为玩家分配阵营
	GetGroup() core.ChessmanGroup                                              //获取玩家所属阵营
	SetIsFirst(isFirst bool)                                                   //设置玩家是否是先手
	GetIsFirst() bool                                                          //获取玩家是否是先手
	SetIsDown(isDown bool)                                                     //设置玩家是否在棋盘俯视图下方
	GetIsDown() bool                                                           //获取玩家位于棋盘俯视图位置

	/**玩家存活的棋子相关**/
	GetOwnChessmen() ([]core.ChessmanCode, error) //获取玩家所有存活的棋子
	AddOwnChessman(code core.ChessmanCode) error  //添加一个存活的棋子
	DelOwnChessman(code core.ChessmanCode) error  //删除一个存活的棋子
	ClearOwnChessman()                            //清除玩家所有存活的棋子
	AddOwnChessmen(codes []core.ChessmanCode)     //批量添加玩家存活的棋子

	/**玩家失去的棋子相关**/
	GetLostChessmen() ([]core.ChessmanCode, error) //获取玩家所有失去的的棋子
	AddLostChessman(code core.ChessmanCode) error  //添加一个失去的棋子
	DelLostChessman(code core.ChessmanCode) error  //删除一个失去的棋子
	ClearLostChessman()                            //清除玩家所有失去的棋子
	AddLostChessmen(codes []core.ChessmanCode)     //批量添加玩家失去的棋子

	/**玩家赢得的棋子相关**/
	GetWonChessmen() ([]core.ChessmanCode, error) //获取玩家所有赢得的的棋子
	AddWonChessman(code core.ChessmanCode) error  //添加一个赢得的棋子
	DelWonChessman(code core.ChessmanCode) error  //删除一个赢得的棋子
	ClearWonChessman()                            //清除玩家所有赢得的棋子
	AddWonChessmen(codes []core.ChessmanCode)     //批量添加玩家赢得的棋子

}
