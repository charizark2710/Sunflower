#include <Arduino.h>
#include <WiFi.h>
#include <PubSubClient.h>
#include <Arduino_JSON.h>
#include <UUID.h>

#define RELAY 12
#define LED 2

const char *ssid = "Trung Thanh";
const char *password = "thanh0106";

// MQTT Broker
const char *mqtt_server = "192.168.1.3"; // 192.168.1.183 //172.24.64.1
const char *mqtt_username = "admin";     // admin
const char *mqtt_password = "admin";     // admin
const int mqtt_port = 1883;              // 1883

const unsigned long interval = 10000; // 10 seconds (in milliseconds)
unsigned long previousMillis = 0;
unsigned long currentMillis;

// Topic to receive messages from server
const String exchangeReceived = "server.*";

// Topic to send messages
const String deviceName = "device_name";
const String exchangeSend = "gateway.";
const String postDevice = ".PostDevice";

WiFiClient espClient;
PubSubClient client(espClient);

UUID uuid;

String correlationId;
String receiveMessage;
String deviceId;

void generateCorrelationId()
{
  uuid.setVariant4Mode();
  uuid.generate();
  correlationId = uuid.toCharArray();
}

String getSendTopic(String topic)
{
  return exchangeSend + deviceName + topic;
}

String getReceiveTopic()
{
  return exchangeReceived;
}

void subscribeServerMessage(bool isPublised)
{
  if (!isPublised)
  {
    Serial.println("Publish message failed");
    return;
  }

  String receiveTopic = getReceiveTopic();
  bool isSubscribe = client.subscribe(receiveTopic.c_str());
  String result = isSubscribe ? "successfully" : "failed";

  Serial.print("Subscribe topic: ");
  Serial.print(receiveTopic);
  Serial.print(" ");
  Serial.println(result);
}

void postDeviceApi()
{
  // Create request body
  JSONVar data;
  JSONVar body;
  body["name"] = WiFi.macAddress();
  body["type"] = "family";
  body["status"] = "Active";
  body["firmware_ver"] = 1;
  body["app_ver"] = 1;

  generateCorrelationId();
  data["CorrelationId"] = uuid.toCharArray();
  data["body"] = body;

  // Convert JSON object to string
  String request = JSON.stringify(data);

  // Publish to topic: gateway.device_name.PostDevice
  Serial.println("Publish message: " + request);
  bool isPublished = client.publish(getSendTopic(postDevice).c_str(), request.c_str());

  // Subscribe messsage from server
  subscribeServerMessage(isPublished);
}

void handlePostDeviceResponse(char *topic, String receiveMessage)
{
  String receiveTopic = "server/+";
  bool isTopicMatched = (String)topic == receiveTopic;

  if (!isTopicMatched)
  {
    return;
  }

  Serial.print("Handle message with topic: ");
  Serial.println(topic);

  JSONVar response = JSON.parse(receiveMessage);
  String correlationIdResponse = response["CorrelationId"];
  if (correlationIdResponse == correlationId)
  {
    int httpCode = response["httpCode"];
    String message = response["message"];
    if (httpCode == 200)
    {
      JSONVar data = response["data"];
      String id = data["Id"];
      deviceId = id;
      Serial.print("Post device successed with id: ");
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
}

// Callback function for received messages
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
  handlePostDeviceResponse(topic, receiveMessage);
}

// Connect to WiFi network
void connectWiFi()
{
  Serial.print("Connecting to WiFi...");

  WiFi.begin(ssid, password);
  while (WiFi.status() != WL_CONNECTED)
  {
    delay(500);
    Serial.print(".");
  }

  Serial.println();
  Serial.print("WiFi connected. IP address: ");
  Serial.println(WiFi.localIP());
}

// Connect to MQTT broker
void connectMQTT()
{
  client.setServer(mqtt_server, mqtt_port);
  client.setBufferSize(1000);
  client.setCallback(callback);

  while (!client.connected())
  {
    String client_id = "esp32-client-";
    client_id += String(WiFi.macAddress());
    Serial.printf("The client %s connects to the public MQTT broker\n", client_id.c_str());
    if (client.connect(client_id.c_str(), mqtt_username, mqtt_password))
    {
      Serial.println("Connected to MQTT broker.");
    }
    else
    {
      Serial.print("Connection failed, rc=");
      Serial.println(client.state());
      delay(5000);
    }
  }
}

void setup()
{
  Serial.begin(115200);
  connectWiFi();
  connectMQTT();

  // Call postDeviceApi
  postDeviceApi();
}

void loop()
{
  currentMillis = millis();
  if (currentMillis - previousMillis > interval)
  {
    if (WiFi.status() == WL_CONNECTED)
    {
      previousMillis = currentMillis;
    }
  }
  client.loop();
  delay(5000);
}

