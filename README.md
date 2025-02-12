### Как запустит приложение?
1) Создайте в корневой директории файл .env
**Пример файла .env**:
```
CFG_PATH=./configs/cfg.yaml 
STORAGE_TYPE=postgres
DB_HOST=db
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
DB_NAME=appdb
DB_SSL=disable
REDIS_ADDR=redis:6379
```
2) Запустите контейнеры через терминал
```sudo docker-compose up --build```