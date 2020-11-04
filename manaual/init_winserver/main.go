package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"proxy-forward/pkg/utils"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func init() {
	if runtime.GOOS == "windows" {
		var st int
		file, _ := exec.LookPath(os.Args[0])
		cmd, err := exec.Command("tasklist").CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}
		procs := strings.Split(string(cmd), "\r\n")
		for _, item := range procs {
			item = strings.TrimSpace(item)
			if strings.Index(item, file) > -1 {
				st += 1
			}
		}
		fmt.Printf("[%s]=process:%d=\r\n", time.Now().Format("2006-01-02 15:04:05"), st)
		if st > 1 {
			os.Exit(0)
		}
	}
}

func main() {
	var err error
	err = netsh_port_to_localhost()
	log.Println(err)
}

// netsh 初始化端口 4001 - 4100 => 5001 -5100
func netsh_port_to_localhost() error {
	var (
		cmd []byte
		err error
	)
	if cmd, err = exec.Command("netsh", "interface", "portproxy", "show", "all").CombinedOutput(); err != nil {
		return err
	}
	if cmd, err = utils.GbkToUtf8(cmd); err != nil {
		return err
	}
	portProxy := strings.Split(string(cmd), "\r\n")
	for _, port := range utils.MakeRange(4001, 4100) {
		for _, item := range portProxy {
			if strings.Index(item, strconv.Itoa(port)) > -1 {
				goto endfor
			}
		}
		if cmd, err = exec.Command("netsh", "interface", "portproxy", "add", "v4tov4", fmt.Sprintf("listenport=%d", port), "connectaddress=127.0.0.1", fmt.Sprintf("connectport=%d", port+1000)).CombinedOutput(); err != nil {
			return err
		}
		if cmd, err = utils.GbkToUtf8(cmd); err != nil {
			return err
		}
	endfor:
	}
	return nil
}
