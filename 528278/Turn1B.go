
package main
import (
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"
	"github.com/pkg/profile"
)

func main() {
	// Start CPU profiling
	defer profile.Start(profile.CPUProfile).Stop()

	// Start memory profiling
	defer profile.Start(profile.MemProfile).Stop()

	fmt.Println("Starting application...")
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Channel to control the loop
	loopCh := make(chan struct{})

	// Simple UI application using a range loop
	go func() {
		for {
			select {
			case <-loopCh:
				return
			default:
				// Perform UI updates or computations here
			}
		}