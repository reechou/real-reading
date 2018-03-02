package logic

import (
	"fmt"
	"testing"

	"github.com/jinzhu/now"
	"time"
)

func TestWeek(t *testing.T) {
	now.FirstDayMonday = true

	weekBegin := now.BeginningOfWeek().Unix()
	weekConfig := []int64{1, 3, 5}
	for _, v := range weekConfig {
		fmt.Println(weekBegin + (v-1)*86400)
	}

	fmt.Println(now.EndOfWeek().Unix() + 1)

	fmt.Println("===========")

	fmt.Println(now.New(time.Unix(1519747200, 0)).BeginningOfWeek())

	fmt.Println((1519747200-now.New(time.Unix(1519747200, 0)).BeginningOfWeek().Unix())/86400 + 1)

	fmt.Println(1519920000 + 3*86400)

	fmt.Println(time.Unix(1546185600, 0).Year())

	fmt.Println(int64(time.Unix(1546185600, 0).Month()))
}
