package network

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func RetrieveNetworkStats() ([]NetworkStats, error) {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parseNetworkStats(file)
}

type NetworkStats struct {
	InterfaceName    string
	ReceivedBytes    uint64
	TransmittedBytes uint64
}

func parseNetworkStats(reader io.Reader) ([]NetworkStats, error) {
	scanner := bufio.NewScanner(reader)
	var networkStats []NetworkStats
	for scanner.Scan() {
		lineParts := strings.SplitN(scanner.Text(), ":", 2)
		if len(lineParts) != 2 {
			continue
		}
		fields := strings.Fields(lineParts[1])
		if len(fields) < 16 {
			continue
		}
		interfaceName := strings.TrimSpace(lineParts[0])
		if interfaceName == "lo" {
			continue
		}
		receivedBytes, err := strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("ошибка при разборе полученных байтов интерфейса %s", interfaceName)
		}
		transmittedBytes, err := strconv.ParseUint(fields[8], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("ошибка при разборе переданных байтов интерфейса %s", interfaceName)
		}
		networkStats = append(networkStats, NetworkStats{
			InterfaceName:    interfaceName,
			ReceivedBytes:    receivedBytes,
			TransmittedBytes: transmittedBytes,
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при сканировании /proc/net/dev: %s", err)
	}
	return networkStats, nil
}
