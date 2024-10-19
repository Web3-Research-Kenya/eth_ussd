package handlers

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/gofiber/fiber/v2"
)

var lastFile string

const (
	intro           = "intro.tmpl"
	register        = "register.tmpl"
	learn           = "learn.tmpl"
	end             = "exit.tmpl"
	registerSuccess = "registerSuccess.tmpl"
)

type Data struct {
	PhoneNumber string
}

func CallbackHandler(c *fiber.Ctx) error {

	var text = c.FormValue("text", "")
	var phoneNumber = c.FormValue("phoneNumber")
	var data Data = Data{
		PhoneNumber: phoneNumber,
	}

	if text == "" {
		tmpl, err := template.ParseFiles("templates/intro.tmpl")
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return err
		}
		_, err = c.WriteString(buf.String())

		return err
	}

	index := getLast(text)

	fmt.Println("index: ", index)
	var file string

	if index != "0" {
		switch text {
		case "1":
			file = register
		case "1*1":
			file = registerSuccess
			lastFile = register
		case "1*2":
			file = registerSuccess
			lastFile = register
		case "2":
			file = learn
			lastFile = intro
		case "3":
			file = end
			lastFile = intro
		}
	} else {
		file = lastFile
	}

	var templateFile string = fmt.Sprintf("templates/%s", file)

	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	_, err = c.WriteString(buf.String())
	fmt.Println("last file: ", lastFile)

	return err
}

func getLast(input string) string {
	parts := strings.Split(input, "*")

	result := parts[len(parts)-1]

	return result
}
