/*
   netns-exec /path/to/namespace program arg1 arg2 ...
   Will switch to that namespace and execute the program with the args
*/

#define _GNU_SOURCE
#include <fcntl.h>
#include <sched.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

int main(int argc, char **argv)
{
  if (argc < 2) {
    fprintf(stderr, "usage: %s namespace program args...\n", argv[0]);
    exit(1);
  }

  // acquire a handle to the namespace
  int namespaceFd = open(argv[1], O_RDONLY);
  if (namespaceFd == -1) {
    fprintf(stderr, "error: open %s\n", argv[1]);
    exit(1);
  }

  // switch this thread to the namespace
  if (setns(namespaceFd, 0) == -1) {
    fprintf(stderr, "error: setns\n");
    exit(1);
  }

  // execute the program
  if (execvp(argv[2], &argv[2]) < 0) {
    fprintf(stderr, "error: execvp\n");
    exit(1);
  }
}
