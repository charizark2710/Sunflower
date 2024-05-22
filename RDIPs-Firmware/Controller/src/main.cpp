#include <Arduino.h>
#include "Gyro.cpp"

Gyro gyro;

void setup()
{
  Serial.begin(9600); // Initiate serial communication for printing the results on the Serial monitor
  Wire.begin();       // Initiate the Wire library
  // Set ADXL345 in measuring mode
  Wire.beginTransmission(gyro.ADXL345); // Start communicating with the device
  Wire.write(0x2D);                // Access/ talk to POWER_CTL Register - 0x2D
  // Enable measurement
  Wire.write(8); // (8dec -> 0000 1000 binary) Bit D3 High for measuring enable
  Wire.endTransmission();
  delay(10);
}

void loop()
{
  gyro.getAngle();
}