#include <Arduino.h>
#include <PubSubClient.h>
#include <UUID.h>

extern PubSubClient client; // Declare the PubSubClient instance
extern UUID uuid;

void setupMqtt();

void callback(char *topic, byte *payload, unsigned int length);

String getSendTopic(String functionName);

String getSendMessageToPutDevice(String id, String status);

String getReceiveTopic();

String getDeviceName();

void generateCorrelationId();

String getCorrelationId();

void handleMessageReceived(char *topic, String receiveMessage);

void handlePutDeviceAfterReceived(String messageStr);
