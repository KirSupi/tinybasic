' BAS-3: Сумма последовательности {2, 5, 8, ...}
160 LET sum = 0
170 FOR i = 10 TO 20
180     LET num = (i - 1) * 3 + 2
190     LET sum = sum + num
200 NEXT i
210 PRINT "Сумма элементов:"; sum
