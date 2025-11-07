package brominmax

import "errors"

type AlgoType int

const (
	BrominType AlgoType = iota
	BromaxType
)

type Config struct {
	AlgoType
	TargetThroughput            int
	ProductionThroughput        int
	ConsumerThroughput          int
	Consumers                   int
	Replicas                    int
	MaxBrokers                  int // B
	UnavailabilityTime          int // u
	MaxOpenFiles                int // Hmax, maximum open files handles
	ReplicationLatency          int
	ReplicationLatencyThreshold int // L
	UnavailabilityThreshold     int // U
}

// GetPartitions returns partition, replica, and error
func GetPartitions(config Config) (int, int, error) {
	if config.AlgoType == BrominType {
		return Bromin(config)
	} else if config.AlgoType == BromaxType {
		return Bromax(config)
	}
	return -1, -1, errors.New("invalid algo type")
}

func Bromin(config Config) (int, int, error) {
	pMin := max(
		config.TargetThroughput/config.ProductionThroughput,
		config.TargetThroughput/config.ConsumerThroughput,
		config.Consumers)

	for b := 0; b < config.MaxBrokers; b++ {
		for p := b * config.MaxOpenFiles / config.Replicas; p >= pMin; p-- {
			if p*config.Replicas*config.ReplicationLatency <= b*config.ReplicationLatencyThreshold &&
				p*config.UnavailabilityTime <= b*config.UnavailabilityThreshold {
				return p, b, nil
			}
		}
	}
	return 0, 0, errors.New("unable to find partition")
}
func Bromax(config Config) (int, int, error) {
	pMin := max(
		config.TargetThroughput/config.ProductionThroughput,
		config.TargetThroughput/config.ConsumerThroughput,
		config.Consumers)

	for b := config.MaxBrokers; b >= config.Replicas; b-- {
		for p := b * config.MaxOpenFiles / config.Replicas; p >= pMin; p-- {
			if p*config.Replicas*config.ReplicationLatency <= b*config.ReplicationLatencyThreshold &&
				p*config.UnavailabilityTime <= b*config.UnavailabilityThreshold {
				return p, b, nil
			}
		}
	}
	return 0, 0, errors.New("unable to find partition")
}
