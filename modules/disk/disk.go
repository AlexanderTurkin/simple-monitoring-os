package disk

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func RetrieveDiskStats() ([]DiskStats, error) {
	file, err := os.Open("/proc/diskstats")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parseDiskStats(file)
}

type DiskStats struct {
	DeviceName      string
	ReadOperations  uint64
	WriteOperations uint64
}

func parseDiskStats(reader io.Reader) ([]DiskStats, error) {
	scanner := bufio.NewScanner(reader)
	var diskStats []DiskStats
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 14 {
			continue
		}
		deviceName := fields[2]
		readOps, err := strconv.ParseUint(fields[3], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("ошибка при разборе завершенных чтений %s", deviceName)
		}
		writeOps, err := strconv.ParseUint(fields[7], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("ошибка при разборе завершенных записей %s", deviceName)
		}
		diskStats = append(diskStats, DiskStats{
			DeviceName:      deviceName,
			ReadOperations:  readOps,
			WriteOperations: writeOps,
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при сканировании /proc/diskstats: %s", err)
	}
	return diskStats, nil
}
