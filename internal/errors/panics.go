package errors

import "fmt"

func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func PanicIfErrMsg(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", msg, err))
	}
}
