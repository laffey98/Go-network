package print

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	con_place = "print"
)

func Printerr(err error, place string) {
	printtime()
	fmt.Print(" [print] ")
	fmt.Print("[", place, "] ")
	fmt.Println("error:", err)
}

func Printvar(place string, a ...interface{}) {
	printtime()
	fmt.Print(" [print] ")
	fmt.Print("[", place, "] ")
	for _, data := range a {
		dataJson, err := json.Marshal(data)
		if err != nil {
			fmt.Print("\n	")
			Printerr(err, con_place)
			fmt.Println("json encode failed!")
			return
		}
		fmt.Print("json encode:")
		fmt.Print(string(dataJson), " ")
	}
	fmt.Print("\n")
}

func printtime() {
	now := time.Now()                  //获取当前时间
	timestamp := now.Unix()            //时间戳
	timeObj := time.Unix(timestamp, 0) //将时间戳转为时间格式
	year := timeObj.Year()             //年
	month := timeObj.Month()           //月
	day := timeObj.Day()               //日
	hour := timeObj.Hour()             //小时
	minute := timeObj.Minute()         //分钟
	second := timeObj.Second()         //秒
	fmt.Printf("%d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, minute, second)
}
