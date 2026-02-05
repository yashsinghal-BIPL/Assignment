package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func printType(value any, indent string) {
	v := reflect.ValueOf(value)
	t := reflect.TypeOf(value)

	if v.Kind() == reflect.Map { // checking for the Map type
		fmt.Println(indent, "Type is :", t, "(Map)")
		for _, value := range v.MapKeys() {
			fmt.Println(indent, "Key:", value.Interface())
			printType(v.MapIndex(value).Interface(), indent+"  ")
		}
	} else if v.Kind() == reflect.Slice || v.Kind() == reflect.Array { //checking for slice or array type
		fmt.Println(indent, "Type is :", t, "(Slice)")
		for i := 0; i < v.Len(); i++ {
			fmt.Printf("%sIndex %d:\n", indent, i)
			printType(v.Index(i).Interface(), indent+"  ")
		}
	} else {

		fmt.Printf("%sType is : %v | Value is : %v\n", indent, t, v.Interface())
	}
}

func main() {

	fmt.Println("Printing Entity with their Values and Types")

	dataFromWeb := []byte(`{
			"name" : "Tolexo Online Pvt. Ltd",
			"age_in_years" : 8.5,
			"origin" : "Noida",
			"head_office" : "Noida, Uttar Pradesh",
			"address" : [
			{
			"street" : "91 Springboard",
			"landmark" : "Axis Bank",
			"city" : "Noida",
			"pincode" : 201301,
			"state" : "Uttar Pradesh"
			},
			{
			"street" : "91 Springboard",
			"landmark" : "Axis Bank",
			"city" : "Noida",
			"pincode" : 201301,
			"state" : "Uttar Pradesh"
			}
			],
			"sponsers" : {
			"name" : "One"
			},
			"revenue" : "19.8 million$",
			"no_of_employee" : 630,
			"str_text" : ["one","two"],
			"int_text" : [1,3,4]
			}`)

	var store any

	err := json.Unmarshal(dataFromWeb, &store)

	if err != nil {
		panic(err)
	}

	//fmt.Println(lcoCourse)

	printType(store, "")

}
