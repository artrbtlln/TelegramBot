package telegram

import (
	"TelegramBot/clients/tg"
	"TelegramBot/events"
	"TelegramBot/lib/e"
	"TelegramBot/storage"
	"errors"
)

type Meta struct {
	ChatID   int
	Username string
}
type Processor struct {
	tg      *tg.CLient
	offset  int
	storage storage.Storage
}

var unknownType = errors.New("unknown type")
var unknownMeta = errors.New("unknown meta")

func New(client *tg.CLient, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}
func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("cant't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))
	for _, u := range updates {
		res = append(res, event(u))
	}
	p.offset = updates[len(updates)-1].Id + 1

	return res, nil
}
func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap("cant process msg", unknownType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("cant process msg", err)
	}
	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Wrap("cant procces msg", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("cant get meta", unknownMeta)
	}
	return res, nil
}

func event(upd tg.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.Id,
			Username: upd.Message.From.Username,
		}
	}
	return res
}

func fetchText(upd tg.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}

func fetchType(upd tg.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}
