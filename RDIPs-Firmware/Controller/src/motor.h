
#ifndef MOTOR_H
#define MOTOR_H

#include "gyro.h"

#define ENA 9
#define IN1 7
#define IN2 8
#define SPEED_3 255
#define SPEED_2 170
#define SPEED_1 85

typedef struct Motor
{
    Gyro myGyro;
    void stop()
    {
        digitalWrite(IN1, LOW);
        digitalWrite(IN2, LOW);
    }

    //==Rotate counter-clockwise==//
    // One speed only
    void forward(int speed)
    {
        Serial.println("Forwarding");
        analogWrite(ENA, speed);
        digitalWrite(IN1, HIGH);
        digitalWrite(IN2, LOW);
    }

    //==Rotate clockwise==//
    // One speed only
    void backward(int speed)
    {
        Serial.println("Backwarding");
        analogWrite(ENA, speed);
        digitalWrite(IN1, LOW);
        digitalWrite(IN2, HIGH);
    }

    // Smoothing motor speed by increase it's speed from minimum speed to max speed
    // If direction == 0, the motor rotate forward
    // If direction == 1, the motor rotate backward
    void accelerate(int minSpeed, int maxSpeed, int direction)
    {
        if (minSpeed != maxSpeed)
        {
            for (int i = minSpeed; i <= maxSpeed; i++)
            {
                analogWrite(ENA, i);
            }
            switch (direction)
            {
            case 0:
                Serial.println("forwarding");
                digitalWrite(IN1, HIGH);
                digitalWrite(IN2, LOW);
                break;
            case 1:
                Serial.println("backwarding");
                digitalWrite(IN1, LOW);
                digitalWrite(IN2, HIGH);
                break;
            default:
                break;
            }
        }
        else
        {
            switch (direction)
            {
            case 0:
                forward(minSpeed);
                break;
            case 1:
                backward(minSpeed);
                break;
            default:
                break;
            }
        }
    }

} Motor;

#endif