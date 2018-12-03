package jupyterhub

import (
	"fmt"
	"os"
	"text/tabwriter"
	// "text/tabwriter"
)

type Tokens struct {
	ApiTokens   []ApiToken   `json:"api_tokens"`
	OAuthTokens []OAuthToken `json:"oauth_tokens"`
}

type ApiToken struct {
	ID           string `json:"id"`
	Kind         string `json:"kind"`
	User         string `json:"user"`
	Created      string `json:"created"`
	LastActivity string `json:"last_activity"`
	Note         string `json:"note"`
}

type OAuthToken struct {
	ID           string `json:"id"`
	Kind         string `json:"kind"`
	User         string `json:"user"`
	Created      string `json:"created"`
	LastActivity string `json:"last_activity"`
	OAuthClient  string `json:"oauth_client"`
}

func GetTokens(username string) (token Tokens, err error) {
	_, err = get(fmt.Sprintf("/users/%s/tokens", username), &token)
	return token, err
}

func (tokens *Tokens) Print() {

	w := tabwriter.NewWriter(os.Stdout, 4, 4, 3, ' ', 0)
	fmt.Fprintf(w, "ID\tKind\tCreated\tLast Activity\tClient\n")
	for _, t := range tokens.ApiTokens {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", t.ID, t.Kind, t.Created, t.LastActivity, t.Note)
	}
	for _, t := range tokens.OAuthTokens {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", t.ID, t.Kind, t.Created, t.LastActivity, t.OAuthClient)
	}
	w.Flush()
}
