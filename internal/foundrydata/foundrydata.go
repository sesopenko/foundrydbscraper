package foundrydata

type DBData struct {
	Journals []Journal `json:"journal"`
	Items    []Item    `json:"items"`
	Actors   []Actor   `json:"actors"`
}

type Journal struct {
	Name  string        `json:"name"`
	ID    string        `json:"_id"`
	Pages []JournalPage `json:"pages"`
}

type JournalPage struct {
	Name string   `json:"name"`
	Type string   `json:"type"`
	ID   string   `json:"_id"`
	Text PageText `json:"text"`
}

type PageText struct {
	Content string `json:"content"`
}

type Item struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type ItemSystem struct {
	Description Description `json:"description"`
}

type Description struct {
	Value string `json:"value"`
}

type Actor struct {
	ID     string      `json:"_id"`
	Name   string      `json:"name"`
	Type   string      `json:"type"`
	System ActorSystem `json:"system"`
}

type ActorSystem struct {
	Details SystemDetails `json:"details"`
}

type SystemDetails struct {
	PublicNotes  string      `json:"publicNotes"`
	PrivateNotes string      `json:"privateNotes"`
	Alignment    StringValue `json:"alignment"`
	Level        NumberValue `json:"level"`
}
type NumberValue struct {
	Value int `json:"value"`
}

type StringValue struct {
	Value string `json:"value"`
}
