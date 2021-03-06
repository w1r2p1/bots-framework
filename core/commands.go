package bots

import (
	"fmt"
	"net/url"
)

type CommandAction func(whc WebhookContext) (m MessageFromBot, err error)
type CallbackAction func(whc WebhookContext, callbackURL *url.URL) (m MessageFromBot, err error)

type CommandMatcher func(Command, WebhookContext) bool

const DEFAULT_TITLE = ""
const SHORT_TITLE = "short_title"

//const LONG_TITLE = "long_title"

type Command struct {
	InputTypes     []WebhookInputType // Instant match if != WebhookInputUnknown && == whc.InputTypes()
	Icon           string
	Replies        []Command
	Code           string
	Title          string
	Titles         map[string]string
	ExactMatch     string
	Commands       []string
	Matcher        CommandMatcher
	Action         CommandAction
	CallbackAction CallbackAction
}

func NewCallbackCommand(code string, action CallbackAction) Command {
	return Command{
		Code:           code,
		//Commands: []string{code},
		CallbackAction: action,
	}
}

func (c Command) String() string {
	return fmt.Sprintf("Command{ByCode: '%v', InputTypes: %v, Icon: '%v', Title: '%v', ExactMatch: '%v', len(Command): %v, len(Replies): %v}", c.Code, c.InputTypes, c.Icon, c.Title, c.ExactMatch, len(c.Commands), len(c.Replies))
}

func (whcb *WebhookContextBase) CommandText(title, icon string) string {
	if title != "" {
		title = whcb.Translate(title)
	}
	return CommandTextNoTrans(title, icon)
}

func CommandTextNoTrans(title, icon string) string {
	if title == "" && icon != "" {
		return icon
	} else if title != "" && icon == "" {
		return title
	} else if title != "" && icon != "" {
		return title + " " + icon
	} else {
		return "<NO_TITLE_OR_ICON>"
	}
}

func (c Command) DefaultTitle(whc WebhookContext) string {
	return c.TitleByKey(DEFAULT_TITLE, whc)
}

func (c Command) TitleByKey(key string, whc WebhookContext) string {
	var title string
	if key == DEFAULT_TITLE && c.Title != "" {
		title = c.Title
	} else if val, ok := c.Titles[key]; ok {
		title = val
	}

	if c.Icon == "" {
		if title == "" {
			title = c.Code
		} else {
			title = whc.Translate(title)
		}
	} else {
		if title == "" {
			title = CommandTextNoTrans("", c.Icon)
		} else {
			title = whc.CommandText(title, c.Icon)
		}
	}
	return title
}
