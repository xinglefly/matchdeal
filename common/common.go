package common

import (
	"time"
	"fmt"
)

func FormatTime() string {
	t := time.Unix(time.Now().Unix(), 0)
	strDate := fmt.Sprintf("%d-%d-%d %d:%d:%d\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	return strDate
}
