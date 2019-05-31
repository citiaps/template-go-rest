# template-go-rest
Repositorio base para la creacion de un servicio REST en GO.

El REST tiene:
 * Base de datos en mongo
 * CORS
 * Carga de multiples archivos .env dependiendo del ambiente 
 * Uso de JWT
 * Acceso restringido a servicios con roles
 * Almacenamiento de passwords hasheados con salt dinamico (bcrypt)

## Ejecución

Antes de ejecutar la aplicación, es necesario:
1) Crear base de datos: Para ello es necesario introducir los comandos definidos en [scripts/createdbuser.js](scripts/createdbuser.js).
2) Definir los parámetros en el archivo .env: Copiar [.env.example](.env.example), renombrar a .env y cambiar los valores.

Para correr el proyecto, clonar en cualquier directorio y utilizar el comando:

```bash
go run app.go
```
La descarga de dependencias se realiza de forma automática.

***NOTA:***
Para esto es necesario tener instalado [golang v1.11 o mayor](https://golang.org/doc/install).

## Estructura del proyecto

El proyecto está estructurado de la siguiente forma:

```bash
├── README.md
├── app.go
├── controller
│   ├── authentication.go
│   └── dog.go
├── go.mod
├── go.sum
├── middleware
│   ├── authentication.go
│   └── cors.go
├── model
│   ├── database.go
│   ├── dog.go
│   ├── index.go
│   └── user.go
├── postman
│   ├── Template-GO-REST.postman_collection.json
│   └── envs
│       └── TemplateGoREST\ -\ Development.postman_environment.json
│       └── TemplateGoREST\ -\ Production.postman_environment.json
│       └── TemplateGoREST\ -\ Testing.postman_environment.json
├── scripts
│   └── createdbuser.js
└── util
    ├── env.go
    ├── error.go
    └── response.go
```

#### Controller
Tiene los controladores de cada uno de los recursos del sistema. En este solo se debería manejar la lógica, recepción y respuesta de las consultas.

#### Model
Tiene los modelos de cada uno de los recursos del sistema. En este solo se debería implementar operaciones sobre la base de datos.

Además se tiene la conexión a la base de datos y la creación de indices para cada una de las colecciones.

#### Middleware
Tiene middlewares definidos para la autentificación del usuario.

#### Postman
Documentación del sistema, se tienen archivos json importables a [Postman](https://www.getpostman.com/), se guardan las colecciones y entornos de desarrollo, testing, producción.

#### scripts
Scripts necesarios para la inicialización de la base de datos.

#### util
Paquete que tiene funciones utilizables a lo largo de toda la aplicación.

## Roles y Auth

Los roles son definidos en [controller/authentication.go](controller/authentication.go), como constantes.

Para exigir autenticación en alguna ruta, utilizar el middleware correspondiente ```authNormal.MiddlewareFunc()```:

```golang
dogRouter.POST("", authNormal.MiddlewareFunc(), dogController.Create())

```

Para proteger una ruta dependiendo de rol del usuario utilizar el middleware antes de la autorización.

```golang
dogRouter.DELETE("/:id", middleware.SetRoles(RolAdmin), authNormal.MiddlewareFunc(), dogController.Delete())
```

Si se quiere permitir a mas de un rol, pasar como otro parametro:

```golang
dogRouter.DELETE("/:id", middleware.SetRoles(RolAdmin, RolUser), authNormal.MiddlewareFunc(), dogController.Delete())
```

