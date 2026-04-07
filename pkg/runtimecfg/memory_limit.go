package runtimecfg

import (
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/martin-helmich/prometheus-nginxlog-exporter/log"
)

const cgroupUnlimitedThreshold = int64(1 << 60)

var cgroupMemoryLimitPaths = []string{
	"/sys/fs/cgroup/memory.max",
	"/sys/fs/cgroup/memory/memory.limit_in_bytes",
}

// ConfigureMemoryLimitFromCgroup sets GOMEMLIMIT from container cgroup limits.
func ConfigureMemoryLimitFromCgroup(logger *log.Logger, enabled bool, ratio float64) error {
	if !enabled {
		return nil
	}

	if ratio <= 0 || ratio > 1 {
		return fmt.Errorf("gomemlimit-ratio must be > 0 and <= 1, got %.3f", ratio)
	}

	if value, exists := os.LookupEnv("GOMEMLIMIT"); exists {
		logger.Infof("GOMEMLIMIT is already set to %q; skipping automatic cgroup memory limit detection", value)
		return nil
	}

	limitBytes, ok, err := detectCgroupMemoryLimitBytes()
	if err != nil {
		return err
	}

	if !ok {
		logger.Infof("cgroup memory limit was not detected; running with Go runtime default memory limit")
		return nil
	}

	goMemLimit := int64(float64(limitBytes) * ratio)
	if goMemLimit <= 0 {
		return fmt.Errorf("computed Go memory limit is not valid: %d", goMemLimit)
	}

	previous := debug.SetMemoryLimit(goMemLimit)
	logger.Infof(
		"configured Go memory limit from cgroup: limit_bytes=%d ratio=%.2f gomemlimit_bytes=%d previous_bytes=%d",
		limitBytes,
		ratio,
		goMemLimit,
		previous,
	)

	return nil
}

func detectCgroupMemoryLimitBytes() (int64, bool, error) {
	var firstErr error
	readAnyFile := false

	for _, path := range cgroupMemoryLimitPaths {
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}

			if firstErr == nil {
				firstErr = fmt.Errorf("could not read cgroup memory limit from %s: %w", path, err)
			}

			continue
		}

		readAnyFile = true
		limitBytes, hasLimit, err := parseCgroupMemoryLimitValue(string(data))
		if err != nil {
			if firstErr == nil {
				firstErr = fmt.Errorf("could not parse cgroup memory limit from %s: %w", path, err)
			}
			continue
		}

		if hasLimit {
			return limitBytes, true, nil
		}
	}

	if firstErr != nil {
		return 0, false, firstErr
	}

	if readAnyFile {
		return 0, false, nil
	}

	return 0, false, nil
}

func parseCgroupMemoryLimitValue(value string) (int64, bool, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || trimmed == "max" {
		return 0, false, nil
	}

	limitBytes, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil {
		return 0, false, err
	}

	if limitBytes <= 0 {
		return 0, false, nil
	}

	// cgroup v1 uses a very large sentinel to indicate "unlimited".
	if limitBytes >= cgroupUnlimitedThreshold {
		return 0, false, nil
	}

	return limitBytes, true, nil
}
