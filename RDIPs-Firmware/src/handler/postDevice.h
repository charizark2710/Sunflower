#include <Arduino_JSON.h>
#include "./common/correlationId.h"

String getSendMessageToPostDevice(
    String name,
    String type,
    String status,
    int firmwareVer,
    int appVer
);

void handlePostDeviceResponse(JSONVar response);