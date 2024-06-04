package main

import (
	"fmt"
	"time"

	"sms/modules/cpu"
	"sms/modules/disk"
	"sms/modules/network"
	"sms/modules/ram"
	"sms/ui"
)

func main() {
	ui.StartUI(1*time.Second, getCPUUsage, getMemoryUsage, getDiskUsage, getNetworkUsage)
}

func getCPUUsage() string {
	beforeCPU, err := cpu.RetrieveStats()
	if err != nil {
		return fmt.Sprintf("Err: %s", err)
	}
	time.Sleep(1 * time.Second)
	afterCPU, err := cpu.RetrieveStats()
	if err != nil {
		return fmt.Sprintf("Err: %s", err)
	}
	totalCPU := float64(afterCPU.TotalTime - beforeCPU.TotalTime)
	return fmt.Sprintf(
		"user: %.1f%%, sys: %.1f%%, idle: %.1f%%",
		float64(afterCPU.UserTime-beforeCPU.UserTime)/totalCPU*100,
		float64(afterCPU.SystemTime-beforeCPU.SystemTime)/totalCPU*100,
		float64(afterCPU.IdleTime-beforeCPU.IdleTime)/totalCPU*100,
	)
}

func getMemoryUsage() string {
	memory, err := ram.RetrieveMemoryStats()
	if err != nil {
		return fmt.Sprintf("Err: %s", err)
	}
	return fmt.Sprintf(
		"total: %.1fG, used: %.1fG, free: %.1fG",
		float64(memory.Total)/(1024*1024*1024),
		float64(memory.Used)/(1024*1024*1024),
		float64(memory.Free)/(1024*1024*1024),
	)
}

func getDiskUsage() string {
	diskStats, err := disk.RetrieveDiskStats()
	if err != nil {
		return fmt.Sprintf("Err: %s", err)
	}
	var readTotal, writeTotal uint64
	for _, ds := range diskStats {
		readTotal += ds.ReadOperations
		writeTotal += ds.WriteOperations
	}
	return fmt.Sprintf(
		"Read: %d, Write: %d",
		readTotal,
		writeTotal,
	)
}

func getNetworkUsage() string {
	networkStats, err := network.RetrieveNetworkStats()
	if err != nil {
		return fmt.Sprintf("Err: %s", err)
	}
	var receivedTotal, transmittedTotal uint64
	for _, ns := range networkStats {
		receivedTotal += ns.ReceivedBytes
		transmittedTotal += ns.TransmittedBytes
	}
	return fmt.Sprintf(
		"RX: %.1fMB, TX: %.1fMB",
		float64(receivedTotal)/(1024*1024),
		float64(transmittedTotal)/(1024*1024),
	)
}
