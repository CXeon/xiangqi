package core

type Coordinate struct {
	X int
	Y int
}

type ChessmanCode string

const (
	BingZu     ChessmanCode = "BingZu"     //兵、卒
	Pao        ChessmanCode = "Pao"        //炮
	Ju         ChessmanCode = "Ju"         //车
	Ma         ChessmanCode = "Ma"         //马
	Xiang      ChessmanCode = "Xiang"      //象
	Shi        ChessmanCode = "Shi"        //士
	JiangShuai ChessmanCode = "JiangShuai" //将、帅
)

type ChessmanGroup int

const (
	GroupNone ChessmanGroup = iota
	Group1
	Group2
)
