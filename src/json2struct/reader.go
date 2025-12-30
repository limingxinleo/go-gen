package json2struct

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/marhaupe/json2struct/pkg/editor"
	"github.com/marhaupe/json2struct/pkg/parse"
)

type Reader struct {
}

func NewReader() *Reader {
	return &Reader{}
}

func (r *Reader) ReadFromFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(data)
}

func (r *Reader) PackJson(j string) (string, error) {
	var data interface{}
	err := json.Unmarshal([]byte(j), &data)
	if err != nil {
		return "", err
	}

	bt, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(bt), nil
}

func (r *Reader) ParseFromString(j string) (parse.Node, error) {
	j, err := r.PackJson(j)
	if err != nil {
		return nil, err
	}

	return parse.ParseFromString(j)
}

func (r *Reader) ReadFromEditor() (parse.Node, error) {
	edit := editor.New()
	defer edit.Delete()
	edit.Display()

	userInput, _ := edit.Read()

	userInputNode, err := r.ParseFromString(userInput)
	if err == nil {
		return userInputNode, nil
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("You supplied invalid JSON. Continue editing? (Y/n) ")
		userAnswer, _ := reader.ReadString('\n')
		userAnswer = strings.TrimSpace(userAnswer)
		userWantsFix := len(userAnswer) == 0 || userAnswer[0] == 'y'
		if !userWantsFix {
			return nil, nil
		}
		fmt.Print("\033[1A\033[2K")
		edit.Display()
		userInput, _ = edit.Read()
		isValid := json.Valid([]byte(userInput))
		if isValid {
			return userInputNode, nil
		}
	}
}
