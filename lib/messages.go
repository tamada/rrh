package lib

import "io"

type MessageCenter struct {
	messages []logEntry
}

type logEntry struct {
	message string
	level   LogLevel
}

type LogLevel int

const (
	VERBOSE LogLevel = iota + 1
	INFO
	WARN
	SEVERE
)

func NewMessageCenter() *MessageCenter {
	var center = new(MessageCenter)
	center.messages = []logEntry{}
	return center
}

func (mc *MessageCenter) PushLog(message string) {
	mc.Push(message, INFO)
}

func (mc *MessageCenter) PushVerbose(message string) {
	mc.Push(message, VERBOSE)
}

func (mc *MessageCenter) Push(message string, level LogLevel) {
	mc.messages = append(mc.messages, logEntry{message: message, level: level})
}

func (mc *MessageCenter) PrintVerbose(out io.Writer) {
	mc.Print(out, VERBOSE)
}

func (mc *MessageCenter) Print(out io.Writer, level LogLevel) {
	var messages = mc.FindMessages(level)
	for _, msg := range messages {
		out.Write([]byte(msg))
		out.Write([]byte("\n"))
	}
}

func (mc *MessageCenter) PrintAll(out io.Writer) {
	mc.Print(out, INFO)
}

func (mc *MessageCenter) FindMessages(level LogLevel) []string {
	var messages = []string{}
	for _, msg := range mc.messages {
		if msg.level >= level {
			messages = append(messages, msg.message)
		}
	}
	return messages
}
