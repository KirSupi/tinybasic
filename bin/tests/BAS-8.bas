' BAS-8: Возведение числа N в степень M
540 INPUT "Введите N и M: ", N, M
550 R = 1
560 FOR I = 1 TO M
570   R = R * N
580 NEXT I
590 PRINT N; " в степени "; M; " = "; R