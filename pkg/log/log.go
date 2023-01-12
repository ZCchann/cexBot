package log

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
)

var basedir string

func init() {
	workdir, _ := os.Getwd()
	basedir = path.Base(workdir)
}

func getFileLine(skips ...int) (file string, line int) {
	skip := 3
	if len(skips) > 0 {
		skip = skips[0]
	}
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "???"
		line = 0
	}
	dir := strings.Split(file, basedir+"/")
	return dir[len(dir)-1], line
}

func _print(level string, v ...interface{}) {
	file, line := getFileLine()
	log.Print(fmt.Sprintf("%s %s:%d - ", level, file, line), fmt.Sprintln(v...))
}

//func Debug(v ...interface{}) {
//	if os.Getenv("DEBUG") == "true" {
//		_print("DEBUG", v...)
//	}
//}
//
//func Info(v ...interface{}) {
//	_print("INFO", v...)
//}
//
//func Warning(v ...interface{}) {
//	_print("WARN", v...)
//}
//
//func Error(v ...interface{}) {
//	_print("ERROR", v...)
//}

func Println(v ...interface{}) {
	_print("INFO", v...)
}

func Fatalln(v ...interface{}) {
	file, line := getFileLine(2)
	log.Fatal(fmt.Sprintf("%s %s:%d - ", "FATAL", file, line), fmt.Sprintln(v...))
}
