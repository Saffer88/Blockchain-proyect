package logic

// signals.go

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Def_handler() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT)

	go func() {
		sig := <-sigchan
		fmt.Printf("\n\n[!] Saliendo del pograma: %v\n", sig)
		os.Exit(1)
	}()
}
