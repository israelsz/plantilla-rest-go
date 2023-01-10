package middleware

import (
	"log"
	"os"
	"rest-template/models"
	"rest-template/utils"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// Roles en el sistema, para registrar nuevos roles, hacerlo aca
const (
	RolAdmin = "Admin"
	RolUser  = "User"
)

// AuthorizatorFunc : funcion tipo middleware que define si el usuario esta autorizado a utilizar la siguiente funcion
func AuthorizatorFunc(data interface{}, c *gin.Context) bool {

	/*
		user := data.(map[string]interface{})
		colUser, session := model.GetCollection(model.CollectionNameUser)
		defer session.Close()
		var usuario model.User

		if err := colUser.FindId(bson.ObjectIdHex(user["id"].(string))).One(&usuario); err != nil {
			return false
		}
		roles, exists := c.Get("roles")
		if !exists {
			return true
		}
		for _, r := range roles.([]string) {
			if usuario.Rol == r {
				return true
			}
		}
		return false
	*/
	userData := data.(map[string]interface{})
	log.Println(userData)
	log.Println("Correo: ", userData["email"])

	roles, exists := c.Get("roles")
	if !exists {
		return true
	}
	for _, r := range roles.([]string) {
		if userData["rol"] == r {
			return true
		}
	}
	return false

	//if datosUsuario, ok := data.(map[string]interface{}); ok && datosUsuario["email"] == "admin@a.com" {
	//	return true
	//}

	//return false
}

// UnauthorizedFunc : funcion que se llama en caso de no estar autorizado a accesar al servicio
func UnauthorizedFunc(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"message": message,
	})
}

// PayLoad : funcion que define lo que tendra el jwt que se enviara al realizarse el login
func PayLoad(data interface{}) jwt.MapClaims {
	user := data.(models.User)
	usuario := models.User{Email: user.Email, Name: user.Name, ID: user.ID}
	if v, ok := data.(models.User); ok {
		claim := jwt.MapClaims{
			"user": usuario,
			"rol":  v.Rol,
		}
		log.Printf("%v", claim)
		return claim
	}
	return jwt.MapClaims{}
}

func IdentityHandlerFunc(c *gin.Context) interface{} {
	jwtClaims := jwt.ExtractClaims(c)
	return jwtClaims["user"]
}

// Función que permite hacer login
func LoginFunc(c *gin.Context) (interface{}, error) {
	/*
		var loginVals models.Login
		//Se revisa el json entrante
		if err := c.BindJSON(&loginVals); err != nil {
			return "", jwt.ErrMissingLoginValues
		}
		// Se establece conexion a la base de datos
		//Se busca la coleccion de usuarios y al usuario correspondiente


			colUser, session := models.GetCollection(model.CollectionNameUser)
			defer session.Close()
			var usuario models.User

			if err := colUser.Find(bson.M{"email": loginVals.Email}).One(&usuario); err != nil {
				//return nil, jwt.ErrFailedAuthentication
				return nil, errors.New("Usuario y contraseña incorrectos")
			} else {
				if err := ComparePasswords(usuario.Hash, loginVals.Password); err != nil {
					//return nil, jwt.ErrFailedAuthentication
					return nil, errors.New("Usuario y contraseña incorrectos")
				}
				return usuario, nil
			}
	*/

	var loginVals models.Login
	if err := c.ShouldBind(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	email := loginVals.Email
	password := loginVals.Password

	if email == "admin@a.com" && password == "admin" {
		return models.User{
			Email: email,
			Name:  "ModoKernel",
			Rol:   RolAdmin,
		}, nil
	}

	if email == "user@a.com" && password == "user" {
		return models.User{
			Email: email,
			Name:  "ModoUsuario",
			Rol:   RolUser,
		}, nil
	}

	return nil, jwt.ErrFailedAuthentication
}

// SetRoles : funcion tipo middleware que define los roles que pueden realizar la siguiente funcion
// Se implementa sobre las rutas para definir que rol puede ocupar el servicio
func SetRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set example variable
		c.Set("roles", roles)
		// before request
		c.Next()
	}
}

func LoadJWTAuth() *jwt.GinJWTMiddleware {
	log.Print("LoadJWTAuth\n")
	var key string
	var set bool
	key, set = os.LookupEnv("JWT_KEY")
	if !set {
		key = "string_largo_unico_por_proyecto"
	}

	log.Println("key: " + key)

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm: "test zone",
		Key:   []byte(key),
		//tiempo que define cuanto vence el jwt
		Timeout: time.Hour * 24 * 7, //una semana
		//tiempo maximo para poder refrescar el jwt token
		MaxRefresh: time.Hour * 24 * 7,

		PayloadFunc:     PayLoad,
		IdentityHandler: IdentityHandlerFunc,
		Authenticator:   LoginFunc,
		Authorizator:    AuthorizatorFunc,
		Unauthorized:    UnauthorizedFunc,
		//HTTPStatusMessageFunc: ResponseFunc,
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	// Verificar si existen errores
	utils.Check(err)

	return authMiddleware

}
