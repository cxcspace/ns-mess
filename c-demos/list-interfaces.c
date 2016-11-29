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
  int fileDescriptor;
  char *filePath;
  for (int i = 1; i < argc; i++) {
    filePath = argv[i];
    fileDescriptor = open(filePath, O_RDONLY);
    if (fileDescriptor == -1) {
      fprintf(stderr, "error: open %s\n", filePath);
      return 1;
    }

    printf("switching to %s\n", filePath);

    if (setns(fileDescriptor, 0) == -1) {
      fprintf(stderr, "error: setns\n");
      return 1;
    }

    if (printInterfaces() == -1) {
      fprintf(stderr, "printInterfaces");
      return 1;
    }

    printf("\n");
  }
}
