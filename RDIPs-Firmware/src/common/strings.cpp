#include "strings.h"

JSONVar parseStringToJson(String str){
    Serial.println("parseStringToJson");
    JSONVar response;
    return JSON.parse(str);
}
