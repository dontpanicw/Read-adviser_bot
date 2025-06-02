package main

import (
	"flag"
	"log"
	tgClient "tgBot/clients/telegram"
	event_consumer "tgBot/consumer/event-consumer"
	telegram "tgBot/events/telegram"
	"tgBot/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "telegram_files_storage"
	batchSize   = 100
)

// 7389866096:AAFW-xmO9gtyhJPk8AUr1GcxMsqB4tcBcd4

func main() {
	//token = flags.Get(token)
	token := mustToken()
	host := mustHost()

	//tgClient = telegram.New(token)
	telegramClient := tgClient.New(host, token)

	//fetcher = fetcher.New()

	//processor = processor.New()
	eventProcessor := telegram.New(telegramClient, files.New(storagePath))

	log.Println("Service started")

	// consumer.Start(fetcher, processor)
	consumer := event_consumer.New(eventProcessor, eventProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal(err)
	}

}

func mustToken() string {
	token := flag.String("tg-bot-token", "", "Telegram bot token")

	flag.Parse()

	if *token == "" {
		log.Fatal("Telegram bot token is required")
	}

	return *token

}

func mustHost() string {
	host := flag.String("host", "api.telegram.org", "Telegram bot host")
	flag.Parse()

	return *host
}
