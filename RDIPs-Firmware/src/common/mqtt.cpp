#include "mqtt.h"
#include "strings.h"

// config mqtt
const char *mqtt_server = MQTT_SERVER;
const char *mqtt_username = MQTT_USERNAME;
const char *mqtt_password = MQTT_PASSWORD;
const int mqtt_port = MQTT_PORT;
const int mqtt_buffer_size = MQTT_BUFFER_SIZE;

// Topic to send messages
const String exchange_send = "gateway.";
const String device_name = "device_name";

// Topic to receive from server
const String exchange_received = "server.*";

String putCorrelationId = ""; // default value
String postCorrelationId = "";
String deviceId = "";

void setupMqtt()
{
  client.setServer(mqtt_server, mqtt_port);
  client.setBufferSize(mqtt_buffer_size);
  client.setCallback(callback);

  while (!client.connected())
  {
    String client_id = "esp32-client-";
    Serial.printf("The client %s connects to the public MQTT broker\n", client_id.c_str());
    if (client.connect(getDeviceName().c_str(), mqtt_username, mqtt_password))
    {
      Serial.println("Connected to RabbitMQ");
    }
    else
    {
      Serial.printf("Connection failed, rc= %d\n", client.state());
      delay(2000);
    }
  }
}

void callback(char *topic, byte *payload, unsigned int length)
{
  Serial.println("-----------------------");
  Serial.printf("Message arrived in topic: %s\n", topic);
  Serial.print("Message: ");
  String receiveMessage = "";
  for (int i = 0; i < length; i++)
  {
    receiveMessage.concat((char)payload[i]);
  }
  Serial.println(receiveMessage);
  handleMessageReceived(topic, receiveMessage);
}

String getSendTopic(String functionName)
{
  return exchange_send + device_name + functionName;
}

String getSendMessageToPutDevice(String id, String status)
{
  JSONVar data;
  JSONVar param;
  JSONVar body;

  param["id"] = id;
  body["status"] = status;

  putCorrelationId = generateCorrelationId();
  data["CorrelationId"] = putCorrelationId;
  data["param"] = param;
  data["body"] = body;

  String request = JSON.stringify(data);
  return request;
}

String getReceiveTopic()
{
  return exchange_received;
}

String getDeviceName()
{
  return device_name;
}

String generateCorrelationId()
{
  uuid.setVariant4Mode();
  uuid.generate();
  return uuid.toCharArray();
}

void handleMessageReceived(char *topic, String receiveMessage)
{
  Serial.printf("Handle message with topic is %s\n", topic);
  if ((String)topic == "server/+")
  {
    Serial.println("handleReceivedMessage");
    JSONVar response = parseStringToJson(receiveMessage);
    // check response from server with my correlationId correctly
    String correlationIdResponse = response["CorrelationId"];

    // TOD0: handle other message arrived from server
    if (strcmp(putCorrelationId.c_str(), correlationIdResponse.c_str()) == 0)
    {
      Serial.println("Handle PUT device");
      handlePutDeviceAfterReceived(response);
    }
    if (strcmp(postCorrelationId.c_str(), correlationIdResponse.c_str()) == 0)
    {
      Serial.println("Handle POST device");
      handlePostDeviceResponse(response);
    }
  }
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
  postCorrelationId = generateCorrelationId();
  data["CorrelationId"] = postCorrelationId;
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
}
