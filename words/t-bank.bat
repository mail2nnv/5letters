setLocal EnableExtensions EnableDelayedExpansion

chcp 65001

:: Инициализация переменных по умолчанию
@echo off
for /f "skip=1" %%a in ('wmic os get localdatetime') do (
    set D=%%a
    goto :break
)
: break
set NOW=%D:~0,4%-%D:~4,2%-%D:~6,2%_%D:~8,2%-%D:~10,2%
echo NOW: %NOW%

:: Переход к основной логике скрипта
goto main_logic

:: strlen - быстрая функция расчета длины строки в 13 итераций вместо итерирования по каждому символу
:strlen <resultVar> <stringVar>
(   
    setlocal EnableDelayedExpansion
    (set^ tmp=!%~2!)
    if defined tmp (
        set "len=1"
        for %%P in (4096 2048 1024 512 256 128 64 32 16 8 4 2 1) do (
            if "!tmp:~%%P,1!" NEQ "" ( 
                set /a "len+=%%P"
                set "tmp=!tmp:~%%P!"
            )
        )
    ) ELSE (
        set len=0
    )
)
( 
    endlocal
    set "%~1=%len%"
    exit /b
)


:: enterWord - Функция ввода слова с проверкой с трех попыток
:: В первом параметре принимает номер запрашиваемого слова. Во второй возвращает введеное слово
:enterWord <promptNum> <wordVar>
(
  setlocal EnableDelayedExpansion
  set word=
  set result=1
  for /l %%a in (1,1,3) do (
    set input=
    set /p input="Введите слово %~1 и результат проверки: "
    call :strlen len input
    if "!len!"=="0" (
      :: отказ от ввода, выход
      set result=2
      goto :breakEnterWord
    )
    if "!len!"=="11" (
      :: успешный ввод, выход
      set result=0
      set word=!input!
      goto :breakEnterWord
    )
    :: продолжаем 33 попытки ввода
    echo слово «!input!» длиной !len!, а ожидается 11 символов, 
    echo Например: «слово=20100», где «2»-буква на своем месте, «1»-не на своем, «0»-отсуствует
  )
)
: breakEnterWord
(
  endlocal
  set "%~2=%word%"
  exit /b %result%
)

:main_logic
(
  setlocal EnableDelayedExpansion
  set result=0

  for /l %%a in (1,1,5) do (
    call :enterWord #%%a, word
    if !ERRORLEVEL! neq 0 (
      set result=1
      goto :breakMain
    )
    set word[%%a]=!word!
    echo !word[1]! !word[2]! !word[3]! !word[4]! !word[5]! >> .tbank\%NOW%.txt
    words.exe !word[1]! !word[2]! !word[3]! !word[4]! !word[5]! >> .tbank\%NOW%.txt
    if !ERRORLEVEL! neq 0 (
      set result=1
      echo *** ОШИБКА *** >> .tbank\%NOW%.txt
      goto :breakMain
    )
  )
)
: breakMain
(
  endlocal
  exit /b %result%
)
