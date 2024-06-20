#include <Arduino.h>
#include <motor.h>

void Motor::stop()
{
    {
        digitalWrite(IN1, LOW);
        digitalWrite(IN2, LOW);
    }
}

void Motor::forward(int speed)
{
    Serial.println("Forwarding");
    analogWrite(ENA, speed);
    digitalWrite(IN1, HIGH);
    digitalWrite(IN2, LOW);
}

void Motor::backward(int speed)
{
    Serial.println("Backwarding");
    analogWrite(ENA, speed);
    digitalWrite(IN1, LOW);
    digitalWrite(IN2, HIGH);
}

void Motor::accelerate(int minSpeed, int maxSpeed, int direction)
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