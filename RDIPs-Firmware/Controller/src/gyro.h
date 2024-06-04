#ifndef GYRO_H
#define GYRO_H
#include <Arduino.h>
#include <Wire.h>
#include <math.h>
typedef struct Gyro
{
    int ADXL345 = 0x53;                 // The ADXL345 sensor I2C address
    float X_degree, Y_degree, Z_degree; // Output

    enum MODEL
    {
        HORIZON = 0,
        SPHERE = 1
    };

    // Return angle by model
    // If model is horizon, input 0, return X-axis only
    // If model is sphere, input 1, return X-axis and Y-axis
    void displayAngleByModel(int model)
    {
        switch (model)
        {
        case HORIZON:
            Serial.print("X = ");
            Serial.println(X_degree);
            break;
        case SPHERE:
            Serial.print("X= ");
            Serial.print(X_degree);
            Serial.print("  Y= ");
            Serial.println(Y_degree);
        default:
            break;
        }
    }

    // Convert radians to degrees
    float toDegree(float radians)
    {
        return radians * (180 / M_PI);
    }

    // === Read acceleromter data === //
    void getAngle()
    {
        float X_radians, Y_radians, Z_radians; // Outputs
        Wire.beginTransmission(ADXL345);
        Wire.write(0x32); // Start with register 0x32 (ACCEL_XOUT_H)
        Wire.endTransmission(false);
        Wire.requestFrom(ADXL345, 6, true);           // Read 6 registers total, each axis value is stored in 2 registers
        X_radians = (Wire.read() | Wire.read() << 8); // X-axis value
        X_radians = X_radians / 256;                  // For a range of +-2g, we need to divide the raw values by 256, according to the datasheet
        Y_radians = (Wire.read() | Wire.read() << 8); // Y-axis value
        Y_radians = Y_radians / 256;
        Z_radians = (Wire.read() | Wire.read() << 8); // Z-axis value
        Z_radians = Z_radians / 256;

        X_degree = round(toDegree(X_radians));
        Y_degree = round(toDegree(Y_radians));
        Z_degree = round(toDegree(Z_radians));

        delay(500);
    }
} Gyro;

#endif