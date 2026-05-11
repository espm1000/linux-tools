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
	Line5 string
}

func menuText() Menu {
	return Menu{
		Line1: "0 --- Install Initial Dependencies (Debian)",
		Line2: "1 --- Run Tool",
		Line3: "2 --- Install Dev Tools (Debian only)",
		Line4: "3 --- Install Docker",
		Line5: "4 --- Generate Templates",
	}
}

func DisplayMenu() (string, error) {
	var selection string
	text := menuText()
	fmt.Println(text.Line1)
	fmt.Println(text.Line2)
	fmt.Println(text.Line3)
	fmt.Println(text.Line4)
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
