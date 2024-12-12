@echo off

docker image prune -ya
docker build -t docker-database-volume-backup:dev .
if %errorlevel% neq 0 exit /b %errorlevel%
docker stack rm backup
docker stack deploy -c docker-compose.yml backup
docker service logs backup_docker-database-volume-backup -f
