# AquaHub

Before run app need up database:

migrate -path schema/. -database postgres://postgres:postgres@172.19.0.2:5432/postgres?sslmode=require up

Run app, and browse to http://localhost:8000/swagger/index.html. You will see Swagger 2.0 Api documents as shown below:

![swagger-image](../main/assets/swagger-image.jpg)



