#include <UUID.h>

extern UUID uuid;

String generateCorrelationId(String prefix);

int getCorrelationIndex();

void removeMatchedCorrelationId(String correlationId);