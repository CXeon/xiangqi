package player

import (
	"errors"
	"github.com/CXeon/xiangqi/core"
)

type Player struct {
	ID           int                 //玩家id
	group        core.ChessmanGroup  //阵营
	isFirst      bool                //是否是先手
	isDown       bool                //是否位于棋盘俯视图下方
	ownChessman  []core.ChessmanCode //自己拥有的己方棋子code数组
	lostChessman []core.ChessmanCode //失去的己方棋子code数组
	wonChessman  []core.ChessmanCode //赢得的棋子数组
}

func NewPlayer() *Player {
	return &Player{
		ownChessman:  make([]core.ChessmanCode, 0),
		lostChessman: make([]core.ChessmanCode, 0),
		wonChessman:  make([]core.ChessmanCode, 0),
	}
}

// ReceiveStatement player接收意图
func (p *Player) ReceiveStatement(ch chan Statement, quit chan struct{}) (Statement, error) {

	select {
	case <-quit: // 优先响应退出信号
		return Statement{}, errors.New("receive quit signal")
	case sta, ok := <-ch: // 同时监听数据通道
		if !ok {
			return Statement{}, errors.New("channel is closed")
		}
		return sta, nil
	}

}

// SetGroup 分配阵营
func (p *Player) SetGroup(group core.ChessmanGroup) {
	p.group = group
}

// GetGroup 获取阵营
func (p *Player) GetGroup() core.ChessmanGroup {
	return p.group
}

// SetIsFirst 设置是否先手
func (p *Player) SetIsFirst(isFirst bool) {
	p.isFirst = isFirst
}

// GetIsFirst 查询是否先手
func (p *Player) GetIsFirst() bool {
	return p.isFirst
}

// SetIsDown 设置玩家位于棋盘俯视图的位置
func (p *Player) SetIsDown(isDown bool) {
	p.isDown = isDown
}

// GetIsDown 获取玩家位于棋盘俯视图的位置
func (p *Player) GetIsDown() bool {
	return p.isDown
}

/**玩家存活的棋子相关**/

// GetOwnChessmen 获取玩家所有存活的棋子
func (p *Player) GetOwnChessmen() ([]core.ChessmanCode, error) {
	if p.ownChessman == nil {
		return nil, errors.New("chessman is nil")
	}
	return p.ownChessman, nil
}

// AddOwnChessman 添加一个存活的棋子
func (p *Player) AddOwnChessman(code core.ChessmanCode) error {
	if p.ownChessman == nil {
		return errors.New("chessman is nil")
	}
	p.ownChessman = append(p.ownChessman, code)
	return nil
}

// DelOwnChessman 删除一个存活的棋子
func (p *Player) DelOwnChessman(code core.ChessmanCode) error {
	if p.ownChessman == nil {
		return errors.New("chessman is nil")
	}
	if len(p.ownChessman) == 0 {
		return nil
	}

	index := -1
	for i := 0; i < len(p.ownChessman); i++ {
		if p.ownChessman[i] == code {
			index = i
			break
		}
	}
	if index < 0 {
		return nil
	}
	if index == 0 {
		p.ownChessman = p.ownChessman[1:]
		return nil
	}
	if len(p.ownChessman)-1 == index {
		p.ownChessman = p.ownChessman[:index]
		return nil
	}

	s1 := p.ownChessman[:index]
	s2 := p.ownChessman[index+1:]
	p.ownChessman = append(s1, s2...)
	return nil
}

// ClearOwnChessman 清除玩家所有存活的棋子
func (p *Player) ClearOwnChessman() {
	p.ownChessman = make([]core.ChessmanCode, 0)
	return
}

// AddOwnChessmen 批量添加玩家存活的棋子
func (p *Player) AddOwnChessmen(codes []core.ChessmanCode) {
	p.ownChessman = append(p.ownChessman, codes...)
}

/**玩家失去的棋子相关**/
// GetLostChessmen 获取玩家所有存活的棋子
func (p *Player) GetLostChessmen() ([]core.ChessmanCode, error) {
	if p.lostChessman == nil {
		return nil, errors.New("chessman is nil")
	}
	return p.lostChessman, nil
}

// AddLostChessman 添加一个存活的棋子
func (p *Player) AddLostChessman(code core.ChessmanCode) error {
	if p.lostChessman == nil {
		return errors.New("chessman is nil")
	}
	p.lostChessman = append(p.lostChessman, code)
	return nil
}

// DelLostChessman 删除一个存活的棋子
func (p *Player) DelLostChessman(code core.ChessmanCode) error {
	if p.lostChessman == nil {
		return errors.New("chessman is nil")
	}
	if len(p.lostChessman) == 0 {
		return nil
	}

	index := -1
	for i := 0; i < len(p.lostChessman); i++ {
		if p.lostChessman[i] == code {
			index = i
			break
		}
	}
	if index < 0 {
		return nil
	}
	if index == 0 {
		p.lostChessman = p.lostChessman[1:]
		return nil
	}
	if len(p.lostChessman)-1 == index {
		p.lostChessman = p.lostChessman[:index]
		return nil
	}

	s1 := p.lostChessman[:index]
	s2 := p.lostChessman[index+1:]
	p.lostChessman = append(s1, s2...)
	return nil
}

// ClearLostChessman 清除玩家所有存活的棋子
func (p *Player) ClearLostChessman() {
	p.lostChessman = make([]core.ChessmanCode, 0)
	return
}

// AddLostChessmen 批量添加玩家存活的棋子
func (p *Player) AddLostChessmen(codes []core.ChessmanCode) {
	p.lostChessman = append(p.lostChessman, codes...)
}

/**玩家赢得的棋子相关**/
// GetWonChessmen 获取玩家所有存活的棋子
func (p *Player) GetWonChessmen() ([]core.ChessmanCode, error) {
	if p.wonChessman == nil {
		return nil, errors.New("chessman is nil")
	}
	return p.wonChessman, nil
}

// AddWonChessman 添加一个存活的棋子
func (p *Player) AddWonChessman(code core.ChessmanCode) error {
	if p.wonChessman == nil {
		return errors.New("chessman is nil")
	}
	p.wonChessman = append(p.wonChessman, code)
	return nil
}

// DelWonChessman 删除一个存活的棋子
func (p *Player) DelWonChessman(code core.ChessmanCode) error {
	if p.wonChessman == nil {
		return errors.New("chessman is nil")
	}
	if len(p.wonChessman) == 0 {
		return nil
	}

	index := -1
	for i := 0; i < len(p.wonChessman); i++ {
		if p.wonChessman[i] == code {
			index = i
			break
		}
	}
	if index < 0 {
		return nil
	}
	if index == 0 {
		p.wonChessman = p.wonChessman[1:]
		return nil
	}
	if len(p.wonChessman)-1 == index {
		p.wonChessman = p.wonChessman[:index]
		return nil
	}

	s1 := p.wonChessman[:index]
	s2 := p.wonChessman[index+1:]
	p.wonChessman = append(s1, s2...)
	return nil
}

// ClearWonChessman 清除玩家所有存活的棋子
func (p *Player) ClearWonChessman() {
	p.wonChessman = make([]core.ChessmanCode, 0)
	return
}

// AddWonChessmen 批量添加玩家存活的棋子
func (p *Player) AddWonChessmen(codes []core.ChessmanCode) {
	p.wonChessman = append(p.wonChessman, codes...)
}
