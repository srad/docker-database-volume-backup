@echo off

docker stop wordpress-backup-enhanced-dev
docker rm wordpress-backup-enhanced-dev
docker build -t wordpress-backup-enhanced:dev .
if %errorlevel% neq 0 exit /b %errorlevel%
REM docker run --name wordpress-backup-enhanced-dev -p 8080:8080 -e MYSQL_HOST_FILE=mysql -it wordpress-backup-enhanced:dev
docker stack rm backup
docker stack deploy -c docker-compose.yml backup
docker service logs backup_wordpress-back-enhanced -f