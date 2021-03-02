package main

import (
	"fmt"
	"math/rand"
	"time"
)

// 消息生成器
func genMsg(name string, done chan struct{}) chan string {
	c := make(chan string)

	go func() {
		i := 0
		for {
			select {
			case <-time.After(time.Duration(rand.Intn(5000)) * time.Millisecond):
				c <- fmt.Sprintf("task:%s -- msg %d", name, i)
			case <-done:
				fmt.Println("cleaning up")
				time.Sleep(2 * time.Second)
				fmt.Println("cleaning done")
				done <- struct{}{}
				return
			}

			i++
		}
	}()
	return c
}

// 非阻塞方法
func noBlockWait(c chan string) (string, bool) {

	select {
	case m := <-c:
		return m, true
	default:
		return "", false

	}
}

func fanIn(c1, c2 chan string) chan string {

	c := make(chan string)
	go func() {
		for {
			select {
			case m := <-c1:
				c <- m
			case m := <-c2:
				c <- m
			}
		}
	}()
	return c
}

// 超时等待
func timeoutWait(c chan string, timeout time.Duration) (string, bool) {
	select {
	case m := <-c:
		return m, true
	case <-time.After(timeout):
		return "", false

	}
}
func main() {
	done := make(chan struct{})
	m1 := genMsg("one", done)
	for i := 0; i < 5; i++ {
		if m, ok := timeoutWait(m1, 1*time.Second); ok {
			fmt.Println(m)
		} else {
			fmt.Println("timeout")
		}
	}
	done <- struct{}{}
	<-done

	time.Sleep(time.Second)
}
