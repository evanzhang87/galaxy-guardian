package quota

import (
	"fmt"
	"galaxy-guardian/logger"
	"github.com/shirou/gopsutil/process"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func WatchRss(rssLimitMb int) {
	pid := int32(os.Getegid())
	pcs, err := process.NewProcess(pid)
	count := 0
	if err != nil {
		logger.Logger.Errorf("get process err :%v", err)
		return
	}
	logger.Logger.Info("start to watch rss")
	for {
		mem, err := pcs.MemoryInfo()
		if err != nil {
			logger.Logger.Errorf("get memory info err : %v", err)
		} else {
			rssMb := float64(mem.RSS) / 1024 / 1024
			if rssMb > float64(rssLimitMb) {
				count += 1
				logger.Logger.Warnf("CurrentRSS above limit: %v >= %v Mb", rssMb, rssLimitMb)
			} else {
				count = 0
			}
			if count >= 3 {
				logger.Logger.Errorf("Suicide, CurrentRSS above limit: %v >= %v Mb", rssMb, rssLimitMb)
				os.Exit(1)
			}
		}
		time.Sleep(time.Minute)
	}
}

func WatchSocket(socketLimit int) {
	pid := strconv.Itoa(os.Getegid())
	fdPath := fmt.Sprintf("/proc/%s/fd", pid)
	logger.Logger.Info("start to watch socket")
	for {
		var softs = getFileSoftLink(fdPath)
		var socketNum = 0
		for _, soft := range softs {
			if strings.Contains(soft, "socket") {
				socketNum++
			}
		}
		if socketNum >= socketLimit {
			logger.Logger.Errorf("Suicide, Current socket above limit: %d >= %d", socketNum, socketLimit)
			os.Exit(1)
		}
		time.Sleep(time.Minute)
	}
}

func getFileSoftLink(filePath string) []string {
	var result = make([]string, 0)
	if filePath == "" {
		return result
	}
	var fdResults, err = ioutil.ReadDir(filePath)
	if err != nil {
		return result
	}
	for _, fdResult := range fdResults {
		fi, err := os.Lstat(filePath + `/` + fdResult.Name())
		if err != nil {
			continue
		}
		// True if the file is a symlink.
		if fi.Mode()&os.ModeSymlink != 0 {
			link, err := os.Readlink(filePath + `/` + fi.Name())
			if err != nil {
				continue
			}
			result = append(result, link)
		}
	}
	return result
}
