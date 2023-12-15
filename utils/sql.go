package utils

import (
	"fmt"
	"log"
	"regexp"
)

// ReplaceBinds replaces given pattern in string
func ReplaceBinds(stm string) string {
	regexp, err := regexp.Compile(`:[a-zA-Z_0-9]+`)
	if err != nil {
		log.Fatal(err)
	}
	match := regexp.ReplaceAllString(stm, "?")
	return match
}

// PrintSQL prints SQL/args
func PrintSQL(stm string, args []interface{}, msg string) {
	if msg != "" {
		log.Println(msg)
	} else {
		log.Println("")
	}
	log.Printf("### SQL statement ###\n%s\n\n", stm)
	var values string
	for _, v := range args {
		values = fmt.Sprintf("%s\t'%v'\n", values, v)
	}
	log.Printf("### SQL values ###\n%s\n", values)
}
