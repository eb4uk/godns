package main

import (
	"github.com/eb4uk/godns/settings"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"time"
)

var (
	logger *GoDNSLogger
)

func main() {
	settings.InitializeConfig()
	initLogger()

	server := &Server{
		host:     settings.Config.Server.Host,
		port:     settings.Config.Server.Port,
		rTimeout: 5 * time.Second,
		wTimeout: 5 * time.Second,
	}

	server.Run()

	logger.Info("godns %s start", settings.Config.Version)

	if settings.Config.Debug {
		go profileCPU()
		go profileMEM()
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

forever:
	for {
		select {
		case <-sig:
			logger.Info("signal received, stopping")
			break forever
		}
	}

}

func profileCPU() {
	f, err := os.Create("godns.cprof")
	if err != nil {
		logger.Error("%s", err)
		return
	}

	pprof.StartCPUProfile(f)
	time.AfterFunc(6*time.Minute, func() {
		pprof.StopCPUProfile()
		f.Close()

	})
}

func profileMEM() {
	f, err := os.Create("godns.mprof")
	if err != nil {
		logger.Error("%s", err)
		return
	}

	time.AfterFunc(5*time.Minute, func() {
		pprof.WriteHeapProfile(f)
		f.Close()
	})

}

func initLogger() {
	logger = NewLogger()

	if settings.Config.Log.Stdout {
		logger.SetLogger("console", nil)
	}

	if settings.Config.Log.File != "" {
		config := map[string]interface{}{"file": settings.Config.Log.File}
		logger.SetLogger("file", config)
	}

	logger.SetLevel(settings.Config.Log.LogLevel())
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
