package docker

// Hashmap to store the docker status output only once since it takes 2 seconds to call docker stats ( for correct cpu percent calculation)
type DockerContainerStats struct {
	Stats map[string]StatType
}

type StatType struct {
	CpuPercent    float64
	MemoryPercent float64
}

func NewDockerContainerStats() *DockerContainerStats {
	return &DockerContainerStats{
		Stats: make(map[string]StatType),
	}
}
