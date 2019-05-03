package main

import "github.com/lxn/walk"
import UI "github.com/lxn/walk/declarative"

import (
	"excelFilter/util"
	"log"
	"os"
	"path"
	"strings"
)

var logFile *os.File

// init will change the stdout to logging file
func init() {
	var err error
	logFile, err = os.OpenFile("errors.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalln("打開日誌文件失敗：", err)
	}

	log.SetOutput(logFile)
}

// MyMainWindow represents the MainWindow
type MyMainWindow struct {
	*walk.MainWindow
	filterGB    *walk.GroupBox
	actionGB    *walk.GroupBox
	filterRule  *walk.TextEdit
	openFileBtn *walk.PushButton
	dragFileBtn *walk.PushButton
	filterBtn   *walk.PushButton
	saveFileBtn *walk.PushButton

	filePath string
}

func (mw *MyMainWindow) openFileDialog() bool {
	dlg := new(walk.FileDialog)
	dlg.Title = "選擇檔案"
	dlg.Filter = "Csv files (*.csv)|*.csv|All files (*.*)|*.*"

	if ok, err := dlg.ShowOpen(mw); err != nil {
		log.Println("Error : OpenFileDialog Fail")
		panic(err)
	} else if !ok {
		log.Println("FileDialog: Cancel")
		return false
	}
	mw.filePath = dlg.FilePath
	return true
}

func (mw *MyMainWindow) handleFile(tv *MyTableView, model *TableModel, file string) {
	mw.filePath = file
	tv.handleFile(model, file)

	mw.openFileBtn.SetVisible(false)
	mw.dragFileBtn.SetVisible(false)

	tv.TableView.SetVisible(true)
	mw.filterBtn.SetVisible(true)
	mw.actionGB.SetVisible(true)
	mw.filterGB.SetVisible(true)
}

func (mw *MyMainWindow) showAbout() {
	walk.MsgBox(mw, "關於", "Developed by Pionxzh\r", walk.MsgBoxIconInformation)
}

func main() {
	model := new(TableModel)

	tv := &MyTableView{}
	mw := &MyMainWindow{}

	var openAction, showAboutBoxAction *walk.Action

	MW := UI.MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "CSV文件過濾器",
		MinSize:  UI.Size{Width: 600, Height: 400},
		Size:     UI.Size{Width: 1000, Height: 700},
		Layout:   UI.VBox{},
		OnDropFiles: func(file []string) {
			mw.handleFile(tv, model, file[0])
		},
		MenuItems: []UI.MenuItem{
			UI.Menu{
				Text: "&操作",
				Items: []UI.MenuItem{
					UI.Action{
						AssignTo: &openAction,
						Text:     "開啟檔案",
						Shortcut: UI.Shortcut{Modifiers: walk.ModControl, Key: walk.KeyO},
						OnTriggered: func() {
							if mw.openFileDialog() != true {
								return
							}
							mw.handleFile(tv, model, mw.filePath)
						},
					},
					UI.Action{
						Text:        "離開 (Exit)",
						OnTriggered: func() { mw.Close() },
					},
				},
			},
			UI.Menu{
				Text: "&關於",
				Items: []UI.MenuItem{
					UI.Action{
						AssignTo:    &showAboutBoxAction,
						Text:        "About",
						OnTriggered: mw.showAbout,
					},
				},
			},
		},
		ContextMenuItems: []UI.MenuItem{
			UI.ActionRef{Action: &showAboutBoxAction},
		},
		Children: []UI.Widget{
			UI.HSplitter{
				Children: []UI.Widget{
					UI.PushButton{
						AssignTo: &mw.openFileBtn,
						Text:     "選擇CSV文件",
						MinSize:  UI.Size{Height: 150},
						OnClicked: func() {
							if mw.openFileDialog() != true {
								return
							}
							mw.handleFile(tv, model, mw.filePath)
						},
					},
					UI.PushButton{
						AssignTo: &mw.dragFileBtn,
						Text:     "或直接拖曳檔案",
						MinSize:  UI.Size{Height: 150},
					},
				},
			},
			UI.TableView{
				AssignTo:              &tv.TableView,
				AlwaysConsumeSpace:    false,
				AlternatingRowBGColor: walk.RGB(239, 239, 239),
				Visible:               false,
				MinSize:               UI.Size{Height: 300},
				Model:                 model,
				Columns: []UI.TableViewColumn{
					{Title: "#"},
				},
			},
			UI.HSplitter{
				Children: []UI.Widget{
					UI.GroupBox{
						AssignTo:           &mw.filterGB,
						AlwaysConsumeSpace: false,
						Title:              "過濾規則",
						Visible:            false,
						Layout:             UI.Grid{Rows: 2},
						Children: []UI.Widget{
							UI.Label{Text: "請輸入過濾規則 一行一個\r\n空行會導致全部資料保留"},
							UI.TextEdit{
								AssignTo: &mw.filterRule,
							},
						},
					},
					UI.GroupBox{
						AssignTo:           &mw.actionGB,
						AlwaysConsumeSpace: false,
						Title:              "操作",
						Visible:            false,
						Layout:             UI.Grid{Rows: 2},
						Children: []UI.Widget{
							UI.VSplitter{
								Children: []UI.Widget{
									UI.PushButton{
										AssignTo: &mw.filterBtn,
										Text:     "過濾",
										MinSize:  UI.Size{Height: 150},
										OnClicked: func() {
											// fuck it = =, window use the \r\n
											rules := strings.Split(mw.filterRule.Text(), "\r\n")
											util.Map(rules, func(str string) string { return strings.Trim(str, " ") })

											// filter the data by rule
											result := FilterByListRule(model.source, rules)
											model.items = util.ArrayTo2DSlice(result, util.SeperatorStr)
											model.PublishRowsReset()
											mw.saveFileBtn.SetEnabled(true)

											mw.ShowDialog("提示", "過濾完成")
										},
									},
									UI.PushButton{
										AssignTo: &mw.saveFileBtn,
										Text:     "存檔",
										Enabled:  false,
										MinSize:  UI.Size{Height: 150},
										OnClicked: func() {
											ext := path.Ext(mw.filePath)
											fileName := strings.TrimSuffix(mw.filePath, ext)
											filePath := fileName + "[已過濾]" + ext
											err := util.WriteCsvToFile(filePath, model.items)
											if err != nil {
												mw.ShowDialog("提示", "寫入檔案失敗")
											} else {
												mw.ShowDialog("提示", "已儲存至\r\n"+filePath)
											}
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if _, err := MW.Run(); err != nil {
		logFile.Close()
		log.Fatalln(err)
		os.Exit(1)
	}

	defer logFile.Close()
}
