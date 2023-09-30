#include "postDevice.h"

String deviceId = "";

String getSendMessageToPostDevice(
    String name,
    String type,
    String status,
    int firmwareVer,
    int appVer)
{
  // Create request body
  JSONVar data;
  JSONVar body;
  body["name"] = name;
  body["type"] = type;
  body["status"] = status;
  body["firmware_ver"] = firmwareVer;
  body["app_ver"] = appVer;
  String correlationId = generateCorrelationId("PostDevice");
  data["CorrelationId"] = correlationId;
  data["body"] = body;

  // Convert JSON object to string
  String request = JSON.stringify(data);

  return request;
}

void handlePostDeviceResponse(JSONVar response)
{
  int httpCode = response["httpCode"];
  String message = response["message"];
  if (httpCode == 200)
  {
    JSONVar data = response["data"];
    String id = data["Id"];
    deviceId = id;
    Serial.print("Post device successed with deviceId: ");
    Serial.println(deviceId);
  }
  else if (message != "")
  {
    Serial.print("Post device fail, with error ");
    Serial.print(httpCode);
    Serial.print(", cause is ");
    Serial.println(message);
  }
  correlationIdsSize--;
}