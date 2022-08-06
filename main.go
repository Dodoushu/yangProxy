package main

import (
	"code.byted.org/demo/goPractice/services"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

//func main1() {
//	err := initConfig()
//	if err != nil {
//		log.Fatalf("err : %s", err)
//	}
//	Clean(&service.S)
//}

func formatGroupTopic(format string, id int) string {
	if strings.ContainsRune(format, '%') {
		s := fmt.Sprintf(format, id)
		if strings.Contains(s, "%!") {
			msg := fmt.Sprintf("invalid format %q, after Sprintf we got %q", format, s)
			panic(msg)
		}
		return s
	}
	return format
}

func Clean(s *services.Service) {
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		for _ = range signalChan {
			fmt.Println("\nReceived an interrupt, stopping services...")
			(*s).Clean()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
