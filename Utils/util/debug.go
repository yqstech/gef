package util

import (
	"fmt"
	"log"
	"os"
	"time"
)

func ErrorLog(err error) interface{} {
	if err != nil {
		//输出屏幕中
		fmt.Println(err)
		//记录到文件中
		fileName := "data/log/" + time.Now().Format("2006-01-02_") + "err.log"
		logFile, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		defer logFile.Close()
		debugLog := log.New(logFile, "[Error]", log.Llongfile)
		debugLog.Println(err)
		return err
	} else {
		return nil
	}
}
