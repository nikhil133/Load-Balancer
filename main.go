package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Request is used for making requests to services behind a load balancer.
type Request struct {
	Payload interface{}
	RspChan chan Response
}

// Response is the value returned by services behind a load balancer.
type Response interface{}

// LoadBalancer is used for balancing load between multiple instances of a service.
type LoadBalancer interface {
	Request(payload interface{}) chan Response
	RegisterInstance(chan Request)
}

var etcd map[chan Request]bool

func init() {
	etcd = make(map[chan Request]bool)
}

// MyLoadBalancer is the load balancer you should modify!
type MyLoadBalancer struct {
	workers []TimeService
}

// Request is currently a dummy implementation. Please implement it!
func (lb *MyLoadBalancer) Request(payload interface{}) chan Response {
	ch := make(chan Response, 1)
	var i int
	req := Request{
		Payload: payload,
		RspChan: make(chan Response),
	}

	for {
		if len(etcd) == 0 {
			log.Println("Info: No workers available")
			return ch
		}
		i = rand.Int() % len(lb.workers)
		if up, _ := etcd[lb.workers[i].ReqChan]; !up {
			delete(etcd, lb.workers[i].ReqChan)
			lb.workers = append(lb.workers[:i], lb.workers[i+1:]...)
			continue
		}
		break
	}
	go func(i int) {
		worker := &lb.workers[i]
		log.Println("Worker ", &worker)
		lb.workers[i].ReqChan <- req
	}(i)

	select {
	case x := <-req.RspChan:
		ch <- x
	}
	return ch
}

// RegisterInstance is currently a dummy implementation. Please implement it!
func (lb *MyLoadBalancer) RegisterInstance(ch chan Request) {
	ts := TimeService{ReqChan: ch}
	lb.workers = append(lb.workers, ts)
	etcd[ch] = true
	return
}

/******************************************************************************
 *  STANDARD TIME SERVICE IMPLEMENTATION -- MODIFY IF YOU LIKE                *
 ******************************************************************************/

// TimeService is a single instance of a time service.
type TimeService struct {
	Dead            chan struct{}
	ReqChan         chan Request
	AvgResponseTime float64
	WorkerId        int
}

// Run will make the TimeService start listening to the two channels Dead and ReqChan.
func (ts *TimeService) Run() {
	for {
		select {
		case <-ts.Dead:
			etcd[ts.ReqChan] = false
			return
		case req := <-ts.ReqChan:
			processingTime := time.Duration(ts.AvgResponseTime+1.0-rand.Float64()) * time.Second
			time.Sleep(processingTime)
			req.RspChan <- time.Now()
		}
	}
}

/******************************************************************************
 *  CLI -- YOU SHOULD NOT NEED TO MODIFY ANYTHING BELOW                       *
 ******************************************************************************/

// main runs an interactive console for spawning, killing and asking for the
// time.
func main() {
	rand.Seed(int64(time.Now().Nanosecond()))

	bio := bufio.NewReader(os.Stdin)
	var lb LoadBalancer = &MyLoadBalancer{}

	manager := &TimeServiceManager{}

	for {
		fmt.Printf("> ")
		cmd, err := bio.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading command: ", err)
			continue
		}
		switch strings.TrimSpace(cmd) {
		case "kill":
			manager.Kill()
		case "spawn":
			ts := manager.Spawn()
			lb.RegisterInstance(ts.ReqChan)
			go ts.Run()
		case "time":
			select {
			case rsp := <-lb.Request(nil):
				fmt.Println(rsp)
			case <-time.After(5 * time.Second):
				fmt.Println("Timeout")
			}
		default:
			fmt.Printf("Unknown command: %s Available commands: time, spawn, kill\n", cmd)
		}
	}
}

// TimeServiceManager is responsible for spawning and killing.
type TimeServiceManager struct {
	Instances []TimeService
}

// Kill makes a random TimeService instance unresponsive.
func (m *TimeServiceManager) Kill() {
	if len(m.Instances) > 0 {
		n := rand.Intn(len(m.Instances))
		close(m.Instances[n].Dead)
		m.Instances = append(m.Instances[:n], m.Instances[n+1:]...)
	}
}

// Spawn creates a new TimeService instance.
func (m *TimeServiceManager) Spawn() TimeService {
	ts := TimeService{
		Dead:            make(chan struct{}, 0),
		ReqChan:         make(chan Request, 10),
		AvgResponseTime: rand.Float64() * 3,
	}
	m.Instances = append(m.Instances, ts)

	return ts
}
