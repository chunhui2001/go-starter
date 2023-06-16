package utils

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

var shutdownSignals chan os.Signal
var funcChannel chan func()
var isStarted bool = false
var funcArray []func()
var done chan struct{}

var (
	logger *logrus.Entry
)

func AddShutDownHook(log *logrus.Entry, f func()) {
	logger = log
	if !isStarted {
		shutdownSignals = make(chan os.Signal, 1)
		funcChannel = make(chan func())
		done = make(chan struct{})
		signal.Notify(shutdownSignals, syscall.SIGINT, syscall.SIGTERM)
		start()
		isStarted = true
	}
	funcChannel <- f
}

func WaitShutDown() {
	<-done
}

func executeHooks() {
	for _, f := range funcArray {
		f()
	}
}

func start() {
	go func() {
		shutdown := false
		for !shutdown {
			select {
			case <-shutdownSignals:
				for signal := range shutdownSignals {
					if signal == syscall.SIGTERM || signal == syscall.SIGQUIT {
						logger.Info("kill -15 退出进程")
					} else if signal == syscall.SIGILL {
						logger.Info("kill -4 退出进程")
					} else {
						logger.Info("kill -? 退出进程")
						break
					}
				}
				executeHooks()
				shutdown = true
			case f := <-funcChannel:
				funcArray = append(funcArray, f)
			}
		}
		done <- struct{}{}
	}()
}
