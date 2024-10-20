package main

import (
	"bufio"
	"common"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"colors"
	"netcat_user"
)

var (
	ClientsMap    map[uint]*netcat_user.User
	uncoloredMsgs chan string
)

const (
	joinLeaveMsgColor = colors.FgBrightCyan
	serverMsgColor    = colors.FgBrightYellow
	clientPromptColor = colors.FgBrightGreen
	errorColor        = colors.FgBrightRed

	colorInstructions = "To change your color, you may specify it\n- In English;\n- by its R(ed)G(reen)B(lue) values, each from 0 to 255 and separated like so: '127-127-127'.\nNOTE: You can also use hexadecimal values for RGB coloring: '#00ff00' gives bright green.\n"
	nameInstructions  = "You are free to choose any name you wish.\nThe only rule being characters should only belong to alphanumeric list: letters and numbers.\n"
	errorReading      = "An error occured while reading your prompt, and your connection was terminated to prevent corruption on both parts. Sorry for the inconvenience.\n(Press Return to exit...)"
)

func init() {
	ClientsMap = make(map[uint]*netcat_user.User, 0)
	uncoloredMsgs = make(chan string)
}

func HandleClient(conn net.Conn) {
	defer conn.Close()

	handledClient := netcat_user.NewUser(conn, "", colors.NewRandomFGColorRGB(), getNextAvailableID())
	ClientsMap[handledClient.Id] = handledClient

	go saveToLogs()

	fmt.Println(colors.SprintfANSI(fmt.Sprintf("New client connected: %s", handledClient.Connection.RemoteAddr().String()), serverMsgColor, colors.BgReset))
	uncoloredMsgs <- fmt.Sprintf("[%s] New client connected: %s\n", time.Now().Format("2006-01-02 15:04:05"), handledClient.Connection.RemoteAddr().String())
	displayWelcomeMessage(handledClient.Connection)

	reader := bufio.NewReader(handledClient.Connection)
	ClientsMap[handledClient.Id].Reader = reader

	for {
		for handledClient.Name == "" {
			SendMessage(handledClient.Connection, colors.SprintfANSI("[ENTER YOUR NAME]: ", clientPromptColor, colors.BgReset))
			username, errRName := ClientsMap[handledClient.Id].Reader.ReadString('\n')
			if errRName != nil {
				msg := colors.SprintfANSI(fmt.Sprintf("%s has left the server.", handledClient.Connection.RemoteAddr().String()), serverMsgColor, colors.BgReset)
				fmt.Println(colors.SprintfANSI("["+time.Now().Format("2006-01-02 15:04:05")+"] ", joinLeaveMsgColor, colors.BgReset) + msg)
				uncoloredMsgs <- fmt.Sprintf("[%s] %s has left the server.\n", time.Now().Format("2006-01-02 15:04:05"), handledClient.Connection.RemoteAddr().String())
				broadcastMessage(msg, handledClient.Connection, true)
				SendMessage(ClientsMap[handledClient.Id].Connection, colors.SprintfANSI(errorReading, errorColor, colors.BgReset))
				ClientsMap[handledClient.Id].Connection.Close()
				delete(ClientsMap, handledClient.Id)
				return
			}

			if ok, errName := netcat_user.IsValidUsername(strings.TrimSpace(username), ClientsMap); !ok {
				SendMessage(handledClient.Connection, "ERROR: Username "+errName.Error()+"\n")
				continue
			}

			username = strings.TrimSpace(username)
			ClientsMap[handledClient.Id].Name = username
			broadcastMessage(colors.SprintfANSI("["+time.Now().Format("2006-01-02 15:04:05")+"] ", joinLeaveMsgColor, colors.BgReset)+handledClient.ColoredUsername()+colors.SprintfANSI(" has joined the chat.", joinLeaveMsgColor, colors.BgReset), ClientsMap[handledClient.Id].Connection, true)
			uncoloredMsgs <- fmt.Sprintf("[%s] %s has joined the chat.\n", time.Now().Format("2006-01-02 15:04:05"), ClientsMap[handledClient.Id].Name)
			SendMessage(ClientsMap[handledClient.Id].Connection, colors.SprintfANSI("Welcome to the chat, ", colors.FgBrightYellow, colors.BgReset)+handledClient.ColoredUsername()+colors.SprintfANSI("!\n", colors.FgBrightYellow, colors.BgReset))

			for _, entry := range Logs[:len(Logs)-1] {
				if common.ContainsIP(entry) {
					continue
				}
				SendMessage(ClientsMap[handledClient.Id].Connection, entry)
			}
		}

		SendMessage(ClientsMap[handledClient.Id].Connection, "> ")

		msg, errRMsg := ClientsMap[handledClient.Id].Reader.ReadString('\n')
		if errRMsg != nil {
			broadcastMessage(colors.SprintfANSI("["+time.Now().Format("2006-01-02 15:04:05")+"] ", joinLeaveMsgColor, colors.BgReset)+handledClient.ColoredUsername()+colors.SprintfANSI(" has left the chat.", joinLeaveMsgColor, colors.BgReset), ClientsMap[handledClient.Id].Connection, true)
			uncoloredMsgs <- fmt.Sprintf("[%s] %s has left the chat.\n", time.Now().Format("2006-01-02 15:04:05"), ClientsMap[handledClient.Id].Name)
			SendMessage(ClientsMap[handledClient.Id].Connection, colors.SprintfANSI(errorReading, errorColor, colors.BgReset))
			ClientsMap[handledClient.Id].Connection.Close()
			delete(ClientsMap, handledClient.Id)
			break
		}

		msg = strings.TrimSpace(msg)

		if msg != "" {
			if msg[0] == '/' {
				switch msg[1:] {
				case "quit", "exit", "leave":
					switch confirmation(handledClient, "Are you sure you want to leave the server?", ClientsMap[handledClient.Id].Reader) {
					case -1:
						return

					case 0:
						SendMessage(ClientsMap[handledClient.Id].Connection, colors.SprintfANSI("Aborting.\n", serverMsgColor, colors.BgReset))
						continue

					case 1:
						broadcastMessage(colors.SprintfANSI("["+time.Now().Format("2006-01-02 15:04:05")+"] ", joinLeaveMsgColor, colors.BgReset)+handledClient.ColoredUsername()+colors.SprintfANSI(" has left the chat.", joinLeaveMsgColor, colors.BgReset), ClientsMap[handledClient.Id].Connection, true)
						uncoloredMsgs <- fmt.Sprintf("[%s] %s has left the chat.\n", time.Now().Format("2006-01-02 15:04:05"), ClientsMap[handledClient.Id].Name)
						SendMessage(ClientsMap[handledClient.Id].Connection, colors.SprintfANSI("Your connection has been closed. Bye!\n(Press Return to exit...)", serverMsgColor, colors.BgReset))
						ClientsMap[handledClient.Id].Connection.Close()
						delete(ClientsMap, handledClient.Id)
						return

					case 2:
						SendMessage(ClientsMap[handledClient.Id].Connection, colors.SprintfANSI("Wat\n", serverMsgColor, colors.BgReset))
						continue
					}

				case "color":
					colorPrompt := colors.SprintfANSI(colorInstructions, serverMsgColor, colors.BgReset) + colors.SprintfANSI("[ENTER THE COLOR YOU WANT]: ", clientPromptColor, colors.BgReset)
					SendMessage(ClientsMap[handledClient.Id].Connection, colorPrompt)

					ok := false
					newColor := ""
					var errColor error

					for !ok {
						req, errRColor := ClientsMap[handledClient.Id].Reader.ReadString('\n')
						req = strings.TrimSpace(req)
						if errRColor != nil {
							broadcastMessage(colors.SprintfANSI("["+time.Now().Format("2006-01-02 15:04:05")+"] ", joinLeaveMsgColor, colors.BgReset)+handledClient.ColoredUsername()+colors.SprintfANSI(" has left the chat.", joinLeaveMsgColor, colors.BgReset), ClientsMap[handledClient.Id].Connection, false)
							uncoloredMsgs <- fmt.Sprintf("[%s] %s has left the chat.\n", time.Now().Format("2006-01-02 15:04:05"), ClientsMap[handledClient.Id].Name)
							SendMessage(ClientsMap[handledClient.Id].Connection, colors.SprintfANSI(errorReading, errorColor, colors.BgReset))
							ClientsMap[handledClient.Id].Connection.Close()
							delete(ClientsMap, handledClient.Id)
							return
						}

						newColor, errColor = netcat_user.StrToColor(req)
						if errColor != nil {
							SendMessage(handledClient.Connection, colors.SprintfANSI("Error: "+errColor.Error()+"\n", errorColor, colors.BgReset)+colorPrompt)
							continue
						}

						switch confirmation(handledClient, colors.SprintfANSI("Do you want to use", serverMsgColor, colors.BgReset)+fmt.Sprintf(" %sthis color%s", newColor, colors.ResetColorsTag)+colors.SprintfANSI("?", serverMsgColor, colors.BgReset), ClientsMap[handledClient.Id].Reader) {
						case -1:
							return

						case 0:
							SendMessage(ClientsMap[handledClient.Id].Connection, colors.SprintfANSI("Aborting.\n", serverMsgColor, colors.BgReset))
							ok = true

						case 1:
							oldColor := handledClient.Color
							handledClient.Color = newColor
							ClientsMap[handledClient.Id] = handledClient
							broadcastMessage(fmt.Sprintf("%s%s%d", oldColor, handledClient.Name, colors.FgReset)+colors.SprintfANSI(" changed their color to ", serverMsgColor, colors.BgReset)+handledClient.ColoredUsername(), ClientsMap[handledClient.Id].Connection, false)
							ok = true

						case 2:
							SendMessage(ClientsMap[handledClient.Id].Connection, colors.SprintfANSI("Wat\n", serverMsgColor, colors.BgReset))
							ok = true
						}
					}

				case "rename":
					namePrompt := colors.SprintfANSI(nameInstructions, serverMsgColor, colors.BgReset) + colors.SprintfANSI("[ENTER YOUR NEW NAME]: ", clientPromptColor, colors.BgReset)
					SendMessage(ClientsMap[handledClient.Id].Connection, namePrompt)

					var req string
					ok := false
					var errName, errRName error

					for !ok {
						req, errRName = ClientsMap[handledClient.Id].Reader.ReadString('\n')
						req = strings.TrimSpace(req)
						if errRName != nil {
							broadcastMessage(colors.SprintfANSI("["+time.Now().Format("2006-01-02 15:04:05")+"] ", joinLeaveMsgColor, colors.BgReset)+handledClient.ColoredUsername()+colors.SprintfANSI(" has left the chat.", joinLeaveMsgColor, colors.BgReset), ClientsMap[handledClient.Id].Connection, false)
							uncoloredMsgs <- fmt.Sprintf("[%s] %s has left the chat.\n", time.Now().Format("2006-01-02 15:04:05"), handledClient.Name)
							SendMessage(ClientsMap[handledClient.Id].Connection, colors.SprintfANSI(errorReading, errorColor, colors.BgReset))
							ClientsMap[handledClient.Id].Connection.Close()
							delete(ClientsMap, handledClient.Id)
							return
						}

						ok, errName = netcat_user.IsValidUsername(req, ClientsMap)
						if !ok {
							SendMessage(handledClient.Connection, colors.SprintfANSI("Error: "+errName.Error()+"\n", errorColor, colors.BgReset)+namePrompt)
							continue
						}

						ok = false

						switch confirmation(handledClient, colors.SprintfANSI("Are you sure you want to be named ", serverMsgColor, colors.BgReset)+fmt.Sprintf("%s%s%s", handledClient.Color, req, colors.ResetColorsTag)+colors.SprintfANSI("?", serverMsgColor, colors.BgReset), ClientsMap[handledClient.Id].Reader) {
						case -1:
							return

						case 0:
							SendMessage(ClientsMap[handledClient.Id].Connection, colors.SprintfANSI("Aborting.\n", serverMsgColor, colors.BgReset))
							ok = true

						case 1:
							oldName := handledClient.Name
							handledClient.Name = req
							ClientsMap[handledClient.Id] = handledClient
							broadcastMessage(fmt.Sprintf("%s%s%s", ClientsMap[handledClient.Id].Color, oldName, colors.ResetColorsTag)+colors.SprintfANSI(" is now known as ", serverMsgColor, colors.BgReset)+handledClient.ColoredUsername(), ClientsMap[handledClient.Id].Connection, false)
							ok = true

						case 2:
							SendMessage(ClientsMap[handledClient.Id].Connection, colors.SprintfANSI("Wat\n", serverMsgColor, colors.BgReset))
							ok = true
						}
					}

				default:
					broadcastMessage("["+time.Now().Format("2006-01-02 15:04:05")+"]["+handledClient.ColoredUsername()+"]: "+colors.SprintfANSI(msg, colors.FgBrightWhite, colors.BgReset), ClientsMap[handledClient.Id].Connection, true)
					uncoloredMsgs <- fmt.Sprintf("[%s][%s]: %s\n", time.Now().Format("2006-01-02 15:04:05"), ClientsMap[handledClient.Id].Name, msg)
				}

			} else {
				broadcastMessage("["+time.Now().Format("2006-01-02 15:04:05")+"]["+handledClient.ColoredUsername()+"]: "+colors.SprintfANSI(msg, colors.FgBrightWhite, colors.BgReset), ClientsMap[handledClient.Id].Connection, true)
				uncoloredMsgs <- fmt.Sprintf("[%s][%s]: %s\n", time.Now().Format("2006-01-02 15:04:05"), ClientsMap[handledClient.Id].Name, msg)
			}
		}
	}
}

func broadcastMessage(message string, senderConn net.Conn, deleteLine bool) {
	fmt.Println(message)
	for _, client := range ClientsMap {
		formattedMsg := ""

		if deleteLine {
			formattedMsg = DeleteLinesCode

			if client.Connection == senderConn {
				formattedMsg += MoveCursorUpCode
			} else {
				var saved []byte
				client.Reader.Read(saved)
				fmt.Printf("%s's saved message: %s\n", client.Name, saved)
			}
		}

		formattedMsg += message + "\n"
		if client.Connection != senderConn {
			formattedMsg += "> "
		}

		SendMessage(client.Connection, formattedMsg)
	}
}

func SendMessage(conn net.Conn, message string) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error sending message:", err.Error())
	}
}

func confirmation(client *netcat_user.User, question string, answer *bufio.Reader) int {
	SendMessage(ClientsMap[client.Id].Connection, question+"\n"+colors.SprintfANSI("[YOUR ANSWER]: ", clientPromptColor, colors.BgReset))
	var conf string
	var errConf error

	conf, errConf = answer.ReadString('\n')
	if errConf != nil {
		broadcastMessage(colors.SprintfANSI("["+time.Now().Format("2006-01-02 15:04:05")+"] ", joinLeaveMsgColor, colors.BgReset)+client.ColoredUsername()+colors.SprintfANSI(" has left the chat.", joinLeaveMsgColor, colors.BgReset), ClientsMap[client.Id].Connection, false)
		uncoloredMsgs <- fmt.Sprintf("[%s] %s has left the chat.\n", time.Now().Format("2006-01-02 15:04:05"), ClientsMap[client.Id].Name)
		SendMessage(ClientsMap[client.Id].Connection, colors.SprintfANSI(errorReading, errorColor, colors.BgReset))
		ClientsMap[client.Id].Connection.Close()
		delete(ClientsMap, client.Id)
		return -1
	}

	switch strings.ToLower(strings.TrimSpace(conf)) {
	case "yes", "ye", "y", "yeah", "yup", "yuh", "yuh uh", "ya", "ja", "da", "oui", "affirmative":
		return 1

	case "no", "na", "n", "nah", "nu", "nuh", "nuh uh", "nope", "nop", "nein", "nyet", "non", "negatory":
		return 0

	default:
		return 2
	}
}

/*
func removeClient(arr []netcat_user.User, toRemove *netcat_user.User) []netcat_user.User {
	newArr := make([]netcat_user.User, 0)

	for _, client := range arr {
		if !(&client == toRemove) {
			newArr = append(newArr, client)
		} else {
			client.Connection.Close()
		}
	}

	return newArr
}
*/

func saveToLogs() {
	for msg := range uncoloredMsgs {
		Logs = append(Logs, msg)

		LogWriteMut.Lock()
		n, errW := LogFile.WriteString(msg)
		if errW != nil {
			fmt.Fprintln(os.Stderr, colors.SprintfANSI(fmt.Sprintf("could not write to %s:\n%v", LogFile.Name(), errW), colors.FgRed, colors.BgReset))
		}
		LogWriteMut.Unlock()

		if n < len(msg) {
			fmt.Fprintln(os.Stderr, colors.SprintfANSI("bytes written lesser than message's length", colors.FgRed, colors.BgReset))
		}
	}
}

func displayWelcomeMessage(conn net.Conn) {
	SendMessage(conn, " 	 "+colors.SprintfANSI("________", colors.FgCyan, colors.BgReset)+"\n"+
		"	"+colors.SprintfANSI("|__   __/   _   _   _   _   _   _ __    _____  ____", colors.FgCyan, colors.BgReset)+"\n"+
		"       "+colors.SprintfANSI("___", colors.FgBrightYellow, colors.BgReset)+" "+colors.SprintfANSI("| | | | | | | | | | | \\ | | | |\\ \\  |  __/ | |\\ \\", colors.FgCyan, colors.BgReset)+"  "+colors.SprintfANSI("_", colors.FgBrightYellow, colors.BgReset)+"\n"+
		"      "+colors.SprintfANSI("|", colors.FgBrightYellow, colors.BgReset)+"    "+colors.SprintfANSI("| | | |_| | | | | | |  \\| | | | | | | |_   | | | |", colors.FgCyan, colors.BgReset)+"  "+colors.SprintfANSI("|", colors.FgBrightYellow, colors.BgReset)+"\n"+
		"      "+colors.SprintfANSI("|", colors.FgBrightYellow, colors.BgReset)+"    "+colors.SprintfANSI("| | |  _  | | | | | |     | | | | | |  _/  | |-,<", colors.FgCyan, colors.BgReset)+"   "+colors.SprintfANSI("|", colors.FgBrightYellow, colors.BgReset)+"\n"+
		"       "+colors.SprintfANSI("\\", colors.FgBrightYellow, colors.BgReset)+"   "+colors.SprintfANSI("| | | | | | | |_| | | |\\  | | |_| | | |___ | | | |", colors.FgCyan, colors.BgReset)+" "+colors.SprintfANSI("/", colors.FgBrightYellow, colors.BgReset)+"\n"+
		"	"+colors.SprintfANSI("\\", colors.FgBrightYellow, colors.BgReset)+"  "+colors.SprintfANSI("|/  |/  |/   \\__/\\| |/  \\/  |__/|/  |____/ |/  |/", colors.FgCyan, colors.BgReset)+" "+colors.SprintfANSI("/", colors.FgBrightYellow, colors.BgReset)+"\n"+
		"	 "+colors.SprintfANSI("\\", colors.FgBrightYellow, colors.BgReset)+"       "+colors.SprintfANSI("_______", colors.FgCyan, colors.BgReset)+"                                    "+colors.SprintfANSI("/", colors.FgBrightYellow, colors.BgReset)+"\n"+
		"	  "+colors.SprintfANSI("\\", colors.FgBrightYellow, colors.BgReset)+"     "+colors.SprintfANSI("|_   __/_____   ____    _____  _____", colors.FgCyan, colors.BgReset)+"       "+colors.SprintfANSI("/", colors.FgBrightYellow, colors.BgReset)+"\n"+
		"	   "+colors.SprintfANSI("\\", colors.FgBrightYellow, colors.BgReset)+"      "+colors.SprintfANSI("| |__|  /| | | |\\ \\  |  __/ |  __/", colors.FgCyan, colors.BgReset)+"      "+colors.SprintfANSI("/", colors.FgBrightYellow, colors.BgReset)+"\n"+
		"	    "+colors.SprintfANSI("\\", colors.FgBrightYellow, colors.BgReset)+"     "+colors.SprintfANSI("|  _/| | | | | | | | | |    | |__", colors.FgCyan, colors.BgReset)+"      "+colors.SprintfANSI("/", colors.FgBrightYellow, colors.BgReset)+"\n"+
		"	     "+colors.SprintfANSI("\\", colors.FgBrightYellow, colors.BgReset)+"    "+colors.SprintfANSI("| |  | | | | | |-,<  | |    |  _/", colors.FgCyan, colors.BgReset)+"     "+colors.SprintfANSI("/", colors.FgBrightYellow, colors.BgReset)+"\n"+
		"	      "+colors.SprintfANSI("\\", colors.FgBrightYellow, colors.BgReset)+"   "+colors.SprintfANSI("| |  | |_| | | | | | | |___ | |___", colors.FgCyan, colors.BgReset)+"   "+colors.SprintfANSI("/", colors.FgBrightYellow, colors.BgReset)+"\n"+
		"	       "+colors.SprintfANSI("\\", colors.FgBrightYellow, colors.BgReset)+"  "+colors.SprintfANSI("|/   |__/ \\| |/  |/  |____/ |____/", colors.FgCyan, colors.BgReset)+"  "+colors.SprintfANSI("/", colors.FgBrightYellow, colors.BgReset)+"\n"+
		"		"+colors.SprintfANSI("\\", colors.FgBrightYellow, colors.BgReset)+"          "+colors.SprintfANSI("___   __      _", colors.FgBrightRed, colors.BgReset)+"           "+colors.SprintfANSI("/", colors.FgBrightYellow, colors.BgReset)+"\n"+
		"		 "+colors.SprintfANSI("\\", colors.FgBrightYellow, colors.BgReset)+"        "+colors.SprintfANSI("\\   / \\  /    \\ /", colors.FgBrightRed, colors.BgReset)+"         "+colors.SprintfANSI("/", colors.FgBrightYellow, colors.BgReset)+"\n"+
		"		  "+colors.SprintfANSI("\\", colors.FgBrightYellow, colors.BgReset)+"        "+colors.SprintfANSI("| |   \\ \\    //", colors.FgBrightRed, colors.BgReset)+"         "+colors.SprintfANSI("/", colors.FgBrightYellow, colors.BgReset)+"\n"+
		"		   "+colors.SprintfANSI("\\______", colors.FgBrightYellow, colors.BgReset)+" "+colors.SprintfANSI("| |", colors.FgBrightRed, colors.BgReset)+" "+colors.SprintfANSI("__", colors.FgBrightYellow, colors.BgReset)+" "+colors.SprintfANSI("\\ \\  //", colors.FgBrightRed, colors.BgReset)+" "+colors.SprintfANSI("________/", colors.FgBrightYellow, colors.BgReset)+"\n"+
		"			   "+colors.SprintfANSI("| |     \\ \\//", colors.FgBrightRed, colors.BgReset)+"\n"+
		"			  "+colors.SprintfANSI("/___\\     \\_/", colors.FgBrightRed, colors.BgReset)+"\n")
}

func getNextAvailableID() uint {
	for i := range MAX_CLIENTS {
		if _, ok := ClientsMap[i]; !ok {
			return i
		}
	}

	return MAX_CLIENTS + 1
}
