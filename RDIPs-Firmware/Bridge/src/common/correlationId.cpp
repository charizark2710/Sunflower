#include "correlationId.h"

String correlationIds[CORRELATION_ID_MAX_SIZE] = {"", "", "", "", ""};
unsigned int correlationIndex = 0;

String generateCorrelationId(String prefix)
{
  String correlationId = "";
  uuid.setVariant4Mode();
  uuid.generate();

  correlationId += prefix;
  correlationId += "-";
  correlationId += uuid.toCharArray();

  correlationIds[correlationIndex] = correlationId;
  ++correlationIndex;

  return correlationId;
}

int getCorrelationIndex()
{
  return correlationIndex;
}

void removeMatchedCorrelationId(String correlationId)
{
  int isFounded = 0;

  for (int i = 0; i <= correlationIndex; i++)
  {
    if (strcmp(correlationId.c_str(), correlationIds[i].c_str()) == 0)
    {
      isFounded = 1;
    }

    if (isFounded)
    {
      if (i == correlationIndex && correlationIndex == CORRELATION_ID_MAX_SIZE - 1)
      {
        correlationIds[i] = "";
      }
      else
      {
        correlationIds[i] = correlationIds[i + 1];
      }
    }
  }

  --correlationIndex;
}