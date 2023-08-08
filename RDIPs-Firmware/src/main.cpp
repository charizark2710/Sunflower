#include <Arduino.h>
#include <WiFi.h>
#include <PubSubClient.h>
#include <HttpClient.h>

#define RELAY 12
#define LED 2

const char *ssid = "Ngo Van Nhieu";
const char *password = "01698592252";

// MQTT Broker
const char *mqtt_server = "192.168.1.183"; // 192.168.1.183 //172.24.64.1
const char *topic = "devices.device1.PUT/device";
const char *mqtt_username = "admin"; // admin
const char *mqtt_password = "admin"; // admin
const int mqtt_port = 1883;          // 1883

const unsigned long interval = 10000; // 10 seconds (in milliseconds)
unsigned long previousMillis = 0;

// HttpRequest
const char *uri = "http://192.168.1.183:8080/";
const char *content_type_header = "Content-Type";
const char *content_type_value = "application/json";

WiFiClient espClient;
PubSubClient client(espClient);

void setupWifi()
{
  delay(100);
  Serial.print("\nConnecting to");
  Serial.print(ssid);

  WiFi.begin(ssid, password);
  // this setting for public wifi (no-password)
  // WiFi.begin(ssid);
  WiFi.mode(WIFI_STA);

  while (WiFi.status() != WL_CONNECTED)
  {
    Serial.print(".");
    delay(100);
  }

  Serial.print("\nConnected to");
  Serial.println(ssid);
}

void callback(char *topic, byte *payload, unsigned int length)
{
  Serial.print("Message arrived in topic: ");
  Serial.println(topic);
  Serial.print("Message:");
  for (int i = 0; i < length; i++)
  {
    Serial.print((char)payload[i]);
  }
  Serial.println();
  Serial.println("-----------------------");
}

String setupPutDeviceApi()
{
  // Make HTTP POST request
  HTTPClient http;
  http.begin(uri);
  http.addHeader(content_type_header, content_type_value);
  int httpResponseCode = http.PUT("{\"name\": \"Device 33\",\"type\": \"Family\"}");
  String response = http.getString();
  http.end();
  return response.c_str();
}

void setup()
{
  delay(3000);
  Serial.begin(115200);
  setupWifi();

  client.setServer(mqtt_server, mqtt_port);
  client.setCallback(callback);
  while (!client.connected())
  {
    String client_id = "esp32-client-";
    client_id += String(WiFi.macAddress());
    Serial.printf("The client %s connects to the public MQTT broker\n", client_id.c_str());
    if (client.connect(client_id.c_str(), mqtt_username, mqtt_password))
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
  // Publish and subscribe
  client.subscribe(topic);
  client.publish(topic, "Hi, I'm ESP32 ^^");
}

void loop()
{
  unsigned long currentMillis = millis();
  if (currentMillis - previousMillis > interval)
  {
    if (WiFi.status() == WL_CONNECTED)
    {
      previousMillis = currentMillis;

      // String response = setupPutDeviceApi();
      HTTPClient http;
      http.begin(uri);
      http.addHeader(content_type_header, content_type_value);
      int httpResponseCode = http.PUT("{\"name\": \"Device 33\",\"type\": \"Family\"}");
      String response = http.getString();
      http.end();
      // Publish the API response to an MQTT topic
      // client.publish(topic, response.c_str());
      // if(httpResponseCode > 0)
      client.publish(topic, http.getString().c_str() );
    }
  }
  client.loop();
}
