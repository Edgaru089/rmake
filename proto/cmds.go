package proto

const (
	HashTypeSHA256 int16 = iota
)

const (
	MsgOK              int16 = iota // OK.
	MsgInvalidMsg                   // Invalid message type. Right after this the connection is closed.
	MsgInvalidParam                 // Invalid parameters. Right after this the connection is closed.
	MsgAuthFailed                   // Wrong password / Not authenticated. Right after this the connection is closed.
	MsgUnsupportedHash              // Unsupported hash type.

	MsgErrorMessage // A error with a string message. After this the connection is closed.
)

const (
	// Client sends MsgAuthenticate to authenticate themselves.
	// Salt is joined right after the real password before hashing.
	//
	// [MsgAuthenticate] [HashType int16] [Username string] [Hash []byte]
	MsgAuthenticate int16 = 100 + iota

	// Client sends MsgMakeRequest to request a make.
	// After a MsgMakeRequestOK the client sends files via MsgMakeSendFile and MsgMakeSendFileEnd.
	// and waits for MsgErrorMessage/MsgMakeOutput and MsgMakeOutputEnd.
	//
	// For now the Filename count should be 1.
	// [MsgMakeRequest] [Type string] [FilenameCount int32] [... [Filename string] .. ]
	// [OptionCount int32] [... [OptionKey string] [OptionValue string] ...]
	MsgMakeRequest

	// Server sends MsgMakeRequestOK after the request has been accepted.
	// [MsgMakeRequestOK] [RecipeType string]
	MsgMakeRequestOK
	// Server sends MsgMakeRequestNoRecipe if no matching recipe can be found.
	MsgMakeRequstNoRecipe

	// Client sends MsgMakeSendFile to send a file.
	//
	// [MsgMakeSendFile] [Name string] [Data []byte]
	MsgMakeSendFile
	// Client sends MsgMakeSendFileEnd after it has sent every file in the MakeRequest.
	MsgMakeSendFileEnd

	// Server sends MsgMakeOutput after making the recipe.
	//
	// [MsgMakeOutput] [Name string] [Data []byte]
	MsgMakeOutput
	// Server sends MsgMakeOutputEnd after it has sent every file in the output.
	//
	// After this message the make session is considered complete.
	MsgMakeOutputEnd

	// MsgGoodbye can be sent at any time by any side to close the connection.
	MsgGoodbye
)

// Messages.
// All of them are followed by a single [string].
const (
	// MsgStdmsg carries a message from the rmake daemon to the client.
	//
	// Usually it indicates the next command to be executed.
	MsgStdmsg int16 = 200 + iota
	// MsgStdout carries a stdout message by one of the Commands.
	MsgStdout
	// MsgStderr carries a stderr message by one of the Commands.
	MsgStderr
)
