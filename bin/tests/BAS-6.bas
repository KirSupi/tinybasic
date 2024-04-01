' BAS-6: N первых чисел Фибоначчи
390 LET a = 0
400 LET b = 1
410 PRINT a;
420 FOR i = 2 TO N
430   LET c = a + b
440   PRINT c;
450   LET a = b
460   LET b = c
470 NEXT i