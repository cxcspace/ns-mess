package main

import (
	"fmt"
	"log"
	"os"

	"github.com/containernetworking/cni/pkg/ns"
	"golang.org/x/sys/unix"
)

func mainWithErr() error {
	originalMain := Snap("main start")
	log.Println(originalMain)

	// create a new namespace
	newNS, err := ns.NewNS()
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := newNS.Close(); closeErr != nil {
			log.Fatalf("close ns error: %s", closeErr)
			panic(closeErr)
		}
	}()

	makeMeSomeGoRoutines := make(chan string)
	stopTheGoRoutines := make(chan struct{})

	// in the background, we spin up some goroutines
	go func() {
		for name := range makeMeSomeGoRoutines {
			myName := name
			go func() {
				originalGoRoutine := Snap(myName)
				reportIfNamespaceSwitch(myName, "start", originalMain, originalGoRoutine)
				log.Println(originalGoRoutine)
				<-stopTheGoRoutines
				finalGoRoutine := Snap(myName)
				reportIfNamespaceSwitch(myName, "end", originalGoRoutine, finalGoRoutine)
			}()
		}
	}()

	// spawn a goroutine that sits in that namespace
	err = newNS.Do(func(prevNS ns.NetNS) error {
		originalNewNS := Snap("newns start")
		log.Println(originalNewNS)
		for i := 0; i < 25; i++ {
			makeMeSomeGoRoutines <- fmt.Sprintf("goroutine %2d", i)
		}
		close(stopTheGoRoutines)
		finalNewNS := Snap("newns stop")
		reportIfNamespaceSwitch("newns", "end", originalNewNS, finalNewNS)
		return nil
	})

	log.Println(Snap("main stop "))
	return err
}

func Snap(name string) Snapshot {
	return Snapshot{
		Name:   name,
		Thread: WhatThread(),
		NS:     WhatNS(),
	}
}

func reportIfNamespaceSwitch(name, state string, original, final Snapshot) {
	if final.NS != original.NS {
		log.Printf("%25s: %5s: expected to be in NS %s but am instead in NS %s", name, state, original.NS, final.NS)
	}
}

func (s Snapshot) String() string {
	return fmt.Sprintf("%25s: thread %6s in namespace %10s", s.Name, s.Thread, s.NS)
}

func (s Snapshot) Numbers() string {
	return fmt.Sprintf("thread %s in namespace %s", s.Thread, s.NS)
}

type Snapshot struct {
	Name   string
	Thread string
	NS     string
}

func WhatNS() string {
	inode, err := getInodeCurNetNS()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", inode)
}

func WhatThread() string {
	return fmt.Sprintf("%d", unix.Gettid())
}

func getInodeCurNetNS() (uint64, error) {
	curNS, err := ns.GetCurrentNS()
	if err != nil {
		return 0, err
	}
	defer curNS.Close()
	return getInodeNS(curNS)
}

func getInodeNS(netns ns.NetNS) (uint64, error) {
	return getInodeFd(int(netns.Fd()))
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
