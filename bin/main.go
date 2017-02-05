package main

import "../crawler"
import "fmt"

func main() {
	lines := crawler.ActiveLines()

	fmt.Println("Trams:")
	fmt.Println(lines.Trams)
	fmt.Println("Trolleys")
	fmt.Println(lines.Trolleys)
	fmt.Println("Urban bus lines")
	fmt.Println(lines.Buses)
	fmt.Println("Suburban bus lines")
	fmt.Println(lines.SuburbanBuses)
	fmt.Println("Subway")
	fmt.Println(lines.SubwayLines)
}
