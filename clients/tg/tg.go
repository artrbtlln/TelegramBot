package tg

import (
	"TelegramBot/lib/e"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type CLient struct {
	host     string
	basePath string
	client   http.Client
}

const (
	getUpdateMethod   = "getUpdates"
	errMsg            = "cant do request"
	sendMessageMethod = "sendMessage"
)

func New(host string, token string) *CLient {
	return &CLient{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *CLient) Updates(offset int, limit int) ([]Update, error) {
	var res UpdateResponse
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdateMethod, q)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	return res.Result, nil
}

func (c *CLient) SendMessages(chatId int, text string) error {
	q := url.Values{}
	q.Add("chatId", strconv.Itoa(chatId))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap("cant send message", err)
	}
	return nil

}
func (c *CLient) doRequest(method string, query url.Values) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}
	return body, nil
}
