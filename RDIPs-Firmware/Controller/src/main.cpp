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

  pinMode(ENA, OUTPUT);
  pinMode(IN1, OUTPUT);
  pinMode(IN2, OUTPUT);

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
  // gyro.getAngle();
  // gyro.displayAngleByModel(0);
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
  rotate(45);
  delay(1000);
  rotate(30);
  delay(1000);
  rotate(-30);
  delay(1000);
  rotate(-45);
  delay(1000);
  rotate(0);
  Serial.println("Motor functional");

  return true;
}

//==The more forward the more Degree increase from negative to positive==//
void rotate(int destinationDegrees)
{

  Serial.println("Rolling");

  // Get current degrees
  gyro.getAngle();
  int currentDegrees = (int)gyro.X_degree;
  // Serial.print("currentDegrees = ");
  // Serial.println(currentDegrees);

  int distant = abs(destinationDegrees - currentDegrees); // get absolute distant
  int speed[2];                                           // speed[0] is minimum speed, speed[1] is maximum speed
  int duration = 0;                                       // the duration of the motor while it running
  // int times = 0;

  // The farther the distance is, the faster and longer the motor run
  while (distant != 0)
  {

    distant = abs(destinationDegrees - currentDegrees); // Get distant again to calculate speed and duration
    // Serial.print("Distant: ");
    // Serial.println(distant);

    if (distant > 30)
    {
      speed[0] = SPEED_2;
      speed[1] = SPEED_3;
      duration = 1000;
    }

    else if (distant > 10 && distant <= 20)
    {
      speed[0] = SPEED_1;
      speed[1] = SPEED_2;
      duration = 200;
    }

    else if (distant > 4 && distant <= 10)
    {
      speed[0] = SPEED_1;
      speed[1] = SPEED_2;
      duration = 170;
    }

    else if (distant > 1 && distant <= 4)
    {
      speed[0] = SPEED_1;
      speed[1] = SPEED_2;
      duration = 100;
    }

    else if (distant == 1)
    {
      speed[0] = SPEED_1;
      speed[1] = SPEED_1;
      duration = 160;
    }

    else if (distant == 0)
    {
      motor.stop();
      Serial.println("Finish rotation");
      // Serial.print("Time: ");
      // Serial.println(times);
      break;
    }

    // If destinationn degrees is greater than current degrees, the motor will rotate forward
    if (destinationDegrees > currentDegrees)
    {
      motor.accelerate(speed[0], speed[1], 0);
      delay(duration);
      motor.stop();
      // times += 1;
    }

    // If destinationn degrees is smaller than current degrees, the motor will rotate backward
    else if (destinationDegrees < currentDegrees)
    {
      motor.accelerate(speed[0], speed[1], 1);
      delay(duration);
      motor.stop();
      // times += 1;
    }

    // get current degrees again to re-calculate speed and duration
    gyro.getAngle();
    currentDegrees = (int)gyro.X_degree;
    gyro.displayAngleByModel(0);
  }
}
