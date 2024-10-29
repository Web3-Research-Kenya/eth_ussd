package main

import (
	"fmt"
	"strings"
)

// MenuNode represents a node in the menu tree.
type MenuNode struct {
	name     string
	children map[string]*MenuNode
}

// NewMenuNode creates a new MenuNode.
func NewMenuNode(name string) *MenuNode {
	return &MenuNode{
		name:     name,
		children: make(map[string]*MenuNode),
	}
}

// AddChild adds a child node to the current MenuNode.
func (m *MenuNode) AddChild(number string, node *MenuNode) {
	m.children[number] = node
}

// MenuTree represents the entire menu structure.
type MenuTree struct {
	root *MenuNode
}

// NewMenuTree creates a new MenuTree with a root node.
func NewMenuTree() *MenuTree {
	return &MenuTree{
		root: NewMenuNode("Root"),
	}
}

// AddMenu adds a menu item based on the provided path.
func (mt *MenuTree) AddMenu(path []string, name string) {
	current := mt.root
	for _, p := range path {
		if _, exists := current.children[p]; !exists {
			current.AddChild(p, NewMenuNode(name+" under "+current.name))
		}
		current = current.children[p]
	}
}

// Navigate traverses the menu tree based on the input path string.
func (mt *MenuTree) Navigate(pathStr string) {
	path := strings.Split(pathStr, "*")
	current := mt.root
	stack := []*MenuNode{} // Stack to keep track of navigation history

	for _, p := range path {
		if child, exists := current.children[p]; exists {
			stack = append(stack, current) // Save current position
			current = child
			fmt.Printf("You are now at: %s\n", current.name)
		} else if p == "0" { // Handle 'back' option
			if len(stack) > 0 {
				current = stack[len(stack)-1] // Go back to the previous node
				stack = stack[:len(stack)-1]  // Remove the last node from the stack
				fmt.Printf("Went back to: %s\n", current.name)
			} else {
				fmt.Println("Already at the root, cannot go back further.")
			}
		} else {
			fmt.Printf("Invalid path: %s\n", pathStr)
			return
		}
	}
}

// Main function to demonstrate the usage of the MenuTree and Navigate function.
func main() {
	// Create the menu tree
	menuTree := NewMenuTree()

	// Adding menus based on the provided structure
	menuTree.AddMenu([]string{"1"}, "Create Account")
	menuTree.AddMenu([]string{"1", "1"}, "Phone Number")
	menuTree.AddMenu([]string{"2"}, "Account Details")
	menuTree.AddMenu([]string{"2", "1"}, "Phone Number")
	menuTree.AddMenu([]string{"2", "2"}, "Address")
	menuTree.AddMenu([]string{"3"}, "Send ETH")
	menuTree.AddMenu([]string{"3", "1"}, "Receiver's Phone Number")
	menuTree.AddMenu([]string{"3", "2"}, "Amount in ETH")
	menuTree.AddMenu([]string{"4"}, "Receive ETH")
	menuTree.AddMenu([]string{"4", "1"}, "This is your Account Phone Number")
	menuTree.AddMenu([]string{"5"}, "Buy goods/services")
	menuTree.AddMenu([]string{"5", "1"}, "Till Number")
	menuTree.AddMenu([]string{"5", "2"}, "Amount in ETH")

	// Example of navigating through the menu tree using a long traversal string
	longNavigationString := "1*0" // A long navigation string

	fmt.Printf("\nNavigating: %s\n", longNavigationString)
	menuTree.Navigate(longNavigationString)
}
