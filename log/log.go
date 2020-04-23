package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	callLevel   = 2
	logLevel    uint
	port        string
	prefix      string
	prefixMutex sync.Mutex
	levelMutex  sync.Mutex
)

func Register(level uint, p string) {
	log.SetFlags(0)
	logLevel = level
	port = p
}

func LogSetPrefix(prefixString string) {
	log.SetPrefix(prefix + prefixString)
}

func setPrefix(level string) {
	_, file, line, _ := runtime.Caller(callLevel)
	index := strings.Index(file, "blockchain")
	if index >= 0 {
		file = file[index+len("blockchain "):]
	}
	prefixMutex.Lock()

	// go routine
	routineId := make([]byte, 20)
	runtime.Stack(routineId, false)
	routineId = routineId[10:]
	for k, v := range routineId {
		if v == 0x20 {
			routineId = routineId[:k]
			break
		}
	}

	LogSetPrefix(fmt.Sprintf("[%d][%s][%s][%s]: { %s +%d } ",
		time.Now().UnixNano(), port, string(routineId), level, file, line))
}

func SetCallerLevel(level int) {
	if level != 0 {
		levelMutex.Lock()
	} else {
		levelMutex.Unlock()
	}
	callLevel = level + 2
}

func Funcname(level int) string {
	pc, _, _, _ := runtime.Caller(level + 1)
	return runtime.FuncForPC(pc).Name()
}

func Tracef(format string, v ...interface{}) {
	setPrefix("TRACE")
	log.Printf(format, v...)
	prefixMutex.Unlock()
}

func Traceln(v ...interface{}) {
	setPrefix("TRACE")
	log.Println(v...)
	prefixMutex.Unlock()
}

func Debugf(format string, v ...interface{}) {
	if logLevel < 3 {
		return
	}
	setPrefix("DEBUG")
	log.Printf(format, v...)
	prefixMutex.Unlock()
}

func Debugln(v ...interface{}) {
	if logLevel < 3 {
		return
	}
	setPrefix("DEBUG")
	log.Println(v...)
	prefixMutex.Unlock()
}

func Infof(format string, v ...interface{}) {
	if logLevel < 2 {
		return
	}
	setPrefix("INFO")
	log.Printf(format, v...)
	prefixMutex.Unlock()
}

func Infoln(v ...interface{}) {
	if logLevel < 2 {
		return
	}
	setPrefix("INFO")
	log.Println(v...)
	prefixMutex.Unlock()
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
	prefixMutex.Unlock()
}

func Warnln(v ...interface{}) {
	if logLevel < 1 {
		return
	}
	setPrefix("WARN")
	log.Println(v...)
	prefixMutex.Unlock()
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
	prefixMutex.Unlock()
	os.Exit(1)
}

func Errln(v ...interface{}) {
	if logLevel < 0 {
		return
	}
	setPrefix("ERROR")
	log.Println(v...)
	prefixMutex.Unlock()
	os.Exit(1)
}

func PrintStack() {
	levelMutex.Lock()
	LogSetPrefix("")
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
	levelMutex.Unlock()
}

func NotImplement() {
	PrintStack()
	SetCallerLevel(1)
	Errln("NotImplement")
}
