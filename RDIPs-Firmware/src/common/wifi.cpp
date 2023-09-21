#include <Arduino.h>
#include <WiFi.h>

const char *ssid = WIFI_NAME_SSID;
const char *password = WIFI_PASSWORD;

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
