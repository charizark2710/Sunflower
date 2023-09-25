#include <Arduino.h>
#include "common/mqtt.h"
#include "common/strings.h"
#include "common/wifi.h"

const unsigned long interval = 10000; // 10 seconds (in milliseconds)
unsigned long previousMillis = 0;
unsigned long currentMillis;

const String postDevice = ".PostDevice";
const String putDetailDevice = ".PutDetailDevice";

WiFiClient espClient;
PubSubClient client(espClient);

UUID uuid;

void setup()
{
  Serial.begin(115200);
  setupWifi();
  setupMqtt();
  client.subscribe(getReceiveTopic().c_str());
  client.subscribe(getSendTopic(putDetailDevice).c_str());
  client.subscribe(getSendTopic(postDevice).c_str());

  // Call postDeviceApi
  client.publish(
      getSendTopic(postDevice).c_str(),
      getSendMessageToPostDevice(
          WiFi.macAddress(),
          "family",
          "Active",
          1,
          1)
          .c_str());
}

void loop()
{
  currentMillis = millis();

  // Call api after 10 seconds
  if (currentMillis - previousMillis > interval)
  {
    if (WiFi.status() == WL_CONNECTED)
    {
      // Send message with device id and status
      client.publish(
          getSendTopic(putDetailDevice).c_str(),
          getSendMessageToPutDevice(
              "e4264e43-01c9-49c9-adfa-0b524ea82a5f",
              "Sleep")
              .c_str());

      previousMillis = currentMillis;
    }
  }
  client.loop();
}
