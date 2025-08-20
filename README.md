## 🌿 spy cat agency API (Golang + Gin + PostgreSQL)

проєкт реалізує RESTful API для керування котами, місіями, та цілями згідно технічного завдання. 


## 🚀запуск

```
docker compose -f docker-compose.yml up -d
go mod init sca
go mod tidy
cd cmd
go run .

```



### 📝особливості:

- при кожному старті контейнера БД виконується скрипт init.sql, який створює необхідні таблиці (cats, missions, targets) та базові дані 🌱.  
- для підключення до БД використовуються змінні середовища з файлу .env 🌿
- запити та відповіді логуються в файл server.log



## 📬postman колекція

вся документація для тестування API доступна через Postman Collection у файлі sca_collection.json 📦  



## 🐱породи котів

- при створенні кота поле breed перевіряється через API:  
  [The Cat API](https://api.thecatapi.com/v1/breeds) 🌱  
- породи кешуються та потім перевіряються без додаткових запитів до API