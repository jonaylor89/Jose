package main

import (
	"bufio"
	"fmt"
	"math/rand"
  "os"
  "sort"
	"strings"

	// "strconv"
	"time"
)

const (

	// Process States

	// CREATED : process created
	CREATED = iota

	// RUNNING : process running
	RUNNING

	// WAITING : process waiting
	WAITING

	// BLOCKED : process blocked
	BLOCKED

	// TERMINATED : process terminated
	TERMINATED
)

var (

	// ProcNum : PID for the highest process
	ProcNum int = 0
)

// Process : Running set of code
type Process struct {
	PID     int
	state   int
	runtime int
	memory  int
}

// Scheduler : Module to schedule process to run
type Scheduler struct {
	inMsg     chan *Process
	processes []*Process
}

// CreateProc : create a new process correctly
func CreateProc(runtime int, mem int) *Process {

	ProcNum++

	return &Process{
		PID:     ProcNum,
		state:   CREATED,
		runtime: runtime,
		memory:  mem,
	}
}

func remove(slice []*Process, s int) []*Process {
	return append(slice[:s], slice[s+1:]...)
}

// Run : Start the schedule and process execution
func (s *Scheduler) Run() {
	for {

		// Check for new processes to schedule
		select {
		case x, ok := <-s.inMsg:
			if ok {
        // New process ready to be executed

        s.processes = append(s.processes, x)
        
        // Naive priority algorithm
        sort.Slice(s.processes, func(i, j int) bool { return s.processes[i].runtime < s.processes[j].runtime })
			} else {
				// Channel is closed to execution must exit
				return
			}
		default:
			// No new processes
			break
		}

		for i, curProc := range s.processes {
			curProc.state = RUNNING

      // I'm assuming this will get much more complex beyond just subtracting runtime
      // Fortunately, as of now it is basic round robin execution
			curProc.runtime -= 10

			if curProc.runtime <= 0 {
				s.processes = remove(s.processes, i)
			} else {
				curProc.state = WAITING
      }

		  time.Sleep(200 * time.Millisecond)
		}

	}
}

func main() {

	rand.Seed(time.Now().UnixNano())

	// Message channel between main kernel and scheduler
	ch := make(chan *Process, 10)
	defer close(ch)

	s := Scheduler{
		inMsg:     ch,
		processes: []*Process{},
	}

	// Run the scheduler
	go s.Run()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("OS Shell")
	fmt.Println("---------------------")

	for {

		fmt.Print("==> ")
		text, err := reader.ReadString('\n')
		text = strings.ReplaceAll(text, "\n", "")
		if err != nil {
			fmt.Println("failed to read user input")
		}

		switch text {
		case "new":
			p := CreateProc(rand.Intn(500)+1, rand.Intn(100)+1)
			ch <- p
			fmt.Println("processes: ", len(s.processes), "; queue: ", len(ch))
		case "len":
			fmt.Println("processes: ", len(s.processes), "; queue: ", len(ch))
		case "dump":
			fmt.Println("process dump:")
			for _, proc := range s.processes {
				fmt.Println(*proc)
			}
		case "exit":
			fmt.Println("exiting simulator")
			return
		}
	}
}
