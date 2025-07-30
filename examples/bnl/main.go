package main

import (
	"fmt"

	"github.com/gkoos/skyline/skyline"
)

func main() {
	// Example: products with price (minimize) and battery life (maximize)
	dims := []string{"price", "battery"}

	data := skyline.Dataset{
		{400, 10},
		{500, 12},
		{300, 9},
		{450, 11},
		{420, 15},
		{460, 14},
		{390, 8},
	}

	prefs := skyline.Preference{skyline.Min, skyline.Max}

	result, err := skyline.Skyline(data, dims, prefs, "bnl")
	if err != nil {
		panic(err)
	}
	fmt.Println("Skyline (BNL):", result)
}
