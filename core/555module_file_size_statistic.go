package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

// 获取文件/目录的内容大小
func (this *InternalFunctionSet) Fsize(args []interface{}) int64 {
	size := len(args)
	if size < 1 {
		log.Println("fsize(path[, excludes]) parameter path is required")
		return 0
	}
	dir := args[0].(string)
	var excludes []string
	if size > 1 {
		excludesStr := args[1].(string)
		for _, exclude := range strings.Split(excludesStr, ",") {
			excludes = append(excludes, exclude)
		}
	}
	fss := newFileSizeStatistic(dir)
	fss.excludes = excludes
	return fss.do()
}

// 获取文件/目录的内容大小
func fileSize(dir string) int64 {
	return newFileSizeStatistic(dir).do()
}

type FileSizeStatistic struct {
	sizeChan chan int64
	size     int64
	dir      string
	excludes []string
	wt       *sync.WaitGroup
}

func newFileSizeStatistic(dir string) *FileSizeStatistic {
	return &FileSizeStatistic{
		dir:      dir,
		wt:       &sync.WaitGroup{},
		sizeChan: make(chan int64, 4),
	}
}

func (f *FileSizeStatistic) do() int64 {
	st, err := os.Stat(f.dir)
	if err != nil || st == nil {
		fmt.Println(f.dir, "is not found")
		return 0
	}

	f.wt.Add(1)
	go f.doStatistic(f.dir)

	go func() {
		f.wt.Wait()
		close(f.sizeChan)
	}()

	for fsize := range f.sizeChan {
		f.size += fsize
	}

	return f.size
}

func (f *FileSizeStatistic) doStatistic(dir string) {
	defer f.wt.Done()
	st, err := os.Stat(dir)
	if err != nil {
		fmt.Println(err)
		return
	}
	if st != nil && !st.IsDir() {
		f.sizeChan <- st.Size()
		return
	}

	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, info := range infos {
		if f.excludeCheck(info.Name()) {
			continue
		}

		if info.IsDir() {
			f.wt.Add(1)
			nextDir := pathJoin(dir, info.Name())
			go f.doStatistic(nextDir)
		} else {
			f.sizeChan <- info.Size()
		}
	}
}

func (f *FileSizeStatistic) excludeCheck(name string) bool {
	for _, exclude := range f.excludes {
		if exclude == name {
			return true
		}
	}
	return false
}
