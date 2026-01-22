//go:build windows
// +build windows

// based on https://github.com/shirou/gopsutil
package cpu

import (
	"fmt"
	"unsafe"

	"github.com/Av4JhG/system-monitoring/internal/core"
)

func getCPU() (*cpuData, error) {
	var lpIdleTime core.FILETIME
	var lpKernelTime core.FILETIME
	var lpUserTime core.FILETIME

	r, _, err := core.GetSystemTimes.Call(
		uintptr(unsafe.Pointer(&lpIdleTime)),
		uintptr(unsafe.Pointer(&lpKernelTime)),
		uintptr(unsafe.Pointer(&lpUserTime)))
	if r == 0 {
		return nil, fmt.Errorf("call GetSystemTimes error: %w", err)
	}

	LOT := 0.0000001
	HIT := LOT * 4294967296.0
	idle := (HIT * float64(lpIdleTime.DwHighDateTime)) + (LOT * float64(lpIdleTime.DwLowDateTime))
	user := (HIT * float64(lpUserTime.DwHighDateTime)) + (LOT * float64(lpUserTime.DwLowDateTime))
	kernel := (HIT * float64(lpKernelTime.DwHighDateTime)) + (LOT * float64(lpKernelTime.DwLowDateTime))
	system := kernel - idle

	return &cpuData{
		total:  user + system + idle,
		user:   user,
		system: system,
		idle:   idle,
	}, nil
}
