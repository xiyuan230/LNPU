package service

import (
	"fmt"
	"testing"
)

func TestSSOLogin(t *testing.T) {
	cookie, err := SSOLogin("2012200108", "Sql123..")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cookie)
}
