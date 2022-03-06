package proto

import (
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Edgaru089/rmake/util"
)

var (
	totalMakes int64 // The total number of makes. must be accessed by sync/atomic!
)

type stderrWriter struct {
	msg    int16 // one of MsgStdout/MsgStderr/MsgStdmsg
	writer io.Writer
	mu     *sync.Mutex
}

var _ io.Writer = &stderrWriter{}

func (s *stderrWriter) Write(data []byte) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	err = util.WriteInt16(s.writer, s.msg)
	if err != nil {
		return
	}
	err = util.WriteInt64(s.writer, int64(len(data)))
	if err != nil {
		return
	}

	return s.writer.Write(data)
}

// Make processes and completes the make request.
//
// Make should be called after the server replied MsgMakeRequestOK,
// c.Request.Recipe is set, and the client should now be sending files.
//
// If err!=nil, the connection is closed.
// And if err is not one of ErrInvalidMsg/ErrInvalidParam,
// A MsgErrorMessage is sent.
func (c *Conn) Make() (err error) {
	r := c.Request

	makeid := atomic.AddInt64(&totalMakes, 1)

	// Make a temporary directory
	r.tempdir = filepath.Join(os.TempDir(), "rmaked_"+strconv.FormatInt(makeid, 10))
	err = os.Mkdir(r.tempdir, 0755)
	if err != nil {
		return err
	}
	defer os.RemoveAll(r.tempdir)

	// Read all the files.
readloop:
	for {
		msg, err := util.ReadInt16(c.Underlying)
		if err != nil {
			return err
		}

		switch msg {
		case MsgMakeSendFile:
			// Read the next file.
			filename, err := util.ReadString(c.Underlying)
			if err != nil {
				return err
			}

			// Read the size
			filesize, err := util.ReadInt64(c.Underlying)
			if err != nil {
				return err
			}

			// Open&Write the file
			f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				return err
			}
			_, err = io.CopyN(f, c.Underlying, filesize)
			f.Close()
			if err != nil {
				return err
			}

			log.Printf("[Make %d] Input [%s] %d bytes", makeid, filename, filesize)

		case MsgMakeSendFileEnd:
			break readloop
		default:
			util.WriteInt16(c.Underlying, MsgInvalidMsg)
			return ErrInvalidMsg
		}
	}

	// Guards the MsgStdout/MsgStderr stuff
	var muConn sync.Mutex

	// Now that we have all the files ready,
	// lets kick the make process start.
	for i := range r.Recipe.Command {

		// Build the command string
		icmdbuild := &strings.Builder{}
		err = r.Recipe.cmdtemplates[i].Execute(icmdbuild, r)
		if err != nil {
			return err
		}
		invokecmd := icmdbuild.String()
		log.Printf("[Make %d] Invoking: \"%s\"", makeid, invokecmd)

		// Send a message
		err = util.WriteString(c.Underlying, strings.Join([]string{"Executing: \"", invokecmd, "\""}, ""))
		if err != nil {
			return err
		}

		// Build the context and fields
		fields := strings.Fields(invokecmd)
		ctx, cancel := context.WithDeadline(
			context.Background(),
			time.Now().Add(time.Duration(r.Recipe.MaxTime)*time.Millisecond),
		)
		cmd := exec.CommandContext(ctx, fields[0], fields[1:]...)

		// Build the Stdout/Stderr pipes
		cmd.Stderr = &stderrWriter{msg: MsgStderr, mu: &muConn, writer: c.Underlying}
		cmd.Stdout = &stderrWriter{msg: MsgStdout, mu: &muConn, writer: c.Underlying}

		err = cmd.Run()
		cancel()
		if err != nil {
			return err
		}
	}

	// Grab the output files
	outputs, err := filepath.Glob(filepath.Join(r.tempdir, r.Recipe.OutputGlob))
	if err != nil {
		return err
	}
	for _, o := range outputs {
		stat, err := os.Stat(o)
		if err != nil {
			return err
		}
		filelen := stat.Size()

		file, err := os.Open(o)
		if err != nil {
			return err
		}

		util.WriteInt16(c.Underlying, MsgMakeOutput)
		util.WriteString(c.Underlying, o)
		util.WriteInt64(c.Underlying, filelen)
		_, err = io.CopyN(c.Underlying, file, filelen)
		file.Close()
		if err != nil {
			return err
		}
		log.Printf("[Make %d] Output [%s] %d bytes", makeid, o, filelen)
	}
	util.WriteInt16(c.Underlying, MsgMakeOutputEnd)

	return nil
}
