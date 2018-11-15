package snippets

import (
	"fmt"
	"time"
)

func everySec() {
	c := time.Tick(1 * time.Second)
	for now := range c {
		fmt.Printf("%v\n", now)
	}
}

func afterOneSec() {
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("timed out")
	}
}

func parseJsDt(t string) (time.Time, error) {
	pt, err := time.Parse("Mon Jan 02 2006 15:04:05 GMT-0700 (MST)", t)
	if err != nil {
		return time.Time{}, err
	}
	return pt.UTC(), nil
}

func parseMysqlDt(t string) (time.Time, error) {
	pt, err := time.Parse("2006-01-02 15:04:05", t)
	if err != nil {
		return time.Time{}, err
	}
	return pt.UTC(), nil
}
