package storage

import (
	"bufio"
	"bytes"
	"os"
	"strings"

	"github.com/bluele/gcache"
)

var Storage gcache.Cache

func init() {
	gc := gcache.New(20).
		LRU().
		Build()
	Storage = gc
	open, err := os.Open("/etc/hosts")
	if err != nil {
		return
	}
	defer open.Close()
	reader := bufio.NewReader(open)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			return
		}
		lines := string(bytes.TrimSpace(line))
		if strings.Index(lines, "#") != -1 {
			continue
		}
		split := strings.Split(lines, " ")
		var c []string
		for _, v := range split {
			if v != "" {
				c = append(c, v)
			}
		}

		if len(c) >= 2 {
			gc.Set(c[0], c[1])
		}
	}
}
