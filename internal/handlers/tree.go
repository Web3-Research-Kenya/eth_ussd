package handlers

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

// Template names as constants
const (
	createAccount  = "createAccount.tmpl"
	root           = "root.tmpl"
	end            = "end.tmpl"
	setPin         = "setPin.tmpl"
	createAccEnd   = "createAccEnd.tmpl"
	receiveEth     = "receiveEth.tmpl"
	sendEth        = "sendEth.tmpl"
	buyGoods       = "buyGoods.tmpl"
	amount         = "amount.tmpl"
	accountDetails = "accountDetails.tmpl"
)

type NavigationContext struct {
	Path     []string
	Data     *Data
	Response interface{}
	tree     *MenuTree
}

type NodeFunction func(context *NavigationContext) error

type MenuNode struct {
	tmplName    string
	children    map[string]*MenuNode
	executeFunc NodeFunction
	parent      *MenuNode
}

func NewMenuNode(name string) *MenuNode {
	return &MenuNode{
		tmplName: name,
		children: make(map[string]*MenuNode),
	}
}

type MenuTree struct {
	root *MenuNode
	mu   sync.RWMutex
}

func NewMenuTree() *MenuTree {
	return &MenuTree{
		root: NewMenuNode(root),
	}
}

func (mt *MenuTree) AddNodeToPath(basePath []string, option string, fn NodeFunction, tmplName string) error {
	// mt.mu.Lock()
	// defer mt.mu.Unlock()

	// Navigate to the base path
	current := mt.root
	for _, p := range basePath {
		next, exists := current.children[p]
		if !exists {
			return fmt.Errorf("base path %v does not exist", basePath)
		}
		current = next
	}

	// Create and add the new node
	newNode := NewMenuNode(tmplName)
	newNode.executeFunc = fn
	current.children[option] = newNode
	newNode.parent = current
	return nil
}

// AddNodeDynamic adds a new node during navigation without deadlocking
func (mt *MenuTree) AddNodeDynamic(currentPath []string, option string, fn NodeFunction, tmplName string) error {
	if len(currentPath) == 0 {
		return fmt.Errorf("path cannot be empty")
	}

	// Create new node
	newNode := NewMenuNode(tmplName)
	newNode.executeFunc = fn

	// Navigate to parent node
	current := mt.root
	for _, p := range currentPath {
		next, exists := current.children[p]
		if !exists {
			return fmt.Errorf("path segment %s does not exist", p)
		}
		current = next
	}

	// Add the new node
	current.children[option] = newNode
	newNode.parent = current

	return nil
}

// Navigate function with support for dynamic node addition
func (mt *MenuTree) Navigate(pathStr *string, d *Data) string {
	if pathStr == nil {
		return root
	}

	// mt.mu.Lock() // Use single lock for entire navigation
	// defer mt.mu.Unlock()

	path := strings.Split(*pathStr, "*")
	current := mt.root
	stack := make([]*MenuNode, 0, len(path))

	ctx := &NavigationContext{
		Path: path,
		Data: d,
		tree: mt,
	}

	for _, p := range path {
		if p == "0" {
			if len(stack) > 0 {
				current = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				continue
			}
			return root
		}

		child, exists := current.children[p]
		if !exists {
			log.Printf("Invalid menu option: %s", p)
			return end
		}

		stack = append(stack, current)
		current = child

		if current.executeFunc != nil {
			if err := current.executeFunc(ctx); err != nil {
				log.Printf("Error executing node function: %v", err)
				return end
			}
		}
	}

	*pathStr = ""
	return current.tmplName
}

// Helper function to print the menu structure
func (mt *MenuTree) PrintStructure() {
	mt.mu.RLock()
	defer mt.mu.RUnlock()

	fmt.Println("\nMenu Structure:")
	var printNode func(*MenuNode, string, string)
	printNode = func(node *MenuNode, prefix string, option string) {
		fmt.Printf("%s%s -> %s\n", prefix, option, node.tmplName)
		for opt, child := range node.children {
			printNode(child, prefix+"  ", opt)
		}
	}
	printNode(mt.root, "", "root")
}

// Example usage:
func ExampleUsage() {
	menuTree := NewMenuTree()

	// Create a function that adds a dynamic node
	createDynamicNode := func(ctx *NavigationContext) error {
		// Create a function for the new node
		newNodeFunc := func(ctx *NavigationContext) error {
			fmt.Println("Dynamic node executed!")
			return nil
		}

		// Add the new node
		err := ctx.tree.AddNodeDynamic(
			ctx.Path,       // current path
			"dynamic",      // new option
			newNodeFunc,    // function for the new node
			"dynamic.tmpl", // template
		)
		if err != nil {
			return fmt.Errorf("failed to add dynamic node: %w", err)
		}
		fmt.Println("Dynamic node added successfully!")
		return nil
	}

	// Add initial node
	initialPath := []string{}
	err := menuTree.AddNodeDynamic(
		initialPath,
		"1",
		createDynamicNode,
		"creator.tmpl",
	)
	if err != nil {
		log.Fatalf("Failed to add initial node: %v", err)
	}

	// Print initial structure
	menuTree.PrintStructure()

	// Navigate to trigger dynamic node creation
	path := "1"
	template := menuTree.Navigate(&path, &Data{})
	fmt.Printf("Navigated to template: %s\n", template)

	// Print updated structure
	menuTree.PrintStructure()

	// Try navigating to the newly created node
	path = "1*dynamic"
	template = menuTree.Navigate(&path, &Data{})
	fmt.Printf("Navigated to dynamic node template: %s\n", template)
}

