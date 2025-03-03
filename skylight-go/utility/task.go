package utility

import (
	"time"
)

func StartPeriodTask(interval int, task func()) {
	ticker := time.NewTicker(time.Second * time.Duration(interval))
	defer ticker.Stop()

	for range ticker.C {
		task()
	}
}
