#include <hal_thread.h>
#include <iec61850_model.h>
#include <iec61850_server.h>
#include <mms_value.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>

#include <iec61850_config_file_parser.h>

static int running = 0;

void sigint_handler(int signalId) { running = 0; }

void read_value(IedServer server, DataAttribute *da) {
  MmsValue *value = IedServer_getAttributeValue(server, da);
  if (MmsValue_getType(value) == MMS_BIT_STRING) {
    int size = MmsValue_getBitStringSize(value);
    printf(" - ");
    for (int i = 0; i < size; i++) {
      printf("%i", MmsValue_getBitStringBit(value, i));
    }
    printf(" - %i", MmsValue_getBitStringAsInteger(value));
    printf("\n");
  } else if (MmsValue_getType(value) == MMS_INTEGER) {
    printf(" - %i\n", MmsValue_toInt32(value));
  } else {
    printf("unknown value type: %s!\n", MmsValue_getTypeString(value));
  }
}

int main(int argc, char **argv) {
  int port = 69;
  if (argc > 1) {
    port = atoi(argv[1]);
  }

  char *hostname = "localhost";
  if (argc > 2) {
    hostname = argv[2];
  }

  char *cfg_p = "test.cfg";
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

	printf("%s\n", hostname);
  IedServer server = IedServer_create(model);


  IedServer_start(server, port);
  if (!IedServer_isRunning(server)) {
    printf("Error: starting ied server failed.\n");
    goto cleanup;
  }

	IedServer_setLocalIpAddress(server, hostname);


  int ld_count = IedModel_getLogicalDeviceCount(model);

  // for (int i = 0; i < ld_count; i++) {
  //   char ln_path[100];
  //   char doi_path[100];
  //   char da_path[100];
  //   LogicalDevice *ld = IedModel_getDeviceByIndex(model, i);
  //   if (ld == NULL) {
  //     break;
  //   }
  //   printf("Logical Device %i: %s\n", i, ld->name);
  //   LogicalNode *lns = (LogicalNode *)ld->firstChild;
  //   while (lns != NULL) {
  //     sprintf(ln_path, "%s/%s", ld->name, lns->name);
  //     DataObject *doi = (DataObject *)lns->firstChild;
  //     while (doi != NULL) {
  //       sprintf(doi_path, "%s.%s", ln_path, doi->name);
  //       DataAttribute *da = (DataAttribute *)doi->firstChild;
  //       while (da != NULL) {
  //         sprintf(da_path, "%s.%s", doi_path, da->name);
  //         da = (DataAttribute *)da->sibling;
  //         printf("%s\n", da_path);
  //       }
  //       doi = (DataObject *)doi->sibling;
  //     }
  //     lns = (LogicalNode *)lns->sibling;
  //   }
  // }
  //

  char *lln0_addr = "Device1/LLN0.Mod.q";
  char *mmxu_addr = "Device1/MMXU1.Mod.q";
  char *ctl_addr = "Device1/MMXU1.Mod.ctlModel";
  DataAttribute *lln0_mod_q =
      (DataAttribute *)IedModel_getModelNodeByShortObjectReference(model,
                                                                   lln0_addr);
  if (lln0_mod_q == NULL) {
    printf("%s does not exist!\n", lln0_addr);
    return 1;
  }
  DataAttribute *mmxu_mod_q =
      (DataAttribute *)IedModel_getModelNodeByShortObjectReference(model,
                                                                   mmxu_addr);
  if (mmxu_mod_q == NULL) {
    printf("%s does not exist!\n", mmxu_addr);
    return 1;
  }
  DataAttribute *ctl_mod_q =
      (DataAttribute *)IedModel_getModelNodeByShortObjectReference(model,
                                                                   ctl_addr);
  if (ctl_mod_q == NULL) {
    printf("%s does not exist!\n", ctl_addr);
    return 1;
  }

  uint32_t val = 0;
  MmsValue *bsval = MmsValue_newBitString(8);
  MmsValue *intval = MmsValue_newInteger(32);

  signal(SIGINT, sigint_handler);
  running = 1;

  while (running) {
    IedServer_lockDataModel(server);
    MmsValue_setBitStringFromInteger(bsval, val);
    IedServer_updateAttributeValue(server, lln0_mod_q, bsval);
    IedServer_updateAttributeValue(server, mmxu_mod_q, bsval);
    IedServer_updateAttributeValue(server, ctl_mod_q, intval);

    IedServer_unlockDataModel(server);
    printf("%s:\n", lln0_addr);
    read_value(server, lln0_mod_q);
    printf("%s:\n", mmxu_addr);
    read_value(server, mmxu_mod_q);
    printf("%s:\n", ctl_addr);
    read_value(server, ctl_mod_q);

    val++;
    MmsValue_setInt32(intval, val);
    Thread_sleep(1000);
  }

  IedServer_stop(server);
cleanup:
  IedServer_destroy(server);
  IedModel_destroy(model);
  return 0;
}
