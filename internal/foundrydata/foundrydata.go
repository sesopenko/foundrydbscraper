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
	ID     string     `json:"_id"`
	Name   string     `json:"name"`
	Type   string     `json:"type"`
	System ItemSystem `json:"system"`
}

type ItemByName []Item

func (ibn ItemByName) Swap(i, j int) {
	ibn[i], ibn[j] = ibn[j], ibn[i]
}

func (ibn ItemByName) Less(i, j int) bool {
	return ibn[i].Name < ibn[j].Name
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

type ActorByName []Actor

// Len returns the length of the slice
func (a ActorByName) Len() int {
	return len(a)
}

// Swap swaps the elements with indexes i and j
func (a ActorByName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less returns true if the element with index i should sort before the element with index j
func (a ActorByName) Less(i, j int) bool {
	return a[i].Name < a[j].Name
}

type ActorSystem struct {
	Details SystemDetails `json:"details"`
}

type SystemDetails struct {
	PublicNotes  string      `json:"publicNotes"`
	PrivateNotes string      `json:"privateNotes"`
	Alignment    StringValue `json:"alignment"`
	Level        NumberValue `json:"level"`
	Description  string      `json:"description"`
	Disable      string      `json:"disable"`
	Reset        string      `json:"reset"`
	Routine      string      `json:"routine"`
	Blurb        string      `json:"blurb"`
}
type NumberValue struct {
	Value int `json:"value"`
}

type StringValue struct {
	Value string `json:"value"`
}
