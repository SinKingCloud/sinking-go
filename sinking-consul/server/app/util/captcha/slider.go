package captcha

import (
	"errors"
	images "github.com/wenlng/go-captcha-assets/resources/images_v2"
	"github.com/wenlng/go-captcha-assets/resources/tiles"
	"github.com/wenlng/go-captcha/v2/slide"
	"log"
)

var slideCapt slide.Captcha

func init() {
	builder := slide.NewBuilder(
		slide.WithGenGraphNumber(1),
		slide.WithEnableGraphVerticalRandom(true),
	)
	img, err := images.GetImages()
	if err != nil {
		log.Fatalln(err)
	}
	graphs, err := tiles.GetTiles()
	if err != nil {
		log.Fatalln(err)
	}
	var newGraphs = make([]*slide.GraphImage, 0, len(graphs))
	for i := 0; i < len(graphs); i++ {
		graph := graphs[i]
		newGraphs = append(newGraphs, &slide.GraphImage{
			OverlayImage: graph.OverlayImage,
			MaskImage:    graph.MaskImage,
			ShadowImage:  graph.ShadowImage,
		})
	}
	builder.SetResources(
		slide.WithGraphImages(newGraphs),
		slide.WithBackgrounds(img),
	)
	slideCapt = builder.Make()
}

// GenSlide 生成滑动验证码
func GenSlide() (result *slide.Block, masterImage string, tileImage string, height int, width int, err error) {
	captData, err := slideCapt.Generate()
	if err != nil {
		return nil, "", "", 0, 0, errors.New("生成验证码失败")
	}
	result = captData.GetData()
	if result == nil {
		return nil, "", "", 0, 0, errors.New("创建验证码失败")
	}
	masterImage, err = captData.GetMasterImage().ToBase64()
	if err != nil {
		return nil, "", "", 0, 0, errors.New("生成验证码主图失败")
	}
	tileImage, err = captData.GetTileImage().ToBase64()
	if err != nil {
		return nil, "", "", 0, 0, errors.New("生成验证码副图失败")
	}
	img := slideCapt.GetOptions().GetImageSize()
	height = img.Height
	width = img.Width
	return result, masterImage, tileImage, height, width, err
}

// CheckSlide 判断滑动验证码
func CheckSlide(x int, y int, result *slide.Block) bool {
	if result == nil {
		return false
	}
	return slide.CheckPoint(int64(x), int64(y), int64(result.X), int64(result.Y), 4)
}
