package middleware

import (
	"errors"
	"log"
	"net/http"
	"os"
	"rest-template/controller"
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
	//log.Println("DATOS DE USUARIO: ", userData)
	//log.Println("SOLO ROL: ", userData["rol"])

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
	usuario := models.User{Email: user.Email, Name: user.Name, ID: user.ID, Rol: user.Rol}
	if v, ok := data.(models.User); ok {
		claim := jwt.MapClaims{
			"user": usuario,
			"rol":  v.Rol,
		}
		log.Printf("DENTRO DEL PAYLOAD: %v", claim)
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
	log.Println("VALOR DE C: ", c)
	if err := c.ShouldBind(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	//email := loginVals.Email
	//password := loginVals.Password

	// Se establece conexion a la base de datos
	//Se busca la coleccion de usuarios y al usuario correspondiente
	var user models.User
	//Se trae al usuario buscandolo por su email
	log.Println("Contexto antes de getUser", c)
	controller.GetUserByEmail(c)
	log.Println("Contexto despues de getUser", c)
	//Se escriben los datos del usuario traido desde la base de datos al model
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("No fue posible encontrar al usuario")
		c.AbortWithError(http.StatusBadRequest, err)
		return nil, errors.New("usuario y contraseña incorrectos")
	}

	//Chequear credenciales del usuario
	if err := controller.ComparePasswords(user.Hash, loginVals.Password); err != nil {
		//return nil, jwt.ErrFailedAuthentication
		return nil, errors.New("contraseña incorrecta")
	}

	//Si usuario y contraseña son correctos
	return user, nil

	//return nil, jwt.ErrFailedAuthentication
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
		//TokenLookup: "header: Authorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,

		// Guardar token JWT como cookie en el navegador
		SendCookie:     true,
		SecureCookie:   false, //non HTTPS dev environments
		CookieHTTPOnly: true,  // JS can't modify
		//CookieDomain:   "localhost:8080", Se debe ingresar la URL del host
		CookieName:     "token", // default jwt
		TokenLookup:    "cookie:token",
		CookieSameSite: http.SameSiteDefaultMode, //SameSiteDefaultMode, SameSiteLaxMode, SameSiteStrictMode, SameSiteNoneMode
	})

	// Verificar si existen errores
	utils.Check(err)

	return authMiddleware

}
