package cmd

import (
	"fmt"
	jh "github.com/jdrivas/sponde/jupyterhub"
	t "github.com/jdrivas/sponde/term"
	"github.com/juju/ansiterm"
	"github.com/spf13/viper"
	"os"
)

// Connection proxies for the jh.Connection
// so we can add some display functionality to it.
type Connection struct {
	*jh.Connection
}

// ConnectionList is, well, a list of connections
type ConnectionList []Connection

// Current and default connection state
var currentConnection *Connection
var lastConnection Connection

// var defaultConnection Connection

// List displpays the list of connections and notes the current one.
func (conns ConnectionList) List() {
	if len(conns) > 0 {
		currentName := getCurrentConnection().Name
		w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
		fmt.Fprintf(w, "%s\n", t.Title("\tName\tURL\tToken"))
		for _, c := range conns {
			name := t.Text(c.Name)
			current := ""
			if c.Name == currentName {
				name = t.Highlight("%s", c.Name)
				current = t.Highlight("%s", "*")
			}
			token := c.getSafeToken(false, true)
			fmt.Fprintf(w, "%s\t%s\t%s\n", current, name, t.Text("%s\t%s", c.HubURL, token))
		}
		w.Flush()
	} else {
		fmt.Printf("%s\n", t.Title("There were no connections."))
	}
}

// Deep copy a connection.
func (conn Connection) copy() Connection {
	// Copy the connection
	newConn := conn

	// Deep copy the pointer to jh.Conneciton (Connection.Auth isn't a pointer)
	c := *conn.Connection
	newConn.Connection = &c

	return newConn
}

func getAllConnections() ConnectionList {
	return getAllConnectionsFromConfig()
}

// GetCurrentConnection returns the current connection object for
// the JupyterhHub API
func getCurrentConnection() Connection {
	return currentConnection.copy()
}

// SetCurrentConneciton sets the connection
func setCurrentConnection(conn Connection) {
	if currentConnection != nil {
		if currentConnection.Name != updatedConnectionName {
			lastConnection = *currentConnection
		}
	}
	currentConnection = &conn
}

// Get a named connection
func getConnection(name string) (Connection, bool) {
	return getConnectionFromConfig(name)
}

//
// setConnection sets the connection to the named connection,
// or returns an error if the named connection doesn't exist.
// Connection is not set if there is an error, and whathever
//  connection is current will continued to be used.
func setConnection(name string) (err error) {
	conn, ok := getConnection(name)
	if ok {
		setCurrentConnection(conn)
	} else {
		err = fmt.Errorf("couldn't find connection \"%s\"", name)
	}
	return err
}

// updateCurrentConnection only makes changes to fields
// that are set in the argument
const updatedConnectionName = "Flag Updated Connection"

func updateCurrentConnection(conn Connection) {
	newConn := updateConnection(conn, getCurrentConnection())
	newConn.Name = updatedConnectionName
	setCurrentConnection(newConn)
}

func updateConnection(update, existing Connection) Connection {

	// this should be easier
	conn := existing
	if update.Name != "" {
		conn.Name = update.Name
	}
	if update.HubURL != "" {
		conn.HubURL = update.HubURL
	}
	if update.Token != "" {
		conn.Token = update.Token
	}
	if conn.Auth.ClientID != "" {
		conn.Auth.ClientID = update.Auth.ClientID
	}
	if conn.Auth.ClientSecret != "" {
		conn.Auth.ClientSecret = update.Auth.ClientSecret
	}
	if conn.Auth.RedirectURL != "" {
		conn.Auth.RedirectURL = update.Auth.RedirectURL
	}

	return conn
}

// In addition we allow creation of named connections
// in a configuration file that's managed by Viper.

// getSafeToken returns a token string managed by the showToken state.
// * If neverShowToken is set then always return "****" instead of a token
// * If showToken is set (and neverShoToken is not) then return the actual token, othewise "****"
// * If useEmpty then instad of "****" return ""
func (conn Connection) getSafeToken(useEmpty bool, useShowTokensOnce bool) (token string) {
	token = conn.Token
	if conn.Token == "" {
		token = "<enpty-token>"
	} else {
		token = "****"
		if useEmpty {
			token = ""
		}
		show := getShowTokens()
		if useShowTokensOnce {
			show = show || getShowTokensOnce()
		}
		if !viper.GetBool(neverShowTokensKey) && show {
			token = conn.Token
		}
	}
	return token
}

//
// Config file conneciton management
//

// YAML Variables which show up in viper, but managed here.
const (
	defaultConnectionNameKey   = "defaultConnection"
	defaultConnectionNameValue = "default"
	connectionsKey             = "connections"
	hubURLKey                  = "huburl"
	tokenKey                   = "token"
	neverShowTokensKey         = "neverShowTokens"
	authKey                    = "auth"
	clientIDKey                = "clientID"
	clientSecretKey            = "clientSecret"
	redirectURLKey             = "redirectURL"
)

// Read in the config to get all the named connections
func getAllConnectionsFromConfig() (conns []Connection) {
	connectionsMap := viper.GetStringMap(connectionsKey) // map[string]interface{}
	for name := range connectionsMap {
		conn, ok := getConnectionFromConfig(name)
		if ok {
			conns = append(conns, conn)
		} else {
			cmdError(fmt.Errorf("couldn't create a config for connection \"%s\"", name))
		}
	}
	return conns
}

// Get a connection from the config file.
func getConnectionFromConfig(name string) (conn Connection, ok bool) {
	connKey := fmt.Sprintf("%s.%s", connectionsKey, name)
	if viper.IsSet(connKey) {
		conn = Connection{
			&jh.Connection{
				Name:   name,
				HubURL: viper.GetString(fmt.Sprintf("%s.%s", connKey, hubURLKey)),
				Token:  viper.GetString(fmt.Sprintf("%s.%s", connKey, tokenKey)),
				Auth: jh.Auth{
					ClientSecret: viper.GetString(fmt.Sprintf("%s.%s.%s", connKey, authKey, clientSecretKey)),
					ClientID:     viper.GetString(fmt.Sprintf("%s.%s.%s", connKey, authKey, clientIDKey)),
					RedirectURL:  viper.GetString(fmt.Sprintf("%s.%s.%s", connKey, authKey, redirectURLKey)),
				},
			},
		}
		ok = true
	}
	return conn, ok
}

// initConnections sets up the first current Connection,
// initializes the ShowTokens state, and should be called whenever the Viper config file gets reloaded.
// Since we need at least a URL to break and/or let us know that no token has been set. Also, this value
// will remind that the usual port is 8081.
const defaultHubURL = "http://127.0.0.1:8081"

func initConnections() {

	// Current conenction should be durable during interactive mode
	// reset it to the default ...
	var conn Connection
	if currentConnection == nil {
		// If there is a connection named default, use it ....
		var ok bool
		conn, ok = getConnection(defaultConnectionNameValue)
		if !ok {
			// .. Otherwise, see if there is a _name_ of a defined connection to use as default ...
			defaultName := viper.GetString(defaultConnectionNameKey)
			if defaultName != "" {
				conn, ok = getConnection(defaultName)
				if !ok {
					// ... As a last resort set up a broken empty connection.
					// We won't panic here as we can set it during interactive
					// mode and it will otherwise error.
					conn = Connection{
						Connection: &jh.Connection{
							Name:   defaultConnectionNameValue,
							HubURL: defaultHubURL,
							Token:  "",
							Auth: jh.Auth{
								ClientID:     "",
								ClientSecret: "",
								RedirectURL:  "",
							},
						},
					}
				}
			}
		}
		lastConnection = conn
		setCurrentConnection(conn)
		// or if we've just changed it for one command, reset it to previous.
	} else if getCurrentConnection().Name == updatedConnectionName {
		conn = lastConnection
		setCurrentConnection(conn)
	}

	// This too, should be durable but we need to set it the first time.
	// Set this from the config, if it's not been set before.
	if showTokens == nil {
		showTokens = new(bool)
		setShowTokens(viper.GetBool(showTokensKey))
	}

}

//
// ShowTokens
//

// so we share the name for the key here. The actual variable show be obtained
// through config rather than in viper
const (
	showTokensKey = "showTokens"
)

var showTokens *bool
var showTokensOnceState bool

func toggleShowTokens() {
	*showTokens = !*showTokens
}

// State variable to show connection tokens on the next call
// this will reset to false once the connection has been displayed.
func setShowTokensOnce() {
	showTokensOnceState = true
}

// Reset the ShowTokensOnce state to false.
func resetShowTokensOnce() {
	showTokensOnceState = false
}

func getShowTokensOnce() bool {
	return showTokensOnceState
}

// Current values of the ShowTokensState
func getShowTokens() bool {
	return *showTokens
}

// Set the status of the ShowTokens state.
func setShowTokens(st bool) {
	*showTokens = st
}
