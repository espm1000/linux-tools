package pkg

import (
	"fmt"
	"time"
)

type Menu struct {
	Line1 string
	Line2 string
	Line3 string
}

func menuText() Menu {
	m := Menu{
		Line1: "1 --- Run Tool",
		Line2: "2 --- Generate Configuration Report",
		Line3: "3 --- Exit",
	}

	return m
}

func DisplayMenu() string {
	var selection string
	text := menuText()
	fmt.Println(text.Line1)
	fmt.Println(text.Line2)
	fmt.Println(text.Line3)
	fmt.Print("Selection: ")
	fmt.Scan(&selection)
	time.Sleep(1 * time.Second)
	return selection
}
