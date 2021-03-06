package main

import (
	"fmt"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"github.com/zsdevX/DarkEye/common"
	"github.com/zsdevX/DarkEye/ui"
	"os"
	"runtime"
	"time"
)

var (
	programName = "DarkEye"
)

func main() {
	runApp()
}

func runApp() {
	//加载配置
	_ = loadCfg()
	runtime.GOMAXPROCS(runtime.NumCPU())
	app := widgets.NewQApplication(len(os.Args), os.Args)
	app.SetWindowIcon(gui.NewQIcon5(":/qml/logo.ico"))
	app.SetQuitOnLastWindowClosed(false)
	//初始化窗口
	mainWin := ui.NewMainWindow(nil)
	//初始化数据
	initMainWin(mainWin)
	//托盘图标初始化
	sysTray := NewQSystemTrayIconWithCustomSlot(nil)
	sysTrayDaemon(sysTray, mainWin, app)
	//显示
	sysTray.Show()
	mainWin.Show()
	//通知
	sysTray.TriggerSlot()
	widgets.QApplication_Exec()
}

func initMainWin(mainWin *ui.MainWindow) {
	//FoFa
	LoadFoFa(mainWin)
	//subDomain
	LoadSubDomain(mainWin)
	//Spider
	LoadSpider(mainWin)
}

func sysTrayDaemon(sysTray *QSystemTrayIconWithCustomSlot, mainWin *ui.MainWindow, app *widgets.QApplication) {
	sysTray.SetIcon(gui.NewQIcon5(":/qml/logo.png"))
	sysTray.SetToolTip("白嫖神器")

	sysTrayMenu := widgets.NewQMenu(nil)
	fucker := sysTrayMenu.AddAction("信息收集神器")
	about := sysTrayMenu.AddAction("关于")
	quit := sysTrayMenu.AddAction("退出")
	sysTray.SetContextMenu(sysTrayMenu)

	sysTray.ConnectTriggerSlot(func() {
		sysTray.ShowMessage("信息", common.ProgramVersion, widgets.QSystemTrayIcon__Information, 5000)
	})

	fucker.ConnectTriggered(func(bool) {
		mainWin.Show()
	})

	about.ConnectTriggered(func(bool) {
		information := common.ProgramVersion
		widgets.QMessageBox_Information(nil, "信息", information,
			widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	})

	quit.ConnectTriggered(func(bool) {
		if widgets.QMessageBox_Information(nil, "信息", "客官再见，欢迎白嫖",
			widgets.QMessageBox__Ok, widgets.QMessageBox__Cancel) == widgets.QMessageBox__Ok {
			app.Quit()
		}
	})
}

//外部应使用goroutine调用
func gracefulStop(start, stop *widgets.QPushButton, runCtl chan bool) {
	//终止任务时避免卡顿
	sec := 0
	jumpOut := false
	tick := time.NewTicker(time.Second)
	stop.SetDisabled(true)
	for {
		select {
		case <-runCtl:
			jumpOut = true
		case <-tick.C:
			sec ++
			stop.SetText(fmt.Sprintf("等待%d秒", 60-sec))
		}
		if jumpOut {
			break
		}
	}
	start.SetEnabled(true)
	stop.SetText("停止")
}

type QSystemTrayIconWithCustomSlot struct {
	widgets.QSystemTrayIcon
	_ func() `slot:"triggerSlot"`
}

func logChannel(view *widgets.QTextEdit) (chan string, chan bool) {
	logC := make(chan string, 128)
	runCtl := make(chan bool, 1)
	view.SetReadOnly(true)
	go func() {
		for {
			log := <-logC
			view.Append(log)
		}
	}()
	return logC, runCtl
}
