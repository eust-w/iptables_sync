package ctrl

import (
	"fmt"
	"github.com/eust-w/iptables_sync/rfsnotify"
	"github.com/fsnotify/fsnotify"
	"os"
	"path"
	"strings"
)

var zstackDir = getZstackDir()

//var uiProxyPath = path.Join(zstackDir, "zstack-ui/configs/proxy")
//var uiProxyIptablesPath = "/etc/iptables.rules"
//var uiProxyCmds = []string{
//	"sudo /usr/sbin/nginx -c ~zstack/zstack-ui/configs/nginx.conf -s reload",
////}

var uiProxyPath = "/tmp/lt"
var uiProxyCmds = []string{
	"sudo echo 1 >> lt.record",
}

var uiProxyIptablesPath = path.Join(zstackDir, "zstack-ui/configs/iptables/iptables.rules")
var uiProxyIptablesCmds = []string{
	"sudo echo 2 >> lt.record",
}

func getZstackDir() string {
	bash := Bash{
		Command: "echo ~zstack",
	}
	_, o, e, err := bash.RunWithReturn()
	if err != nil || strings.TrimSpace(o) == "~zstack" {
		fmt.Printf("error to get zstack home dir, stdout: %s, stderr: %s, error: %s", o, e, err)
	}
	return strings.TrimSpace(o)
}

func UiProxyMonitor(peer string) {
	uiProxyCmds = append(uiProxyCmds, fmt.Sprintf("sudo iptables-save > %s", uiProxyIptablesPath))
	//uiProxyCmds = append(uiProxyCmds, fmt.Sprintf("sudo scp %s root@%s:%s", uiProxyIptablesPath, peer, uiProxyIptablesPath))
	uiProxyIptablesCmds = append(uiProxyIptablesCmds, fmt.Sprintf("sudo iptables-restore < %s", uiProxyIptablesPath))
	fmt.Println("uiProxyIptablesCmds:", uiProxyIptablesCmds)
	err := createDirIfNotExist(path.Dir(uiProxyIptablesPath))
	if err != nil {
		fmt.Printf("create dir %s failed: %s", path.Dir(uiProxyIptablesPath), err)
	}
	syncDirAndRunCmds(uiProxyPath, peer, uiProxyCmds, true)
	go dirSyncAndRunCmdsMonitor(uiProxyPath, peer, uiProxyCmds)
	go dirSyncAndRunCmdsMonitor(path.Dir(uiProxyIptablesPath), peer, uiProxyIptablesCmds)
}

func createDirIfNotExist(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0755)
		return err
	}
	return nil
}

func syncDirAndRunCmds(dir, peer string, cmds []string, skipNewer bool) {
	doRsyncToRemote(dir, peer, skipNewer)
	if len(cmds) > 0 {
		fmt.Println("run cmd:cmds", cmds)
		runCmds(cmds)
	}
}

func dirSyncAndRunCmdsMonitor(dir, peer string, cmds []string) {
	watcher, err := rfsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("fsnotify: %s", err)
		return
	}
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		fmt.Printf("syncing dir from %s to %s, and run cmds %v", dir, peer, cmds)
		for true {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					done <- true
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write ||
					event.Op&fsnotify.Remove == fsnotify.Remove ||
					event.Op&fsnotify.Rename == fsnotify.Rename {
					fmt.Printf("syncing dir from %s to %s", dir, peer)
					go syncDirAndRunCmds(dir, peer, cmds, false)
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					fmt.Printf("new create event, watch[%s] and syncing to %s", event.Name, peer)
					err = watcher.AddRecursive(dir)
					if err != nil {
						fmt.Printf("fsnotify: %s", err)
					}
					fmt.Printf("syncing dir from %s to %s", dir, peer)
					go syncDirAndRunCmds(dir, peer, cmds, false)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					done <- true
					return
				}
				fmt.Printf("fsnotify: %s", err)
			}
		}
	}()
	err = watcher.AddRecursive(dir)
	if err != nil {
		fmt.Printf("fsnotify: %s", err)
	}
	<-done
}

func doRsyncToRemote(fpath, peer string, skipNewer bool) {
	opt := "-rtogz"
	if skipNewer {
		opt += "u"
	}

	bash := Bash{
		Command: fmt.Sprintf(
			`/usr/bin/rsync --delete %s -e "ssh -o BatchMode=yes -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null" %s %s`,
			opt,
			fpath,
			fmt.Sprintf("root@%s:%s", peer, path.Dir(fpath))),
	}

	_, stdout, stderr, err := bash.RunWithReturn()
	if err != nil {
		fmt.Printf("rsync %s failed: %s", fpath, stdout+stderr)
	}
}

func runCmds(cmds []string) {
	for _, cmd := range cmds {
		bash := Bash{
			Command: fmt.Sprintf(cmd),
		}
		_, stdout, stderr, err := bash.RunWithReturn()
		if err != nil {
			fmt.Printf("\nrun cmd: %s failed: %s\n", cmd, stdout+stderr)
		} else {
			fmt.Printf("\nrun cmd: %s success: %s\n", cmd, stdout+stderr)
		}
	}
	return
}
