package ram

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func RetrieveMemoryStats() (*MemoryStats, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parseMemoryStats(file)
}

type MemoryStats struct {
	Total, Used, Buffers, Cached, Free, Available, Active, Inactive,
	SwapTotal, SwapUsed, SwapCached, SwapFree, Mapped, Shmem, Slab,
	PageTables, Committed, VmallocUsed uint64
	MemAvailableEnabled bool
}

func parseMemoryStats(reader io.Reader) (*MemoryStats, error) {
	scanner := bufio.NewScanner(reader)
	var memory MemoryStats
	fieldMap := map[string]*uint64{
		"MemTotal":     &memory.Total,
		"MemFree":      &memory.Free,
		"MemAvailable": &memory.Available,
		"Buffers":      &memory.Buffers,
		"Cached":       &memory.Cached,
		"Active":       &memory.Active,
		"Inactive":     &memory.Inactive,
		"SwapCached":   &memory.SwapCached,
		"SwapTotal":    &memory.SwapTotal,
		"SwapFree":     &memory.SwapFree,
		"Mapped":       &memory.Mapped,
		"Shmem":        &memory.Shmem,
		"Slab":         &memory.Slab,
		"PageTables":   &memory.PageTables,
		"Committed_AS": &memory.Committed,
		"VmallocUsed":  &memory.VmallocUsed,
	}
	for scanner.Scan() {
		line := scanner.Text()
		colonIndex := strings.IndexRune(line, ':')
		if colonIndex < 0 {
			continue
		}
		fieldName := line[:colonIndex]
		if fieldPointer := fieldMap[fieldName]; fieldPointer != nil {
			value := strings.TrimSpace(strings.TrimRight(line[colonIndex+1:], "kB"))
			if parsedValue, err := strconv.ParseUint(value, 10, 64); err == nil {
				*fieldPointer = parsedValue * 1024
			}
			if fieldName == "MemAvailable" {
				memory.MemAvailableEnabled = true
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при сканировании /proc/meminfo: %s", err)
	}

	memory.SwapUsed = memory.SwapTotal - memory.SwapFree

	if memory.MemAvailableEnabled {
		memory.Used = memory.Total - memory.Available
	} else {
		memory.Used = memory.Total - memory.Free - memory.Buffers - memory.Cached
	}

	return &memory, nil
}
