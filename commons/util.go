package commons

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Schedule executes callback func with some period.
// no delays force callback to execute first time immediately
// Function returns channel to stop scheduled execution
func Schedule(timeout time.Duration, noDelay bool, callback func()) chan<- struct{} {
	ticker := time.NewTicker(timeout)
	quit := make(chan struct{})
	go func() {
		if noDelay {
			callback()
		}
		for {
			select {
			case <-ticker.C:
				callback()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	return quit
}

// Retry executes callback func until it executes successfully
func Retry(attempts int, timeout time.Duration, callback func() error) (err error) {
	for i := 0; i <= attempts-1; i++ {
		err = callback()
		if err == nil {
			return nil
		}

		//time.Sleep(timeout)
		<-time.After(timeout)
		log.Println("retrying...")
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}

// ShutdownHook adds function to be performed on app shutdown
func ShutdownHook(hook func() error) {
	c := make(chan os.Signal, 1) // Create a channel accepting os.Signal
	// Bind a given os.Signal to the channel we just created
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // Register os.Interrupt, syscall.SIGTERM

	go func() {
		// Start an anonymous func running in a goroutine
		<-c // that will block until a message is received on
		e := hook()
		if nil != e {
			log.Println("Shutdown hook error: ", e)
		}

		os.Exit(1)
	}()
}

// GetLocalIP returns first non-loopback IP address
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// KeySet returns array of map keys
func KeySet(m map[string]interface{}) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}
