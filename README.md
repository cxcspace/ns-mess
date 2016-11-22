# ns-mess

```bash
vagrant up
vagrant ssh
sudo su
ns-mess
```


success looks like
```
root@ubuntu-xenial:/# ns-mess
      main   start in namespace f0000075
     newns   start in namespace f0000fc1
     newns     end in namespace f0000fc1
      main waiting in namespace f0000075
      main     end in namespace f0000075
```


failure looks like
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
