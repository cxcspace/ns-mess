package main

import (
	"fmt"
	"os"

	"github.com/containernetworking/cni/pkg/ns"
	"golang.org/x/sys/unix"
)

func main() {
	fmt.Println(Snap())
}

func Snap() Snapshot {
	return Snapshot{
		Thread: WhatThread(),
		NS:     WhatNS(),
	}
}

func (s Snapshot) String() string {
	return fmt.Sprintf("%s: %s", s.Thread, s.NS)
}

type Snapshot struct {
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
