/*
   inherit /path/to/namespace1 /path/to/namespace2

   Will switch to the first namespace, then spawn a thread which will switch
   to the second namespace.
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


int getThreadID() {
  return syscall(SYS_gettid);
}

long int getInodeOfCurrentNetNS() {
  int threadID = getThreadID();
  char myNSPath[100];
  sprintf(myNSPath, "/proc/self/task/%d/ns/net", threadID);

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
  printf("\tswitched to %s\n", namespacePath);
}

void report(char* msg) {
  fprintf(stdout, "%30s: on thread %d in netns %lx\n", msg, getThreadID(), getInodeOfCurrentNetNS());
}

void* threadWorker(char* namespacePath) {
  report("started thread");
  switchToNamespace(namespacePath);
  report("thread switched");
  return NULL;
}

void launchThreadAndWait(char* threadWorkerArg) {
  pthread_t thread;
  pthread_create(&thread, NULL,
        (void *(*) (void *)) threadWorker,
        (void *) threadWorkerArg
    );

  void* res;
  if (pthread_join(thread, &res) != 0) {
    perror("error: join");
    exit(1);
  }
}


int main(int argc, char **argv) {
  if (argc != 3) {
    fprintf(stderr, "usage: %s nspath1 nspath2\n", argv[0]);
    exit(1);
  }

  report("main started");
  switchToNamespace(argv[1]);
  report("main, after first switch");
  printf("\tcreating new thread...\n");
  launchThreadAndWait(argv[2]);
  report("returned to main");

  return 0;
}
