#include <Arduino.h>
#include <WiFi.h>
#include <PubSubClient.h>
#include <ArduinoJson.h>

#define RELAY 12
#define LED 2

// WiFi
const char *ssid = "Ngo Van Nhieu";
const char *password = "01698592252";

// MQTT Broker
const char *mqtt_server = "192.168.1.183";
const char *mqtt_username = "admin";
const char *mqtt_password = "admin";
const int mqtt_port = 1883;

const unsigned long interval = 10000; // 10 seconds (in milliseconds)
unsigned long previousMillis = 0;
unsigned long currentMillis;

// Topic to receive messages from server
const String correlationId = "392abe2a-33b3-4ec7-88e8-629c1c7efd7d"; //which is gen by uuid library
const String exchangeReceived = "server.*";

// Topic to send messages
const String deviceName = "device_name";
const String exchangeSend = "gateway.";
const String putDetailDevice = ".PutDetailDevice";
const String getDetailDevice = ".GetDetailDevice";

WiFiClient espClient;
PubSubClient client(espClient);

String receiveMessage;

void setupWifi()
{
  delay(100);
  Serial.print("\nConnecting to ");
  Serial.print(ssid);

  WiFi.begin(ssid, password);
  WiFi.mode(WIFI_STA);

  while (WiFi.status() != WL_CONNECTED)
  {
    Serial.print(".");
    delay(100);
  }

  Serial.print("\nConnected to ");
  Serial.println(ssid);
}

String getSendTopic()
{
  return exchangeSend + deviceName + putDetailDevice;
}

String getReceiveTopic()
{
  return exchangeReceived;
}

void parseStringToJson(String messageStr)
{
  Serial.println("parseStringToJson");
  DynamicJsonDocument doc(256);
  DeserializationError error = deserializeJson(doc, messageStr.c_str());

  if (error)
  {
    Serial.print("Deserialization failed: ");
    Serial.println(error.c_str());
    return;
  }

  int httpCode = doc["httpCode"];
  const char *message = doc["message"];

  Serial.print("Canh ");
  Serial.println(httpCode);
  int success = 200;
  if (httpCode == success)
  {
    Serial.println("Update Device successfully!");
    // Discard the perious log
    Serial.println("Discard the perious log");
    Serial.println("Save new log");
    // Save new log
  }
  else if ((String)message != "")
  {
    Serial.print("Update Device fail, with error ");
    Serial.print(httpCode);
    Serial.print(" , and cause is ");
    Serial.println(message);
  }
}

void handleMessage(char *topic, String receiveMessage)
{
  Serial.println("Handle message ");
  Serial.println(topic);
  String server = "server/";

  if ((String)topic == server)
  {
    parseStringToJson(receiveMessage);
  }
}

void callback(char *topic, byte *payload, unsigned int length)
{
  Serial.println("-----------------------");
  Serial.print("Message arrived in topic: ");
  Serial.println(topic);
  Serial.print("Message:");
  receiveMessage = "";
  for (int i = 0; i < length; i++)
  {
    receiveMessage.concat((char)payload[i]);
    Serial.print((char)payload[i]);
  }
  Serial.println();
  handleMessage(topic, receiveMessage);
}

void setupClient()
{
  client.setServer(mqtt_server, mqtt_port);
  client.setCallback(callback);

  while (!client.connected())
  {
    String client_id = "esp32-client-";
    Serial.printf("The client %s connects to the public MQTT broker\n", client_id.c_str());
    if (client.connect(deviceName.c_str(), mqtt_username, mqtt_password))
    {
      Serial.println("Connected to RabbitMQ");
    }
    else
    {
      Serial.print("Connection failed, rc=");
      Serial.println(client.state());
      delay(2000);
    }
  }
}

String getSendMessage(String correlationId, String id, String status)
{
  StaticJsonDocument<256> doc;
  doc["correlationId"] = correlationId;

  JsonObject param = doc.createNestedObject("param");
  param["id"] = id;

  JsonObject body = doc.createNestedObject("body");
  body["status"] = status;

  String message;
  serializeJson(doc, message);
  return message;
}

void setup()
{
  delay(3000);
  Serial.begin(115200);
  setupWifi();
  setupClient();

  client.subscribe(getReceiveTopic().c_str());

  client.subscribe(getSendTopic().c_str());
}

void loop()
{
  currentMillis = millis();
  if (currentMillis - previousMillis > interval)
  {
    if (WiFi.status() == WL_CONNECTED)
    {
      //Send message with device id and status
      String sendMessage = getSendMessage(correlationId, "e4264e43-01c9-49c9-adfa-0b524ea82a5f", "Sleep");
      client.publish(getSendTopic().c_str(), sendMessage.c_str());
      previousMillis = currentMillis;
    }
  }
  client.loop();
}
