package handlers

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/gofiber/fiber/v2"
)

var menuTree *MenuTree

type Data struct {
	PhoneNumber string
}

func init() {
	menuTree = NewMenuTree()

	menuTree.AddMenu([]string{"1"}, createAccount)
	menuTree.AddMenu([]string{"1", "1"}, phoneNumber)
	menuTree.AddMenu([]string{"2"}, accountDetails)
	menuTree.AddMenu([]string{"2", "1"}, phoneNumber)
	menuTree.AddMenu([]string{"2", "2"}, accountDetails)
	menuTree.AddMenu([]string{"3"}, sendEth)
	menuTree.AddMenu([]string{"3", "1"}, phoneNumber)
	menuTree.AddMenu([]string{"3", "2"}, amount)
	menuTree.AddMenu([]string{"4"}, recieveEth)
	menuTree.AddMenu([]string{"4", "1"}, amount)
	menuTree.AddMenu([]string{"5"}, buyGoods)
	menuTree.AddMenu([]string{"5", "1"}, buyGoods)
	menuTree.AddMenu([]string{"5", "2"}, amount)

}

func CallbackHandler(c *fiber.Ctx) error {

	var text = c.FormValue("text", "")
	var phoneNumber = c.FormValue("phoneNumber")
	var data Data = Data{
		PhoneNumber: phoneNumber,
	}

	if text == "" {
		return render(root, c, data)
	}

	fmt.Println("TEXT: ", text)

	template := menuTree.Navigate(&text)

	return render(template, c, data)

}

func render(fileName string, c *fiber.Ctx, data Data) error {

	var templateFile string = fmt.Sprintf("templates/%s", fileName)

	tmpl, err := template.ParseFiles(templateFile)
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
