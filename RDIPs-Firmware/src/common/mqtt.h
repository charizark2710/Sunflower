#include <PubSubClient.h>
#include <Arduino_JSON.h>

extern PubSubClient client; // Declare the PubSubClient instance

void setupMqtt();

void callback(char *topic, byte *payload, unsigned int length);

String getSendTopic(String functionName);

String getReceiveTopic();

String getDeviceName();

void handleMessageReceived(char *topic, String receiveMessage);

String getResponseMethod(String response);