package main

import (
	"fmt"
	"reflect"
	"time"
)

func main() {
	var x int
	x = 5

	xt := reflect.TypeOf(x)
	fmt.Println("Type of x:", xt.Name())
	fmt.Println("Kind of x:", xt.Kind())

	xpt := reflect.TypeOf(&x)
	fmt.Println("Type of &x:", xpt.Name())
	fmt.Println("Kind of &x:", xpt.Kind())

	fmt.Println("Original value of x:", x)
	valPtr := reflect.ValueOf(&x)
	val := valPtr.Elem()
	val.SetInt(3)
	fmt.Println("Overriden value of x:", x)

	// create new int value
	var y int
	intType := reflect.TypeOf(y)
	newPtrVal := reflect.New(intType)
	newPtrVal.Elem().SetInt(20)
	newVal := newPtrVal.Elem().Int()
	fmt.Println(newVal)

	testFun := makeTimedFunction(func() {
		fmt.Println("It works")
	}).(func())
	testFun()

	fmt.Printf("Got new struct: %+v\n", makeStruct())
}

func makeTimedFunction(f interface{}) interface{} {
	rf := reflect.TypeOf(f)
	if rf.Kind() != reflect.Func {
		panic("Expects a function")
	}
	vf := reflect.ValueOf(f)

	wrapperF := reflect.MakeFunc(rf, func(in []reflect.Value) []reflect.Value {
		start := time.Now()
		out := vf.Call(in)
		end := time.Now()
		fmt.Printf("calling took %v\n", end.Sub(start))
		return out
	})
	return wrapperF.Interface()
}

func makeStruct() interface{} {
	var i int
	var s string
	it := reflect.TypeOf(i)
	st := reflect.TypeOf(s)

	fields := []reflect.StructField{
		{
			Name: "Name",
			Type: st,
		},
		{
			Name: "Age",
			Type: it,
		},
	}
	structType := reflect.StructOf(fields)
	structVal := reflect.New(structType)

	structVal.Elem().FieldByName("Name").SetString("Bob")
	structVal.Elem().FieldByName("Age").SetInt(50)

	return structVal.Interface()
}
