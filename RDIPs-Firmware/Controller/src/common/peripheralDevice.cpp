#include "Arduino.h"
#include "Wire.h"
#include "peripheralDevice.h"

bool checkI2Cstatus(int sensorAdress)
{
    Wire.beginTransmission(sensorAdress);
    Wire.write(sensorAdress);
    byte error = Wire.endTransmission();

    Wire.requestFrom(sensorAdress, 1);
    if (Wire.available())
    {
        int data = Wire.read();
        Serial.print("Sensor response: ");
        Serial.println(data, HEX);
        return true;
    }
    else
    {
        TWCR = 0;
        Serial.println("Sensor not connected or failed to respond.");
        return false;
    }
}