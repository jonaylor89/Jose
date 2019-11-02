package memory

const (
	pageLength = 256
)

var (
	// pages : pages in secondary memory
	virtualMemory [1024]Page

)

// RAM : virtual physical memory
type RAM struct {
	// frames : basically a cache of pages for the simulator because of the lack of hardware
	frames [256]Page
}

// Page : a page of memory
type Page struct {
	PID    int // Process ID of the process using this page
	length int
}

// Mutex : Mutex lock
type Mutex struct {
	locked bool
}

func (m *Mutex) acquire() bool {
	if m.locked {
		return false
	} 

	m.locked = true
	return true
}

func (m *Mutex) release() {
	m.locked = false
}
