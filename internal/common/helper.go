package common

import (
	"fmt"
	"strings"
)

func Boom(err error, msg ...string) {
	if err != nil {
		fmt.Println(strings.Join(msg, " "))
		panic(err)
	}
}

func Assert(ok bool, msg ...string) {
	if !ok {
		panic("assertion failed: " + strings.Join(msg, " "))
	}
}
