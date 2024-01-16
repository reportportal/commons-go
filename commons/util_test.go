package commons

import (
	"errors"
	"log"
	"net"
	"reflect"
	"regexp"
	"sort"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestKeySet(t *testing.T) {
	mp := map[string]interface{}{
		"one": struct{}{},
		"two": struct{}{},
	}

	actual := KeySet(mp)

	// make sure they are sorted before validation
	sort.Strings(actual)

	expected := []string{"one", "two"}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Arrays are not equal. Expected:%s, Actual:%s", expected, actual)
	}
}

func TestSchedule(t *testing.T) {
	i := 0
	quite := Schedule(1*time.Second, true, func() {
		i++
	})
	time.Sleep(2 * time.Second)
	quite <- struct{}{}

	if i == 0 {
		t.Errorf("Incorrect execution count: %d", i)
	}
}

func TestRetryAttempts(t *testing.T) {
	RegisterTestingT(t)

	i := 0
	err := Retry(2, 1*time.Second, func() error {
		i++
		return errors.New("some error")
	})
	log.Println(err)
	log.Println(i)
	if 2 != i {
		t.Errorf("Incorrect attempts count: %d", i)
	}
}

func TestGetLocalIP(t *testing.T) {
	ip := net.ParseIP(GetLocalIP())
	if ip.IsLoopback() {
		t.Errorf("IP is loopback: %s", ip.String())
	}

	if !regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`).MatchString(ip.String()) {
		t.Errorf("Incorrect IP format: %s", ip.String())
	}

	print(ip.String())
	print()
}
