/*
   Run with 1 or more args
   Each arg is a filepath to a namespace
   Will switch to that namespace and print all interfaces it finds
*/

#define _GNU_SOURCE
#include <fcntl.h>
#include <sched.h>
#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>

#include <ifaddrs.h>

int printInterfaces()
{
     struct ifaddrs *ifaddr, *ifa;
     int n;

     if (getifaddrs(&ifaddr) == -1) {
         return -1;
     }

     for (ifa = ifaddr, n = 0; ifa != NULL; ifa = ifa->ifa_next, n++) {
         if (ifa->ifa_addr == NULL)
             continue;
         printf("%20s\n", ifa->ifa_name);
     }

     freeifaddrs(ifaddr);
     return 0;
}

int main(int argc, char **argv)
{
  int n;
  int fd;
  for (n = 1; n < argc; n++) {
     fd = open(argv[n], O_RDONLY);
     if (fd == -1) {
       fprintf(stderr, "error: open %s\n", argv[n]);
       exit(1);
     }

     fprintf(stderr, "switching to %s\n", argv[n]);

     if (setns(fd, 0) == -1) {
       fprintf(stderr, "error: setns\n");
       exit(1);
     }

     fprintf(stderr, "interfaces:\n");

     if (printInterfaces() == -1) {
       fprintf(stderr, "printInterfaces");
       exit(1);
     }

     fprintf(stderr, "\n");
  }
}
