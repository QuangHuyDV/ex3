package main

import (
	// "fmt"
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

var mu sync.Mutex
var wg sync.WaitGroup

//ex1
func chanRoutine1() {
	wg.Add(1)
    log.Print("hello 1")
    go func() {
        time.Sleep(1 * time.Second)
        log.Print("hello 3")
		wg.Done()
    }()
	wg.Wait()
    log.Print("hello 2")
}

func chanRoutine2() {
    log.Print("hello 1")
	mu.Lock()
    go func() {
		time.Sleep(1 * time.Second)
        log.Print("hello 3")
		mu.Unlock()
    }()
	mu.Lock()	
	log.Print("hello 2")
	mu.Unlock()
}

func chanRoutine3() {
	ch := make(chan string,3)
	ch <- "hello 1"
	ch <- "hello 3"
	ch <- "hello 2"
    log.Print(<-ch)
    go func() {
		time.Sleep(1 * time.Second)
        log.Print(<-ch)
	}()
	log.Print(<-ch)
	time.Sleep(2*time.Second)
}

//ex2
func add1(X map[string]string) {
	a := "student"
	for i := 0; i < 1000; i++ {
		mu.Lock()
		add := a + strconv.Itoa(i)
		X[add] = strconv.Itoa(i)
 		mu.Unlock()
	}
	wg.Done()
}

func add2(X map[string]string) {
	a := "class"
	for i := 0; i < 1000; i++ {
		mu.Lock()
		add := a + strconv.Itoa(i)
		X[add] = strconv.Itoa(i)
 		mu.Unlock()
	}
	wg.Done()
}

func add3(X map[string]string) {
	a := "School"
	for i := 0; i < 1000; i++ {
		mu.Lock()
		add := a + strconv.Itoa(i)
		X[add] = strconv.Itoa(i)
 		mu.Unlock()
	}
	wg.Done()

}

//ex3
func errFunc() {
	m := make(map[int]int)
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		defer wg.Done()
		go func() {
			for j := 1; j < 10000; j++ {
				mu.Lock()
				defer mu.Unlock()
				if _, ok := m[j]; ok {
					delete(m, j)
					continue
				}
				m[j] = j * 10
			}
			wg.Wait()
		}()
	}
	for i := range m {
		fmt.Printf(" gia tri %v\n", m[i])
	}
	log.Print("done")
}

type Line struct {
	line_number int
	data string
}

//ex4
func ex4() {
	ch := make(chan string,10)
	finish := make(chan bool)
	f, err := os.Open("file.txt")
	if err != nil {
		log.Fatal(err)
	}
	num:= 0
	li := []*Line{}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		ch <- scanner.Text()
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go achan(finish,ch,scanner)
			wg.Done()
		}
		time.Sleep(100*time.Millisecond)
		a := Line{line_number: num, data: scanner.Text()}
		new := append(li, &a)
		num ++
		for i := range new {
			fmt.Printf("%v giá trị là: %v\n",new[i].line_number,new[i].data)
		}	
	}

	fmt.Println("xong")
}

func achan(finish chan bool,ch chan string, scanner *bufio.Scanner) {
	for range scanner.Text() {
		fmt.Printf("Giá trị %v đã lấy ra\n", <-ch)
	}
	finish <- false
	close(finish)
	close(ch)
}

func main() {
	//ex1
	// fmt.Println("Sử dụng WaitGroup: ")
	// chanRoutine1()
	// time.Sleep(1*time.Second)
	// fmt.Println("Sử dụng Mutex: ")
	// chanRoutine2()
	// time.Sleep(1*time.Second)
	// fmt.Println("Sử dụng chan: ")
	// chanRoutine3()


	//ex2
	// X :=  make(map[string]string)
	// wg.Add(3)
	// go add1(X)
	// go add2(X)
	// go add3(X)
	// wg.Wait()
	// count := 0
	// for key, value := range X {
	// 	fmt.Printf("Key: %v, value: %v\n", key,value)
	// 	count ++
	// 	// if count == 15 {
	// 	// 	break
	// 	// }
	// }
	// fmt.Printf("\nTong gia tri da duoc them: %v",count)


	//ex3
	// errFunc()


	//ex4
	ex4()
}