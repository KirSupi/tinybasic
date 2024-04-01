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