package ui

import (
	"time"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
)

func StartUI(updateInterval time.Duration, getCPUUsage, getMemoryUsage, getDiskUsage, getNetworkUsage func() string) {
	err := ui.Main(func() {
		cpuLabel := ui.NewLabel("")
		memoryLabel := ui.NewLabel("")
		diskLabel := ui.NewLabel("")
		networkLabel := ui.NewLabel("")

		cpuGroup := ui.NewGroup("CPU Usage")
		cpuGroup.SetChild(cpuLabel)

		memoryGroup := ui.NewGroup("Memory Usage")
		memoryGroup.SetChild(memoryLabel)

		diskGroup := ui.NewGroup("Disk Usage")
		diskGroup.SetChild(diskLabel)

		networkGroup := ui.NewGroup("Network Usage")
		networkGroup.SetChild(networkLabel)

		box := ui.NewVerticalBox()
		box.Append(cpuGroup, false)
		box.Append(memoryGroup, false)
		box.Append(diskGroup, false)
		box.Append(networkGroup, false)

		window := ui.NewWindow("System Stats", 400, 400, false)
		window.SetChild(box)
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()

		go func() {
			for {
				cpuUsage := getCPUUsage()
				memoryUsage := getMemoryUsage()
				diskUsage := getDiskUsage()
				networkUsage := getNetworkUsage()

				ui.QueueMain(func() {
					cpuLabel.SetText(cpuUsage)
					memoryLabel.SetText(memoryUsage)
					diskLabel.SetText(diskUsage)
					networkLabel.SetText(networkUsage)
				})

				time.Sleep(updateInterval)
			}
		}()
	})
	if err != nil {
		panic(err)
	}
}
