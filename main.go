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

func mainAPI(c *gin.Context) {
	isDataAvailable := make(chan bool, 1)
	wg := &sync.WaitGroup{}
	funcDBLength := funcDB("getlength", 0)
	fmt.Println("funcDBLength : ", funcDBLength)
	wg.Add(1)

	var i = 0
	for i < funcDBLength.(int)+1 {
		go Thread1(c, i, wg, isDataAvailable)
		time.Sleep(1 * time.Second)
		i++
	}
	wg.Wait()
}

func Thread1(c *gin.Context, i int, wg *sync.WaitGroup, isDataAvailable chan bool) {
	funcData := funcDB("getData", i)

	// fmt.Printf("funcData value: %v\n", funcData)
	// fmt.Printf("funcData type : %T", funcData)

	val, ok := funcData.(int)
	if ok && val == -1 {
		isDataAvailable <- false
		go Thread2("", wg, isDataAvailable, c)
	} else {
		data := Person{Name: funcData.(Person).Name, Age: funcData.(Person).Age, Year: funcData.(Person).Year}
		fmt.Println("data : ", data)

		p1, err := json.Marshal(data)
		if err != nil {
			fmt.Println("error while marshaling 1 : ", err)
		}

		fmt.Println("marshalled data p1 : ", string(p1))
		isDataAvailable <- true
		go Thread2(string(p1), wg, isDataAvailable, c)
	}
}

func Thread2(data string, wg *sync.WaitGroup, isDataAvailable chan bool, c *gin.Context) {
	for {
		select {

		case out := <-isDataAvailable:
			fmt.Println("isDataAvailable: ", out)
			if out {
				var unMarshalledPerson Person
				if len(data) == 0 {
					fmt.Println("Error: Empty JSON input")
				}
				err := json.Unmarshal([]byte(data), &unMarshalledPerson)
				if err != nil {
					fmt.Println("error while Unmarshaling 2: ", err)
				}
				fmt.Println("unmarshalled data : ", unMarshalledPerson)
				//TODO: send the recived data to Thread 3 over grpc by goroutine
			} else {
				fmt.Println("Data is not available")
			}
		}
	}

}
