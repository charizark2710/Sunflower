#include <map>
#include "mqtt.h"
#include "strings.h"
#include "handler/postDevice.h"
#include "handler/putDetailDevice.h"
#include "correlationId.h"

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

String getReceiveTopic()
{
  return exchange_received;
}

String getDeviceName()
{
  return device_name;
}

String getResponseMethod(String response)
{
  return strtok((char *)response.c_str(), "-");;
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
    String responseMethod = getResponseMethod(correlationIdResponse);

    std::map<std::string, int> methodMap;
    methodMap["PostDevice"] = 1;
    methodMap["PutDetailDevice"] = 2;

    if (methodMap.find(responseMethod.c_str()) != methodMap.end())
    {
      int command = methodMap[responseMethod.c_str()];
      switch (command)
      {
      case 1:
        Serial.println("Handle POST device");
        handlePostDeviceResponse(response);
        break;
      case 2:
        Serial.println("Handle PUT device");
        handlePutDeviceAfterReceived(response);
        break;
      }
      removeMatchedCorrelationId(correlationIdResponse);
    }
  }
}
