package jupyterhub

import (
	"fmt"
	"os"
	"text/tabwriter"
)

type Users []User

type User struct {
	Name         string            `json:"name"`
	Kind         string            `json:"kind"`
	Admin        bool              `json:"admin"`
	Created      string            `json:"created"`
	LastActivity string            `json:"last_activity"`
	ServerURL    string            `json:"server"`
	Pending      string            `json:"pending"`
	Servers      map[string]Server `json:"servers"`
}

type Server struct {
	Name         string      `json:"name"`
	LastActivity string      `json:"last_activity"`
	Started      string      `json:"started"`
	Pending      string      `json:"pending"`
	Ready        bool        `json:"ready"`
	State        StateValues `json:"state"`
	URL          string      `json:"url"`
	ProgressURL  string      `json:"progress_url"`
}

type StateValues struct {
	PodName string `json:"progress_url"`
}

// GetUser retruns a users information
func GetUser(username string) (user User, err error) {
	err = Get(fmt.Sprintf("%s%s", "/users/", username), &user)
	return user, err
}

// GetUsers returns a list of logged in JupyterHub users.
func GetUsers() (users Users, err error) {
	err = Get("/users", &users)
	return users, err
}

func ListUsers(users Users) {
	w := tabwriter.NewWriter(os.Stdout, 4, 4, 3, ' ', 0)
	fmt.Fprintf(w, "Name\tAdmin\tCreated\tServer\tLast\n")
	for _, u := range users {
		fmt.Fprintf(w, "%s\t%t\t%s\t%s\t%s\n", u.Name, u.Admin, u.Created, u.ServerURL, u.LastActivity)
	}
	w.Flush()
}

func DescribeUsers(users Users) {
	for _, u := range users {
		w := tabwriter.NewWriter(os.Stdout, 4, 4, 3, ' ', 0)
		fmt.Fprintf(w, "Name\tKind\tAdmin\tServer\n")
		fmt.Fprintf(w, "%s\t%s\t%t\t%s\n", u.Name, u.Kind, u.Admin, u.ServerURL)
		w.Flush()
		fmt.Println()
		w = tabwriter.NewWriter(os.Stdout, 4, 4, 3, ' ', 0)
		fmt.Fprintf(w, "Created\tLast Activity\tPending\n")
		pending := "<empty>"
		if u.Pending != "" {
			pending = u.Pending
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", u.Created, u.LastActivity, pending)
		w.Flush()
		fmt.Printf("\nServers\n")
		for _, s := range u.Servers {
			w = tabwriter.NewWriter(os.Stdout, 4, 4, 3, ' ', 0)
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
}
