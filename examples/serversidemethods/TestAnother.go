package main

import (
	"fmt"
	"net/http"
)

type TestAnother struct {
}

func (teststruct TestAnother) GetSomething(w http.ResponseWriter, r *http.Request) {
	fmt.Println("TestAnother - GetSomething")
}

func (teststruct TestAnother) CallSomething(w http.ResponseWriter, r *http.Request) {
	fmt.Println("TestAnother - CallSomething")
}
