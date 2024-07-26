package main

import "github.com/rivo/tview"

// Center returns a new primitive which shows the provided primitive in its
// center, given the provided primitive's size.
func Center(width, height int, title string, p tview.Primitive) *tview.Flex {
	flx := tview.NewFlex().
	AddItem(nil, 0, 1, false).
	AddItem(tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(p, height, 1, true).
		AddItem(nil, 0, 1, false), width, 1, true).
	AddItem(nil, 0, 1, false)

	if title != "" {
		flx.SetTitle(title)
	}
	flx.SetBorder(true)

	return flx
}