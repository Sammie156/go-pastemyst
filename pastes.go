package gopastemyst

// Currently only checking some fields.

// TODO: Add more fields based on the responses
// Check https://docs.beta.myst.rs/pastes
type Pasty struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Language string `json:"language"`
}

// TODO: Continue from here tomorrow and make a working
//
//	Struct for each Paste
type Paste struct {
	ID string `json:"id"`
}

// TODO: Work on this after defining the structure of a Paste
// func (c *Client) GetPaste(ctx context.Context, pasteID string) (*Pasty, error) {
// 	//url := fmt.Sprintf("%s/paste/%s", c.baseURL, pasteID)

// 	return ()
// }
