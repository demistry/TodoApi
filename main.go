package main

import (
	"github.com/Demistry/ToDoAPi/handler"
	"github.com/gin-gonic/gin"
	"path"
	"path/filepath"
)


func main(){
	g := gin.Default()

	g.NoRoute(func(c * gin.Context) {
		dir, file := path.Split(c.Request.RequestURI)
		ext := filepath.Ext(file)
		if file == "" || ext == "" {
			c.File("./ui/dist/ui/index.html")
		} else {
			c.File("./ui/dist/ui/" + path.Join(dir, file))
		}
	})
	g.GET("/todo", handler.GetTodohandler)
	g.POST("/todo", handler.AddTodoHandler)
	g.DELETE("/todo/:id", handler.DeleteTodoHandler)
	g.PUT("/todo", handler.CompleteTodoHandler)

	err := g.Run(":8080")
	if err != nil{
		panic(err)
	}
}
