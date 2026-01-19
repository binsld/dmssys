package main

import (
    "bufio"
    "fmt"
    "net"
	"time"
)
var tasks []time.Time
var nearest int

func main() {
	nearest := -1
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
	nearest := -1
	for {
		time.Sleep(1 * time.Second)
	/*fmt.Print("Timers: ")
	for _, task := range tasks {
        fmt.Print(" ", task.Sub(time.Now()))
    }
	fmt.Println()*/

		if nearest != -1 {
			//fmt.Print(nearest)
			fmt.Print("Nearest timer:", tasks[nearest].Sub(time.Now()))
		}
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
			if message[4:9] == "timer" {

				parsedTime, err := time.Parse("02-Jan-2006 15:04:05", message[10:])
    			if err != nil {
					fmt.Println("Time parsing error:", err)
        			return
    			}
				tasks = append(tasks, parsedTime)
				fmt.Println(tasks)
				fmt.Println(nearest)
				if nearest == -1 {
					nearest = 0
					fmt.Println("Nearest = 0")
				} else if parsedTime.Before(tasks[nearest]) {
					nearest = len(tasks) - 1
				}


			}
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

