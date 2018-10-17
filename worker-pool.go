package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Pool interface {
	Start()
	Assign(Job)
	Shutdown()
}

type Job interface {
	Perform()
}

type FixedPool struct {
	Capacity int
	queue    chan Job
	kill     []chan bool
}

func (fp *FixedPool) Start() {
	fp.kill = make([]chan bool, 0, fp.Capacity)
	fp.queue = make(chan Job)
	for i := 0; i < fp.Capacity; i++ {
		quit := make(chan bool)
		fp.kill = append(fp.kill, quit)
		go func(quit chan bool, queue chan Job) {
			for {
				select {
				case job := <-queue:
					job.Perform()
				case <-quit:
					return
				}
			}
		}(quit, fp.queue)
	}
}

func (fp *FixedPool) Assign(job Job) {
	fp.queue <- job
}

func (fp *FixedPool) Shutdown() {
	for i := 0; i < fp.Capacity; i++ {
		fp.kill[i] <- true
		close(fp.kill[i])
	}
	close(fp.queue)
}

type Response struct {
	Content string
	Err     error
}

type Downloader struct {
	URL    string
	Result chan Response
}

func (d Downloader) Perform() {
	log.Println("Downloading: ", d.URL)
	// mimic long running task
	time.Sleep(1 * time.Second)
	response := Response{Content: "", Err: nil}
	resp, err := http.Get(d.URL)
	if err != nil {
		response.Err = err
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		response.Err = err
		return
	}
	response.Content = string(body)
	log.Println("Download complete: ", d.URL)
	d.Result <- response
}

func main() {
	done := make(chan bool)
	results := make(chan Response)
	go func(results chan Response) {
		for response := range results {
			if response.Err != nil {
				fmt.Println("Failed to fetch response: ", response.Err)
			} else {
				fmt.Println("Got response: ", response.Content[0:20])
			}
		}
		done <- true
	}(results)
	pool := &FixedPool{Capacity: 4}
	pool.Start()
	links := []string{
		"https://www.google.com",
		"http://example.com",
		"https://www.google.com",
		"http://example.com",
		"https://www.google.com",
		"http://example.com",
		"https://www.google.com",
		"http://example.com",
		"https://www.google.com",
		"http://example.com",
		"https://www.google.com",
		"http://example.com",
	}
	for _, link := range links {
		d := Downloader{URL: link, Result: make(chan Response)}
		pool.Assign(d)
		go func(queue, res chan Response) {
			queue <- (<-res)
		}(results, d.Result)
	}
	pool.Shutdown()
	time.Sleep(1 * time.Second)
	close(results)
	<-done
}
