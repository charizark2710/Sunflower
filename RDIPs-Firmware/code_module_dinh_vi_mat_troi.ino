#define LDR1 2
#define LDR2 3
#define LDR3 4
#define LDR4 5
#define LDR5 6

int R1[5], R2[5], R3[5], R4[5], R5[5], C[5];

void setup() {
  pinMode(A1, INPUT);
  pinMode(A2, INPUT);
  pinMode(A3, INPUT);
  pinMode(A4, INPUT);
  pinMode(A5, INPUT);
  pinMode(LDR1, OUTPUT);
  pinMode(LDR2, OUTPUT);
  pinMode(LDR3, OUTPUT);
  pinMode(LDR4, OUTPUT);
  pinMode(LDR5, OUTPUT);

  Serial.begin(9600);
}

void loop() {
  readSensorValues();
  findMaxValues();
  printMaxPositions();
  delay(1000);
}

void readSensorValues() {
  readSensorValuesForArray(R1, LDR1);
  readSensorValuesForArray(R2, LDR2);
  readSensorValuesForArray(R3, LDR3);
  readSensorValuesForArray1(R4, LDR4);
  readSensorValuesForArray1(R5, LDR5);

}

void readSensorValuesForArray(int array[], int pin) {
  digitalWrite(pin, HIGH);
  for (int i = 0; i < 3; i++) {
    array[i] = analogRead(A1 + i);  
  }
  digitalWrite(pin, LOW);
}

void readSensorValuesForArray1(int array1[], int pin1)
{
  digitalWrite(pin1, HIGH);
  for (int i = 0; i<2; i++){
    array1[i] = analogRead(A4 + i);
  }
  digitalWrite(pin1, LOW);
}



void findMaxValues() {
  for (int i = 0; i < 5; i++) {
    C[i] = findMaxValueInArray(getRowByIndex(i));
  }
}

int* getRowByIndex(int index) {
  switch (index) {
    case 0:
      return R1;
    case 1:
      return R2;
    case 2:
      return R3;
    case 3:
      return R4;
    case 4:
      return R5;
    default:
      return NULL;
  }
}

int findMaxValueInArray(int array[]) {
  int maxValue = array[0];
  for (int i = 1; i < 5; i++) {
    if (array[i] > maxValue) {
      maxValue = array[i];
    }
  }
  return maxValue;
}

void printMaxPositions() {
  int maxC = findMaxValueInArray(C);
  for (int i = 0; i < 5; i++) {
    if (C[i] == maxC) {
      Serial.print("Horizontal Position: ");
      Serial.println(i + 1);
      printVerticalPosition(getRowByIndex(i), maxC);
    }
  }
}

void printVerticalPosition(int array[], int maxValue) {
  for (int i = 0; i < 5; i++) {
    if (array[i] == maxValue) {
      Serial.print("Vertical Position: ");
      Serial.println(i + 1);
    }
  }
}
