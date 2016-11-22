package main

import (
	"fmt"
	"log"
	"os"

	"github.com/containernetworking/cni/pkg/ns"
	"golang.org/x/sys/unix"
)

const nGoRoutines = 50

// Snapshot is a string representation of the namespace inode
type Snapshot string

func SnapshotNS() Snapshot {
	inode, err := getInode(getCurrentThreadNetNSPath())
	if err != nil {
		panic(err)
	}
	return Snapshot(fmt.Sprintf("%x", inode))
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

func checkUnexpectedNS(msg string, expected, actual Snapshot) bool {
	if actual != expected {
		log.Printf("MESSY %15s: expected %s, actual %s", msg, expected, actual)
		return false
	}
	return true
}

func reportNamespace(routine, state string, snapshot Snapshot) {
	log.Printf("%10s %7s in namespace %s", routine, state, snapshot)
}

func mainWithErr() error {
	hostNS := SnapshotNS()
	reportNamespace("main", "start", hostNS)

	workQueue := make(chan string)
	done := make(chan struct{})
	results := make(chan bool, 2*nGoRoutines)

	// in the background, we will spin up go routines, on demand
	// the call stack that spawned them never changed namespaces
	// so we would hope that the code inside the goroutine
	// would only ever run on the host namespace
	go func() {
		for item := range workQueue {
			myName := item
			go func() {
				originalNS := SnapshotNS()
				results <- checkUnexpectedNS(myName+" start", hostNS, originalNS)

				<-done

				finalNS := SnapshotNS()
				results <- checkUnexpectedNS(myName+"   end", originalNS, finalNS)
			}()
		}
	}()

	// separately, we create a new namespace
	newNS, err := ns.NewNS()
	if err != nil {
		return err
	}

	// and do some work inside this new namespace
	err = newNS.Do(func(prevNS ns.NetNS) error {
		newNS := SnapshotNS()
		reportNamespace("newns", "start", newNS)

		// queue up some work
		for i := 0; i < nGoRoutines; i++ {
			workQueue <- fmt.Sprintf("goroutine %2d", i)
		}
		close(done)

		reportNamespace("newns", "end", SnapshotNS())
		return nil
	})
	if err != nil {
		return err
	}

	reportNamespace("main", "waiting", SnapshotNS())

	allGood := true
	for i := 0; i < nGoRoutines; i++ {
		allGood = allGood && <-results
	}
	reportNamespace("main", "end", SnapshotNS())

	if allGood {
		return nil
	} else {
		return fmt.Errorf("at least one goroutine saw the wrong namespace")
	}
}

func main() {
	log.SetFlags(0)
	if err := mainWithErr(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
