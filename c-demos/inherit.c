/*
   inherit.c

   Will switch to the namespace provided as an arg, and will then spawn a thread.
*/

#define _GNU_SOURCE
#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>
#include <pthread.h>
#include <sys/stat.h>
#include <sys/syscall.h>

void switchToNamespace(char* namespacePath) {
  int namespaceFd = open(namespacePath, O_RDONLY);
  if (namespaceFd < 0) {
    perror("error: open namespace");
    exit(1);
  }

  if (setns(namespaceFd, 0) < 0) {
    perror("error: setns");
    exit(1);
  }
}

int getThreadID() {
  return syscall(SYS_gettid);
}

long int getInodeOfCurrentNetNS() {
  char myNSPath[100];
  sprintf(myNSPath, "/proc/self/task/%d/ns/net", getThreadID());

  int currentNS = open(myNSPath, O_RDONLY);
  if (currentNS < 0) {
      perror("error: open namespace");
      exit(1);
  }

  struct stat nsStat;
  if (fstat (currentNS, &nsStat) < 0) {
      perror("error: stat namespace");
      exit(1);
  }
  close(currentNS);

  return nsStat.st_ino;
}

void report(char* msg) {
  fprintf(stdout, "%20s: on thread %d in netns %ld\n", msg, getThreadID(), getInodeOfCurrentNetNS());
}

void* threadWorker(void* _) {
  report("in new thread");
  return NULL;
}

void launchThreadAndWait() {
  pthread_t thread;
  pthread_create(&thread, NULL,
        (void *(*) (void *)) threadWorker,
        (void *) NULL
    );

  void* res;
  // wait for the thread to complete
  if (pthread_join(thread, &res) != 0) {
    perror("error: join");
    exit(1);
  }
}


int main(int argc, char **argv) {
  if (argc != 2) {
    fprintf(stderr, "usage: %s nspath1\n", argv[0]);
    exit(1);
  }

  report("main started");

  switchToNamespace(argv[1]);
  printf("switched to %s\n", argv[1]);
  report("main, after switch");

  printf("creating new thread...\n");
  launchThreadAndWait();
}
