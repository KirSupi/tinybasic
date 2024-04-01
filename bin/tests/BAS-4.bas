' BAS-4: Сумма элементов последовательности {100, 97, 94, ..., 1
220 LET sum = 0
230 FOR i = 10 TO 20
240     LET num = 100 - (i - 1) * 3
250     LET sum = sum + num
260 NEXT i
270 PRINT "Сумма элементов:"; sum