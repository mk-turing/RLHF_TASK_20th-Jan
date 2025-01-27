package main
import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/benbjohnson/phantom"
	"github.com/go-kit/kit/log/term"
	"github.com/matryer/freeticks"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func main() {
	// Set GOMAXPROCS to utilize all available CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Run the benchmark suite
	benchmarkSuite()
}

func benchmarkSuite() {
	// Initialize the logger for terminal visualization
	var logger term.Logger

	// Create a new stress test with 10 seconds duration and 500ms ramp-up time
	test := phantom.NewTest(10*time.Second, 500*time.Millisecond)

	// Add an HTTP request load to the stress test
	test.AddLoad(phantom.HTTPLoad{
		URL:         "http://your-api-endpoint.com/resource",
		Method:      "GET",
		Concurrency: phantom.CappedConcurrency(10, 500), // Limit concurrency to 500
	})

	// Define the metrics we want to track
	var totalErrors uint64
	var cpuUsagePercent float64
	var memUsageMB float64

	// Run the stress test
	done := make(chan struct{})
	go func() {
		err := test.Run()
		if err != nil {
			logger.Log("ERROR", err)
		}
		close(done)
	}()

	// Ticker to update system resource usage and visualize metrics every second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var ft *freeticks.Ticks
	ft, err := freeticks.New(1, 3) // Initialize visualization with 1 second interval and 3 lines
	if err != nil {
		log.Fatal(err)
	}

	defer ft.Close()

	for {
		select {
		case <-ticker.C:
			// Update CPU usage
			stats, err := cpu.Times(false)
			if err == nil {
				cpuUsagePercent = stats[0].User + stats[0].System
			}

			// Update Memory usage
			vmem, err := mem.VirtualMemory()
			if err == nil {
				memUsageMB = float64(vmem.Used) / 1024 / 1024
			}

			// Log and visualize metrics
			logger.Log("TEST", "Status", "CONCURRENCY", test.Concurrency(), "ERRORS", totalErrors, "CPU", cpuUsagePercent, "MEM", memUsageMB)