package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

func init() {
	log.SetFlags(log.Lmicroseconds | log.Ltime)
}

func setPrefix(level string) {
	_, file, line, _ := runtime.Caller(2)
	file = file[strings.Index(file, "blockchain")+len("blockchain"):]
	log.SetPrefix(fmt.Sprintf("[%s]: %s:%d ", level, file, line))
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

func Warnf(format string, v ...interface{}) {
	setPrefix("WARN")
	log.Printf(format, v...)
}

func Warnln(v ...interface{}) {
	setPrefix("WARN")
	log.Println(v...)
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
