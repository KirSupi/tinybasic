' BAS-1: Среднее арифметическое последовательности {0, 2, 4, 6, ...}
10 LET sum = 0
20 FOR i = 10 TO 20
30     LET sum = sum + (i * 2)
40 NEXT i
50 LET average = sum / 11
60 PRINT "Среднее арифметическое:"; average

' BAS-2: 10 случайных чисел от 0 до 9 и их среднее арифметическое
70 LET sum = 0
80 FOR i = 1 TO 10
90     LET num = INT(RND(1) * 10)
100    LET sum = sum + num
110    PRINT num;
120 NEXT i
130 PRINT
140 LET average = sum / 10
150 PRINT "Среднее арифметическое:"; average

' BAS-3: Сумма последовательности {2, 5, 8, ...}
160 LET sum = 0
170 FOR i = 10 TO 20
180     LET num = (i - 1) * 3 + 2
190     LET sum = sum + num
200 NEXT i
210 PRINT "Сумма элементов:"; sum

' BAS-4: Сумма элементов последовательности {100, 97, 94, ..., 1
220 LET sum = 0
230 FOR i = 10 TO 20
240     LET num = 100 - (i - 1) * 3
250     LET sum = sum + num
260 NEXT i
270 PRINT "Сумма элементов:"; sum

' BAS-5: Элементы и сумма отрицательных элементов последовательности
280 A = 0
290 B = 4
300 SUM = 0
310 FOR I = 3 TO 10
320  N = B - 3*A
330  A = B
340  B = N
350  PRINT N
360  IF N < 0 THEN SUM = SUM + N
370 NEXT I
380 PRINT "Сумма отрицательных элементов: "; SUM

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

' BAS-7: Факториал числа N
480 INPUT "Введите N: ", N
490 F = 1
500 FOR I = 1 TO N
510  F = F * I
520 NEXT I
530 PRINT "Факториал: "; F

' BAS-8: Возведение числа N в степень M
540 INPUT "Введите N и M: ", N, M
550 R = 1
560 FOR I = 1 TO M
570   R = R * N
580 NEXT I
590 PRINT N; " в степени "; M; " = "; R