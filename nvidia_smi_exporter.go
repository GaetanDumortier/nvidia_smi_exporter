package main

import (
    "bytes"
    "encoding/csv"
    "fmt"
    "net/http"
    "log"
    "os"
    "os/exec"
    "strings"
    "strconv"
)

func metrics(response http.ResponseWriter, request *http.Request) {
    out, err := exec.Command(
        "nvidia-smi",
        "--query-gpu=name,index,count,temperature.gpu,utilization.gpu,utilization.memory,memory.total,memory.free,memory.used,power.draw",
        "--format=csv,noheader,nounits").Output()

    if err != nil {
        fmt.Printf("%s\n", err)
        return
    }

    csvReader := csv.NewReader(bytes.NewReader(out))
    csvReader.TrimLeadingSpace = true
    records, err := csvReader.ReadAll()

    if err != nil {
        fmt.Printf("%s\n", err)
        return
    }

    metricList := []string {
        "count",
        "temperature.gpu",
        "utilization.gpu",
        "utilization.memory", "memory.total", "memory.free", "memory.used", "power.draw"}

    result := ""
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