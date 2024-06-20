#ifndef CHECK_PERIPHERAL_DEVICE_STATUS_H
#define CHECK_PERIPHERAL_DEVICE_STATUS_H
bool checkI2Cstatus(int sensorAdrress);
bool checkAnalogConnectionStatus(int analogPin);
bool checkDigitalConnectionStatuts(int digitalPin);

#endif //CHECK_PERIPHERAL_Device_H