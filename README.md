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

## Нагрузка
```
ghz --insecure --proto=api/shortener/shortener.proto --call=shortener.URLShortener.Create -D create_requests.json -n 2000 -c 20 -r 200 0.0.0.0:50051
```
```
ghz --insecure --proto=api/shortener/shortener.proto --call=shortener.URLShortener.Get -D get_requests.json -n 2000 -c 20 -r 200 0.0.0.0:50051
```

