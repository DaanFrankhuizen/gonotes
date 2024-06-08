package main

import (
	"fmt"
	"os"
)

type Note struct {
	Title string
	Body  []byte
}

func (n *Note) save() error {
	filename := n.Title + ".txt"
	return os.WriteFile(filename, n.Body, 0600)
}

func loadNote(title string) (*Note, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return &Note{Title: title, Body: body}, nil
}

func main() {
	n1 := &Note{Title: "testNote", Body: []byte("This is a test note.")}
	n1.save()
	n2, _ := loadNote("testNote")
	fmt.Println(string(n2.Body))
}
