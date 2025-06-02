package telegram

import (
	"errors"
	"tgBot/clients/telegram"
	"tgBot/events"
	"tgBot/lib/e"
	"tgBot/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	update, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap(err, "failed to fetch updates")
	}

	if len(update) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(update))

	for _, update := range update {
		res = append(res, event(update))
	}

	p.offset = update[len(update)-1].ID + 1

	return res, nil

}

func (p *Processor) Process(event events.Event) (err error) {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap(errors.New("unknown event type"), "can't process message")
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap(err, "can't process message")
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Wrap(err, "can't process message")
	}

	return nil

}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap(errors.New("unknown meta type"), "can't get meta")
	}

	return res, nil
}

func event(update telegram.Update) events.Event {
	updType := fetchType(update)
	res := events.Event{
		Type: updType,
		Text: fetchText(update),
	}
	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   update.Message.Chat.ID,
			Username: update.Message.From.Username,
		}
	}

	return res

}

func fetchType(update telegram.Update) events.Type {
	if update.Message == nil {
		return events.Unknown
	}
	return events.Message
}
func fetchText(update telegram.Update) string {
	if update.Message == nil {
		return ""
	}
	return update.Message.Text
}
