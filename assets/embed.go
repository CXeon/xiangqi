package assets

import (
	_ "embed"
	_ "image/png"
)

var (
	//go:embed board.png
	Board []byte //棋盘文件

	//go:embed black_jiang.png
	BlackJiang []byte //黑将文件

	//go:embed black_ju.png
	BlackJu []byte //黑车文件

	//go:embed black_ma.png
	BlackMa []byte //黑马文件

	//go:embed black_pao.png
	BlackPao []byte //黑炮文件

	//go:embed black_shi.png
	BlackShi []byte //黑士文件

	//go:embed black_xiang.png
	BlackXiang []byte //黑象文件

	//go:embed black_zu.png
	BlackZu []byte //黑卒文件

	//go:embed red_shuai.png
	RedJiang []byte //红帅文件

	//go:embed red_ju.png
	RedJu []byte //红车文件

	//go:embed red_ma.png
	RedMa []byte //红马文件

	//go:embed red_pao.png
	RedPao []byte //红炮文件

	//go:embed red_shi.png
	RedShi []byte //红士文件

	//go:embed red_xiang.png
	RedXiang []byte //红象文件

	//go:embed red_bing.png
	RedBing []byte //红兵文件

	//go:embed win.png
	Win []byte //胜利界面

	//go:embed once_again.png
	OnceAgain []byte //再来一局按钮
)
