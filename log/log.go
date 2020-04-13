package log

import (
	"log"
	"os"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ltime)
	log.SetPrefix("[info]: ")
}

func Infof(format string, v ...interface{}) {
	log.Printf("[INFO]: ")
	log.Printf(format, v...)
}

func Infoln(v ...interface{}) {
	log.Printf("[INFO]: ")
	log.Println(v...)
}

func Warnf(format string, v ...interface{}) {
	log.Printf("[WARN]: ")
	log.Printf(format, v...)
}

func Warnln(v ...interface{}) {
	log.Printf("[WARN]: ")
	log.Println(v...)
}

func Errf(format string, v ...interface{}) {
	log.Printf("[ERROR]: ")
	log.Printf(format, v...)
	os.Exit(1)
}

func Errln(v ...interface{}) {
	log.Printf("[ERROR]: ")
	log.Println(v...)
	os.Exit(1)
}
