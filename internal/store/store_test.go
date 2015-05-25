package store

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func PrintSync(ch chan SM) {
	for {
		select {
		case c, ok := <-ch:
			if ok {
				fmt.Printf("k: %s -> %s\n", c.uri, string(c.body))
			} else {
				return
			}
		}
	}
}
func SumSync(ch chan SM) int {
	rv := 0
	for {
		select {
		case c, ok := <-ch:
			if ok {
				nrv, _ := strconv.ParseInt(string(c.body), 16, 64)
				rv += int(nrv)
			} else {
				return rv
			}
		}
	}
}
func DumpArr(arr [][]byte) {
	for i, v := range arr {
		fmt.Printf("%d : %v\n", i, string(v))
	}
}
func TestOsterone(t *testing.T) {
	PutMessage("a/b/c", []byte("v(a/b/c)"))
	PutMessage("a/b/c/1", []byte("v(a/b/c/1)"))
	PutMessage("a/b/d", []byte("v(a/b/d)"))
	rc := make(chan SM, 3)
	go GetMatchingMessage("a/*/1", rc)
	PrintSync(rc)
}

func TestIcle(t *testing.T) {
	datasetvector := []struct {
		URI  string
		Data string
	}{
		{"tstes/a/b/c", "1"},
		{"tstes/a/b/d", "2"},
		{"tstes/a/b/c/1", "4"},
		{"tstes/x/b/c/1", "8"},
		{"tstes/foo/c/1", "10"},
		{"tstes/foo/c/2", "20"},
	}
	testvector := []struct {
		QRY      string
		Expected int
	}{
		{"tstes/a/b/c", 1},
		{"tstes/a/b/+", 1 + 2},
		{"tstes/a/b/*", 1 + 2 + 4},
		{"tstes/+/c/+", 0x10 + 0x20},
		{"*/1", 0x4 + 0x8 + 0x10},
		{"+/*", 0x3F},
		{"*/+", 0x3F},
	}
	for _, v := range datasetvector {
		PutMessage(v.URI, []byte(v.Data))
	}
	fmt.Println("============= PUT COMPLETE ================")
	for i, v := range testvector {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("===== TESTING [", i, "] ", v.QRY, " ================")
		rc := make(chan SM, 3)
		go GetMatchingMessage(v.QRY, rc)
		got := SumSync(rc)
		time.Sleep(100 * time.Millisecond)
		if got != v.Expected {
			fmt.Printf("For test vector %d expected %d, got %d\n", i, v.Expected, got)
			t.FailNow()
		} else {
			fmt.Printf("Test vector ok\n")
		}
	}
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Done")

}
