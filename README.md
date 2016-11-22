# ns-mess

```bash
vagrant up
vagrant ssh
sudo su
ns-mess
```

typical output:
```
root@ubuntu-xenial:/# ns-mess
               main start in namespace   f0000075
              newns start in namespace   f00005cd
             main waiting in namespace   f0000075
               main stop  in namespace   f0000075

root@ubuntu-xenial:/# ns-mess
               main start in namespace   f0000075
              newns start in namespace   f0000601
MESSY goroutine  6 start: expected to be in NS f0000075 but am instead in NS f0000601
MESSY goroutine  5 start: expected to be in NS f0000075 but am instead in NS f0000601
MESSY goroutine 12 start: expected to be in NS f0000075 but am instead in NS f0000601
MESSY goroutine 11 start: expected to be in NS f0000075 but am instead in NS f0000601
             main waiting in namespace   f0000075
MESSY goroutine  6 end: expected to be in NS f0000601 but am instead in NS f0000075
MESSY goroutine  5 end: expected to be in NS f0000601 but am instead in NS f0000075
MESSY goroutine 12 end: expected to be in NS f0000601 but am instead in NS f0000075
MESSY goroutine 14 start: expected to be in NS f0000075 but am instead in NS f0000601
               main stop  in namespace   f0000601

root@ubuntu-xenial:/# ns-mess
               main start in namespace   f0000075
              newns start in namespace   f0000635
MESSY goroutine  4 start: expected to be in NS f0000075 but am instead in NS f0000635
MESSY goroutine  3 start: expected to be in NS f0000075 but am instead in NS f0000635
MESSY goroutine  6 start: expected to be in NS f0000075 but am instead in NS f0000635
MESSY goroutine 12 start: expected to be in NS f0000075 but am instead in NS f0000635
MESSY goroutine 11 start: expected to be in NS f0000075 but am instead in NS f0000635
             main waiting in namespace   f0000075
MESSY goroutine 24 end: expected to be in NS f0000075 but am instead in NS f0000635
MESSY goroutine 49 start: expected to be in NS f0000075 but am instead in NS f0000635
MESSY goroutine  0 end: expected to be in NS f0000075 but am instead in NS f0000635
MESSY goroutine  4 end: expected to be in NS f0000635 but am instead in NS f0000075
MESSY goroutine  3 end: expected to be in NS f0000635 but am instead in NS f0000075
MESSY goroutine 11 end: expected to be in NS f0000635 but am instead in NS f0000075
MESSY goroutine 12 end: expected to be in NS f0000635 but am instead in NS f0000075
MESSY goroutine 16 end: expected to be in NS f0000075 but am instead in NS f0000635
               main stop  in namespace   f0000635
```
