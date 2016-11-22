package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/containernetworking/cni/pkg/ns"
	"golang.org/x/sys/unix"
)

func mainWithErr() error {
	originalMain := Snap("main start")
	log.Println(originalMain)

	newNS, err := ns.NewNS()
	if err != nil {
		return err
	}

	makeMeSomeGoRoutines := make(chan string)
	stopTheGoRoutines := make(chan struct{})
	wg := sync.WaitGroup{}

	// in the background, we spin up some goroutines
	go func() {
		for name := range makeMeSomeGoRoutines {
			wg.Add(1)
			myName := name
			go func() {
				originalGoRoutine := Snap(myName)
				reportIfNamespaceSwitch(myName, "start", originalMain, originalGoRoutine)
				<-stopTheGoRoutines
				finalGoRoutine := Snap(myName)
				reportIfNamespaceSwitch(myName, "end", originalGoRoutine, finalGoRoutine)
				wg.Done()
			}()
		}
	}()

	// spawn a goroutine that sits in that namespace
	err = newNS.Do(func(prevNS ns.NetNS) error {
		originalNewNS := Snap("newns start")
		log.Println(originalNewNS)
		for i := 0; i < 50; i++ {
			makeMeSomeGoRoutines <- fmt.Sprintf("goroutine %2d", i)
		}
		close(stopTheGoRoutines)
		finalNewNS := Snap("newns stop")
		reportIfNamespaceSwitch("newns", "end", originalNewNS, finalNewNS)
		return nil
	})

	log.Println(Snap("main waiting"))
	wg.Wait()
	log.Println(Snap("main stop "))
	return err
}

func Snap(name string) Snapshot {
	inode, err := getInode(getCurrentThreadNetNSPath())
	if err != nil {
		panic(err)
	}
	return Snapshot{
		Name: name,
		NS:   fmt.Sprintf("%x", inode),
	}
}

func reportIfNamespaceSwitch(name, state string, original, final Snapshot) {
	if final.NS != original.NS {
		log.Printf("MESSY %s %s: expected to be in NS %s but am instead in NS %s", name, state, original.NS, final.NS)
	}
}

func (s Snapshot) String() string {
	return fmt.Sprintf("%25s in namespace %10s", s.Name, s.NS)
}

type Snapshot struct {
	Name string
	NS   string
}

func getCurrentThreadNetNSPath() string {
	// /proc/self/ns/net returns the namespace of the main thread, not
	// of whatever thread this goroutine is running on.  Make sure we
	// use the thread's net namespace since the thread is switching around
	return fmt.Sprintf("/proc/%d/task/%d/ns/net", os.Getpid(), unix.Gettid())
}

func getInode(path string) (uint64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	return getInodeFd(int(file.Fd()))
}

func getInodeFd(fd int) (uint64, error) {
	stat := &unix.Stat_t{}
	err := unix.Fstat(fd, stat)
	return stat.Ino, err
}

func main() {
	log.SetFlags(0)
	if err := mainWithErr(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
