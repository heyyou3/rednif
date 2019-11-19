// Demo code for the TreeView primitive.
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Show a navigable tree view of the current directory.
func main() {
	app := tview.NewApplication()
	rootDir := "."
	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	// A helper function which adds the files and directories of the given path
	// to the given target node.
	add := func(target *tview.TreeNode, path string) {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			node := tview.NewTreeNode(fmt.Sprintf("[%s] %s", file.ModTime().Format(time.RFC3339), file.Name())).
				SetReference(filepath.Join(path, file.Name())).
				SetSelectable(true)
			if file.IsDir() {
				node.SetColor(tcell.ColorGreen)
			}
			target.AddChild(node)
		}
	}

	// Add the current directory to the root node.
	add(root, rootDir)

	pathToMove := tview.NewInputField().
		SetLabel("directory path: ").
		SetFieldWidth(20).
		SetAcceptanceFunc(tview.InputFieldMaxLength(100))

	grid := tview.NewGrid().
		SetRows(0, 2).
		SetColumns(0, 0).
		SetBorders(true).
		AddItem(pathToMove, 1, 0, 1, 2, 0, 0, false).
		AddItem(tree, 0, 0, 1, 1, 0, 0, true)

	fileActionMenu := tview.NewList().
		AddItem("名前変更", "", 'r', nil).
		AddItem("削除", "", 'd', nil).
		AddItem("複製", "", 'c', nil).
		AddItem("閉じる", "", 'q', func() {
			grid.RemoveItem(app.GetFocus())
			app.SetFocus(tree)
		})

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return
		}
		path := reference.(string)
		_, err := ioutil.ReadDir(path)
		if err == nil {
			children := node.GetChildren()
			if len(children) <= 0 {
				add(node, path)
			} else {
				node.SetExpanded(!node.IsExpanded())
			}
		} else {
			grid.AddItem(fileActionMenu, 0, 1, 1, 1, 0, 0, true)
			app.SetFocus(fileActionMenu)
		}
	})

	tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Name() {
		// quit action
		case "Rune[Q]":
			app.Stop()
			os.Exit(0)
		// go to file path input field
		case "Rune[g]":
			app.SetFocus(pathToMove)
		// display directory control menu
		case "Rune[m]":
			grid.AddItem(fileActionMenu, 0, 1, 1, 1, 0, 0, true)
			app.SetFocus(fileActionMenu)
		}
		return event
	})

	pathToMove.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		// move TreeView
		case tcell.KeyCtrlG:
			app.SetFocus(tree)
		// display input directory
		case tcell.KeyEnter:
			inputText := pathToMove.GetText()
			root = tview.NewTreeNode(inputText).
				SetColor(tcell.ColorRed)
			tree = tree.
				SetRoot(root).
				SetCurrentNode(root)
			add(root, inputText)
			app.SetAfterDrawFunc(func(screen tcell.Screen) {
				tree.Draw(screen)
			})
		}
		return event
	})

	if err := app.SetRoot(grid, true).SetFocus(grid).Run(); err != nil {
		panic(err)
	}
}
