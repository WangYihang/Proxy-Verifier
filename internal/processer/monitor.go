package processer

import (
	"time"

	"github.com/WangYihang/Proxy-Verifier/internal"
)

func Monitor() {
	for {
		// Sleep for interval seconds
		time.Sleep(time.Duration(internal.Options.MonitorInterval) * time.Second)
		// Take a snapshot
		internal.State.Snapshot()
	}
}
