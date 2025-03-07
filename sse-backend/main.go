package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type SSEValue struct {
	LabelKey   string  `json:"labelKey"`
	LabelValue float64 `json:"labelValue"`
}

type SSEData struct {
	CreatedAt string     `json:"createdAt"`
	Value     []SSEValue `json:"value"`
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	for i := 1; i <= 120; i++ {
		cpuUsage, _ := cpu.Percent(0, false)
		memStat, _ := mem.VirtualMemory()

		memoryUsageVal, err := strconv.ParseFloat(fmt.Sprintf("%.2f", memStat.UsedPercent), 64)
		if err != nil {
			fmt.Fprintf(w, "data: {\"error\": \"Failed to encode JSON\"}\n\n")
			flusher.Flush()
			continue
		}

		cpuUsageVal, err := strconv.ParseFloat(fmt.Sprintf("%.2f", cpuUsage[0]), 64)
		if err != nil {
			fmt.Fprintf(w, "data: {\"error\": \"Failed to encode JSON\"}\n\n")
			flusher.Flush()
			continue
		}

		data := SSEData{
			CreatedAt: time.Now().Format("15:04:05"),
			Value: []SSEValue{
				{
					LabelKey:   "memory",
					LabelValue: memoryUsageVal,
				},
				{
					LabelKey:   "cpu",
					LabelValue: cpuUsageVal,
				},
			},
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Fprintf(w, "data: {\"error\": \"Failed to encode JSON\"}\n\n")
			flusher.Flush()
			continue
		}
		fmt.Fprintf(w, "data: %s\n\n", jsonData)

		flusher.Flush()
		time.Sleep(1 * time.Second)
	}

	fmt.Fprintln(w, "event: close\ndata: Connection closed")
	flusher.Flush()
}

func main() {

	http.HandleFunc("/events", eventsHandler)

	fmt.Println("SSE server running on http://localhost:8001")
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		panic(err)
	}

}
