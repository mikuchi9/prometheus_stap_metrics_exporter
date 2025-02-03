package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	processInfo_running_time = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "process_info_running_time",
			Help: "Running time of processes that fork or clone.",
		},
		[]string{"name", "pid", "ppid", "niceness"},
	)

	processInfo_memory_usage_size = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "process_info_memory_usage_size_mib",
			Help: "Memory size of processes that fork or clone.",
		},
		[]string{"name", "pid", "ppid", "niceness", "type"},
	)

	processInfo_memory_usage_rss = prometheus.NewGaugeVec( // resident set size
		prometheus.GaugeOpts{
			Name: "process_info_memory_usage_rss_mib",
			Help: "Resident Set Size of processes that fork or clone.",
		},
		[]string{"name", "pid", "ppid", "niceness", "type"},
	)

	processInfo_memory_usage_shr = prometheus.NewGaugeVec( // Shared Mem size
		prometheus.GaugeOpts{
			Name: "process_info_memory_usage_shr_mib",
			Help: "Shared memory size of processes that fork or clone.",
		},
		[]string{"name", "pid", "ppid", "niceness", "type"},
	)

	processInfo_memory_usage_txt = prometheus.NewGaugeVec( // Data+Stack size
		prometheus.GaugeOpts{
			Name: "process_info_memory_usage_txt_mib",
			Help: "Data and stack size of processes that fork or clone.",
		},
		[]string{"name", "pid", "ppid", "niceness", "type"},
	)

	processInfo_memory_usage_data = prometheus.NewGaugeVec( // Code size
		prometheus.GaugeOpts{
			Name: "process_info_memory_usage_data_mib",
			Help: "Code size of processes that fork or clone.",
		},
		[]string{"name", "pid", "ppid", "niceness", "type"},
	)
)

type metricStr []string

func (s metricStr) getProperValue(fieldNum int) float64 {

	s_M := strings.TrimSpace(strings.Split(s[fieldNum], ":")[1])
	s_number := strings.Replace(s_M, "M", "", -1)
	s_f64, _ := strconv.ParseFloat(s_number, 64)

	return s_f64
}

func init() {
	// Register the metric with Prometheus
	prometheus.MustRegister(processInfo_running_time)
	prometheus.MustRegister(processInfo_memory_usage_size)
	prometheus.MustRegister(processInfo_memory_usage_rss)
	prometheus.MustRegister(processInfo_memory_usage_shr)
	prometheus.MustRegister(processInfo_memory_usage_txt)
	prometheus.MustRegister(processInfo_memory_usage_data)

}

func updateMetrics() {

	
	cmd := exec.Command("staprun", "custom_metrics.ko")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error creating stdout pipe: %v\n", err)
		return 
	}
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting command: %v\n", err)
		return
	}
	scanner := bufio.NewScanner(stdout)
	fmt.Println("Metrics collector is running...")

	for scanner.Scan() {
		process_info := scanner.Text()
		pr_info_data := strings.Split(process_info, "@")[1] // selecting data not name
		pr_info_fields := strings.Split(pr_info_data, "|") // split by '|'
		running_time, _ := strconv.ParseFloat(strings.Split(pr_info_fields[4], "=")[1], 64)

		name := strings.Split(pr_info_fields[0], "=")[1]
		pid := strings.Split(pr_info_fields[1], "=")[1]
		ppid := strings.Split(pr_info_fields[2], "=")[1]
		niceness := strings.Split(pr_info_fields[3], "=")[1]
		mem_usage := strings.Split(pr_info_fields[5], "=")[1]
		
		var mem_usage_fields metricStr = strings.Split(mem_usage, ",")

		mem_usage_size := mem_usage_fields.getProperValue(0)
		mem_usage_rss := mem_usage_fields.getProperValue(1)
		mem_usage_shr := mem_usage_fields.getProperValue(2)
		mem_usage_txt := mem_usage_fields.getProperValue(3)
		mem_usage_data := mem_usage_fields.getProperValue(4)

		processInfo_running_time.With(
			prometheus.Labels{"name": name,
							  "pid": pid, 
							  "ppid": ppid,
							  "niceness": niceness,
							}).Set(running_time)

		processInfo_memory_usage_size.With(
			prometheus.Labels{"name": name,
							  "pid": pid, 
							  "ppid": ppid,
							  "niceness": niceness,
							  "type": "size",
							}).Set(mem_usage_size)
		
		processInfo_memory_usage_rss.With(
			prometheus.Labels{"name": name,
							  "pid": pid, 
							  "ppid": ppid,
							  "niceness": niceness,
							  "type": "rss",
							}).Set(mem_usage_rss)

		processInfo_memory_usage_shr.With(
			prometheus.Labels{"name": name,
							  "pid": pid, 
							  "ppid": ppid,
							  "niceness": niceness,
							  "type": "shr",
							}).Set(mem_usage_shr)

		processInfo_memory_usage_txt.With(
			prometheus.Labels{"name": name,
							  "pid": pid, 
							  "ppid": ppid,
							  "niceness": niceness,
							  "type": "txt",
							}).Set(mem_usage_txt)

		processInfo_memory_usage_data.With(
			prometheus.Labels{"name": name,
							  "pid": pid, 
							  "ppid": ppid,
							  "niceness": niceness,
							  "type": "data",
							}).Set(mem_usage_data)
	}
	cmd.Wait()

}

func main() {
	// Start the Prometheus HTTP server
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		for {
			updateMetrics()
		}
	}()
	
	http.ListenAndServe(":9100", nil)
}
