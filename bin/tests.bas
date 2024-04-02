' BAS-1: Среднее арифметическое последовательности {0, 2, 4, 6, ...}
00 PRINT "BAS-1: Среднее арифметическое последовательности {0, 2, 4, 6, ...}"
10 LET S = 0
20 LET C = 0
30 FOR I = 10 TO 20
40     PRINT 2 * (I - 1) ;
50     LET S = S + (2 * (I - 1))
60 NEXT
70 PRINT
80 LET A = S / 11
90 PRINT "Среднее арифметическое:"; A

' BAS-2: 10 случайных чисел от 0 до 9 и их среднее арифметическое
100 PRINT "BAS-2: 10 случайных чисел от 0 до 9 и их среднее арифметическое"
110 LET S = 0
120 FOR I = 1 TO 10
130     LET N = RND(10)
140     LET S = S + N
150     PRINT N;
160 NEXT
170 PRINT
180 LET A = S / 10
190 PRINT "Среднее арифметическое:"; A

' BAS-3: Сумма последовательности {2, 5, 8, ...}
200 PRINT "BAS-3"
210 LET S = 0
220 FOR I = 10 TO 20
230     LET N = (I - 1) * 3 + 2
240     PRINT N ;
250     LET S = S + N
260 NEXT
270 PRINT
280 PRINT "Сумма элементов:"; S

' BAS-4: Сумма элементов последовательности {100, 97, 94, ..., 1
300 PRINT "BAS-4: Сумма элементов последовательности {100, 97, 94, ..., 1}"
310 LET S = 0
320 FOR i = 10 TO 20
330     LET N = 100 - (I - 1) * 3
340     PRINT N ;
350     LET S = S + N
360 NEXT
370 PRINT "Сумма элементов:"; S

' ' BAS-5: Элементы и сумма отрицательных элементов последовательности
400 PRINT "BAS-5: Элементы и сумма отрицательных элементов последовательности"
410 LET A = 0
420 LET B = 4
430 LET S = 0
440 LET N = 10
450 FOR I = 3 TO N
460    LET N = B - 3*A
470    LET A = B
480    LET B = N
490    PRINT N;
500    IF N < 0 THEN LET S = S + N
510 NEXT
520 PRINT "Сумма отрицательных элементов: "; S

' BAS-6: N первых чисел Фибоначчи
600 PRINT "BAS-6: N первых чисел Фибоначчи"
610 LET A = 0
620 LET B = 1
630 LET N = 0
640 PRINT "Введите число N: "
650 INPUT N
660 PRINT B;
670 FOR I = 2 TO N
680   LET C = A + B
690   PRINT C;
710   LET A = B
720   LET B = C
730 NEXT
740 PRINT

' BAS-7: Факториал числа N
800 PRINT "BAS-7: Факториал числа N"
810 LET N = 0
820 PRINT "Введите число N: "
830 INPUT N
840 LET F = 1
850 FOR I = 1 TO N
860  LET F = F * I
870 NEXT
880 PRINT "Факториал: "; F

' BAS-8: Возведение числа N в степень M
900 PRINT "BAS-8: Возведение числа N в степень M"
910 LET N = 0
920 LET M = 0
930 PRINT "Введите N и M: "
940 INPUT N, M
950 LET S = N
960 FOR I = 1 TO M
970     LET S = S * N
980 NEXT
990 PRINT N; " в степени "; M; " = "; S