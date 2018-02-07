package controller

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"os"
	"time"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/reechou/holmes"
	"github.com/reechou/real-reading/config"
	"golang.org/x/image/font"
)

type GraduationInfo struct {
	Name          string
	TimeInfo      string
	CourseName    string
	CourseNum     string
	GraduationNum string
}

type GraduationManager struct {
	cfg  *config.Config
	font *truetype.Font
	bm   image.Image
}

func NewGraduationManager(cfg *config.Config) (*GraduationManager, error) {
	gm := &GraduationManager{cfg: cfg}
	var err error
	gm.font, err = gm.getFontFamily()
	if err != nil {
		holmes.Error("get font family error: %v", err)
		return nil, err
	}
	gm.bm, err = imaging.Open(gm.cfg.Graduation.ByzsPic)
	if err != nil {
		holmes.Error("open byzs file failed")
		return nil, err
	}

	return gm, nil
}

func (gm *GraduationManager) HandleUserImage(info *GraduationInfo) (string, error) {
	start := time.Now()
	defer func() {
		holmes.Debug("Handle graduation pic[%v] use_time[%v]", info, time.Now().Sub(start))
	}()

	dst := imaging.Clone(gm.bm)
	// 图片按比例缩放
	//dst := imaging.Resize(bm, 744, 463, imaging.Lanczos)
	gm.writeOnImage(dst, info)

	graduationName := fmt.Sprintf("%s.jpg", info.GraduationNum)
	fileName := fmt.Sprintf("%s/%s", gm.cfg.Graduation.TmpPath, graduationName)
	err := imaging.Save(dst, fileName)
	if err != nil {
		return "", err
	}

	return graduationName, nil
}

func (gm *GraduationManager) CheckUserImage(graduationNum string) (string, bool) {
	graduationName := fmt.Sprintf("%s.jpg", graduationNum)
	fileName := fmt.Sprintf("%s/%s", gm.cfg.Graduation.TmpPath, graduationName)
	_, err := os.Stat(fileName)
	if err != nil && os.IsNotExist(err) {
		return "", false
	}
	return graduationName, true
}

func (gm *GraduationManager) writeOnImage(target *image.NRGBA, info *GraduationInfo) error {
	c := freetype.NewContext()

	c.SetDPI(256)
	c.SetClip(target.Bounds())
	c.SetDst(target)
	c.SetHinting(font.HintingFull)

	// 设置文字颜色、字体、字大小
	c.SetSrc(image.NewUniform(color.RGBA{R: 0, G: 0, B: 0, A: 255}))
	c.SetFontSize(6)
	c.SetFont(gm.font)

	pt := freetype.Pt(60, 295)
	_, err := c.DrawString(info.Name, pt)
	if err != nil {
		holmes.Error("draw error: %v \n", err)
		return err
	}

	c.SetFontSize(5.5)
	pt = freetype.Pt(300, 295)
	_, err = c.DrawString(info.TimeInfo, pt)
	if err != nil {
		holmes.Error("draw error: %v \n", err)
		return err
	}

	c.SetFontSize(6)
	pt = freetype.Pt(230, 365)
	_, err = c.DrawString(info.CourseNum, pt)
	if err != nil {
		holmes.Error("draw error: %v \n", err)
		return err
	}

	pt = freetype.Pt(310, 365)
	_, err = c.DrawString(info.CourseName, pt)
	if err != nil {
		holmes.Error("draw error: %v \n", err)
		return err
	}

	pt = freetype.Pt(150, 855)
	_, err = c.DrawString(info.GraduationNum, pt)
	if err != nil {
		holmes.Error("draw error: %v \n", err)
		return err
	}
	return nil
}

func (gm *GraduationManager) getFontFamily() (*truetype.Font, error) {
	fontBytes, err := ioutil.ReadFile(gm.cfg.Graduation.Font)
	if err != nil {
		holmes.Error("read file error: %v", err)
		return &truetype.Font{}, err
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		holmes.Error("parse font error: %v", err)
		return &truetype.Font{}, err
	}

	return f, err
}
