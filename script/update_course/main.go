package main

import (
	"github.com/reechou/real-reading/config"
	"github.com/reechou/real-reading/script/update_course/logic"
)

func main() {
	logic.NewUpdateCourseLogic(config.NewConfig()).Run()
}
