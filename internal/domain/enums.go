package domain

type Language string

func (l Language) String() string {
	return string(l)
}

const (
	LanguageEN Language = "en"
	LanguageRU Language = "ru"
	LanguageUA Language = "ua"
)

type ChatService string

func (c ChatService) String() string {
	return string(c)
}

const (
	ChatServiceTelegram ChatService = "telegram"
)

type StreamingService string

func (s StreamingService) String() string {
	return string(s)
}

const (
	StreamingServiceTwitch StreamingService = "twitch"
)
