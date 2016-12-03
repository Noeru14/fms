/************/
/* LOG FUNC */
/************/

package logger

import (
	"log"
	"os"
)

var l = log.New(os.Stdout, "[INFO]\t", log.Ldate | log.Ltime)
var lError = log.New(os.Stdout, "[ERROR]\t", log.Ldate | log.Ltime)

func Print(msg ...interface{}) {
	l.Println(msg...)
}

func PrintError(msg ...interface{}) {
	lError.Println(msg...)
}