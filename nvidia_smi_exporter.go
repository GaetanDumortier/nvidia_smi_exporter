package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	GpuLostOutput    = "GPU is lost"        // Substring to look for in output when GPU is broken
	GpuMetric        = "nvidia_smi_failure" // Metrics name
	GpuLost          = "gpu_lost"           // Metrics value when gpu is lost/broken
	NvidiaSmiFailure = "nvidia-smi_failure" // Metrics value when nvidia-smi command failed
)

func metrics(response http.ResponseWriter, request *http.Request) {
	out, err := exec.Command(
		"nvidia-smi",
		"--query-gpu=name,index,count,temperature.gpu,utilization.gpu,utilization.memory,memory.total,memory.free,memory.used,power.draw",
		"--format=csv,noheader,nounits").Output()
	outStr := strings.ToLower(string(out))
	result := ""

	if err != nil {
		// Check the output of nvidia-smi for "gpu is lost" to determine GPU failures
		if strings.Contains(outStr, strings.ToLower(GpuLostOutput)) {
			fmt.Println("ALERT: nvidia-smi reports a GPU is lost/broken")
			result = fmt.Sprintf("%s=%s", GpuMetric, GpuLost)
		} else {
			fmt.Println("ALERT: Error while running nvidia-smi command!")
			result = fmt.Sprintf("%s=%s", GpuMetric, NvidiaSmiFailure)
		}
		// fmt.Printf("%s\n", err)
		fmt.Fprintf(response, result)
		return
	}

	csvReader := csv.NewReader(bytes.NewReader(out))
	csvReader.TrimLeadingSpace = true
	records, err := csvReader.ReadAll()

	if err != nil {
		fmt.Println("Error csv")
		fmt.Printf("%s\n", err)
		return
	}

	metricList := []string{
		"count",
		"temperature.gpu",
		"utilization.gpu",
		"utilization.memory", "memory.total", "memory.free", "memory.used", "power.draw"}

	for _, row := range records {
		name := fmt.Sprintf("%s[%s]", row[0], row[1])
		for idx, value := range row[2:] {
			floatVal, _ := strconv.ParseFloat(value, 64)

			result = fmt.Sprintf(
				"%s%s{gpu=\"%s\"} %.2f\n",
				result,
				strings.Replace(metricList[idx], ".", "_", -1),
				strings.Replace(name, ".", "_", -1),
				floatVal)
		}
	}

	fmt.Fprintf(response, result)
}

func main() {
	addr := ":9101"
	if len(os.Args) > 1 {
		addr = ":" + os.Args[1]
	}

	http.HandleFunc("/metrics/", metrics)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
