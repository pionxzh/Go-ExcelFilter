package main

import (
	"github.com/lxn/walk"
)

// MyDialog represents
type MyDialog struct {
	*walk.Dialog
	ui myDialogUI
}

// myDialogUI represents the alert dialog
type myDialogUI struct {
	text  *walk.TextLabel
	okBtn *walk.PushButton
}

/*
// 這東西太坑了...調不好
func (w *MyDialog) alertDialog(owner walk.Form, message string) (err error) {
	if w.Dialog, err = walk.NewDialog(owner); err != nil {
		return err
	}

	succeeded := false
	defer func() {
		if !succeeded {
			w.Dispose()
		}
	}()

	w.SetName("Dialog")
	if err := w.SetClientSize(walk.Size{Width: 300, Height:140}); err != nil {
		return err
	}
	if err := w.SetTitle(`提示`); err != nil {
		return err
	}

	// textLabel
	if w.ui.text, err = walk.NewTextLabelWithStyle(w, win.SS_CENTER); err != nil {
		return err
	}

	if err := w.ui.text.SetBounds(walk.Rectangle{X: 40, Y: 50, Width: 220, Height: 30}); err != nil {
		return err
	}

	if err := w.ui.text.SetAlignment(walk.Alignment2D(walk.AlignHCenterVCenter)); err != nil {
		return err
	}

	if err := w.ui.text.SetText(message); err != nil {
		return err
	}

	// okBtn
	if w.ui.okBtn, err = walk.NewPushButton(w); err != nil {
		return err
	}

	if err := w.ui.okBtn.SetBounds(walk.Rectangle{X: 120, Y: 100, Width: 60, Height: 24}); err != nil {
		return err
	}

	if err := w.ui.okBtn.SetText("OK"); err != nil {
		return err
	}

	w.ui.okBtn.Clicked().Attach(func() {
		w.Accept()
	})

	w.Run()

	succeeded = true

	return nil
}
*/

// ShowDialog will generate the dialog to show message
func (mw *MyMainWindow) ShowDialog(title, message string) {
	walk.MsgBox(mw, title, message, walk.MsgBoxIconInformation)
}

/*
func ShowDialog(owner walk.Form, message string) (int, error) {
	var dlg *walk.Dialog
	var acceptPB *walk.PushButton

	return UI.Dialog{
		AssignTo:      &dlg,
		Title:         "提示",
		DefaultButton: &acceptPB,
		MinSize:       UI.Size{Width: 250, Height: 140},
		Layout:        UI.VBox{},
		Children: []UI.Widget{
			UI.Label{
				Text:      message,
				Alignment: UI.AlignHCenterVCenter,
			},
			UI.Composite{
				Layout: UI.HBox{},
				Children: []UI.Widget{
					UI.PushButton{
						AssignTo: &acceptPB,
						Text:     "OK",
						OnClicked: func() {
							dlg.Accept()
						},
					},
				},
			},
		},
	}.Run(owner)
}
*/
