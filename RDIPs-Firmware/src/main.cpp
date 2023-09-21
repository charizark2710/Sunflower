#include <Arduino.h>
#include <WiFi.h>
#include <PubSubClient.h>
#include <Arduino_JSON.h>
#include <UUID.h>

#include "common/mqtt.h"
#include "common/strings.h"
#include "common/wifi.h"

const unsigned long interval = 10000; // 10 seconds (in milliseconds)
unsigned long previousMillis = 0;
unsigned long currentMillis;

const String putDetailDevice = ".PutDetailDevice";

WiFiClient espClient;
PubSubClient client(espClient);

UUID uuid;

void setup()
{
  Serial.begin(115200);
  setupWifi();
  setupMqtt();
  generateCorrelationId();
  client.subscribe(getReceiveTopic().c_str());
  client.subscribe(getSendTopic(putDetailDevice).c_str());
}

void loop()
{
  currentMillis = millis();

  //Call api after 10 seconds
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
