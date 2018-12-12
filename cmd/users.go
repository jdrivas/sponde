package cmd

import (
	"fmt"
	"net/http"
	"os"
	"sort"

	jh "github.com/jdrivas/sponde/jupyterhub"
	t "github.com/jdrivas/sponde/term"
	"github.com/juju/ansiterm"
	"github.com/spf13/cobra"
)

// User is a proxy for the jupyterhub/User.
type User jh.User

// UserList is a proxy for jh.UserList
type UserList jh.UserList

// ByName implements sort inteface for []jh.User
type ByName UserList

// Len implements sort interface
func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

// List prints a consice one line at a time reprsentation of
// users.
func (ul UserList) List() {
	users := jh.UserList(ul)
	if len(users) > 0 {
		w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
		fmt.Fprintf(w, "%s\n", t.Title("Name\tAdmin\tGroups\tCreated\tPending\tServer\tLast"))
		sort.Sort(ByName(users))
		for _, u := range users {
			serverURL := "<empty>"
			if u.ServerURL != "" {
				serverURL = u.ServerURL
			}
			fmt.Fprintf(w, "%s\n", t.SubTitle("%s\t%t\t%v\t%s\t%s\t%s\t%s", u.Name, u.Admin, u.Groups, u.Created, u.Pending, serverURL, u.LastActivity))
		}
		w.Flush()
	} else {
		fmt.Printf("There were no users.")
	}
}

// Describe prints all of the infomration there is about each user in the list.
// These are sorted by UserName and the servers are sorted by Name (this last
// implemented with sort.Stings()
func (ul UserList) Describe() {
	users := jh.UserList(ul)
	sort.Sort(ByName(users))
	for _, u := range users {
		w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
		fmt.Fprintf(w, "Name\tKind\tAdmin\tServer\tCreated\tLast Activity\tPending\n")
		pending := checkForEmptyString(u.Pending)
		serverURL := checkForEmptyString(u.ServerURL)
		fmt.Fprintf(w, "%s\t%s\n", t.Highlight("%s ", u.Name), t.Text("%s\t%t\t%s\t%s\t%s\t%s", u.Kind, u.Admin, serverURL, u.Created, u.LastActivity, pending))
		w.Flush()
		fmt.Println()
		if len(u.Servers) == 0 {
			fmt.Printf("No Servers\n")
		} else {
			fmt.Printf("Servers\n")
			w = ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
			fmt.Fprintf(w, "%s\n", t.Title("Name\tPdd\tReady\tPending\tStarted\tLast Activity"))
			var serverNames []string
			for k := range u.Servers {
				serverNames = append(serverNames, k)
			}
			sort.Strings(serverNames)
			for _, sn := range serverNames {
				s := u.Servers[sn]
				name := "<empty>"
				if s.Name != "" {
					name = s.Name
				}
				pending := "<empty>"
				if s.Pending != "" {
					pending = u.Pending
				}
				fmt.Fprintf(w, "%s\n", t.Text("%s\t%s\t%t\t%s\t%s\t%s", name, s.State.PodName, s.Ready, pending, s.Started, s.LastActivity))
			}
			w.Flush()
		}
		fmt.Println()
	}
}

// For the doUsers below.
func listUsers(u UserList, resp *http.Response, err error) {
	List(u, resp, err)
}
func describeUsers(u UserList, resp *http.Response, err error) {
	Describe(u, resp, err)
}

// doUsers is a command handler that will print a list of all users on the hub
// if no arguments are provided, or treat arguments as user names and print a list of users
// found on the Hub with details, and the names of users not found on the hub.
func doUsers(listFunc func(UserList, *http.Response, error)) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {

		var users jh.UserList
		var badNames []string
		var resp *http.Response
		var err error

		if len(args) > 0 {
			users, badNames, resp, err = jh.GetUsers(args)
		} else {
			users, resp, err = jh.GetAllUsers()
		}

		// Display users
		listFunc(UserList(users), resp, err)

		// Print an extra line if you have both
		if len(users) > 0 && len(badNames) > 0 {
			fmt.Println("")
		}
		// Displpay bad names if you have them
		if len(badNames) > 0 {
			// TODO: Pluralize
			fmt.Printf("There were %d user names not found on the Hub:\n", len(badNames))
			for _, n := range badNames {
				fmt.Printf("%s\n", n)
			}
		}
	}
}

// UpdatedUser is a proxy to add methods to jupyterhub/UpdatedUser
type UpdatedUser jh.UpdatedUser

// List displays the newely updated user details.
func (u UpdatedUser) List() {
	user := jh.UpdatedUser(u)
	w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
	fmt.Fprintf(w, "%s\n", t.Title("Name\tAdmin"))
	fmt.Fprintf(w, "%s\n", t.Text("%s\t%t", user.Name, user.Admin))
	w.Flush()
}

// Tokens is a proxy to add methhods to jupyterhub/Tokens
type Tokens jh.Tokens

// List tokens display a lit of API and OAuthTokens assocaited with a user.
func (ts Tokens) List() {
	tokens := jh.Tokens(ts)
	if len(tokens.APITokens) > 0 || len(tokens.OAuthTokens) > 0 {
		w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
		fmt.Fprintf(w, "%s\n", t.Title("ID\tKind\tCreated\tExpires\tLast Activity\tNote (OAuth client)"))
		for _, tk := range tokens.APITokens {
			tk.Expires = checkForEmptyString(tk.Expires)
			fmt.Fprintf(w, "%s\n", t.Text("%s\t%s\t%s\t%s\t%s\t%s", tk.ID, tk.Kind, tk.Created, tk.Expires, tk.LastActivity, tk.Note))
		}
		for _, tk := range tokens.OAuthTokens {
			tk.Expires = checkForEmptyString(tk.Expires)
			fmt.Fprintf(w, "%s\n", t.Text("%s\t%s\t%s\t%s\t%s\t%s", tk.ID, tk.Kind, tk.Created, tk.Expires, tk.LastActivity, tk.OAuthClient))
		}
		w.Flush()
	} else {
		fmt.Printf("No users tokens.\n")
	}
}

// APIToken is a proxy to add methods to jupyterhub/APIToken
type APIToken jh.APIToken

// Describe displays a single APIToken
func (tk APIToken) Describe() {
	token := jh.APIToken(tk)
	w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
	fmt.Fprintf(w, "%s\n", t.Title("ID\tKind\tCreated\tExpires\tLast Activity\tNote (OAuth client)"))
	token.Expires = checkForEmptyString(token.Expires)
	fmt.Fprintf(w, "%s\n", t.Text("%s\t%s\t%s\t%s\t%s\t%s", token.ID, token.Kind, token.Created, token.Expires, token.LastActivity, token.Note))
	w.Flush()
}
