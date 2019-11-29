package handler


import (
	"encoding/json"
	"github.com/Demistry/ToDoAPi/todo"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
)



func GetTodohandler(c *gin.Context){
	c.JSON(http.StatusOK, todo.GetAllList())
}

func AddTodoHandler(c *gin.Context){
	todoItem, statusCode, err := convertHTTPBodyToTodo(c.Request.Body)
	if err != nil{
		c.JSON(statusCode, err)
		return
	}
	c.JSON(statusCode, gin.H{"id":todo.Add(todoItem.Description)})
}

func CompleteTodoHandler(c * gin.Context){
	todoItem, statusCode, err := convertHTTPBodyToTodo(c.Request.Body)
	if err != nil{
		c.JSON(statusCode, err)
		return
	}

	if todo.Complete(todoItem.ID) != nil{
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "")
}

func DeleteTodoHandler(c *gin.Context){
	todoID := c.Param("id")
	if err := todo.Delete(todoID); err != nil{
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "")
}

func convertHTTPBodyToTodo(httpBody io.ReadCloser) (todo.Todo, int, error) {
	body, err := ioutil.ReadAll(httpBody)
	if err != nil {
		return todo.Todo{}, http.StatusInternalServerError, err
	}
	defer httpBody.Close()
	return convertJSONBodyToTodo(body)
}

func convertJSONBodyToTodo(jsonBody []byte) (todo.Todo, int, error) {
	var todoItem todo.Todo
	err := json.Unmarshal(jsonBody, &todoItem)
	if err != nil {
		return todo.Todo{}, http.StatusBadRequest, err
	}
	return todoItem, http.StatusOK, nil
}


func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message": "not found"}`))
}

func home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "hello world, welcome to this rest api"}`))
}
