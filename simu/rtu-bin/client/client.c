#include "linked_list.h"
#include "mms_value.h"
#include <hal_thread.h>
#include <iec61850_client.h>
#include <signal.h>
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

  ClientDataSet cds = NULL;
  char lds_path[100];

  LinkedList ld = IedConnection_getLogicalDeviceList(conn, &err)->next;
  while (ld != NULL) {
    printf("Logical Device: %s\n", (char *)ld->data);
    LinkedList lds =
        IedConnection_getLogicalDeviceDataSets(conn, &err, (char *)ld->data)
            ->next;
    while (lds != NULL) {
      snprintf(lds_path, sizeof(lds_path), "%s/%s", (char *)ld->data,
               (char *)lds->data);
      cds = IedConnection_readDataSetValues(conn, &err, lds_path, NULL);
      if (cds == NULL) {
        printf("Failed to read DataSet %s: %s\n", lds_path,
               IedClientError_toString(err));
      }
      lds = lds->next;
    }
    ld = ld->next;
  }

  while (running) {
    IedConnection_readDataSetValues(conn, &err, lds_path, cds);
    MmsValue *cds_arr = ClientDataSet_getValues(cds);
    int size = MmsValue_getArraySize(cds_arr);
    printf("DataSet %s:\n", lds_path);
    for (int i = 0; i < size; i++) {
      MmsValue *val = MmsValue_getElement(cds_arr, i);
      printf("%s : ", MmsValue_getTypeString(val));
      if (MmsValue_getType(val) == MMS_BIT_STRING) {
        for (int j = 0; j < MmsValue_getBitStringSize(val); j++) {
          printf("%i ", MmsValue_getBitStringBit(val, j));
        }
        printf(" - %i", MmsValue_getBitStringAsInteger(val));
        printf("\n");
      } else if (MmsValue_getType(val) == MMS_INTEGER) {
        printf("%i\n", MmsValue_toInt32(val));
      } else {
        printf(" unknown type\n");
      }
    }
    Thread_sleep(1000);
  }
cleanup:
  IedConnection_close(conn);
  IedConnection_destroy(conn);
  return 0;
}
