# golang_api2

Basado en la API: https://github.com/diegochavezcarro/golang_api y https://github.com/ydhnwb/golang_api

Pasos a seguir:

Levantar un MySQL(en una imagen de Docker):

docker run -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=password -e MYSQL_DATABASE=golang_api -d mysql:5.7.25

Ejecutar la API:

go run server.go

Entrar al mysql:

docker exec -it mysql bash

mysql -ppassword

use golang_api;

Debido al Automigrate configurado en gorm,se crean las tablas al inicio: 

show tables;

Al inicio estaran vacias:

select * from courses;

Luego registramos un usuario:

tiene validaciones, entre ellas, requiere un mail valido:

Utilizar POSTMAN y ejecutar en el primer request:

'{ "name":"diego", "email":"bla@bla.com", "password":"123" }' http://localhost:8080/api/auth/register

Luego de este register podemos usar el mismo request, modificando los datos para registrar nuestro usuario.

Desde Mysql, veremos los usuarios creados:

select * from users;

Para ver cursos disponibles:

Se debera realizar el request mediante el metodo POST utilizando POSTMAN a la siguiente URL:

http://localhost:8080/api/courses

Utilizando basic Autenticatión (debemos usar un usuario creado previamente)

