#include <Arduino_JSON.h>
#include "./common/correlationId.h"

String getSendMessageToPutDevice(String id, String status);

void handlePutDeviceAfterReceived(JSONVar response);