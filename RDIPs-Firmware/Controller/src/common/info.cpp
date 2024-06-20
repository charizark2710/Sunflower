#include "info.h"
#include "../src/ArduinoUniqueID.h"
ArduinoUniqueID uniqueID;

char* getNanoID()
{   
    const size_t uniqueIDSize = sizeof(uniqueID.id);  // Size of the UniqueID array
    
    // Allocate a buffer to hold the hex string: 2 chars per byte + 1 for null terminator
    static char hexString[uniqueIDSize * 2 + 1];
    memset(hexString, 0, sizeof(hexString));  // Initialize the buffer to zero
    
    // Convert each byte of the uniqueID array to hex and store in the buffer
    for (size_t i = 0; i < uniqueIDSize; ++i) {
        // Use snprintf to convert each byte to hex and append to the string
        snprintf(hexString + (i * 2), 3, "%02X", UniqueID[i]);
    }
    
    return hexString;  // Return the static buffer
}