package logic

import (
	"github.com/reechou/real-reading/config"
	"github.com/reechou/real-reading/models"
)

type UpdateCourseLogic struct {
	cfg *config.Config
}

func NewUpdateCourseLogic(cfg *config.Config) *UpdateCourseLogic {
	ucl := &UpdateCourseLogic{cfg: cfg}
	ucl.init()

	return ucl
}

func (ucl *UpdateCourseLogic) init() {
	models.InitDB(ucl.cfg)
}

func (ucl *UpdateCourseLogic) Run() {

}
