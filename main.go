// Demo code for a timer based update
package main

import (
	"fmt"
	"slices"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	refreshInterval      = 500 * time.Millisecond
	rootPage             = "Root"
	timerPage            = "Timer"
	addTimerModelPage    = "modal-add-timer"
	deleteTimerModelPage = "modal-delete-timer"
	debugMode            = true
)

var (
	app    *tview.Application
	pages  *tview.Pages
	timers []*Timer
)

func updateTime(table *tview.Table) {
	ticker := time.NewTicker(refreshInterval)
	for range ticker.C {
		app.QueueUpdateDraw(func() {
			table.Clear()
			for i, timer := range timers {
				table.SetCell(i, 0, tview.NewTableCell(timer.Name))
				timeCell := tview.NewTableCell(timer.ElapsedTime().Round(time.Second).String())
				if timer.Running {
					timeCell.SetBackgroundColor(tcell.ColorGreen)
				} else {
					timeCell.SetBackgroundColor(tcell.ColorRed)
				}
				table.SetCell(i, 1, timeCell)
			}
		})
	}
}

// Returns a new primitive which puts the provided primitive in the center and
// sets its size to the given width and height.
func createCustomModal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}

func returnToPreviousPage() {
	lastPageName, _ := pages.GetFrontPage()
	if lastPageName != rootPage {
		pages.RemovePage(lastPageName)
	}
}

func debugPrint(msg string) {
	if debugMode {
		pages.SetBorder(true)
		pages.SetTitle(msg)
	}
}

func main() {
	app = tview.NewApplication()
	pages = tview.NewPages()
	pages.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			returnToPreviousPage()
			return nil
		}
		return event
	})

	timers = []*Timer{{Name: "Amazing"}}
	timers[0].Start()
	timerTable := tview.NewTable().SetSelectable(true, false)
	timerTable.SetBorder(true)
	go updateTime(timerTable)
	frame := tview.NewFrame(timerTable).
		SetBorders(0, 0, 1, 0, 0, 0).
		AddText("<esc>   Go to previous page", true, tview.AlignLeft, tcell.ColorGrey).
		AddText("<a>     Add", true, tview.AlignLeft, tcell.ColorGrey).
		AddText("<d>     Delete", true, tview.AlignLeft, tcell.ColorGrey).
		AddText("<enter> Start/Stop", true, tview.AlignLeft, tcell.ColorGrey).
		AddText("-> Timer ", false, tview.AlignLeft, tcell.ColorDarkGreen)
	frame.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'a':
				inputField := tview.NewInputField().
					SetLabel("Timer's name: ").
					SetFieldWidth(20)
				inputField.SetDoneFunc(func(key tcell.Key) {
					switch key {
					case tcell.KeyEnter:
						defer returnToPreviousPage()
						newTimerName := inputField.GetText()
						if newTimerName == "" {
							return
						}
						newTimer := Timer{Name: newTimerName}
						newTimer.Start()
						timers = append(timers, &newTimer)
					}
				})
				pages.AddAndSwitchToPage(addTimerModelPage, createCustomModal(inputField, 30, 20), false)
				return nil
			case 'd':
				row, _ := timerTable.GetSelection()
				timerName := timerTable.GetCell(row, 0).Text
				timerIndex := slices.IndexFunc(timers, func(timer *Timer) bool {
					return timerName == timer.Name
				})
				if timerIndex == -1 {
					return nil
				}

				modal := tview.NewModal().
					SetText(fmt.Sprintf("Are you sure you want to delete '%v'?", timerName)).
					AddButtons([]string{"No", "Yes"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						defer returnToPreviousPage()
						if buttonLabel == "Yes" {
							timers = slices.Delete(timers, timerIndex, timerIndex+1)
						}
					})
				pages.AddAndSwitchToPage(deleteTimerModelPage, modal, false)
				return nil
			}
		case tcell.KeyEnter:
			row, _ := timerTable.GetSelection()
			timerName := timerTable.GetCell(row, 0).Text
			timerIndex := slices.IndexFunc(timers, func(timer *Timer) bool {
				return timerName == timer.Name
			})
			if timerIndex == -1 {
				return nil
			}
			timer := timers[timerIndex]
			if !timer.Running {
				timer.Start()
			} else {
				timer.Pause()
			}
			return nil
		}
		return event
	})

	list := tview.NewList().
		AddItem("Timer", "", 't', func() {
			pages.AddPage(timerPage, frame, true, true)
		}).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		}).
		ShowSecondaryText(false)

	// add root page
	pages.AddPage(rootPage, list, true, true)

	err := app.SetRoot(pages, true).Run()
	if err != nil {
		panic(err)
	}
}
