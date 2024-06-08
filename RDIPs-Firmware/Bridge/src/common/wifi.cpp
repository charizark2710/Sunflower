#include "wifi.h"

const char *ssid = "TrungAnh1";
const char *password = "05113617364";

void setupWifi()
{
    delay(100);
    Serial.printf("\nConnecting to %s\n", &ssid);

    WiFi.begin(ssid, password);
    WiFi.mode(WIFI_STA);

    while (WiFi.status() != WL_CONNECTED)
    {
        Serial.print(".");
        delay(100);
    }
    Serial.printf("\nConnected to %s\n", &ssid);
}
