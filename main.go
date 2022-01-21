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

//ex1: In ra các message theo thứ tự. -- In ra message 3 trước message 2. Sử dụng 3 cách để làm( gợi ý: sử dụng mutex, chan, waitGroup)
//In ra các message theo thứ tự

// su dung waitgroup
func chanRoutine1() {
	wg.Add(1)
    log.Print("hello 1")
    go func() {
        time.Sleep(1 * time.Second)
        log.Print("hello 3")
		wg.Done()
    }()
    log.Print("hello 2")
	wg.Wait()
}

// su dung mutex
func chanRoutine2() {
	mu.Lock()
    log.Print("hello 1")
    go func() {
        time.Sleep(1 * time.Second)
        log.Print("hello 3")
		mu.Unlock()
    }()
    log.Print("hello 2")
	mu.Lock()

}

//su dung chan
func chanRoutine3() {
	ch := make(chan bool)
    log.Print("hello 1")
    if ch != nil {
		go func() {
			time.Sleep(1 * time.Second)
			log.Print("hello 3")
		}()
	}
	log.Print("hello 2")
}

//-- In ra message 3 trước message 2

// su dung WaitGroup
func chanRoutine4() {
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

// su dung mutex
func chanRoutine5() {
    log.Print("hello 1")
    go func() {
		time.Sleep(1 * time.Second)
        log.Print("hello 3")
		mu.Unlock()
    }()
	mu.Lock()	
	log.Print("hello 2")
}

// su dung chan
func chanRoutine6() {
	ch := make(chan bool)
    log.Print("hello 1")
    go func() {
        time.Sleep(1 * time.Second)
        log.Print("hello 3")
    }()
    if ch != nil {
		time.Sleep(2*time.Second)
	}
	log.Print("hello 2")
}

//ex2: tạo 1 biến X map[string]string và 3 goroutine cùng thêm dữ liệu vào X. Mỗi goroutine thêm 1000 key khác nhau. Sao cho quá trình đủ 15 key không mất mát dữ liệu.
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

//ex3: Tìm lỗi.
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
	log.Print("done")
	fmt.Println("Lỗi: Do nhiều dữ liệu chuyển vào m cùng lúc\n -> thêm waitgroup tách các go với nhau, thêm khóa để không bị lỗi khi truyền giá trị")
}

type Line struct {
	line_number int
	data string
}

//ex4: bài tập worker pool: tạo bằng tay file dưới. file.txt sau đó đọc từng dòng file này nạp dữ liệu vào 1 buffer channel có size 10, Điều kiện đọc file từng dòng. Chỉ được sử dụng 3 go routine. Kết quả xử lý xong ỉn ra màn hình + từ xong
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
	// time.Sleep(2*time.Second)
	// fmt.Println("----")
	// chanRoutine4()
	// time.Sleep(2*time.Second)

	// fmt.Println("Sử dụng Mutex: ")
	// chanRoutine2()
	// time.Sleep(2*time.Second)
	// fmt.Println("----")
	// chanRoutine5()
	// time.Sleep(2*time.Second)

	// fmt.Println("Sử dụng chan: ")
	// chanRoutine3()
	// time.Sleep(2*time.Second)
	// fmt.Println("----")
	// chanRoutine6()


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
	// ex4()
}