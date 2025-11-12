package gopastemyst

import "time"

// ----- PASTES TYPES -----

type Pasty struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Language string `json:"language"`
}

type Paste struct {
	// Non-Nullable fields
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresIn string    `json:"expiresIn"`
	Pinned    bool      `json:"pinned"`
	Private   bool      `json:"private"`
	Stars     int       `json:"stars"`
	Tags      []string  `json:"tags"`
	Pasties   []Pasty   `json:"pasties"`

	// Nullable fields
	OwnerID   *string    `json:"ownerID"`
	EditedAt  *time.Time `json:"editedAt"`
	DeletesAt *time.Time `json:"deletesAt"`
}

type PasteDiff struct {
	CurrentPaste Paste `json:"currentPaste"`
	NewPaste     Paste `json:"newPaste"`
	OldPaste     Paste `json:"oldPaste"`
}

type PastyStats struct {
	Bytes int `json:"bytes"`
	Lines int `json:"lines"`
	Words int `json:"words"`
}

type Stats struct {
	Bytes   int                   `json:"bytes"`
	Lines   int                   `json:"lines"`
	Pasties map[string]PastyStats `json:"pasties"`
	Words   int                   `json:"words"`
}

type CreatePastyOptions struct {
	Title    string `json:"title,omitempty"`
	Content  string `json:"content"`
	Language string `json:"language,omitempty"`
}

type CreatePasteOptions struct {
	Title     string               `json:"title,omitempty"`
	ExpiresIn string               `json:"expiresIn,omitempty"`
	Anonymous bool                 `json:"anonymous,omitempty"`
	Private   bool                 `json:"private,omitempty"`
	Pinned    bool                 `json:"pinned,omitempty"`
	Encrypted bool                 `json:"encrypted,omitempty"`
	Tags      []string             `json:"tags,omitempty"`
	Pasties   []CreatePastyOptions `json:"pasties"`
}

type LanguageStats struct {
	Aliases            []string `json:"aliases"`
	CodemirrorMimeType string   `json:"codemirrorMimeType"`
	CodemirrorMode     string   `json:"codemirrorMode"`
	Color              string   `json:"color"`
	Extensions         []string `json:"extensions"`
	Name               string   `json:"name"`
	TmScope            string   `json:"tmScope"`
	Wrap               bool     `json:"wrap"`
}

type PasteLanguageStats struct {
	Language   LanguageStats `json:"language"`
	Percentage float64       `json:"percentage"`
}

type CompactPasteHistory struct {
	EditedAt time.Time `json:"editedAt"`
	ID       string    `json:"id"`
}

// ----- USER TYPES -----

type User struct {
	AvatarID      string    `json:"avatarId"`
	CreatedAt     time.Time `json:"createdAt"`
	ID            string    `json:"id"`
	IsAdmin       bool      `json:"isAdmin"`
	IsContributor bool      `json:"isContributor"`
	Username      string    `json:"username"`
}

type GetUserPasteOptions struct {
	Page     int
	PageSize int
	Tag      string
}
