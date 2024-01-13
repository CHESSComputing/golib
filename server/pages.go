package server

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

//
// helper functions
//

// ErrorPage returns error page
func ErrorPage(fsys fs.FS, msg string, err error) string {
	log.Printf("ERROR: %v\n", err)
	tmpl := MakeTmpl(fsys, "Error")
	tmpl["Content"] = strings.ToTitle(msg)
	return TmplPage(fsys, "error.tmpl", tmpl)
}

// HeaderPage returns header page
func HeaderPage(fsys fs.FS) string {
	tmpl := MakeTmpl(fsys, "Header")
	return TmplPage(fsys, "header.tmpl", tmpl)
}

// FooterPage returns footer page
func FooterPage(fsys fs.FS) string {
	tmpl := MakeTmpl(fsys, "Footer")
	return TmplPage(fsys, "footer.tmpl", tmpl)
}

// helper function to make initial template struct
func MakeTmpl(fsys fs.FS, title string) TmplRecord {
	tmpl := make(TmplRecord)
	tmpl["Title"] = title
	tmpl["StartTime"] = time.Now().Unix()
	return tmpl
}

// ErrorTmpl provides error template message
func ErrorTmpl(fsys fs.FS, msg string, err error) string {
	tmpl := MakeTmpl(fsys, "Status")
	tmpl["Content"] = template.HTML(fmt.Sprintf("<div>%s</div>\n<br/><h3>ERROR</h3>%v", msg, err))
	content := TmplPage(fsys, "error.tmpl", tmpl)
	return content
}

// SuccessTmpl provides success template message
func SuccessTmpl(fsys fs.FS, msg string) string {
	tmpl := MakeTmpl(fsys, "Status")
	tmpl["Content"] = template.HTML(fmt.Sprintf("<h3>SUCCESS</h3><div>%s</div>", msg))
	content := TmplPage(fsys, "success.tmpl", tmpl)
	return content
}

// FAQPage provides FAQ page
func FAQPage(fsys fs.FS) string {
	tmpl := MakeTmpl(fsys, "FAQ")
	return TmplPage(fsys, "faq.tmpl", tmpl)
}

/*
// Memory structure keeps track of server memory
type Memory struct {
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

// Mem structure keeps track of virtual/swap memory of the server
type Mem struct {
	Virtual Memory
	Swap    Memory
}
*/

// MetricsPage represents metrics page
func MetricsPage(fsys fs.FS) TmplRecord {
	// get cpu and mem profiles
	m, _ := mem.VirtualMemory()
	s, _ := mem.SwapMemory()
	l, _ := load.Avg()
	c, _ := cpu.Percent(time.Millisecond, true)
	process, perr := process.NewProcess(int32(os.Getpid()))

	// get unfinished queries
	tmpl := MakeTmpl(fsys, "Metrics")
	tmpl["NGo"] = runtime.NumGoroutine()
	//     virt := Memory{Total: m.Total, Free: m.Free, Used: m.Used, UsedPercent: m.UsedPercent}
	//     swap := Memory{Total: s.Total, Free: s.Free, Used: s.Used, UsedPercent: s.UsedPercent}
	tmpl["Memory"] = m.UsedPercent
	tmpl["Swap"] = s.UsedPercent
	tmpl["Load1"] = l.Load1
	tmpl["Load5"] = l.Load5
	tmpl["Load15"] = l.Load15
	tmpl["CPU"] = c
	if perr == nil { // if we got process info
		conn, err := process.Connections()
		if err == nil {
			tmpl["Connections"] = conn
		}
		openFiles, err := process.OpenFiles()
		if err == nil {
			tmpl["OpenFiles"] = openFiles
		}
	}
	tmpl["Uptime"] = time.Since(Time0).Seconds()
	tmpl["GetRequests"] = TotalGetRequests
	tmpl["PostRequests"] = TotalPostRequests
	return tmpl
}
