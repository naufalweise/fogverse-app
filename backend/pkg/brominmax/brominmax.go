package brominmax

import "errors"

type AlgoType int

const (
	BrominType AlgoType = iota
	BromaxType
)

type Config struct {
	AlgoType
	T    int
	Tp   int
	Tc   int
	C    int
	R    int
	B    int // B
	Uu   int // u
	U    int // U
	Hmax int // Hmax, maximum open files handles
	Lr   int
	L    int // L
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
		config.T/config.Tp,
		config.T/config.Tc,
		config.C)

	for b := config.R; b <= config.B; b++ {
		for p := b * config.Hmax / config.R; p >= pMin; p-- {
			if p*config.R*config.Lr <= b*config.L &&
				p*config.Uu <= b*config.U {
				return p, b, nil
			}
		}
	}
	return 0, 0, errors.New("unable to find partition")
}
func Bromax(config Config) (int, int, error) {
	pMin := max(
		config.T/config.Tp,
		config.T/config.Tc,
		config.C)

	for b := config.B; b >= config.R; b-- {
		for p := b * config.Hmax / config.R; p >= pMin; p-- {
			if p*config.R*config.Lr <= b*config.L &&
				p*config.Uu <= b*config.U {
				return p, b, nil
			}
		}
	}
	return 0, 0, errors.New("unable to find partition")
}
