#include <UUID.h>

extern UUID uuid;

void setCorrelationIds();

String generateCorrelationId(String prefix);

int getCorrelationIndex();

void increaseCorrelationIndex();

void decreaseCorrelationIndex();

void removeMatchedCorrelationId(String correlationId);