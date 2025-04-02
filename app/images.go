package app

import (
	"bytes"
	"github.com/CXeon/xiangqi/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
)

var (
	ScreenWidth  = 576
	ScreenHeight = 672
)

var (
	//棋盘
	ebitenBoardImage *ebiten.Image

	//黑车
	ebitenBlackJuImage      *ebiten.Image
	ebitenBlackJuAlphaImage *image.Alpha

	//黑马
	ebitenBlackMaImage      *ebiten.Image
	ebitenBlackMaAlphaImage *image.Alpha

	//黑象
	ebitenBlackXiangImage      *ebiten.Image
	ebitenBlackXiangAlphaImage *image.Alpha

	//黑士
	ebitenBlackShiImage      *ebiten.Image
	ebitenBlackShiAlphaImage *image.Alpha

	//黑将
	ebitenBlackJiangImage      *ebiten.Image
	ebitenBlackJiangAlphaImage *image.Alpha

	//黑炮
	ebitenBlackPaoImage      *ebiten.Image
	ebitenBlackPaoAlphaImage *image.Alpha

	//黑卒
	ebitenBlackZuImage      *ebiten.Image
	ebitenBlackZuAlphaImage *image.Alpha

	//红车
	ebitenRedJuImage      *ebiten.Image
	ebitenRedJuAlphaImage *image.Alpha

	//红马
	ebitenRedMaImage      *ebiten.Image
	ebitenRedMaAlphaImage *image.Alpha

	//红相
	ebitenRedXiangImage      *ebiten.Image
	ebitenRedXiangAlphaImage *image.Alpha

	//红仕
	ebitenRedShiImage      *ebiten.Image
	ebitenRedShiAlphaImage *image.Alpha

	//红帅
	ebitenRedShuaiImage      *ebiten.Image
	ebitenRedShuaiAlphaImage *image.Alpha

	//红炮
	ebitenRedPaoImage      *ebiten.Image
	ebitenRedPaoAlphaImage *image.Alpha

	//红兵
	ebitenRedBingImage      *ebiten.Image
	ebitenRedBingAlphaImage *image.Alpha
)

func init() {

	//加载棋盘到内存
	imgBoard, _, err := image.Decode(bytes.NewReader(assets.Board))
	if err != nil {
		log.Fatal(err)
	}
	ebitenBoardImage = ebiten.NewImageFromImage(imgBoard)

	{
		//加载黑方车到内存
		imgBlackJu, _, err := image.Decode(bytes.NewReader(assets.BlackJu))
		if err != nil {
			log.Fatal(err)
		}
		ebitenBlackJuImage = ebiten.NewImageFromImage(imgBlackJu)
		bBlackJu := imgBlackJu.Bounds()
		ebitenBlackJuAlphaImage = image.NewAlpha(bBlackJu)
		for j := bBlackJu.Min.Y; j < bBlackJu.Max.Y; j++ {
			for i := bBlackJu.Min.X; i < bBlackJu.Max.X; i++ {
				ebitenBlackJuAlphaImage.Set(i, j, bBlackJu.At(i, j))
			}
		}
	}

	{
		//加载黑方马到内存
		imgBlackMa, _, err := image.Decode(bytes.NewReader(assets.BlackMa))
		if err != nil {
			log.Fatal(err)
		}
		ebitenBlackMaImage = ebiten.NewImageFromImage(imgBlackMa)
		bBlackMa := imgBlackMa.Bounds()
		ebitenBlackMaAlphaImage = image.NewAlpha(bBlackMa)
		for j := bBlackMa.Min.Y; j < bBlackMa.Max.Y; j++ {
			for i := bBlackMa.Min.X; i < bBlackMa.Max.X; i++ {
				ebitenBlackMaAlphaImage.Set(i, j, bBlackMa.At(i, j))
			}
		}
	}

	{
		//加载黑方象到内存
		imgBlackXiang, _, err := image.Decode(bytes.NewReader(assets.BlackXiang))
		if err != nil {
			log.Fatal(err)
		}
		ebitenBlackXiangImage = ebiten.NewImageFromImage(imgBlackXiang)
		bBlackXiang := imgBlackXiang.Bounds()
		ebitenBlackXiangAlphaImage = image.NewAlpha(bBlackXiang)
		for j := bBlackXiang.Min.Y; j < bBlackXiang.Max.Y; j++ {
			for i := bBlackXiang.Min.X; i < bBlackXiang.Max.X; i++ {
				ebitenBlackXiangAlphaImage.Set(i, j, bBlackXiang.At(i, j))
			}
		}
	}

	{
		//加载黑方士到内存
		imgBlackShi, _, err := image.Decode(bytes.NewReader(assets.BlackShi))
		if err != nil {
			log.Fatal(err)
		}
		ebitenBlackShiImage = ebiten.NewImageFromImage(imgBlackShi)
		bBlackShi := imgBlackShi.Bounds()
		ebitenBlackShiAlphaImage = image.NewAlpha(bBlackShi)
		for j := bBlackShi.Min.Y; j < bBlackShi.Max.Y; j++ {
			for i := bBlackShi.Min.X; i < bBlackShi.Max.X; i++ {
				ebitenBlackShiAlphaImage.Set(i, j, bBlackShi.At(i, j))
			}
		}
	}

	{
		//加载黑方将到内存
		imgBlackJiang, _, err := image.Decode(bytes.NewReader(assets.BlackJiang))
		if err != nil {
			log.Fatal(err)
		}
		ebitenBlackJiangImage = ebiten.NewImageFromImage(imgBlackJiang)
		bBlackJiang := imgBlackJiang.Bounds()
		ebitenBlackJiangAlphaImage = image.NewAlpha(bBlackJiang)
		for j := bBlackJiang.Min.Y; j < bBlackJiang.Max.Y; j++ {
			for i := bBlackJiang.Min.X; i < bBlackJiang.Max.X; i++ {
				ebitenBlackJiangAlphaImage.Set(i, j, imgBlackJiang.At(i, j))
			}
		}
	}

	{
		//加载黑方炮到内存
		imgBlackPao, _, err := image.Decode(bytes.NewReader(assets.BlackPao))
		if err != nil {
			log.Fatal(err)
		}
		ebitenBlackPaoImage = ebiten.NewImageFromImage(imgBlackPao)
		bBlackPao := imgBlackPao.Bounds()
		ebitenBlackPaoAlphaImage = image.NewAlpha(bBlackPao)
		for j := bBlackPao.Min.Y; j < bBlackPao.Max.Y; j++ {
			for i := bBlackPao.Min.X; i < bBlackPao.Max.X; i++ {
				ebitenBlackPaoAlphaImage.Set(i, j, bBlackPao.At(i, j))
			}
		}
	}

	{
		//加载黑方卒到内存
		imgBlackZu, _, err := image.Decode(bytes.NewReader(assets.BlackZu))
		if err != nil {
			log.Fatal(err)
		}
		ebitenBlackZuImage = ebiten.NewImageFromImage(imgBlackZu)
		bBlackZu := imgBlackZu.Bounds()
		ebitenBlackZuAlphaImage = image.NewAlpha(bBlackZu)
		for j := bBlackZu.Min.Y; j < bBlackZu.Max.Y; j++ {
			for i := bBlackZu.Min.X; i < bBlackZu.Max.X; i++ {
				ebitenBlackZuAlphaImage.Set(i, j, bBlackZu.At(i, j))
			}
		}
	}

	{
		//加载红方车到内存
		imgRedJu, _, err := image.Decode(bytes.NewReader(assets.RedJu))
		if err != nil {
			log.Fatal(err)
		}
		ebitenRedJuImage = ebiten.NewImageFromImage(imgRedJu)
		bRedJu := imgRedJu.Bounds()
		ebitenRedJuAlphaImage = image.NewAlpha(bRedJu)
		for j := bRedJu.Min.Y; j < bRedJu.Max.Y; j++ {
			for i := bRedJu.Min.X; i < bRedJu.Max.X; i++ {
				ebitenRedJuAlphaImage.Set(i, j, bRedJu.At(i, j))
			}
		}
	}

	{
		//加载红方马到内存
		imgRedMa, _, err := image.Decode(bytes.NewReader(assets.RedMa))
		if err != nil {
			log.Fatal(err)
		}
		ebitenRedMaImage = ebiten.NewImageFromImage(imgRedMa)
		bRedMa := imgRedMa.Bounds()
		ebitenRedMaAlphaImage = image.NewAlpha(bRedMa)
		for j := bRedMa.Min.Y; j < bRedMa.Max.Y; j++ {
			for i := bRedMa.Min.X; i < bRedMa.Max.X; i++ {
				ebitenRedMaAlphaImage.Set(i, j, bRedMa.At(i, j))
			}
		}
	}

	{
		//加载红方相到内存
		imgRedXiang, _, err := image.Decode(bytes.NewReader(assets.RedXiang))
		if err != nil {
			log.Fatal(err)
		}
		ebitenRedXiangImage = ebiten.NewImageFromImage(imgRedXiang)
		bRedXiang := imgRedXiang.Bounds()
		ebitenRedXiangAlphaImage = image.NewAlpha(bRedXiang)
		for j := bRedXiang.Min.Y; j < bRedXiang.Max.Y; j++ {
			for i := bRedXiang.Min.X; i < bRedXiang.Max.X; i++ {
				ebitenRedXiangAlphaImage.Set(i, j, bRedXiang.At(i, j))
			}
		}
	}

	{
		//加载红方仕到内存
		imgRedShi, _, err := image.Decode(bytes.NewReader(assets.RedShi))
		if err != nil {
			log.Fatal(err)
		}
		ebitenRedShiImage = ebiten.NewImageFromImage(imgRedShi)
		bRedShi := imgRedShi.Bounds()
		ebitenRedShiAlphaImage = image.NewAlpha(bRedShi)
		for j := bRedShi.Min.Y; j < bRedShi.Max.Y; j++ {
			for i := bRedShi.Min.X; i < bRedShi.Max.X; i++ {
				ebitenRedShiAlphaImage.Set(i, j, bRedShi.At(i, j))
			}
		}
	}

	{
		//加载红方帅到内存
		imgRedShuai, _, err := image.Decode(bytes.NewReader(assets.RedJiang))
		if err != nil {
			log.Fatal(err)
		}
		ebitenRedShuaiImage = ebiten.NewImageFromImage(imgRedShuai)
		bRedShuai := imgRedShuai.Bounds()
		ebitenRedShuaiAlphaImage = image.NewAlpha(bRedShuai)
		for j := bRedShuai.Min.Y; j < bRedShuai.Max.Y; j++ {
			for i := bRedShuai.Min.X; i < bRedShuai.Max.X; i++ {
				ebitenRedShuaiAlphaImage.Set(i, j, imgRedShuai.At(i, j))
			}
		}
	}

	{
		//加载黑方炮到内存
		imgRedPao, _, err := image.Decode(bytes.NewReader(assets.RedPao))
		if err != nil {
			log.Fatal(err)
		}
		ebitenRedPaoImage = ebiten.NewImageFromImage(imgRedPao)
		bRedPao := imgRedPao.Bounds()
		ebitenRedPaoAlphaImage = image.NewAlpha(bRedPao)
		for j := bRedPao.Min.Y; j < bRedPao.Max.Y; j++ {
			for i := bRedPao.Min.X; i < bRedPao.Max.X; i++ {
				ebitenRedPaoAlphaImage.Set(i, j, bRedPao.At(i, j))
			}
		}
	}

	{
		//加载红方兵到内存
		imgRedBing, _, err := image.Decode(bytes.NewReader(assets.RedBing))
		if err != nil {
			log.Fatal(err)
		}
		ebitenRedBingImage = ebiten.NewImageFromImage(imgRedBing)
		bRedBing := imgRedBing.Bounds()
		ebitenRedBingAlphaImage = image.NewAlpha(bRedBing)
		for j := bRedBing.Min.Y; j < bRedBing.Max.Y; j++ {
			for i := bRedBing.Min.X; i < bRedBing.Max.X; i++ {
				ebitenRedBingAlphaImage.Set(i, j, bRedBing.At(i, j))
			}
		}
	}
}
