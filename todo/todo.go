package todo

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
)


var (
	once sync.Once
	mtx sync.RWMutex
	list []Todo
)


type Todo struct{
	ID string `json:"id"`
	Description string `json:"description"`
	Complete bool	`json:"complete"`
}

func init() {
	once.Do(initialiseList)
}

func initialiseList() {
	list = []Todo{}
}

func GetAllList() []Todo{
	return list
}

func Add(message string)string{
	t := newTodo(message)
	mtx.Lock()
	list = append(list, t)
	mtx.Unlock()
	return t.ID
}

func Delete(id string) error{
	index, err := findTodoLocation(id)
	if err != nil{
		fmt.Println(err)
		return err
	}
	mtx.Lock()
	list = append(list[:index], list[index + 1:]...)
	mtx.Unlock()
	return nil
}

func Complete(id string) error{
	location, err := findTodoLocation(id)
	if err != nil{
		return err
	}
	mtx.Lock()
	list[location].Complete = true
	mtx.Unlock()
	return nil
}

func newTodo(message string)Todo{
	return Todo{generateRandomString(),
					message,
					false}
}

func findTodoLocation(id string)(int, error){
	mtx.RLock()
	defer mtx.RUnlock()

	for i,todo := range list{
		if todo.ID == id{
			return i, nil
		}
	}
	errorStr := "Could not find todo with ID " + id
	return 0, errors.New(errorStr)
}

func generateRandomString() string{
	alphabets := []rune("1234567890qwertyuiopasdfghjklzxcvbnm")
	b := make([]rune, 10)
	for i := range b{
		b[i] = alphabets[rand.Intn(len(alphabets))]
	}
	return string(b)
}