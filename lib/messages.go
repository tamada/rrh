package lib

import "io"

/*
MessageCenter manages log entries.
*/
type MessageCenter struct {
	messages []logEntry
}

type logEntry struct {
	message string
	level   LogLevel
}

/*
LogLevel represent the level for logging.
*/
type LogLevel int

/*
   VERBOSE shows verbose level of log level.
*/
const (
	VERBOSE LogLevel = iota + 1
	/*
	   INFO shows info level of log level.
	*/
	INFO
	/*
	   WARN shows warning level of log level.
	*/
	WARN
	/*
	   SEVERE shows severe level of log level.
	*/
	SEVERE
)

/*
NewMessageCenter creates an instance of MessageCenter.
*/
func NewMessageCenter() *MessageCenter {
	var center = new(MessageCenter)
	center.messages = []logEntry{}
	return center
}

/*
PushLog append given message to receiver MessageCenter as INFO level.
*/
func (mc *MessageCenter) PushLog(message string) {
	mc.Push(message, INFO)
}

/*
PushVerbose append given message to receiver MessageCenter as VERBOSE level.
*/
func (mc *MessageCenter) PushVerbose(message string) {
	mc.Push(message, VERBOSE)
}

/*
Push append given message to receiver MessageCenter as given level.
*/
func (mc *MessageCenter) Push(message string, level LogLevel) {
	mc.messages = append(mc.messages, logEntry{message: message, level: level})
}

/*
PrintVerbose prints messages of VERBOSE, INFO, WARN, and SEVERE log level
in the receiver MessageCenter to the given writer.
*/
func (mc *MessageCenter) PrintVerbose(out io.Writer) {
	mc.Print(out, VERBOSE)
}

/*
PrintLog prints messages of INFO, WARN, and SEVERE log level
in the receiver MessageCenter to the given writer.
*/
func (mc *MessageCenter) PrintLog(out io.Writer) {
	mc.Print(out, INFO)
}

/*
Print prints greater than given level messages
in the receiver MessageCenter to the given writer.
*/
func (mc *MessageCenter) Print(out io.Writer, level LogLevel) {
	var messages = mc.FindMessages(level)
	for _, msg := range messages {
		out.Write([]byte(msg))
		out.Write([]byte("\n"))
	}
}

/*
FindMessages returns messages in the receiver MessageCenter with
greater than given log level.
*/
func (mc *MessageCenter) FindMessages(level LogLevel) []string {
	var messages = []string{}
	for _, msg := range mc.messages {
		if msg.level >= level {
			messages = append(messages, msg.message)
		}
	}
	return messages
}
