package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 0. temp child process
	child := exec.Command("sleep", "5")
	if err := child.Start(); err != nil {
		fmt.Printf("Failed to start child process: %v\n", err)
		return
	}
	fmt.Printf("Child process started (PID: %d)\n", child.Process.Pid)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan error, 1)
	go func() {
		done <- child.Wait()
	}()

	select {
	case <-sigCh: // close signal
		// 1. child shutdown
		fmt.Println("Sending SIGTERM to child...")
		child.Process.Signal(syscall.SIGTERM)

		// 2. child가 완전히 종료 될 때까지 2초 정도 시간 제공
		select {
		case err := <-done:
			// 정상 종료
			fmt.Printf("Child terminated safely. (Result: %v)\n", err)
		case <-time.After(2 * time.Second):
			// 응답이 없으면 강제 종료
			fmt.Println("Grace period exceeded. Escalating to SIGKILL!")
			child.Process.Kill()
		}

	case err := <-done:
		// 외부 시그널 없이 자식 프로세스가 먼저 스스로 종료된 경우
		fmt.Printf("Child exited unexpectedly or finished first. (Result: %v)\n", err)
	}
	fmt.Println("End")
}
