package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	type Moive struct {
		Title  string
		Year   int  `json:"released"`
		Color  bool `json:"color,omitempty"`
		Actors []string
	}

	var moives = []Moive{
		{Title: "中国", Year: 1942, Color: false,
			Actors: []string{"Humphrey Bogart", "Ingrid Bergman"}},
		{Title: "Cool Hand Luke", Year: 1967, Color: true,
			Actors: []string{"Paul Newman"}},
		{Title: "Bullitt", Year: 1968, Color: true,
			Actors: []string{"Steve McQueen", "Jacqueline Bisset"}},
	}

	data, err := json.MarshalIndent(moives, "", "	")

	if err != nil {
		log.Fatalf("json marshaling failed: %s", err)
	}

	fmt.Printf("%s\n", data)
}
