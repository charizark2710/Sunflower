#include "correlationId.h"

String generateCorrelationId(String prefix)
{
  char *correlationId;
  uuid.setVariant4Mode();
  uuid.generate();

  strcpy(correlationId, prefix.c_str());
  strcpy(correlationId, "-");
  strcpy(correlationId, uuid.toCharArray());

  correlationIds[correlationIdsSize] = correlationId;
  correlationIdsSize++;

  Serial.print("correlationIdsSize: ");
  Serial.println(correlationIdsSize);

  return correlationId;
}

// TODO: implement remove correlationId from correlationIds array
void removeMatchedCorrelationId(String correlationId) {

}