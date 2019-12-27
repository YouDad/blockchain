package log

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "[info]: ", log.Ltime|log.Lmicroseconds|log.Lshortfile)

func Print(v ...interface{}) {
	logger.Print(v...)
}

func Printf(format string, v ...interface{}) {
	logger.Printf(format, v...)
}

func Println(v ...interface{}) {
	logger.Println(v...)
}

func Fatal(v ...interface{}) {
	logger.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	logger.Fatalf(format, v...)
}

func Fatalln(v ...interface{}) {
	logger.Fatalln(v...)
}

func Panic(v ...interface{}) {
	logger.Panic(v...)
}

func Panicf(format string, v ...interface{}) {
	logger.Panicf(format, v...)
}

func Panicln(v ...interface{}) {
	logger.Panicln(v...)
}
