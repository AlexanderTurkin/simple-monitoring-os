package cpu

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func RetrieveStats() (*CPUStats, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parseStats(file)
}

type CPUStats struct {
	UserTime, NiceTime, SystemTime, IdleTime, IOWaitTime, IRQTime, SoftIRQTime, StealTime, GuestTime, GuestNiceTime, TotalTime uint64
	NumCPUs, NumFields                                                                                                         int
}

type cpuFieldMapping struct {
	fieldName string
	fieldPtr  *uint64
}

func parseStats(reader io.Reader) (*CPUStats, error) {
	scanner := bufio.NewScanner(reader)
	var stats CPUStats

	fieldMappings := []cpuFieldMapping{
		{"user", &stats.UserTime},
		{"nice", &stats.NiceTime},
		{"system", &stats.SystemTime},
		{"idle", &stats.IdleTime},
		{"iowait", &stats.IOWaitTime},
		{"irq", &stats.IRQTime},
		{"softirq", &stats.SoftIRQTime},
		{"steal", &stats.StealTime},
		{"guest", &stats.GuestTime},
		{"guest_nice", &stats.GuestNiceTime},
	}

	if !scanner.Scan() {
		return nil, fmt.Errorf("не удалось прочитать /proc/stat")
	}

	valueStrings := strings.Fields(scanner.Text())[1:]
	stats.NumFields = len(valueStrings)
	for i, valStr := range valueStrings {
		val, err := strconv.ParseUint(valStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("ошибка при разборе %s из /proc/stat", fieldMappings[i].fieldName)
		}
		*fieldMappings[i].fieldPtr = val
		stats.TotalTime += val
	}

	stats.TotalTime -= stats.GuestTime
	stats.TotalTime -= stats.GuestNiceTime

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu") && unicode.IsDigit(rune(line[3])) {
			stats.NumCPUs++
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при сканировании /proc/stat: %s", err)
	}

	return &stats, nil
}
