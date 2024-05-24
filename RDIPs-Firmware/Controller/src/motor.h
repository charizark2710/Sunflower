
#ifndef MOTOR_H
#define MOTOR_H

#include "gyro.h"

#define IN1 7
#define IN2 8
#define MAX_SPEED 255
#define MIN_SPEED 0

typedef struct Motor
{
    Gyro myGyro;
    void stop()
    {
        digitalWrite(IN1, LOW);
        digitalWrite(IN2, LOW);
    }
    //==Rotate counter-clockwise==//
    void forward(int speed)
    {
        speed = constrain(speed, MIN_SPEED, MAX_SPEED);
        digitalWrite(IN1, HIGH);
        analogWrite(IN2, 255 - speed);
    }
    //==Rotate clockwise==//
    void backward(int speed)
    {
        speed = constrain(speed, MIN_SPEED, MAX_SPEED);
        digitalWrite(IN1, LOW);
        analogWrite(IN2, speed);
    }
    //==The more forward the more Degree increase from negative to positive==//
    void rotate(float destinationDegrees)
    {
        float currentDegrees = myGyro.X_degree;
        while (currentDegrees != destinationDegrees)
        {
            if (destinationDegrees > currentDegrees)
            {
                forward(MAX_SPEED);
            }
            else
            {
                backward(MAX_SPEED);
            }
        } 
        stop();
    }
} Motor;

#endif