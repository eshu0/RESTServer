package main

import (
	"fmt"
	"net/http"
)

type TestStruct struct {
}

func (teststruct TestStruct) GetSomething(w http.ResponseWriter, r *http.Request) {
	fmt.Println("TestStruct - GetSomething")
}

func (teststruct TestStruct) CallSomething(w http.ResponseWriter, r *http.Request) {
	fmt.Println("TestStruct - CallSomething")
}
