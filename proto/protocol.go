package proto

import (
	"errors"
	"io"
	"regexp"
	"text/template"
	"time"
)

const (
	// The default port for rmaked to listen on
	DefaultPort = 25386
)

var (
	ErrInvalidMsg   = errors.New("invalid message code")
	ErrInvalidParam = errors.New("invalid message parameters")
)

// Recipe contains instructions for what to do for a request.
type Recipe struct {
	Name     string // Short name for the recipe.
	FileGlob string // The recipe is selected automatically if the glob matches the file uploaded.

	Command    []string // What to execute in order in the working directory.
	OutputGlob string   // Matches output files to be sent to the client.

	MaxTime int // Maximum time, in milliseconds, for the task to finish.

	// private fields
	maxtime            time.Duration
	regfile, regoutput *regexp.Regexp
	cmdtemplates       []*template.Template
}

// Request is a make request.
// The request is fed to Recipe.Command via template.
type Request struct {
	Recipe   *Recipe           // The recipe used, nil on the client.
	Filename string            // The first filename.
	Type     string            // The recipe type, if specified.
	Option   map[string]string // Other options, if specified.

	tempdir string // The temporary directory used to build.
}

// Conn is a rmake connection.
type Conn struct {
	Underlying io.ReadWriteCloser // The underlying connection.

	Authed bool   // Whether or not the connection is authencated.
	User   string // The user name if authencated.

	Request *Request // The current request on the connection, or nil for none.
}
