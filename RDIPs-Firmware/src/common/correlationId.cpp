#include "correlationId.h"

String correlationIds[CORRELATION_ID_MAX_SIZE];
unsigned int correlationIndex = 0;

void setCorrelationIds()
{
  correlationIds[0] = "";
  correlationIds[1] = "";
  correlationIds[2] = "";
  correlationIds[3] = "";
  correlationIds[4] = "";
}

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

void increaseCorrelationIndex()
{
  correlationIndex++;
}

void decreaseCorrelationIndex()
{
  correlationIndex--;
}

void removeMatchedCorrelationId(String correlationId)
{
  int isFounded = 0;
  int index = 0;

  for (int i = 0; i <= correlationIndex; i++)
  {
    if (strcmp(correlationId.c_str(), correlationIds[i].c_str()) == 0)
    {
      isFounded = 1;
      index = i;
      break;
    }
  }

  if (isFounded)
  {
    for (int i = index; i <= correlationIndex; i++)
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
    decreaseCorrelationIndex();
  }
}