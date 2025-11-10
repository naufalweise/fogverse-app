package brominmax

import "errors"

type AlgoType string

const (
	TypeBromin AlgoType = "bromin"
	TypeBromax          = "bromax"
)

type Config struct {
	AlgoType `json:"algorithm" default:"bromin"`
	T        float64 `json:"T" default:"100e6"`    // Target throughput, in byte/s.
	L        float64 `json:"L" default:"200e-3"`   // Replication latency threshold, ins.
	U        float64 `json:"U" default:"2000e-3"`  // Application unavailability threshold, in s.
	Tp       float64 `json:"Tp" default:"10e6"`    // Production throughput, in byte/s.
	Tc       float64 `json:"Tc" default:"20e6"`    // Consumption throughput, in byte/s.
	Hmax     int     `json:"Hmax" default:"30000"` // Hmax, maximum open files handles
	Lr       float64 `json:"lr" default:"1e-3"`    // Replication latency, in s.
	Uu       float64 `json:"u" default:"5e-3"`     // Observed unavailability time, in s.
	C        int     `json:"c" default:"100"`      // Number of consumers
	R        int     `json:"r" default:"3"`        // Replication factor
	B        int     `json:"B" default:"10"`       // Number of brokers available
}

// GetPartitions returns partition, replica, and error
func GetPartitions(config Config) (int, int, error) {
	if config.AlgoType == TypeBromin {
		return Bromin(config)
	} else if config.AlgoType == TypeBromax {
		return Bromax(config)
	}
	return -1, -1, errors.New("invalid algo type")
}

// helper function, get minimum allowed partition according to constraint/objective
func getPmin(config Config) int {
	return max(
		int(config.T/config.Tp),
		int(config.T/config.Tc),
		config.C)
}

// check if p and b satisfy the algo constraint
func checkFeasible(config Config, p int, b int) bool {
	return float64(p*config.R)*config.Lr <= float64(b)*config.L &&
		float64(p)*config.Uu <= float64(b)*config.U
}

// Bromin returns the partitions and brokers
func Bromin(config Config) (int, int, error) {
	pMin := getPmin(config)

	for b := config.R; b <= config.B; b++ {
		for p := b * config.Hmax / config.R; p >= pMin; p-- {
			if checkFeasible(config, p, b) {
				return p, b, nil
			}
		}
	}
	return 0, 0, errors.New("unable to find partition")
}

// Bromax returns the partitions and brokers
func Bromax(config Config) (int, int, error) {
	pMin := getPmin(config)

	for b := config.B; b >= config.R; b-- {
		for p := b * config.Hmax / config.R; p >= pMin; p-- {
			if checkFeasible(config, p, b) {
				return p, b, nil
			}
		}
	}
	return 0, 0, errors.New("unable to find partition")
}
