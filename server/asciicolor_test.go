package main

import (
	"fmt"
	"testing"

	"colors"
)

func TestAscii(t *testing.T) {
	fmt.Println("	 " + colors.SprintfANSI("________", colors.FgCyan, colors.BgReset) + "\n" +
		"	" + colors.SprintfANSI("|__   __/   _   _   _   _   _   _ __    _____  ____", colors.FgCyan, colors.BgReset) + "\n" +
		"       " + colors.SprintfANSI("___", colors.FgYellow, colors.BgReset) + " " + colors.SprintfANSI("| | | | | | | | | | | \\ | | | |\\ \\  |  __/ | |\\ \\", colors.FgCyan, colors.BgReset) + "  " + colors.SprintfANSI("_", colors.FgYellow, colors.BgReset) + "\n" +
		"      " + colors.SprintfANSI("|", colors.FgYellow, colors.BgReset) + "    " + colors.SprintfANSI("| | | |_| | | | | | |  \\| | | | | | | |_   | | | |", colors.FgCyan, colors.BgReset) + "  " + colors.SprintfANSI("|", colors.FgYellow, colors.BgReset) + "\n" +
		"      " + colors.SprintfANSI("|", colors.FgYellow, colors.BgReset) + "    " + colors.SprintfANSI("| | |  _  | | | | | |     | | | | | |  _/  | |-,<", colors.FgCyan, colors.BgReset) + "   " + colors.SprintfANSI("|", colors.FgYellow, colors.BgReset) + "\n" +
		"       " + colors.SprintfANSI("\\", colors.FgYellow, colors.BgReset) + "   " + colors.SprintfANSI("| | | | | | | |_| | | |\\  | | |_| | | |___ | | | |", colors.FgCyan, colors.BgReset) + " " + colors.SprintfANSI("/", colors.FgYellow, colors.BgReset) + "\n" +
		"	" + colors.SprintfANSI("\\", colors.FgYellow, colors.BgReset) + "  " + colors.SprintfANSI("|/  |/  |/   \\__/\\| |/  \\/  |__/|/  |____/ |/  |/", colors.FgCyan, colors.BgReset) + " " + colors.SprintfANSI("/", colors.FgYellow, colors.BgReset) + "\n" +
		"	 " + colors.SprintfANSI("\\", colors.FgYellow, colors.BgReset) + "       " + colors.SprintfANSI("_______", colors.FgCyan, colors.BgReset) + "                                    " + colors.SprintfANSI("/", colors.FgYellow, colors.BgReset) + "\n" +
		"	  " + colors.SprintfANSI("\\", colors.FgYellow, colors.BgReset) + "     " + colors.SprintfANSI("|_   __/_____   ____    _____  _____", colors.FgCyan, colors.BgReset) + "       " + colors.SprintfANSI("/", colors.FgYellow, colors.BgReset) + "\n" +
		"	   " + colors.SprintfANSI("\\", colors.FgYellow, colors.BgReset) + "      " + colors.SprintfANSI("| |__|  /| | | |\\ \\  |  __/ |  __/", colors.FgCyan, colors.BgReset) + "      " + colors.SprintfANSI("/", colors.FgYellow, colors.BgReset) + "\n" +
		"	    " + colors.SprintfANSI("\\", colors.FgYellow, colors.BgReset) + "     " + colors.SprintfANSI("|  _/| | | | | | | | | |    | |__", colors.FgCyan, colors.BgReset) + "      " + colors.SprintfANSI("/", colors.FgYellow, colors.BgReset) + "\n" +
		"	     " + colors.SprintfANSI("\\", colors.FgYellow, colors.BgReset) + "    " + colors.SprintfANSI("| |  | | | | | |-,<  | |    |  _/", colors.FgCyan, colors.BgReset) + "     " + colors.SprintfANSI("/", colors.FgYellow, colors.BgReset) + "\n" +
		"	      " + colors.SprintfANSI("\\", colors.FgYellow, colors.BgReset) + "   " + colors.SprintfANSI("| |  | |_| | | | | | | |___ | |___", colors.FgCyan, colors.BgReset) + "   " + colors.SprintfANSI("/", colors.FgYellow, colors.BgReset) + "\n" +
		"	       " + colors.SprintfANSI("\\", colors.FgYellow, colors.BgReset) + "  " + colors.SprintfANSI("|/   |__/ \\| |/  |/  |____/ |____/", colors.FgCyan, colors.BgReset) + "  " + colors.SprintfANSI("/", colors.FgYellow, colors.BgReset) + "\n" +
		"		" + colors.SprintfANSI("\\", colors.FgYellow, colors.BgReset) + "          " + colors.SprintfANSI("___   __      _", colors.FgBrightRed, colors.BgReset) + "           " + colors.SprintfANSI("/", colors.FgYellow, colors.BgReset) + "\n" +
		"		 " + colors.SprintfANSI("\\", colors.FgYellow, colors.BgReset) + "        " + colors.SprintfANSI("\\   / \\  /    \\ /", colors.FgBrightRed, colors.BgReset) + "         " + colors.SprintfANSI("/", colors.FgYellow, colors.BgReset) + "\n" +
		"		  " + colors.SprintfANSI("\\", colors.FgYellow, colors.BgReset) + "        " + colors.SprintfANSI("| |   \\ \\    //", colors.FgBrightRed, colors.BgReset) + "         " + colors.SprintfANSI("/", colors.FgYellow, colors.BgReset) + "\n" +
		"		   " + colors.SprintfANSI("\\______", colors.FgYellow, colors.BgReset) + " " + colors.SprintfANSI("| |", colors.FgBrightRed, colors.BgReset) + " " + colors.SprintfANSI("__", colors.FgBrightYellow, colors.BgReset) + " " + colors.SprintfANSI("\\ \\  //", colors.FgBrightRed, colors.BgReset) + " " + colors.SprintfANSI("________/", colors.FgYellow, colors.BgReset) + "\n" +
		"			   " + colors.SprintfANSI("| |     \\ \\//", colors.FgBrightRed, colors.BgReset) + "\n" +
		"			  " + colors.SprintfANSI("/___\\     \\_/", colors.FgBrightRed, colors.BgReset) + "\n")

	want := true
	has := true

	if want != has {
		t.Fatalf("lol")
	}
}
