// Demo code for the TreeView primitive.
package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

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
			node := tview.NewTreeNode(file.Name()).
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

	flex := tview.NewFlex().
		AddItem(tree, 0, 1, true)

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
			fileActionMenu := tview.NewList().
				AddItem("名前変更", "", 'r', nil).
				AddItem("削除", "", 'd', nil).
				AddItem("複製", "", 'c', nil).
				AddItem("閉じる", "", 'q', func() {
					flex.RemoveItem(app.GetFocus())
					app.SetFocus(tree)
				})
			flex.AddItem(fileActionMenu, 0, 1, false)
			app.SetFocus(fileActionMenu)
		}
	})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Name() {
		// quit action
		case "Rune[Q]":
			app.Stop()
			os.Exit(0)
		}
		return event
	})

	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}
