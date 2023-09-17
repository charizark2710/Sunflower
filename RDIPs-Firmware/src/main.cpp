#include <Arduino.h>
#include <WiFi.h>
#include <PubSubClient.h>
#include <Arduino_JSON.h>
#include <UUID.h>

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
const String exchangeReceived = "server.*";
String correlationId;

// Topic to send messages
const String deviceName = "device_name";
const String exchangeSend = "gateway.";
const String putDetailDevice = ".PutDetailDevice";
const String getDetailDevice = ".GetDetailDevice";

WiFiClient espClient;
PubSubClient client(espClient);

UUID uuid;

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
  JSONVar response;
  response = JSON.parse(messageStr);

  // check response from server with my correlationId correctly
  String correlationIdResponse = response["CorrelationId"];
  if (correlationIdResponse == correlationId)
  {
    int httpCode = response["httpCode"];
    String message = response["message"];
    if (httpCode == 200)
    {
      Serial.println("Update Device successfully!");
      // Discard the perious log, handle later
      Serial.println("Discard the perious log");
      // Save new log, handle later
      Serial.println("Save new log");
    }
    else if (message != "")
    {
      Serial.print("Update Device fail, with error ");
      Serial.print(httpCode);
      Serial.print(" , and cause is ");
      Serial.println(message);
    }
  }
}

void handleMessage(char *topic, String receiveMessage)
{
  Serial.print("Handle message with topic is ");
  Serial.println(topic);
  String server = "server/+";

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
  String receiveMessage = "";
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

void generateCorrelationId()
{
  uuid.setVariant4Mode();
  uuid.generate();
  correlationId = uuid.toCharArray();
}

String getSendMessage(String id, String status)
{
  JSONVar data;
  JSONVar param;
  JSONVar body;

  param["id"] = id;
  body["status"] = status;

  data["CorrelationId"] = correlationId;
  data["param"] = param;
  data["body"] = body;

  String request = JSON.stringify(data);
  return request;
}

void setup()
{
  Serial.begin(115200);
  setupWifi();
  setupClient();
  generateCorrelationId();
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
      // Send message with device id and status
      String sendMessage = getSendMessage("e4264e43-01c9-49c9-adfa-0b524ea82a5f", "Sleep");
      client.publish(getSendTopic().c_str(), sendMessage.c_str());
      previousMillis = currentMillis;
    }
  }
  client.loop();
}
