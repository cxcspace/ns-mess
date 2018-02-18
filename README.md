# FIXED in Go 1.10

Go 1.10 runtime now does the right thing with Goroutines and Threads,
[see here](https://golang.org/doc/go1.10#runtime).

I'm keeping this repo around as an artifact.

You can play with different Go
versions by adjusting the environment variable in the [`Vagrantfile`](Vagrantfile).

# ns-mess
Investigations into Go programs that switch their Linux [namespace](http://man7.org/linux/man-pages/man7/namespaces.7.html).

## A minimal reproduction
```bash
vagrant up
vagrant ssh
sudo su
ns-mess
```

Sometimes it succeeds:
```
root@ubuntu-xenial:/# ns-mess
      main   start in namespace f0000075
     newns   start in namespace f0000fc1
     newns     end in namespace f0000fc1
      main waiting in namespace f0000075
      main     end in namespace f0000075
```

But often it fails:
```
root@ubuntu-xenial:/# ns-mess
      main   start in namespace f0000075
     newns   start in namespace f000187d
MESSY goroutine  1 start: expected f0000075, actual f000187d
MESSY goroutine  9 start: expected f0000075, actual f000187d
MESSY goroutine  8 start: expected f0000075, actual f000187d
MESSY goroutine 11 start: expected f0000075, actual f000187d
MESSY goroutine 10 start: expected f0000075, actual f000187d
MESSY goroutine 12 start: expected f0000075, actual f000187d
MESSY goroutine 15 start: expected f0000075, actual f000187d
MESSY goroutine 14 start: expected f0000075, actual f000187d
MESSY goroutine 17 start: expected f0000075, actual f000187d
MESSY goroutine 19 start: expected f0000075, actual f000187d
MESSY goroutine 18 start: expected f0000075, actual f000187d
     newns     end in namespace f000187d
MESSY goroutine 25   end: expected f0000075, actual f000187d
      main waiting in namespace f0000075
      main     end in namespace f000187d
error: at least one goroutine saw the wrong namespace
```

## Further reading
- https://github.com/containernetworking/cni/issues/262
- https://github.com/vishvananda/netns/issues/17
- https://github.com/docker/libnetwork/issues/1113
- https://github.com/weaveworks/weave/issues/2388#issuecomment-228365069
- https://gist.github.com/sykesm/020a27341cca5169250990d250b25db4
