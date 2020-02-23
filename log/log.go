package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

var (
	callLevel = 2
)

func init() {
	log.SetFlags(log.Lmicroseconds | log.Ltime)
}

func setPrefix(level string) {
	_, file, line, _ := runtime.Caller(callLevel)
	file = file[strings.Index(file, "blockchain")+len("blockchain")+1:]
	log.SetPrefix(fmt.Sprintf("[%s]: { %s +%d } ", level, file, line))
}

func SetCallerLevel(level int) {
	callLevel = level + 2
}

func Infof(format string, v ...interface{}) {
	setPrefix("INFO")
	log.SetPrefix("[INFO]: ")
	log.Printf(format, v...)
}

func Infoln(v ...interface{}) {
	setPrefix("INFO")
	log.Println(v...)
}

func Warn(err error) {
	if err != nil {
		SetCallerLevel(1)
		Warnln(err)
		SetCallerLevel(0)
	}
}

func Warnf(format string, v ...interface{}) {
	setPrefix("WARN")
	log.Printf(format, v...)
}

func Warnln(v ...interface{}) {
	setPrefix("WARN")
	log.Println(v...)
}

func Err(err error) {
	if err != nil {
		SetCallerLevel(1)
		Errln(err)
	}
}

func Errf(format string, v ...interface{}) {
	setPrefix("ERROR")
	log.Printf(format, v...)
	os.Exit(1)
}

func Errln(v ...interface{}) {
	setPrefix("ERROR")
	log.Println(v...)
	os.Exit(1)
}

func NotImplement() {
	SetCallerLevel(1)
	Errln("NotImplement")
}
