package commitinfo

import (
	"fmt"
	"path/filepath"
	"time"
)

type CommitInfo struct {
	Author    string
	Filename  string
	Timestamp time.Time
}

func (c *CommitInfo) print() {
	fmt.Printf("CommitInfo[Author: %s, Filename: %s, Timestamp: %s]\n",
		c.Author, filepath.Base(c.Filename),
		c.Timestamp.Format(time.RFC3339))
}

func (c *CommitInfo) logPrint() {
	fmt.Printf("%s %s %s\n", c.Author, filepath.Base(c.Filename),
		c.Timestamp.Format(time.RFC3339))
}
