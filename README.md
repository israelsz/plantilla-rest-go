# template-go-rest2
Repositorio base para la creacion de un servicio REST en GO. Basado en template-go-rest

El REST tiene:
 * Base de datos en mongoDB
 * CORS
 * Carga de multiples archivos .env dependiendo del ambiente 
 * Uso de JWT en cookies
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

## Requisitos
* Golang v1.13 o superior.
* MongoDB 3.6 o superior.

## Estructura del proyecto

El proyecto está estructurado de la siguiente forma:

```bash
├── README.md
├── main.go
├── go.mod
├── go.sum
├── controller
│   ├── catController.go
│   └── userController.go
├── middleware
│   ├── authentication.go
│   └── cors.go
├── models
│   ├── cat.go
│   ├── login.go
│   └── user.go
├── routes
│   ├── auth.go
│   ├── routes.go
│   ├── catRouter.go
│   └── userRouter.go
├── scripts
│   └── createdbuser.js
├── services
│   ├── catService.go
│   ├── userServie.go
└── utils
    ├── env.go
    ├── logger.go
    └── hash.go
```

#### Controller
Tiene los controladores de cada uno de los recursos del sistema. En este solo se debería manejar la lógica, recepción y respuesta de las consultas.

#### Model
Tiene los modelos de cada uno de los recursos del sistema. En este solo se debería implementar estructuras de las entidades a utilizar.

#### Routes
Tiene las rutas de los end-points de la aplicación, divididas según entidad. Las rutas son registradas en la aplicación en el archivo routes.go

#### Middleware
Tiene middlewares definidos para la autentificación del usuario. Como la implementación de JWT y la configuración de CORS.

#### Scripts
Scripts necesarios para la inicialización de la base de datos.

#### Services
Tiene la lógica de operaciones en la base de datos. En este solo se deberían implementar operaciones para usar con la base de datos

#### Utils
Paquete que tiene funciones utilizables a lo largo de toda la aplicación.

## Roles y Auth

Los roles son definidos en las rutas [middleware/authentication.go](middleware/authentication.go), como constantes al inicio del archivo.

Para exigir autenticación en alguna ruta, utilizar el middleware correspondiente ```middleware.LoadJWTAuth().MiddlewareFunc()```:


Por ejemplo, para proteger una ruta que solo pueda ser utilizada por un usuario logueado
```golang
catGroup.POST("/", middleware.LoadJWTAuth().MiddlewareFunc(), controller.CreateCat)
```

Para proteger una ruta dependiendo de rol del usuario utilizar el middleware antes de la autorización.

```golang
catGroup.DELETE("/:id", middleware.SetRoles(middleware.RolAdmin), middleware.LoadJWTAuth().MiddlewareFunc(), controller.DeleteCat)
```

Si se quiere permitir a mas de un rol, se puede agregar como otro parametro:

```golang
catGroup.PUT("/:id", middleware.SetRoles(middleware.RolAdmin, middleware.RolUser), middleware.LoadJWTAuth().MiddlewareFunc(), controller.UpdateCat)
```

Hay más ejemplos disponibles al interior del directorio [routes](routes/)

## Para registrar un nuevo servicio
1. Crear modelos como struct en el directorio models
2. Crear controlador que reciba y gestione los parametros https
3. Crear servicio que se comunique con la base de datos, enlazarla al controlador
4. Crear ruta que ejecute la funcionalidad creada en el controlador, agregar restricciones de autenticación y roles si se desea