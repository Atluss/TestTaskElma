// Package cpuStatus work with Linux CPU and broadcast func for clients
package cpuStatus

import (
	"fmt"
	"github.com/Atluss/TestTaskElma/pkg/v1/server/wsServer"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

var Broadcast = make(chan CPULoad) // broadcast CPU load channel

type CPULoad struct {
	CPU string
}

// getCPUSample CPU status in time
func getCPUSample() (idle, total uint64) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return
		}
	}
	return
}

// RunCPUBroadcast run broadcast CPU
func RunCPUBroadcast() {
	go getCpuLoad(Broadcast)
	go handleCPUBroadcast(Broadcast)
}

// getCpuLoad measure CPU load after 3 seconds
func getCpuLoad(broadcast chan<- CPULoad) {
	for {
		idle0, total0 := getCPUSample()
		time.Sleep(3 * time.Second)
		idle1, total1 := getCPUSample()

		idleTicks := float64(idle1 - idle0)
		totalTicks := float64(total1 - total0)
		cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks
		//log.Printf("%9.2f %%", cpuUsage)
		cp := CPULoad{CPU: fmt.Sprintf("%9.2f", cpuUsage)}
		broadcast <- cp
	}
}

// handleCPUBroadcast send CPU load all accepted connections
func handleCPUBroadcast(broadcast <-chan CPULoad) {
	for {
		msg := <-broadcast
		banKeys := wsServer.BanKeys.CloneMe()
		for client := range wsServer.Clients {

			if _, ok := banKeys[client.Key.Key]; ok {
				client.WriteClose(websocket.ClosePolicyViolation, "banned")
				continue
			}

			err := client.Conn.WriteJSON(msg)
			if err != nil {
				client.Stop(err)
			}
		}
	}
}
