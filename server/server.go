package main

import (
	"common"

	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/jroimartin/gocui"
)

var (
	clientCounter   uint
	MainGUI         *gocui.Gui
	Logs            []string
	LogFile         *os.File
	LogWriteMut     sync.Mutex
	receivedSignals chan os.Signal
)

const (
	DeleteLinesCode  = "\033[2K\r"
	MoveCursorUpCode = "\033[A"

	// Change this value to allow more or less users on the server.
	MAX_CLIENTS uint = 10
)

/*
func init() {
	// disable input buffering
	//exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	//exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	var prev2 byte
	var prev byte
	var b [1]byte
	for {
		os.Stdin.Read(b[:])
		//fmt.Println(prev2, prev, b[0])
		switch prev {
		case 91:
			if prev2 == 27 {
				switch b[0] {
				case 65:
					fmt.Println("\nArrow key UP")
				case 66:
					fmt.Println("\nArrow key DOWN")
				case 67:
					fmt.Println("\nArrow key RIGHT")
				case 68:
					fmt.Println("\nArrow key LEFT")
				default:
					fmt.Println("\nUnknown :", b[0])
				}
			} else {
				fmt.Println("\nI got the byte", b, "("+string(b[:])+")")
			}
		case 194:
			fmt.Println("\nI got the 8-bit character :", string(rune(b[0])))
		case 195:
			fmt.Println("\nI got the 8-bit composed character :", string(rune(prev)+rune(b[0])-131))
		default:
			if b[0] == 9 || b[0] == 10 || (b[0] >= 32 && b[0] <= 127) && !(prev == 27 && b[0] == 91) {
				fmt.Println("\nI got the byte", b, "("+string(b[:])+")")
			}
		}
		prev2 = prev
		prev = b[0]
	}
}
*/

func main() {
	port := "8080"

	if len(os.Args) == 2 {
		if !common.IsNumeric(os.Args[1]) {
			fmt.Fprintln(os.Stderr, "Usage: go run . [<port>]")
			os.Exit(1)
		}

		port = os.Args[1]
	}

	clear := exec.Command("clear")
	clear.Stdout = os.Stdout
	clear.Run()

	localIp, errIP := getLocalIP()
	if errIP != nil {
		fmt.Fprintln(os.Stderr, "An unexpected error has occured retrieving local IP:", errIP.Error())
		os.Exit(1)
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()

	receivedSignals = make(chan os.Signal)

	signal.Notify(receivedSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-receivedSignals

		switch sig {
		case syscall.SIGINT:
			fmt.Print("\nCtrl+C detected. Do you want to close the server?\nThis will disconnect all users.\n> ")

			var conf string
			fmt.Scanf("%s", &conf)

			switch strings.ToLower(conf) {
			case "yes", "ye", "y", "yeah", "ya", "ja", "da", "oui":
				fmt.Println("Alright, closing all connections...")
				for _, client := range ClientsMap {
					SendMessage(client.Connection, "\rThe server has been closed, your connection has been terminated. Press a key to close nc...\n")
					client.Connection.Close()
				}
				fmt.Println("Done. Bye!")
				os.Exit(0)

			default:
				fmt.Println("Aborting interruption.")
				return
			}
		case syscall.SIGTERM:
			fmt.Println("Process terminated. Closing all connections...")
			for _, client := range ClientsMap {
				SendMessage(client.Connection, "\rThe server has been closed, your connection has been terminated. Press a key to close nc...\n")
				client.Connection.Close()
			}
			fmt.Println("Done. Bye!")
			os.Exit(0)
		}
	}()

	var errFile error
	Logs = make([]string, 0)

	LogFile, errFile = os.Create("log" + time.Now().Format("2006-01-02_15:04:05") + ".txt")
	if errFile != nil {
		fmt.Fprintf(os.Stderr, "Could not create log file:\n%v\n", errFile)
		os.Exit(1)
	}

	// TODO : Guh??
	// MainGUI, err := gocui.NewGui(gocui.Output256)
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, colors.SprintfANSI(err.Error(), colors.FgRed, colors.BgReset))
	// }
	// defer MainGUI.Close()

	fmt.Printf("Chat server started on %s. Listening on %s\n", localIp, port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			continue
		}
		if clientCounter == MAX_CLIENTS {
			SendMessage(conn, fmt.Sprintf("The maximum amount of users (%d) has been reached. Please try connecting again later.\nPress [ENTER] to close nc...", MAX_CLIENTS))
			conn.Close()
			continue
		}

		clientCounter++
		go HandleClient(conn)
	}
}

func getLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	return strings.Split(conn.LocalAddr().String(), ":")[0], nil
}
