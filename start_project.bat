@echo off
:: Step 5: Start RabbitMQ server
:: Step 4: Kill Erlang processes (if needed)
echo Terminating Erlang processes (if any)...
taskkill /F /IM erl.exe /T > NUL 2>&1
taskkill /F /IM beam.smp.exe /T > NUL 2>&1

echo Starting RabbitMQ server...
start /B cmd /c "cd /d C:\Program Files\RabbitMQ Server\rabbitmq_server-4.0.2\sbin && rabbitmq-server > NUL 2>&1"

:: Step 1: Start Go backend
echo Starting Go backend...
start /B cmd /c "cd /d C:\Users\J.A.R.V.I.S\Desktop\dvm\back-end && go run main.go > NUL 2>&1"

:: Step 2: Start Frontend (Angular)
echo Starting frontend...
start /B cmd /c "cd /d C:\Users\J.A.R.V.I.S\Desktop\dvm\front-end && npm start > NUL 2>&1"

:: Step 3: Start Django server
echo Starting Django server...
start /B cmd /c "cd /d C:\Users\J.A.R.V.I.S\Desktop\dvm && python manage.py runserver > NUL 2>&1"


:: Wait for 8 seconds before starting consumer
timeout /t 8 /nobreak > NUL

:: Step 6: Start Go consumer
echo Starting Go consumer...
start /B cmd /c "cd /d C:\Users\J.A.R.V.I.S\Desktop\dvm\back-end\consumer && go run consumer.go > NUL 2>&1"

echo All services started successfully!
pause
