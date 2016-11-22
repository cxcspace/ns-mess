package main

import (
	"fmt"
	"log"
	"os"

	"github.com/containernetworking/cni/pkg/ns"
	"golang.org/x/sys/unix"
)

func mainWithErr() error {
	log.Println(Snap("main start"))

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

	// spawn a goroutine that sits in that namespace
	err = newNS.Do(func(prevNS ns.NetNS) error {
		log.Println(Snap("newns goroutine"))
		return nil
	})

	log.Println(Snap("main end"))
	return err
}

func Snap(name string) Snapshot {
	return Snapshot{
		Name:   name,
		Thread: WhatThread(),
		NS:     WhatNS(),
	}
}

func (s Snapshot) String() string {
	return fmt.Sprintf("%20s: thread %6s in namespace %10s", s.Name, s.Thread, s.NS)
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
	if err := mainWithErr(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
