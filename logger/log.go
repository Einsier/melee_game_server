package logger

/**
*@Author Sly
*@Date 2022/1/20
*@Version 1.0
*@Description:
 */

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

//初始化几个日志打印对象,并且加入到loggers中
var (
	//LstdFlags 指示打印日期和分钟,Lshortfile指示打印具体的文件名和行
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile) //打印Error有关
	infoLog  = log.New(os.Stdout, "\033[34m[info ]\033[0m ", log.LstdFlags|log.Lshortfile) //打印Info有关
	loggers  = []*log.Logger{errorLog, infoLog}                                            //两个*logger都放到loggers中
	mu       sync.Mutex
)

// log methods		暴露这几个方法,作为模块中的全局变量,方便以后修改向其他地方输出/调用
var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)

const (
	LogInfoLevel = iota
	LogErrorLevel
	LogDisabledLevel
)

// SetLogLevel 设置日志等级,LogInfoLevel为所有,LogErrorLevel为只打印错误,LogDisabledLevel为忽略全部
func SetLogLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	if LogErrorLevel < level {
		errorLog.SetOutput(ioutil.Discard)
	}
	if LogInfoLevel < level {
		infoLog.SetOutput(ioutil.Discard)
	}
}
