package common

//容器参数
type ContainerOption struct {
	Image          string
	Cmd            []string
	SourceFilePath string
	TargetFilePath string
}

//容器结果
type ContainerResult struct {
	ID            string   `json:"container_id"`
	Result        string   `json:"code_result"`
	ExecutionTime int64    `json:"execution_time"`
	Stats         []*Stats `json:"stats"`
}

type Stats struct {
	CurrentTime   int64   `json:"current_time"`
	MemoryPercent float64 `json:"memory_percent"`
	CPUPercent    float64 `json:"cpu_percent"`
}
