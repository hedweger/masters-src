#include "linked_list.h"
#include "mms_value.h"
#include <hal_thread.h>
#include <iec61850_client.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>

static int running = 0;

void sigint_handler(int signalId) { running = 0; }

int main(int argc, char **argv) {
  
  int port = 69;
  if (argc > 1) {
    port = atoi(argv[1]);
  }

  char *hostname = "localhost";
  if (argc > 2) {
    hostname = argv[2];
  }

  IedClientError err;
  IedConnection conn = IedConnection_create();
conn:
  IedConnection_connect(conn, &err, hostname, port);
  printf("Connecting to %s:%i\n", hostname, port);
  if (err != IED_ERROR_OK) {
    printf("Failed to connect. Trying again in 10 seconlds!\n");
    Thread_sleep(10000);
    goto conn;
  }

  signal(SIGINT, sigint_handler);
  printf("Connected!\n");
  running = 1;

  ClientDataSet cds;
  char lds_path[100];

  while (running) {
    if (cds == NULL) {
      LinkedList ld = IedConnection_getLogicalDeviceList(conn, &err);
      ld = ld->next;
      while (ld != NULL) {
        LinkedList lds = IedConnection_getLogicalDeviceDataSets(
            conn, &err, (char *)ld->data);
        lds = lds->next;
        while (lds != NULL) {
          strcpy(lds_path, (char *)lds->data);
          cds = IedConnection_readDataSetValues(conn, &err, (char *)lds->data,
                                                NULL);
        }
        lds = lds->next;
      }
      ld = ld->next;
    }
    IedConnection_readDataSetValues(conn, &err, lds_path, cds);
    MmsValue *cds_arr = ClientDataSet_getValues(cds);
    int size = MmsValue_getArraySize(cds_arr);
    printf("DataSet %s:\n", lds_path);
    for (int i = 0; i < size; i++) {
      MmsValue *val = MmsValue_getElement(cds_arr, i);
      printf("%s\n", MmsValue_getTypeString(val));
    }
    Thread_sleep(10000);
  }

cleanup:
  IedConnection_close(conn);
  IedConnection_destroy(conn);
  return 0;
}
