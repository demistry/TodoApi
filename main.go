package main

import (
	"github.com/Demistry/ToDoAPi/handler"
	"github.com/auth0-community/go-auth0"
	"github.com/gin-gonic/gin"
	jose "gopkg.in/square/go-jose.v2"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

var (
	audience string
	domain   string
)

func main() {
	setAuth0Variables()
	g := gin.Default()

	g.NoRoute(func(c *gin.Context) {
		dir, file := path.Split(c.Request.RequestURI)
		ext := filepath.Ext(file)
		if file == "" || ext == "" {
			c.File("./ui/dist/ui/index.html")
		} else {
			c.File("./ui/dist/ui/" + path.Join(dir, file))
		}
	})

	authorized := g.Group("/")
	authorized.Use(authRequired())
	authorized.GET("/todo", handler.GetTodohandler)
	authorized.POST("/todo", handler.AddTodoHandler)
	authorized.DELETE("/todo/:id", handler.DeleteTodoHandler)
	authorized.PUT("/todo", handler.CompleteTodoHandler)

	err := g.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func setAuth0Variables() {
	audience = os.Getenv("AUTH0_API_IDENTIFIER")
	domain = os.Getenv("AUTH0_DOMAIN")
}

func authRequired() gin.HandlerFunc { //middleware
	return func(c *gin.Context) {
		var auth0Domain = "https://" + domain + "/"
		client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: auth0Domain + ".well-known/jwks.json"}, nil)
		configuration := auth0.NewConfiguration(client, []string{audience}, auth0Domain, jose.RS256)
		validator := auth0.NewValidator(configuration, nil)

		_, err := validator.ValidateRequest(c.Request)

		if err != nil {
			log.Println(err)
			terminateWithError(http.StatusUnauthorized, "token is not valid", c)
			return
		}
		c.Next()
	}
}

func terminateWithError(statusCode int, message string, c *gin.Context) {
	c.JSON(statusCode, gin.H{"error": message})
	c.Abort()
}
