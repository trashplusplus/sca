### 🌿 spy cat agency API (Golang + Gin + PostgreSQL)

проєкт реалізує RESTful API для керування котами, місіями, та цілями згідно технічного завдання. 

тестове завдання для DevelopsToday.


### 🚀запуск

```
git clone https://github.com/trashplusplus/sca.git
docker compose -f docker-compose.yml up -d (або docker-compose)
go mod init sca
go mod tidy
cd cmd
go run .

```



### 📝особливості:

- при кожному старті контейнера БД виконується скрипт init.sql, який створює необхідні таблиці (cats, missions, targets) та базові дані 🌱.  
- для підключення до БД використовуються змінні середовища з файлу .env 🌿
- запити та відповіді логуються в файл server.log



### 📬postman колекція

вся документація для тестування API доступна через Postman Collection у файлі sca_collection.json 📦  



## 🐱породи котів

- при створенні кота поле breed перевіряється через API:  
  [The Cat API](https://api.thecatapi.com/v1/breeds) 🌱  
- породи кешуються та потім перевіряються без додаткових запитів до API

## демонстрація успішного деплою (як приклад власний Raspberry Pi)

<img width="1288" height="188" alt="image" src="https://github.com/user-attachments/assets/14eb6ee1-a803-4034-9fc6-bacf6c64e0fa" />
<img width="1579" height="672" alt="image" src="https://github.com/user-attachments/assets/e460ef48-c7e9-49d9-b9f9-9e1f873dd765" />




