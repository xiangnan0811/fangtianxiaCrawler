package engine

type ConcurrentEngine struct {
	Scheduler        Scheduler
	WorkerCount      int
	ItemChan         chan Item
	RequestProcessor Processor
}

type Processor func(r Request) (ParseResult, error)

type Scheduler interface {
	ReadyNotifier
	Submit(Request)
	WorkerChan() chan Request
	Run()
}
type ReadyNotifier interface {
	WorkerReady(chan Request)
}

func (c *ConcurrentEngine) Run(seeds ...Request) {
	out := make(chan ParseResult)
	c.Scheduler.Run()

	for i := 0; i < c.WorkerCount; i++ {
		c.createWorker(c.Scheduler.WorkerChan(), out, c.Scheduler)
	}

	for _, r := range seeds {
		if isDuplicated(r.Url) {
			continue
		}
		c.Scheduler.Submit(r)
	}

	for {
		result := <-out
		for _, item := range result.Items {
			item := item
			go func() {
				c.ItemChan <- item
			}()
		}

		for _, request := range result.Requests {
			if isDuplicated(request.Url) {
				continue
			}
			c.Scheduler.Submit(request)
		}
	}
}

func (c *ConcurrentEngine) createWorker(in chan Request, out chan ParseResult, ready ReadyNotifier) {
	go func() {
		for {
			ready.WorkerReady(in)
			request := <-in
			result, err := c.RequestProcessor(request)
			if err != nil {
				continue
			}
			out <- result
		}
	}()
}

var visitedUrls = make(map[string]bool)

func isDuplicated(url string) bool {
	if visitedUrls[url] {
		return true
	}
	visitedUrls[url] = true
	return false
}
