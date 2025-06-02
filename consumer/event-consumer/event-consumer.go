package event_consumer

import (
	"log"
	"tgBot/events"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("fetcher.Fetch error: %s", err)

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Printf("handleEvents error: %s", err)

			continue
		}

	}
}

// сделать параллельную обработку sync.WaitGroup()
func (c Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event: %#v", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("can't handle event: %s", err)

			continue
		}
	}
	return nil
}
