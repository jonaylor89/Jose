package main

import (
	// "bufio"
	// "fmt"
	"log"
	// "os"
	// "strconv"
	// "strings"
	// "time"

	ui "github.com/gizak/termui/v3"

	"github.com/jonaylor89/John_Naylor_CMSC312_2019/config"
	"github.com/jonaylor89/John_Naylor_CMSC312_2019/cpu"
	"github.com/jonaylor89/John_Naylor_CMSC312_2019/sched"
	"github.com/jonaylor89/John_Naylor_CMSC312_2019/memory"
	"github.com/jonaylor89/John_Naylor_CMSC312_2019/tui"
)

func main() {

	conf := config.ReadConfig("config.yml")

	// Message channel to scheduler
	ch := make(chan *sched.Process, conf.ProcChanSize)
	defer close(ch)

	// Initialize resources
	cpu1 := cpu.InitCPU(conf.CPU.ClockSpeed1)
	// cpu2 := cpu.InitCPU(conf.CPU.ClockSpeed2)
	mem := memory.InitMemory(conf.Memory.PageSize, conf.Memory.TotalRam)

	// Initialize Scheduler
	s := sched.InitScheduler(cpu1, mem, ch, conf.MinimumFreeFrames, conf.Sched.TimeQuantum)

	// Run the scheduler
	go s.RunRoundRobin()
	// go k.RunFirstComeFirst

	// Initialize the TUI
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// Point the widgets to the scheduler
	tui.InitWidgets(s)

	// Render initial state to the terminal
	tui.Render()

	// Start the tui event loop
	tui.EventLoop(ch)

}
