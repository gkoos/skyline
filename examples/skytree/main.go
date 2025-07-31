package main

import (
	"fmt"
	"math/rand"

	"github.com/gkoos/skyline/skyline"
)

func main() {
	// Example: products with price (minimize) and battery life (maximize)
	dims := []string{"price", "battery"}
	prefs := skyline.Preference{skyline.Min, skyline.Max}

	// Create a dataset with 1000 items, with known skyline points
	// For this synthetic example, we add 5 skyline points that are not dominated by each other
	// and 95 random dominated points
	data := skyline.Dataset{
		{100, 16}, // lowest price, lowest battery
		{105, 17}, // second lowest price, second lowest battery
		{110, 18}, // middle price, middle battery
		{115, 19}, // second highest price, second highest battery
		{120, 20}, // highest price, highest battery
	}

	// Generate points that are strictly dominated by at least one skyline point
	for i := range 995 {
		// Pick a skyline point to dominate this random point
		base := data[i%5]
		// Make price higher and battery lower than the base
		price := base[0] + 1 + rand.Float64()*20
		battery := base[1] - 1 - rand.Float64()*5
		data = append(data, skyline.Point{price, battery})
	}

	result, err := skyline.Skyline(data, dims, prefs, "skytree")
	if err != nil {
		panic(err)
	}
	fmt.Println("Skyline (SkyTree):", result)

	// Expected result (skyline points):
	// [[100 16] [105 17] [110 18] [115 19] [120 20]]
}
