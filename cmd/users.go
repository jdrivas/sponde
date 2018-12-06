package cmd

import (
	"fmt"
	"os"

	jh "github.com/jdrivas/sponde/jupyterhub"
	t "github.com/jdrivas/sponde/term"
	"github.com/juju/ansiterm"
	"github.com/spf13/cobra"
)

// ListUsers prints a consice one line at a time reprsentation of
// users.
func listUsers(users jh.UserList) {
	w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
	fmt.Fprintf(w, "%s\n", t.Title("Name\tAdmin\tCreated\tServer\tLast"))
	for _, u := range users {
		serverURL := "<empty>"
		if u.ServerURL != "" {
			serverURL = u.ServerURL
		}
		// fmt.Fprintf(w, "%s\t%t\t%s\t%s\t%s\n", u.Name, u.Admin, u.Created, serverURL, u.LastActivity)
		fmt.Fprintf(w, "%s\n", t.SubTitle("%s\t%t\t%s\t%s\t%s", u.Name, u.Admin, u.Created, serverURL, u.LastActivity))
	}
	w.Flush()
}

// DescribeUsers prints all of the infomration there is about each user.
func describeUsers(users jh.UserList) {
	for _, u := range users {
		w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
		fmt.Fprintf(w, "Name\tKind\tAdmin\tServer\tCreated\tLast Activity\tPending\n")
		pending := checkForEmptyString(u.Pending)
		serverURL := checkForEmptyString(u.ServerURL)
		fmt.Fprintf(w, "%s\t%s\n", t.Highlight("%s ", u.Name), t.Text("%s\t%t\t%s\t%s\t%s", u.Kind, u.Admin, serverURL, u.Created, u.LastActivity, pending))
		w.Flush()
		fmt.Println()
		if len(u.Servers) == 0 {
			fmt.Printf("No Servers\n")
		} else {
			fmt.Printf("Servers\n")
			for _, s := range u.Servers {
				w = ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
				fmt.Fprintf(w, "Name\tReady\tPending\tStarted\tLast Activity\n")
				name := "<empty>"
				if s.Name != "" {
					name = s.Name
				}
				pending := "<empty>"
				if s.Pending != "" {
					pending = u.Pending
				}
				fmt.Fprintf(w, "%s\t%t\t%s\t%s\t%s\n", name, s.Ready, pending, s.Started, s.LastActivity)
				w.Flush()
			}
		}
		fmt.Println()
	}
}

// doUsers is a command handler that will print a list of all users on the hub
// if no arguments are provided, or treat arguments as user names and print a list of users
// found on the Hub with details, and the names of users not found on the hub.
func doUsers(listFunc func(jh.UserList), logError func(error)) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		var users jh.UserList
		var badNames []string
		var err error
		if len(args) > 0 {
			users, badNames, err = jh.GetUsers(args)
		} else {
			users, err = jh.GetAllUsers()
		}
		if err == nil {
			if len(users) > 0 {
				listFunc(users)
			}
			if len(users) > 0 && len(badNames) > 0 {
				fmt.Println("")
			}
			if len(badNames) > 0 {
				// TODO: Pluralize
				fmt.Printf("There were %d user names not found on the Hub:\n", len(badNames))
				for _, n := range badNames {
					fmt.Printf("%s\n", n)
				}
			}
		} else {
			logError(err)
		}
	}
}

func listTokens(tokens jh.Tokens) {

	w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
	fmt.Fprintf(w, "%s\n", t.Title("ID\tKind\tCreated\tLast Activity\tClient\n"))
	for _, tk := range tokens.APITokens {
		fmt.Fprintf(w, "%s\n", t.Text("%s\t%s\t%s\t%s\t%s", tk.ID, tk.Kind, tk.Created, tk.LastActivity, tk.Note))
	}
	for _, tk := range tokens.OAuthTokens {
		fmt.Fprintf(w, "%s\n", t.Text("%s\t%s\t%s\t%s\t%s", tk.ID, tk.Kind, tk.Created, tk.LastActivity, tk.OAuthClient))
	}
	w.Flush()
}
