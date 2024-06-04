#include <Arduino.h>
#include <Wire.h>
#include "gyro.h"
#include "motor.h"

Gyro gyro;
Motor motor;

bool checkDevice();
void rotate(int destinationDegrees);

void setup()
{
  Serial.begin(9600); // Initiate serial communication for printing the results on the Serial monitor
  Wire.begin();       // Initiate the Wire library
  // Set ADXL345 in measuring mode
  Wire.beginTransmission(gyro.ADXL345); // Start communicating with the device
  Wire.write(0x2D);                     // Access/ talk to POWER_CTL Register - 0x2D
  // Enable measurement
  Wire.write(8); // (8dec -> 0000 1000 binary) Bit D3 High for measuring enable
  Wire.endTransmission();
  delay(10);
  checkDevice();
}

void loop()
{
  if (checkDevice == false)
  {
    exit(0);
  }
}

bool checkDevice()
{
  Serial.println("=== Initiate Horizon Model ===");
  Serial.println("Checking on gyro sensor");
  gyro.getAngle();
  if (gyro.X_degree == NAN)
  {
    Serial.println("ERROR on gyro sensor");
    Serial.println("Check wiring again!!");
    Serial.println("End initiate");
    return false;
  }
  else
  {
    gyro.displayAngleByModel(0);
    Serial.println("Gyro functional");
  }
  Serial.println("Checking on motor module");
  rotate(30);
  rotate(-30);
  rotate(0);
  Serial.println("Motor functional");

  return true;
}

//==The more forward the more Degree increase from negative to positive==//
void rotate(int destinationDegrees)
{
  Serial.println("Rolling");

  gyro.getAngle();
  int currentDegrees = (int)gyro.X_degree;
  Serial.print("currentDegrees = ");
  Serial.println(currentDegrees);

  while (abs(destinationDegrees - currentDegrees) != 0)
  {
    if (destinationDegrees > currentDegrees)
    {
      motor.forward(MAX_SPEED);
      delay(abs(destinationDegrees - currentDegrees) * 100);
      motor.stop();
    }
    else
    {
      motor.backward(MAX_SPEED);
      delay(abs(destinationDegrees - currentDegrees) * 100);
      motor.stop();
    }

    gyro.getAngle();
    currentDegrees = (int)gyro.X_degree;
    gyro.displayAngleByModel(0);
  }
}
