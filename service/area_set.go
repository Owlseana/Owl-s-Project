package service

import (
	"time"
)

func TimeinBeiJing() time.Time {
	now := time.Now()

	loc, _ := time.LoadLocation("Asia/Shanghai")

	// 将当前时间转换为北京时间
	beijingTime := now.In(loc)
	return beijingTime
}
