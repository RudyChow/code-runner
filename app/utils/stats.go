package utils

import (
	"strings"

	"github.com/RudyChow/code-runner/app/common"
	"github.com/docker/docker/api/types"
)

//转换docker的统计数据
func ParseStat(stats []*types.StatsJSON) ([]*common.Stats, int64) {
	var minTime, maxTime int64
	result := make([]*common.Stats, 0)

	for _, stat := range stats {
		//没时间的话就886
		currTime := stat.PreRead.UnixNano() / 1e6
		if currTime <= 0 {
			continue
		}

		record := &common.Stats{}
		record.CurrentTime = currTime

		//记录时间
		if minTime == 0 {
			minTime = currTime
		}
		if maxTime == 0 {
			maxTime = currTime
		}
		if currTime <= minTime {
			minTime = currTime
		} else if currTime >= maxTime {
			maxTime = currTime
		}

		//cpu
		record.CPUPercent = CalculateCPUPercent(stat)
		//内存
		record.MemoryPercent = CalculateMemoryPercent(stat)

		result = append(result, record)
	}

	return result, maxTime - minTime
}

//计算内存百分比
func CalculateMemoryPercent(v *types.StatsJSON) float64 {
	var result = 0.0
	if v.MemoryStats.Limit != 0 {
		result = float64(v.MemoryStats.Usage) / float64(v.MemoryStats.Limit) * 100.0
	}
	return result
}

//计算cpu百分比
func CalculateCPUPercent(v *types.StatsJSON) float64 {
	var (
		cpuPercent = 0.0
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float64(v.CPUStats.CPUUsage.TotalUsage) - float64(v.PreCPUStats.CPUUsage.TotalUsage)
		// calculate the change for the entire system between readings
		systemDelta = float64(v.CPUStats.SystemUsage) - float64(v.PreCPUStats.SystemUsage)
	)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(v.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}
	return cpuPercent
}

//计算io读写
func CalculateBlockIO(v *types.StatsJSON) (blkRead uint64, blkWrite uint64) {
	for _, bioEntry := range v.BlkioStats.IoServiceBytesRecursive {
		switch strings.ToLower(bioEntry.Op) {
		case "read":
			blkRead = blkRead + bioEntry.Value
		case "write":
			blkWrite = blkWrite + bioEntry.Value
		}
	}
	return
}

//计算network
func calculateNetwork(v *types.StatsJSON) (float64, float64) {
	var rx, tx float64

	for _, v := range v.Networks {
		rx += float64(v.RxBytes)
		tx += float64(v.TxBytes)
	}
	return rx, tx
}
