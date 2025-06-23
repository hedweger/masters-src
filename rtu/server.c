#include "iec61850_model.h"
#include "mms_value.h"
#include <hal_thread.h>
#include <iec61850_server.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>

#include <iec61850_config_file_parser.h>

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

  char *cfg_p = "../test.cfg";
  if (argc > 3) {
    cfg_p = argv[3];
  }
  char abs_path[100];
  realpath(cfg_p, abs_path);

  IedModel *model = ConfigFileParser_createModelFromConfigFileEx(abs_path);
  if (model == NULL) {
    fprintf(stderr, "Error: failed to parse: %s\n", abs_path);
    return 1;
  }

  IedServer server = IedServer_create(model);

  IedServer_start(server, port);
  if (!IedServer_isRunning(server)) {
    printf("Error: starting ied server failed.\n");
    goto cleanup;
  }

  int ld_count = IedModel_getLogicalDeviceCount(model);

  for (int i = 0; i < ld_count; i++) {
    char ln_path[100];
    char doi_path[100];
    char da_path[100];
    LogicalDevice *ld = IedModel_getDeviceByIndex(model, i);
    if (ld == NULL) {
      break;
    }
    printf("Logical Device %i: %s\n", i, ld->name);
    LogicalNode *lns = (LogicalNode *)ld->firstChild;
    while (lns != NULL) {
      sprintf(ln_path, "%s/%s", ld->name, lns->name);
      DataObject *doi = (DataObject *)lns->firstChild;
      while (doi != NULL) {
        sprintf(doi_path, "%s.%s", ln_path, doi->name);
        DataAttribute *da = (DataAttribute *)doi->firstChild;
        while (da != NULL) {
          sprintf(da_path, "%s.%s", doi_path, da->name);
          da = (DataAttribute *)da->sibling;
          printf("%s\n", da_path);
        }
        doi = (DataObject *)doi->sibling;
      }
      lns = (LogicalNode *)lns->sibling;
    }
  }
  char *addr = "Device1/DGEN1.Mod.q";
  DataAttribute *lln0_mod_q =
      (DataAttribute *)IedModel_getModelNodeByShortObjectReference(model, addr);
  if (lln0_mod_q == NULL) {
    printf("%s does not exist!\n", addr);
    return 1;
  }
  uint32_t val = 0;
  MmsValue *bsval = MmsValue_newBitString(8);

  signal(SIGINT, sigint_handler);
  running = 1;
  printf("IED server launch succesfully\n");
  while (running) {
    IedServer_lockDataModel(server);
    MmsValue_setBitStringFromInteger(bsval, val);
    IedServer_updateAttributeValue(server, lln0_mod_q, bsval);

    MmsValue *r_bval = IedServer_getAttributeValue(server, lln0_mod_q);
    int size = MmsValue_getBitStringSize(r_bval);
    printf("Size: %i\n", size);
    for (int i = 0; i < size; i++) {
      printf("%i", MmsValue_getBitStringBit(r_bval, i));
    }
    printf(" - %i", MmsValue_getBitStringAsInteger(r_bval));
    printf("\n");

    IedServer_unlockDataModel(server);
    val++;
    Thread_sleep(2000);
  }

  IedServer_stop(server);
cleanup:
  IedServer_destroy(server);
  IedModel_destroy(model);
  return 0;
}

