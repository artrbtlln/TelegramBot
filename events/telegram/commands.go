package telegram

import (
	"TelegramBot/storage"
	"context"
	"errors"
	"log"
	"net/url"
	"strings"
)

const (
	GetCmd   = "/get"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatId int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("get %s from %s", text, username)

	if isAddCmd(text) {
		return p.createPage(chatId, text, username)
	}
	switch text {
	case GetCmd:
		return p.getPage(chatId, username)
	case HelpCmd:
		return p.sendHelp(chatId)
	case StartCmd:
		return p.sendHello(chatId)
	default:
		return p.tg.SendMessages(chatId, msgUnknownCommand)

	}
}
func (p *Processor) createPage(chatId int, pageUrl string, username string) error {
	page := &storage.Page{
		URL:      pageUrl,
		UserName: username,
	}
	isExist, err := p.storage.IsExists(context.Background(), page)
	if err != nil {
		return nil
	}
	if isExist {
		return p.tg.SendMessages(chatId, msgAlreadyExists)
	}

	if err := p.storage.Create(context.Background(), page); err != nil {
		return err
	}
	if err := p.tg.SendMessages(chatId, msgSaved); err != nil {
		return err
	}
	return nil
}

func (p *Processor) getPage(chatId int, username string) error {
	page, err := p.storage.Get(context.Background(), username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessages(chatId, msgNoSavedPages)
	}
	if err := p.tg.SendMessages(chatId, page.URL); err != nil {
		return err
	}
	return nil
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessages(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessages(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isUrl(text)
}

func isUrl(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
