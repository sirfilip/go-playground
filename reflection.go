package main

import (
	"fmt"
	"reflect"
)

type Thing struct {
	Name string
	Age  int
}

func explore(thing interface{}) {
	v := reflect.Indirect(reflect.ValueOf(thing))
	fmt.Println(v.Type())
	for i := 0; i < v.NumField(); i++ {
		switch v.Field(i).Kind() {
		case reflect.Int:
			fmt.Println(v.Type().Field(i).Name, v.Field(i).Interface().(int))
			v.Field(i).SetInt(42)
			fmt.Println(v.Type().Field(i).Name, v.Field(i).Interface().(int))
		case reflect.String:
			fmt.Println(v.Type().Field(i).Name, v.Field(i).Interface().(string))
			v.Field(i).SetString("A New Name")
			fmt.Println(v.Type().Field(i).Name, v.Field(i).Interface().(string))
		}
	}

}

func main() {
	thingPtr := &Thing{"Swamp", 24}

	explore(thingPtr)
}
