package telegram

import (
	"errors"
	"log"
	"net/url"
	"strings"
	"tgBot/storage"
)

const (
	rndCmd = "/rnd"
	help   = "/help"
	start  = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command: %s, from %s", text, username)

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case rndCmd:
		return p.sendRandom(chatID, username)
	case help:
		return p.sendHelp(chatID)
	case start:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMassage(chatID, msgUnknownCommand)

	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMassage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}
	if err := p.tg.SendMassage(chatID, msgSaved); err != nil {
		return err
	}
	return nil
}

func (p *Processor) sendRandom(chatID int, username string) error {
	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMassage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMassage(chatID, page.URL); err != nil {
		return err
	}
	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMassage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMassage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
