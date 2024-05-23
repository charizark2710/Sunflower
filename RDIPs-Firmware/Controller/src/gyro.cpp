#include <Arduino.h>
#include <Wire.h> // Wire library - used for I2C communication
#include <math.h> // Call out Pi
#include <gyro.h>

struct Gyro
{
    int ADXL345 = 0x53;        // The ADXL345 sensor I2C address
    float X_out, Y_out, Z_out; // Outputs
    float toDegree(float radians)
    {
        return radians * (180 / M_PI);
    }

    void getAngle()
    {
        // === Read acceleromter data === //
        Wire.beginTransmission(ADXL345);
        Wire.write(0x32); // Start with register 0x32 (ACCEL_XOUT_H)
        Wire.endTransmission(false);
        Wire.requestFrom(ADXL345, 6, true);       // Read 6 registers total, each axis value is stored in 2 registers
        X_out = (Wire.read() | Wire.read() << 8); // X-axis value
        X_out = X_out / 256;                      // For a range of +-2g, we need to divide the raw values by 256, according to the datasheet
        Y_out = (Wire.read() | Wire.read() << 8); // Y-axis value
        Y_out = Y_out / 256;
        Z_out = (Wire.read() | Wire.read() << 8); // Z-axis value
        Z_out = Z_out / 256;

        float X_degree = toDegree(X_out);
        float Y_degree = toDegree(Y_out);
        float Z_degree = toDegree(Z_out);

        Serial.print("Xa= ");
        Serial.print(X_out);
        Serial.print("   Xd= ");
        Serial.print(X_degree);
        Serial.print("   Ya= ");
        Serial.print(Y_out);
        Serial.print("   Yd= ");
        Serial.print(Y_degree);
        Serial.print("   Za= ");
        Serial.print(Z_out);
        Serial.print("   Zd= ");
        Serial.println(Z_degree);
        delay(500);
    }
};