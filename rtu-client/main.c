#include <hal_thread.h>
#include <iec61850_client.h>
#include <stdio.h>
#include <stdlib.h>

static int running = 0;

void sigint_handler(int signalId) { running = 0; }

int main(int argc, char **argv) {
  char *hostname = "localhost";
  int tcpPort = 102;
  IedClientError error;
  IedConnection con = IedConnection_create();
  printf("Connecting to %s:%i\n", hostname, tcpPort);
  IedConnection_connect(con, &error, hostname, tcpPort);
  if (error == IED_ERROR_OK) {
    printf("Connected\n");
    running = 1;
    while (running == 1) {
      MmsValue *value = IedConnection_readObject(
          con, &error, "simpleIOGenericIO/GGIO1.AnIn1.mag.f", IEC61850_FC_MX);
      if (value != NULL) {
        if (MmsValue_getType(value) == MMS_FLOAT) {
          float fval = MmsValue_toFloat(value);
          printf("read float value: %f\n", fval);
        } else if (MmsValue_getType(value) == MMS_DATA_ACCESS_ERROR) {
          printf("Failed to read value (error code: %i)\n",
                 MmsValue_getDataAccessError(value));
        }
        MmsValue_delete(value);
      }
    }
  } else {
    printf("Failed to connect to %s:%i\n", hostname, tcpPort);
    Thread_sleep(60000);
  }
  IedConnection_destroy(con);
  return 0;
}
