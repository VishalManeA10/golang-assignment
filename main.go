package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Person struct {
	Name `json:"Name"`
	Age  `json:"Age"`
	Year `json:"Year"`
}

type Name string
type Age float64
type Year int

type Info interface {
	getInfo() string
}

func (p Name) getInfo() Name {
	return p
}
func (p Age) getInfo() Age {
	return p
}
func (p Year) getInfo() Year {
	return p
}

func funcDB(action string, index int) interface{} {
	personList := []Person{{"vishal", 24.5, 1996}, {"Amit", 25.5, 1992}, {"Manuwela", 26.5, 1991}, {"Anaya", 27.5, 1990}, {"Rahul", 28.5, 1990}}

	if index >= len(personList) {
		return -1
	}
	if action == "getlength" {
		return len(personList)
	} else {
		name := personList[index].Name.getInfo()
		age := personList[index].Age.getInfo()
		year := personList[index].Year.getInfo()
		return Person{Name: name, Age: age, Year: year}
	}
}

func main() {
	r := gin.Default()
	r.GET("/persons", mainAPI)
	r.Run("localhost:8080")
}

func mainAPI(ctx *gin.Context) {
	isDataAvailable := make(chan bool, 1)
	wg := &sync.WaitGroup{}
	i := 0

	funcDBLength := funcDB("getlength", 0)
	fmt.Println("funcDBLength : ", funcDBLength)
	wg.Add(1)

	for i < funcDBLength.(int)+1 {
		go Thread1(ctx, i, wg, isDataAvailable)
		time.Sleep(10 * time.Millisecond)
		i++
	}
	wg.Wait()
}

func Thread1(ctx *gin.Context, i int, wg *sync.WaitGroup, isDataAvailable chan bool) {
	funcData := funcDB("getData", i)

	// fmt.Printf("funcData value: %v\n", funcData)
	// fmt.Printf("funcData type : %T", funcData)

	val, exist := funcData.(int)
	fmt.Println("exist : ", exist)
	if exist && val == -1 {
		isDataAvailable <- false
		go Thread2("", wg, isDataAvailable, ctx)
	} else {
		data := Person{Name: funcData.(Person).Name, Age: funcData.(Person).Age, Year: funcData.(Person).Year}
		fmt.Println("funcDB data : ", data)

		p1, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Marshaling error : ", err)
		}

		fmt.Println("marshalled data : ", string(p1))
		isDataAvailable <- true
		go Thread2(string(p1), wg, isDataAvailable, ctx)
	}
}

func Thread2(data string, wg *sync.WaitGroup, isDataAvailable chan bool, ctx *gin.Context) {
	for {
		select {
		case out := <-isDataAvailable:
			fmt.Println("isDataAvailable: ", out)
			if out {
				var unMarshalledPerson Person
				if len(data) == 0 {
					fmt.Println("Error: Empty JSON Received")
				}
				err := json.Unmarshal([]byte(data), &unMarshalledPerson)
				if err != nil {
					fmt.Println("Unmarshaling error : ", err)
				}
				fmt.Println("unmarshalled data : ", unMarshalledPerson)
				//TODO: send the recived data to Thread 3 over grpc by goroutine
			} else {
				fmt.Println("Data is not available")
			}
		}
	}
}
