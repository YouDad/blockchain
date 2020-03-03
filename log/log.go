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
	logLevel  int
	prefix    string
)

func Register(level int, port string) {
	log.SetFlags(log.Lmicroseconds | log.Ltime)
	logLevel = level
	prefix = fmt.Sprintf("[%s]", port)
}

func LogSetPrefix(prefixString string) {
	log.SetPrefix(prefix + prefixString)
}

func setPrefix(level string) {
	_, file, line, _ := runtime.Caller(callLevel)
	file = file[strings.Index(file, "blockchain")+len("blockchain")+1:]
	LogSetPrefix(fmt.Sprintf("[%s]: { %s +%d } ", level, file, line))
}

func SetCallerLevel(level int) {
	callLevel = level + 2
}

func Debugf(format string, v ...interface{}) {
	if logLevel < 3 {
		return
	}
	setPrefix("DEBUG")
	log.Printf(format, v...)
}

func Debugln(v ...interface{}) {
	if logLevel < 3 {
		return
	}
	setPrefix("DEBUG")
	log.Println(v...)
}

func Infof(format string, v ...interface{}) {
	if logLevel < 2 {
		return
	}
	setPrefix("INFO")
	log.Printf(format, v...)
}

func Infoln(v ...interface{}) {
	if logLevel < 2 {
		return
	}
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
	if logLevel < 1 {
		return
	}
	setPrefix("WARN")
	log.Printf(format, v...)
}

func Warnln(v ...interface{}) {
	if logLevel < 1 {
		return
	}
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
	if logLevel < 0 {
		return
	}
	setPrefix("ERROR")
	log.Printf(format, v...)
	os.Exit(1)
}

func Errln(v ...interface{}) {
	if logLevel < 0 {
		return
	}
	setPrefix("ERROR")
	log.Println(v...)
	os.Exit(1)
}

func PrintStack() {
	LogSetPrefix("")
	log.SetFlags(0)
	log.Println("")
	LogSetPrefix("[STACK]: ")
	for i := 0; i < 5; i++ {
		pc, file, line, ok := runtime.Caller(2 + 4 - i)
		if !ok {
			continue
		}
		function := runtime.FuncForPC(pc)
		keyword := "blockchain"
		index := strings.Index(file, keyword)
		if index == -1 {
			keyword = "github.com"
			index = strings.Index(file, keyword)
		}
		file = file[index+len(keyword)+1:]
		log.Printf("{ %s +%d } [ %s ]", file, line, function.Name())
	}
}

func NotImplement() {
	PrintStack()
	SetCallerLevel(1)
	Errln("NotImplement")
}
