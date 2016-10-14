package timestamp

import (
	"time"

	"github.com/qgweb/glib/convert"
)

//获取天的时间戳
func GetDayUnix(day int) int64 {
	d := time.Now().AddDate(0, 0, day).Format("20060102")
	a, _ := time.ParseInLocation("20060102", d, time.Local)
	return a.Unix()
}

//获取天的时间戳
func GetDayUnixStr(day int) string {
	return convert.ToString(GetDayUnix(day))
}

//获取小时的时间戳
func GetＨourUnix(hour int) int64 {
	d := time.Now().Add(time.Hour * time.Duration(hour)).Format("2006010215")
	a, _ := time.ParseInLocation("2006010215", d, time.Local)
	return a.Unix()
}

//获取小时的时间戳
func GetＨourUnixStr(hour int) string {
	return convert.ToString(GetＨourUnix(hour))
}

//获取月的时间戳
func GetＭonthUnix(month int) int64 {
	d := time.Now().AddDate(0, month, 0).Format("200601")
	a, _ := time.ParseInLocation("200601", d, time.Local)
	return a.Unix()
}

//获取月的时间戳
func GetMonthUnixStr(month int) string {
	return convert.ToString(GetＭonthUnix(month))
}

// 获取当前时间戳
func GetNow() int64 {
	return time.Now().Unix()
}
