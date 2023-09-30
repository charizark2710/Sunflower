#include <UUID.h>

extern UUID uuid;

extern String correlationIds[5];
extern int correlationIdsSize;

String generateCorrelationId(String prefix);

void removeMatchedCorrelationId(String correlationId);