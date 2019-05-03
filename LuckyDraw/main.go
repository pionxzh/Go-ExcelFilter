package main

import "github.com/lxn/walk"
import UI "github.com/lxn/walk/declarative"

import (
	"excelFilter/util"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
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
	drawList 	*walk.TextEdit
	luckyList	*walk.TextEdit
	openFileBtn *walk.PushButton
	dragFileBtn *walk.PushButton
	drawBtn		*walk.PushButton
	resetBtn	*walk.PushButton
	winnerNum	*walk.NumberEdit

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

func (mw *MyMainWindow) handleFile(file string) {
	csvData, err := util.ReadCsvData(file)
	if err != nil {
		log.Fatalln(err)
	}

	data := util.TwoDimensionTo1DArray(csvData, ",")
	str := strings.Join(data, "\r\n")
	mw.drawList.SetText(str)
}

func (mw *MyMainWindow) showAbout() {
	walk.MsgBox(mw, "關於", "Developed by Pionxzh\r", walk.MsgBoxIconInformation)
}

func filterEmptyString(data []string) []string{
	result := util.Filter(data, func(str string) bool {
		if str != "" {
			return true
		}

		return false
	})

	return result
}

func main() {
	rand.Seed(time.Now().UnixNano())

	mw := &MyMainWindow{}
	var openAction, showAboutBoxAction *walk.Action

	MW := UI.MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "抽獎系統",
		MinSize:  UI.Size{Width: 600, Height: 400},
		Size:     UI.Size{Width: 1000, Height: 700},
		Layout:   UI.VBox{},
		OnDropFiles: func(file []string) {
			mw.handleFile(file[0])
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
							mw.handleFile(mw.filePath)
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
					UI.GroupBox{
						Title:              "抽獎欄",
						Layout:             UI.Grid{Columns: 12},
						Children: []UI.Widget{
							UI.HSplitter{
								ColumnSpan: 12,
								Children: []UI.Widget{
									UI.PushButton{
										AssignTo: &mw.openFileBtn,
										Text:     "開啟CSV文件",
										MinSize:  UI.Size{Height: 150},
										OnClicked: func() {
											if mw.openFileDialog() != true {
												return
											}
											mw.handleFile(mw.filePath)
										},
									},
									UI.PushButton{
										AssignTo: &mw.dragFileBtn,
										Text:     "或直接拖曳檔案",
										MinSize:  UI.Size{Height: 150},
									},
								},
							},
							UI.Composite{
								Layout: UI.Grid{Columns: 2},
								Children: []UI.Widget{
									UI.Label{
										Text: "中獎人數:",
									},
									UI.NumberEdit{
										AssignTo: &mw.winnerNum,
										Value: float64(3),
									},
								},
							},
							UI.TextEdit{
								AssignTo: &mw.drawList,
								Text: "王大明\r\n李鐵花\r\n野比大雄\r\nSteve Jobs\r\nSakura\r\n",
								ColumnSpan: 12,
							},
						},
					},
					UI.GroupBox{
						Title:              "中獎欄",
						Layout:             UI.Grid{Rows: 2},
						Children: []UI.Widget{
							UI.TextEdit{
								AssignTo: &mw.luckyList,
							},
						},
					},

				},
			},
			UI.PushButton{
				AssignTo: &mw.drawBtn,
				Text:     "抽獎",
				MinSize:  UI.Size{Height: 100},
				OnClicked: func() {
					winnerNum := int(mw.winnerNum.Value())
					// fuck it = =, window use the \r\n
					drawPool := strings.Split(mw.drawList.Text(), "\r\n")
					drawPool = filterEmptyString(drawPool)

					if winnerNum > len(drawPool) {
						mw.ShowDialog("提示", "中獎人數不得大於抽獎人數")
						return
					}

					mw.luckyList.SetText("")
					for i := 0; i < winnerNum; i++ {
						index := rand.Intn(len(drawPool))
						mw.luckyList.AppendText(drawPool[index] + "\r\n")
						drawPool = util.RemoveArrayItem(drawPool, index)
					}

					mw.ShowDialog("提示", "抽獎完成")
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
