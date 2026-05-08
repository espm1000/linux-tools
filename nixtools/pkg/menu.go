package pkg

import (
	"fmt"
	"os"
	"time"
)

type Menu struct {
	Line1 string
	Line2 string
	Line3 string
	Line4 string
}

func menuText() Menu {
	return Menu{
		Line1: "1 --- Run Tool",
		Line2: "2 --- Install Dev Tools (Debian only)",
		Line3: "3 --- Exit",
	}
}

func DisplayMenu() (string, error) {
	var selection string
	text := menuText()
	fmt.Println(text.Line1)
	fmt.Println(text.Line2)
	fmt.Println(text.Line3)
	fmt.Print("Selection: ")
	if _, err := fmt.Scan(&selection); err != nil {
		return "", err
	}
	time.Sleep(500 * time.Millisecond)
	return selection, nil
}

func Exit() {
	os.Exit(0)
}
