package main

import (
	"bufio"
	"dd-logstats/engine"
	"dd-logstats/ui"
	"flag"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/golang/glog"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
)

const (
	logEntryBufferSize = 1024
	hitTreshold        = 50
)

var flagAlarmTimeFrame time.Duration
var flagAlarmThreshold uint64
var flagStatsPeriod time.Duration
var flagServeAddr string

func init() {
	flag.DurationVar(&flagAlarmTimeFrame, "alarm-timeframe", time.Duration(2*time.Minute), "consider the total hits for this period of time when raising an alarm")
	flag.Uint64Var(&flagAlarmThreshold, "alarm-threshold", 100, "raise an alarm when the average hit count between [now - alarm_timeframe; now] exceeds the given value")
	flag.DurationVar(&flagStatsPeriod, "stats-period", time.Duration(10*time.Second), "statistics aggregation period")
	flag.StringVar(&flagServeAddr, "serve", ":8080", "address to serve requests")
}

func readLogs(out chan *engine.LogEntry) {
	glog.Info("Reading log lines from standard input")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		obj, err := engine.NewLogEntry(line)
		if err != nil {
			glog.Errorf("Error decoding line from stdin: %s", err)
			continue
		}
		out <- obj
	}
	err := scanner.Err()
	if err != nil {
		glog.Errorf("Error reading line from stdin (%s)", err)
	}
}

func runTrackers(logs chan *engine.LogEntry, context ui.Renderer, quit, done chan bool) {

	stats := engine.NewStats()
	uiState := ui.State{}
	ht := engine.NewHitTracker(uint64(flagAlarmTimeFrame/flagStatsPeriod), flagAlarmThreshold)
	ticker := time.NewTicker(flagStatsPeriod)

	for {
		select {

		case <-ticker.C:
			ht.AddHits(stats.TotalHits)
			stats.Finalize()
			// Update state
			uiState.Update(stats, ht.IsAboveThreshold(), ht.AverageHitCount())
			// Render template
			err := context.Render(&uiState)
			if err != nil {
				glog.Errorf("Error rendering template: %s", err)
			}
			stats = engine.NewStats()

		case l := <-logs:
			stats.AddLogEntry(l)

		case <-quit:
			done <- true
		}
	}
}

type appHandler struct {
	ui.Renderer
	Handle func(ui.Renderer, http.ResponseWriter, *http.Request) (int, error)
}

func logger(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		glog.Infoln(r.RemoteAddr, r.Method, r.RequestURI)
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (ah appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := ah.Handle(ah.Renderer, w, r)
	if err != nil {
		glog.Errorf("HTTP %d: %q", status, err)
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(status), status)
		default:
			http.Error(w, http.StatusText(status), status)
		}
	}
}

func indexHandler(renderer ui.Renderer, w http.ResponseWriter, r *http.Request) (int, error) {
	return w.Write([]byte(renderer.Result()))
}

func main() {

	// force logging to stderr (glog)
	flag.Set("logtostderr", "true")
	flag.Parse()

	glog.Infoln("Hits statistics are refreshed every", flagStatsPeriod)
	glog.Infoln("Alarm is raised when average hit count over the last", flagAlarmTimeFrame, "exceeds", flagAlarmThreshold, "hits")

	logs := make(chan *engine.LogEntry, logEntryBufferSize)
	quit, done := make(chan bool), make(chan bool)

	context, err := ui.NewRenderer("ui/assets", flagStatsPeriod, flagAlarmThreshold, flagAlarmTimeFrame)
	if err != nil {
		glog.Errorf("Failed to setup template renderer (%s)", err)
		os.Exit(1)
	}

	go runTrackers(logs, context, quit, done)
	go readLogs(logs)

	mux := web.New()
	mux.Use(logger)
	mux.Get("/", appHandler{context, indexHandler})
	mux.Get("/assets/*", http.FileServer(http.Dir("ui")))

	glog.Infoln("User interface served on", flagServeAddr, "(HTTP)")

	graceful.AddSignal(syscall.SIGTERM, syscall.SIGINT)
	err = graceful.ListenAndServe(flagServeAddr, mux)
	if err != nil {
		glog.Errorf("Unable to start webserver (%s)", err)
	}
	glog.Infoln("Shutting down. Waiting for goroutines to end.")
	close(quit)
	<-done
}
