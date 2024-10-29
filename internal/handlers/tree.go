package handlers

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

const (
	createAccount  = "createAccount.tmpl"
	root           = "root.tmpl"
	end            = "end.tmpl"
	phoneNumber    = "phoneNumber.tmpl"
	recieveEth     = "recieveEth.tmpl"
	sendEth        = "sendEth.tmpl"
	buyGoods       = "buyGoods.tmpl"
	amount         = "amount.tmpl"
	accountDetails = "accountDetails.tmpl"
)

// MenuNode represents a node in the menu tree.
type MenuNode struct {
	tmplName string
	children map[string]*MenuNode
}

// NewMenuNode creates a new MenuNode.
func NewMenuNode(name string) *MenuNode {
	return &MenuNode{
		tmplName: name,
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
	mu   sync.RWMutex
}

// NewMenuTree creates a new MenuTree with a root node.
func NewMenuTree() *MenuTree {
	return &MenuTree{
		root: NewMenuNode(root),
	}
}

func (mt *MenuTree) AddMenu(path []string, tmplName string) {
	mt.mu.Lock()
	defer mt.mu.Unlock()
	current := mt.root
	for _, p := range path {
		if _, exists := current.children[p]; !exists {
			current.AddChild(p, NewMenuNode(tmplName))
		}
		current = current.children[p]
	}
}

// Navigate traverses the menu tree based on the input path string.
func (mt *MenuTree) Navigate(pathStr *string) string {
	mt.mu.RLock()
	defer mt.mu.RUnlock()
	path := strings.Split(*pathStr, "*")
	current := mt.root
	stack := []*MenuNode{} // Stack to keep track of navigation history

	for _, p := range path {
		fmt.Printf("%s -> ", p)
		if child, exists := current.children[p]; exists {

			stack = append(stack, current) // Save current position
			current = child
			fmt.Printf("You are now at: %s\n", current.tmplName)
		} else if p == "0" { // Handle 'back' option
			if len(stack) > 0 {
				fmt.Printf("You are now at: %s\n", current.tmplName)
				current = stack[len(stack)-1] // Go back to the previous node
				stack = stack[:len(stack)-1]  // Remove the last node from the stack
			} else {
				return root
			}
		} else {
			log.Panicf("Invalid path: %s\n", *pathStr)
			return end
		}
	}
	pathStr = nil
	return current.tmplName
}
