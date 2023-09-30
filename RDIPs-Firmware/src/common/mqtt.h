#include <PubSubClient.h>
#include <UUID.h>
#include <Arduino_JSON.h>
#include "handler/postDevice.h"
#include "handler/putDetailDevice.h"

extern PubSubClient client; // Declare the PubSubClient instance

extern UUID uuid;

void setupMqtt();

void callback(char *topic, byte *payload, unsigned int length);

String getSendTopic(String functionName);

String getReceiveTopic();

String getDeviceName();

void handleMessageReceived(char *topic, String receiveMessage);

