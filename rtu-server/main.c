#include <hal_thread.h>
#include <iec61850_server.h>
#include <stdio.h>
#include <stdlib.h>

#include <iec61850_config_file_parser.h>

static int running = 0;

void sigint_handler(int signalId) { running = 0; }

int main(int argc, char **argv) {
  int tcpPort = 102;
  IedModel *model = ConfigFileParser_createModelFromConfigFileEx(
      "/Users/th/workspace/masters/rtu-server/test.cfg");
  if (model == NULL) {
    printf("Error parsing config file!\n");
    return 1;
  }

  IedServer iedServer = IedServer_create(model);
  DataAttribute *anIn1_mag_f =
      (DataAttribute *)IedModel_getModelNodeByShortObjectReference(
          model, "GenericIO/GGIO1.AnIn1.mag.f");
  DataAttribute *anIn1_t =
      (DataAttribute *)IedModel_getModelNodeByShortObjectReference(
          model, "GenericIO/GGIO1.AnIn1.t");
  printf("Starting server...\n");
  IedServer_start(iedServer, tcpPort);
  if (!IedServer_isRunning(iedServer)) {
    printf("Starting server failed! Exit.\n");
    IedServer_destroy(iedServer);
    exit(-1);
  }
  running = 1;
  float val = 0.f;
  MmsValue *floatValue = MmsValue_newFloat(val);
  printf("Started!\n");
  while (running) {
    if (anIn1_mag_f != NULL) {
      MmsValue_setFloat(floatValue, val);
      IedServer_lockDataModel(iedServer);
      MmsValue_setUtcTimeMs(anIn1_t->mmsValue, Hal_getTimeInMs());
      IedServer_updateAttributeValue(iedServer, anIn1_mag_f, floatValue);
      IedServer_unlockDataModel(iedServer);
      val += 0.1f;
      printf("%f\n", val);
      Thread_sleep(100);
    }
  }
  MmsValue_delete(floatValue);
  IedServer_stop(iedServer);
  IedServer_destroy(iedServer);
  IedModel_destroy(model);
  return 0;
}
