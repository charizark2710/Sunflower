// R=row, C=column
int R1[5], R2[5], R3[5], R4[5], R5[5], C[5];


#define LDR 2
#define LDR2 3
#define LDR3 4
#define LDR4 5
#define LDR5 6

int last;

int maxC;

void MaxROW1();
void MaxROW2();
void MaxROW3();
void MaxROW4();
void MaxROW5();

void MaxCOL();

void locateSun1();
void locateSun2();
void locateSun3();
void locateSun4();
void locateSun5();

void setup() {

  pinMode(A1, INPUT);
  pinMode(A2, INPUT);
  pinMode(A3, INPUT);
  pinMode(A4, INPUT);
  pinMode(A5, INPUT);
  pinMode(LDR, OUTPUT);
  pinMode(LDR2, OUTPUT);
  pinMode(LDR3, OUTPUT);
  pinMode(LDR4, OUTPUT);
  pinMode(LDR5, OUTPUT);
  digitalWrite(LDR, LOW);  //mac dinh quang tro tat
  digitalWrite(LDR2, LOW);
  digitalWrite(LDR3, LOW);
  digitalWrite(LDR4, LOW);
  digitalWrite(LDR5, LOW);
  Serial.begin(9600);  //khoi tao baud rate
}

void loop() {

  for (int i = 0; i <= 4; i++) {
    R1[i] = analogRead(A1);
    R2[i] = analogRead(A2);
    R3[i] = analogRead(A3);
    R4[i] = analogRead(A4);
    R5[i] = analogRead(A5);
  }
  if (millis() - last >= 1000) {
    MaxROW1(R1[0]);
    MaxROW2(R2[0]);
    MaxROW3(R3[0]);
    MaxROW4(R4[0]);
    MaxROW5(R5[0]);

    MaxCOL(C[0]);

    locateSun1(C[0]);
    locateSun2(C[1]);
    locateSun3(C[2]);
    locateSun4(C[3]);
    locateSun5(C[4]);
  }
}

//tim vi tri hang ngang cua quang tro
void MaxROW1(int R1[0]) {
  if (millis() - last >= 1000) {
    int maxR1 = R1[0];
    for (int i = 1; i <= 4; i++) {
      if (maxR1 < R1[i]) {
        maxR1 = R1[i];
      }
    }
    C[0] = maxR1;
  }
}

void MaxROW2(int R2[0]) {
  if (millis() - last >= 1000) {
    int maxR2 = R2[0];
    for (int i = 1; i <= 4; i++) {
      if (maxR2 < R2[i]) {
        maxR2 = R2[i];
      }
    }
    C[1] = maxR2;
  }
}

void MaxROW3(int R3[0]) {
  if (millis() - last >= 1000) {
    int maxR3 = R3[0];
    for (int i = 1; i <= 4; i++) {
      if (maxR3 < R3[i]) {
        maxR3 = R3[i];
      }
    }
    C[2] = maxR3;
  }
}

void MaxROW4(int R4[0]) {
  if (millis() - last >= 1000) {
    int maxR4 = R4[0];
    for (int i = 1; i <= 4; i++) {
      if (maxR4 < R4[i]) {
        maxR4 = R4[i];
      }
    }
    C[3] = maxR4;
  }
}

void MaxROW5(int R5[0]) {
  if (millis() - last >= 1000) {
    int maxR5 = R5[0];
    for (int i = 1; i <= 4; i++) {
      if (maxR5 < R5[i]) {
        maxR5 = R5[i];
      }
    }
    C[4] = maxR5;
  }
}

//xac dinh vi tri hang doc
void MaxCOL(int C[0]) {
  if (millis() - last >= 1000) {
    int maxC = C[0];
    for (int i = 1; i <= 4; i++) {
      if (maxC < C[i]) {
        maxC = C[i];
      }
    }
  }
}

//tim vi tri mat troi
void locateSun1(int C[0]) {
  if (millis() - last >= 1000) {
    maxC = C[0];
    if (maxC == C[0]) {
      for (int i = 0; i <= 4; i++) {
        if (maxC == R1[i]) {
          Serial.println("Vertical Position: 1");
          Serial.println("Horizontal Position: ");
          Serial.println(i + 1);
        }
      }
    }
  }
}

void locateSun2(int C[1]) {
  if (millis() - last >= 1000) {
    if (maxC == C[1]) {
      for (int i = 0; i <= 4; i++) {
        if (maxC == R2[i]) {
          Serial.println("Vertical Position: 2");
          Serial.println("Horizontal Position: ");
          Serial.println(i + 1);
        }
      }
    }
  }
}

void locateSun3(int C[2]) {
  if (millis() - last >= 1000) {
    if (maxC == C[2]) {
      for (int i = 0; i <= 4; i++) {
        if (maxC == R3[i]) {
          Serial.println("Vertical Position: 3");
          Serial.println("Horizontal Position: ");
          Serial.println(i + 1);
        }
      }
    }
  }
}

void locateSun4(int C[3]) {
  if (millis() - last >= 1000) {
    if (maxC == C[3]) {
      for (int i = 0; i <= 4; i++) {
        if (maxC == R4[i]) {
          Serial.println("Vertical Position: 4");
          Serial.println("Horizontal Position: ");
          Serial.println(i + 1);
        }
      }
    }
  }
}

void locateSun5(int C[4]) {
  if (millis() - last >= 1000) {
    if (maxC == C[4]) {
      for (int i = 0; i <= 4; i++) {
        if (maxC == R5[i]) {
          Serial.println("Vertical Position: 5");
          Serial.println("Horizontal Position: ");
          Serial.println(i + 1);
        }
      }
    }
  }
}
