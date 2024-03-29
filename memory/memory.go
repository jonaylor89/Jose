package memory

import (
	"math"

	"github.com/hashicorp/golang-lru"
)

var (
	pageNum = 0
)

type Memory struct {
	// PageSize : length of page contents in Mb as a power of 2
	PageSize int

	// TotalRam : Total amount of physical memory in the simulator in Mb as a power of 2
	TotalRam int

	// PageTable : map page id to index of frame
	PageTable map[int]int

	// VirtualMemory : pages in secondary memory
	VirtualMemory []*Page

	// PhysicalMemory : Memory in RAM
	PhysicalMemory []*Page

	// Cache : Cache of pages
	Cache *lru.ARCCache
}

// Page : a page of memory
type Page struct {
	PageID   int    // ID of page
	ProcID   int    // Process ID of the process using this page
	contents []byte // Contents of the page of memory
}

// InitMemory : create new memory unit
func InitMemory(pageSize int, totalRam int, cacheSize int) *Memory {

	cache, _ := lru.NewARC(cacheSize)

	return &Memory{
		PageSize:       pageSize,
		TotalRam:       totalRam,
		PageTable:      make(map[int]int),
		VirtualMemory:  make([]*Page, 0),
		PhysicalMemory: make([]*Page, 0, totalRam/pageSize),
		Cache:          cache,
	}
}

// GetPage : get a page of memory
func (m *Memory) Get(pageNum int) *Page {

	// Check if the page is in the cache
	if result, ok := m.Cache.Get(pageNum); ok {
		return result.(*Page)
	}

	// Check for page in PhysicalMemory
	if val, ok := m.PageTable[pageNum]; ok {
		return m.PhysicalMemory[val]
	}

	// Otherwise, look through virtual memory for the page
	for i, page := range m.VirtualMemory {
		if page.PageID == pageNum {

			// Move the page to physical memory once found
			m.moveToPhysicalMemory(page, i)

			// Add page to the cache
			m.Cache.Add(pageNum, page)

			return page
		}
	}

	// Page doesn't exist
	return nil
}

// AddPage : Add pages of memory to memory pool, return PageIDs
func (m *Memory) Add(requirement int, pid int) []int {

	var pageIds []int

	numOfPages := int(math.Ceil(float64(requirement) / float64(m.PageSize)))

	for i := 0; i < numOfPages; i++ {

		pageNum++

		p := &Page{
			PageID:   pageNum,
			ProcID:   pid,
			contents: make([]byte, 0, 30),
		}

		pageIds = append(pageIds, pageNum)

		// Append new page to virtual memory
		m.VirtualMemory = append(m.VirtualMemory, p)
	}

	// return pageIds for the process to keep track of
	return pageIds
}

// moveToPhysicalMemory puts pages into RAM and adds the entry to the PageTable
func (m *Memory) moveToPhysicalMemory(p *Page, indexInVm int) {

	// Remove page from virtual memory
	m.VirtualMemory = remove(m.VirtualMemory, indexInVm)

	// if there is an empty space, put page in empty space
	if cap(m.PhysicalMemory)-len(m.PhysicalMemory) > 0 {
		m.PhysicalMemory = append(m.PhysicalMemory, p)

		// Always add new entry to page table and remove old entry if replaced
		m.PageTable[p.PageID] = len(m.PhysicalMemory) - 1
	}

	// if there isn't an empty space, run a replace procedure

	// Find victim page
	i, victimPage := m.findVictim(p.ProcID)
	if i == -1 {
		return
	}

	// Fill victim page's spot
	m.PhysicalMemory[i] = p

	// Always add new entry to page table and remove old entry if replaced
	m.PageTable[p.PageID] = i
	delete(m.PageTable, victimPage.PageID)

	// move victim page to virtual memory
	m.VirtualMemory = append(m.VirtualMemory, victimPage)
}

// findVictim : find a process to replace the current one with
func (m *Memory) findVictim(procID int) (int, *Page) {

	// Literally just find the first page with the same process ID lol
	for _, v := range m.PageTable {
		if m.PhysicalMemory[v].ProcID == procID {
			return v, m.PhysicalMemory[v]
		}
	}

	return -1, nil
}

// RemovePages : remove all pages associated with a pid
func (m *Memory) RemovePages(pid int) {

	// Remove pages from physical memory
	for i := len(m.PhysicalMemory) - 1; i >= 0; i-- {
		page := m.PhysicalMemory[i]
		if page.ProcID == pid {
			m.PhysicalMemory = remove(m.PhysicalMemory, i)

			// Update page table
			delete(m.PageTable, page.PageID)
		}
	}

	// Remove pages from virtual memory
	for i := len(m.VirtualMemory) - 1; i >= 0; i-- {
		page := m.VirtualMemory[i]
		if page.ProcID == pid {

			m.VirtualMemory = remove(m.VirtualMemory, i)
		}
	}
}

func remove(slice []*Page, s int) []*Page {
	slice[s] = slice[len(slice)-1] // Copy last element to index i.
	// slice[len(slice)-1] = nil   	// Erase last element (write zero value)
	slice = slice[:len(slice)-1] // Truncate slice.

	return slice
}
