package middleware

import (
	"errors"
	"log"
	"net/http"
	"os"
	"rest-template/models"
	"rest-template/services"
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

// AuthorizatorFunc : funcion tipo middleware que define si el usuario esta autorizado a utilizar un servicio
func AuthorizatorFunc(data interface{}, c *gin.Context) bool {

	//Se consiguen los datos entrantes a verificar
	userData := data.(map[string]interface{})

	// Se consiguen los roles registrados para la ruta a verificar
	roles, exists := c.Get("roles")
	if !exists {
		return true
	}
	for _, r := range roles.([]string) {
		//Si el usuario tienea algun rol vinculado a la ruta, se le permite su acceso a ella
		if userData["rol"] == r {
			return true
		}
	}
	// En caso contrario, se le deniega el permiso
	return false
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
	//Se fijan los campos que contendra el token jwt insertos
	usuario := models.User{Email: user.Email, Name: user.Name, ID: user.ID, Rol: user.Rol}
	if v, ok := data.(models.User); ok {
		claim := jwt.MapClaims{
			"user": usuario,
			"rol":  v.Rol,
		}
		return claim
	}
	return jwt.MapClaims{}
}

// Función que retorna las claims registradas en la función de Payload
func IdentityHandlerFunc(c *gin.Context) interface{} {
	jwtClaims := jwt.ExtractClaims(c)
	return jwtClaims["user"] //Retrona la claim registrada para usuario en payload.
}

// Función que permite hacer login en la aplicación y conseguir un token jwt
func LoginFunc(c *gin.Context) (interface{}, error) {
	var loginValues models.Login
	// Se asocian los valores entrantes por contexto al modelo de Login creado
	if err := c.ShouldBind(&loginValues); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	email := loginValues.Email
	password := loginValues.Password

	// Se establece conexion a la base de datos
	// y se trae al usuario buscandolo por su email
	user, err := services.GetUserByEmailService(email)
	// Si hubo algun error
	if err != nil {
		log.Println("No fue posible encontrar al usuario")
		c.AbortWithError(http.StatusBadRequest, err)
		return nil, errors.New("usuario y contraseña incorrectos")
	}

	//Chequear credenciales del usuario
	if err := utils.ComparePasswords(user.Hash, password); err != nil {
		//return nil, jwt.ErrFailedAuthentication
		return nil, errors.New("contraseña incorrecta")
	}

	//Retorna al usuario
	return user, nil
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

// Función que retorna una struct del middleware
func LoadJWTAuth() *jwt.GinJWTMiddleware {
	var key string
	var set bool
	//Se carga la key de jwt seteada desde las variables de entorno
	key, set = os.LookupEnv("JWT_KEY")
	if !set {
		//Si no estaba seteada, se fija una por default
		key = "string_largo_unico_por_proyecto"
	}
	//Se crea el middleware
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
	if err != nil {
		log.Println("Hubo un error al cargar el middleware")
	}

	return authMiddleware

}
