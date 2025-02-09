# URL-shortener

## Запуск с postgreSQL
1. Указать в .env значение STORAGE_TYPE=postgres
2. Запуск
```
docker compose --profile postgres up
```

## Запуск с in-memory хранилищем
1. Указать в .env значение STORAGE_TYPE=memory 
2. Запуск
```
docker compose up
```
