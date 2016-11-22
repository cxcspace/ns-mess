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
      main   start in namespace f0000075
     newns   start in namespace f00010c5
MESSY goroutine  7 start: expected f0000075, actual f00010c5
MESSY goroutine 12 start: expected f0000075, actual f00010c5
MESSY goroutine 33 start: expected f0000075, actual f00010c5
MESSY goroutine 41 start: expected f0000075, actual f00010c5
MESSY goroutine 40 start: expected f0000075, actual f00010c5
MESSY goroutine 45 start: expected f0000075, actual f00010c5
MESSY goroutine 49 start: expected f0000075, actual f00010c5
MESSY goroutine 48 start: expected f0000075, actual f00010c5
MESSY goroutine 44 start: expected f0000075, actual f00010c5
     newns     end in namespace f00010c5
      main waiting in namespace f0000075
MESSY goroutine  1   end: expected f0000075, actual f00010c5
MESSY goroutine  7   end: expected f00010c5, actual f0000075
MESSY goroutine 12   end: expected f00010c5, actual f0000075
MESSY goroutine 33   end: expected f00010c5, actual f0000075
MESSY goroutine 41   end: expected f00010c5, actual f0000075
MESSY goroutine  0   end: expected f0000075, actual f00010c5
      main     end in namespace f0000075
```
