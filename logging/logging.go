package logging

import (
	"errors"
	"fmt"
	"log"
	"os"
)

const Reset = "\033[0m"
const Red = "\033[31m"
const Green = "\033[32m"
const Yellow = "\033[33m"
const Blue = "\033[34m"
const Purple = "\033[35m"
const Cyan = "\033[36m"
const Gray = "\033[37m"
const White = "\033[97m"
const BrightBlack = "\u001b[30;1m"
const BrightRed = "\u001b[31;1m"
const BrightGreen = "\u001b[32;1m"
const BrightYellow = "\u001b[33;1m"
const BrightBlue = "\u001b[34;1m"
const BrightMagenta = "\u001b[35;1m"
const BrightCyan = "\u001b[36;1m"
const BrightWhite = "\u001b[37;1m"

var Mylog log.Logger

// Initialise MyLog.
// Override standard log defaults.
// Write to LogFile instead of terminal.
func LogInit() {
	Mylog = *log.Default()

	f, err := os.OpenFile("./logging.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	Mylog.SetOutput(f)
}

// Log a message only to the log file
// Adds a colour that can be used to identify the caller easily.
//
//	startColour: selects a colour from utils.colour rendered using ANSI codes.
//	message: the message to display
//	returns: the formatted message for further reporting such as browser output
func TraceLog(startColour string, message string) string {
	formattedMessage := "[TRACE] " + startColour + message + Reset

	// log the message in the log file
	Mylog.Output(2, message)

	// pass the message back to the caller
	return formattedMessage
}

// convenience function to include a variable in a TraceLog message
//
//	startColour: the colour of the message
//	message: the message, which must contain one formatting symbol
//	details: the variable to be printed
func TraceLogf(startColour string, message string, details ...any) string {
	return TraceLog(startColour, fmt.Sprintf(message, details...))
}

// Log a message both to the log file and to the terminal.
// Adds a colour that can be used to identify the caller easily.
//
//	startColour: selects a colour from utils.colour rendered using ANSI codes.
//	message: the message to display
//	returns: the formatted message for further reporting such as browser output
func TraceInfo(startColour string, message string) string {
	formattedMessage := TraceLog(startColour, message)

	fmt.Println(formattedMessage) // send the message to the terminal

	return formattedMessage // pass the message back to the caller
}

// Log an error to the console and to the log file
func TraceError(message string) error {
	formattedMessage := "[ERROR] " + BrightRed + message + Reset

	// log the message in the log file
	Mylog.Output(2, message)

	// (Development only) send the message to the terminal
	fmt.Println(formattedMessage)

	// pass the message back to the caller
	return errors.New(formattedMessage)
}

// convenience function to include variables in an TraceError message
//
//	startColour: the colour of the message
//	message: the message, containing formatting symbols using fmt.Printf conventions
//	details: the variables to be printed
func TraceErrorf(message string, details ...any) error {
	return TraceError(fmt.Sprintf(message, details...))
}

// convenience function to include variables in a TraceInfo message
//
//	startColour: the colour of the message
//	message: the message, containing formatting symbols using fmt.Printf conventions
//	details: the variables to be printed
func TraceInfof(startColour string, message string, details ...any) string {
	return TraceInfo(startColour, fmt.Sprintf(message, details...))
}

// convenience function to create a component of a trace message.
// Does not add a prefix or newline, so composite messages can
// be constructed from it.
//
//	startColour: the colour of the message
//	message: the message, containing formatting symbols using fmt.Printf conventions
//	details: the variables to be printed
func TraceInfoPart(startColour string, message string, details ...any) string {
	formattedMessage := Reset + startColour + message + Reset
	return fmt.Sprintf(formattedMessage, details...)
}
