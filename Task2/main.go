package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// 1、
func test1(a *int) {
	*a += 10
}

func test2(a []int) {
	for index := range a {
		a[index] *= 2
	}
}

// 2.1
func test3() {

	go func() {
		jishu := []int{1, 3, 5, 7, 9}
		for _, value := range jishu {
			fmt.Println(value)
		}
	}()

	go func() {
		oushu := []int{2, 4, 6, 8, 10}
		for _, value := range oushu {
			fmt.Println(value)
		}
	}()

	time.Sleep(3000)
}

// 2.2
var wg sync.WaitGroup

func taskMethod(i int) {
	defer wg.Done()
	fmt.Println("task i=", i)
	i *= 1000
	time.Sleep(time.Duration(i))
}

func timeUsed(start time.Time) {
	tm := time.Since(start)
	fmt.Println("time use", tm)
}

// 3.1
type Shape interface {
	Area() float32
	Perimeter() float32
}

type Rectangle struct {
	width  float32
	height float32
}

type Circle struct {
	r float32
}

func (r Rectangle) Area() float32 {
	return r.width * r.height
}

func (r Rectangle) Perimeter() float32 {
	return 2 * (r.width + r.height)
}

func (c Circle) Area() float32 {
	var pai float32 = 3.14
	return float32(2) * c.r * pai
}

func (c Circle) Perimeter() float32 {
	var pai float32 = 3.14
	return pai * c.r * pai * c.r
}

// 3.2
type Person struct {
	name string
	age  int
}

type Employee struct {
	Person            // 组合 Person 结构体（继承其 Name、Age 字段）
	employeeID string // Employee 新增字段：员工ID
}

func (e Employee) PrintInfo() {
	fmt.Printf("员工ID：%s\n", e.employeeID)
	fmt.Printf("姓名：%s\n", e.name)
	fmt.Printf("年龄：%d\n", e.age)
}

// 4.1
// 4.2 在定义的时候加上缓冲区就行 chanel := make(chan int, 10)
func pushData(ch chan<- int) {
	for i := 1; i < 11; i++ {
		ch <- i
		fmt.Printf("发送: %d\n", i)
	}
	// close(ch)
}

func getData(ch <-chan int) {
	for value := range ch {
		fmt.Printf("通道收到消息: %d\n", value)
	}
}

// 5.1
type SafeCount struct {
	mu    sync.Mutex
	count int
}

func (s *SafeCount) increment() {
	s.mu.Lock()
	s.count++
	s.mu.Unlock()
}

func (c *SafeCount) GetCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

func main() {
	// s := SafeCount{}
	// for i := 0; i < 10; i++ {
	// 	go func ()  {
	// 		for i := 0; i < 1000; i++ {
	// 			s.increment()
	// 		}
	// 	}()
	// }
	//  // 等待一段时间确保所有goroutine完成
	// time.Sleep(time.Second)

	// // 输出最终计数
	// fmt.Printf("Final count: %d\n", s.GetCount())

	//5.2
	var counter int64
	//协程同步等待组（等待10个协程全部完成）
	var wg sync.WaitGroup

	// 3. 启动10个协程
	for i := 0; i < 10; i++ {
		wg.Add(1) // 每启动一个协程，等待组计数+1
		go func() {
			defer wg.Done() // 协程结束时，等待组计数-1
			// 每个协程执行1000次原子递增
			for j := 0; j < 1000; j++ {
				// 原子操作：对 counter 执行 +1，返回递增后的值（此处忽略返回值）
				atomic.AddInt64(&counter, 1)
			}
		}()
	}

	// 4. 等待所有协程完成（阻塞直到等待组计数为0）
	wg.Wait()

	// 5. 原子读取最终结果（保证读取操作的原子性，避免中间状态）
	finalCount := atomic.LoadInt64(&counter)
	fmt.Printf("最终计数器值：%d\n", finalCount)
}
