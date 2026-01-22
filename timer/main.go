package main

import (
    "bufio"
    "fmt"
    "net"
	"time"
	"sync"
)
var tasks_time []time.Time
var tasks_description []string
var nearest int
var mu sync.Mutex

func main() {
	mu.Lock()
	nearest = -1
	mu.Unlock()
	
    go handleTimer()

	listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        fmt.Println("Unable to open port:", err)
        return
    }
    defer listener.Close() 

    fmt.Println("Server started on 8080", nearest)

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Connect acception errpr:", err)
            continue
        }

        go handleConnection(conn)
    }
}

// Print timers
func handleTimer() {
	mu.Lock()
	nearest = -1
	mu.Unlock()
	
	for {
		time.Sleep(1 * time.Second)
	/*fmt.Print("Timers: ")
	for _, task := range tasks {
        fmt.Print(" ", task.Sub(time.Now()))
    }
	fmt.Println()*/

		mu.Lock()
		now := time.Now()
		if nearest != -1 {
			if tasks_time[nearest].Before(now) {
				fmt.Printf("Timer \"%s\" reached", tasks_description[nearest])
			
				nearest = -1
				for i, _  := range tasks_time {
					if tasks_time[i].After(now) {
						if nearest == -1 {
							nearest = i
						} else if tasks_time[i].Before(tasks_time[nearest]) {
							nearest = i
						}
					}
				}
			}
			//fmt.Print(nearest)
			//fmt.Print("Nearest timer:", tasks_time[nearest].Sub(time.Now()))
		}
		fmt.Println(tasks_time, nearest)
		fmt.Println(tasks_description)
		mu.Unlock()
	}
}

func handleConnection(conn net.Conn) {
    fmt.Println("New connection:", conn.RemoteAddr())
    defer conn.Close()

    scanner := bufio.NewScanner(conn)

    for scanner.Scan() {
        message := scanner.Text()
        fmt.Println("New message:", message)
		if message[:3] == "set" {
			
			divider := findNthSpace(message, 3)
			
			parsedTime, err := time.Parse("02.01.2006 15:04:05", message[4:divider])
			
			if err != nil {
				fmt.Println("Time parsing error:", err)
				return
			}

			fmt.Println("New task named \"", message[divider+1:], "\" at", parsedTime)
			
			mu.Lock()
			tasks_time = append(tasks_time, parsedTime)
			tasks_description = append(tasks_description, message[divider+1:])

			if nearest == -1 {
				nearest = 0
				fmt.Println("Nearest = 0")
			} else if parsedTime.Before(tasks_time[nearest]) {
				nearest = len(tasks_time) - 1
			}
			mu.Unlock()
		}
		if message[:6] == "cancel" {
			
		}
        // _, err := conn.Write([]byte(""))
        // if err != nil {
        //     fmt.Println("", err)
        // }
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("Data read error:", err)
    }
}

