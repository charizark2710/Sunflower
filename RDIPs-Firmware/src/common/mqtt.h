#include <PubSubClient.h>
#include <UUID.h>
#include <Arduino_JSON.h>

extern PubSubClient client; // Declare the PubSubClient instance

extern UUID uuid;

void setupMqtt();

void callback(char *topic, byte *payload, unsigned int length);

String getSendTopic(String functionName);

String getSendMessageToPutDevice(String id, String status);

String getReceiveTopic();

String getDeviceName();

String generateCorrelationId();

void handleMessageReceived(char *topic, String receiveMessage);

void handlePutDeviceAfterReceived(JSONVar response);

String getSendMessageToPostDevice(
    String name,
    String type,
    String status,
    int firmwareVer,
    int appVer
);

void handlePostDeviceResponse(JSONVar response);