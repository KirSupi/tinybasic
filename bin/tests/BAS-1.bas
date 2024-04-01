' BAS-1: Среднее арифметическое последовательности {0, 2, 4, 6, ...}
10 LET sum = 0
20 FOR i = 10 TO 20
30     LET sum = sum + (i * 2)
40 NEXT i
50 LET average = sum / 11
60 PRINT "Среднее арифметическое:"; average