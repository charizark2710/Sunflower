#include "putDetailDevice.h"
#include "./common/correlationId.h"

String getSendMessageToPutDevice(String id, String status)
{
  JSONVar data;
  JSONVar param;
  JSONVar body;

  param["id"] = id;
  body["status"] = status;

  String correlationId = generateCorrelationId("PutDetailDevice");
  data["CorrelationId"] = correlationId;
  data["param"] = param;
  data["body"] = body;

  String request = JSON.stringify(data);
  return request;
}

void handlePutDeviceAfterReceived(JSONVar response)
{
  int httpCode = response["httpCode"];
  String message = response["message"];
  if (httpCode == 200)
  {
    Serial.println("Update Device successfully! - Discard the perious logs - Save new logs");
  }
  else
  {
    Serial.printf("Update Device fail! With error %d\n", httpCode);
    Serial.println("Cause is " + message);
  }
}
