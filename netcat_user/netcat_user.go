package netcat_user

import (
	"bufio"
	"colors"
	"common"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type User struct {
	Id          uint
	Connection  net.Conn
	Name, Color string
	Reader      *bufio.Reader
}

func NewUser(c net.Conn, name, color string, id uint) *User {
	usr := new(User)
	usr.Id = id
	usr.Connection = c
	usr.Name = name
	usr.Color = color
	return usr
}

func IsValidUsername(name string, userMap map[uint]*User) (bool, error) {
	if !common.IsAlphaNum(name) {
		return false, errors.New("contains forbidden characters")
	}

	for _, user := range userMap {
		if user.Name == name {
			return false, errors.New("already taken")
		}
	}

	if len(name) == 0 || name == "\n" {
		return false, errors.New("not valid")
	}

	return true, nil
}

func IsValidColor(req string) (string, error) {
	if req == "" {
		return "none", errors.New("empty request forbidden")
	}

	if req[0] == '#' {
		if !common.IsHexadecimal(req[1:]) {
			return "none", errors.New("non-hexadecimal value after '#'")
		}

		return "hex", nil
	}

	if common.IsAlphabetic(req) {
		if strings.ToLower(req) == "black" {
			return "none", errors.New("black forbidden")
		}

		if _, exists := colors.EnglishColors[strings.ToLower(req)]; !exists {
			return "none", fmt.Errorf("color \"%s\" not supported", req)
		}

		return "eng", nil
	}

	var sep rune
	for _, r := range req {
		if r < '0' || r > '9' {
			sep = r
			break
		}
	}

	values := strings.Split(req, string(sep))
	if len(values) != 3 {
		return "none", errors.New("wrong decimal format")
	}

	if common.IsNumeric(values[0]) && common.IsNumeric(values[1]) && common.IsNumeric(values[2]) {
		for i, val := range values {
			num, _ := strconv.Atoi(val)
			if num < 0 || num > 255 {
				return "none", fmt.Errorf("value #%d out of range", i+1)
			}
		}
		return "dec", nil
	}

	return "none", errors.New("nuh uh")
}

func StrToColor(req string) (string, error) {
	what, err := IsValidColor(req)

	switch what {
	case "eng":
		switch req {
		case "red", "green", "blue", "yellow", "cyan", "magenta", "white":
			return colors.NewFGColorANSI(req)
		default:
			return colors.RGBValuesToColor(colors.EnglishColors[strings.ToLower(req)])
		}

	case "dec":
		var sep rune
		var vals []int
		for _, r := range req {
			if r < '0' || r > '9' {
				sep = r
				break
			}
		}
		values := strings.Split(req, string(sep))
		for _, val := range values {
			num, _ := strconv.Atoi(val)
			vals = append(vals, num)
		}
		return colors.NewFGColorRGB(vals[0], vals[1], vals[2])

	case "hex":
		var vals []int
		hexs := common.StrChunk(strings.ToLower(req[1:]), 2)
		for _, hex := range hexs {
			vals = append(vals, common.AtoiBase(hex, "0123456789abcdef"))
		}
		return colors.NewFGColorRGB(vals[0], vals[1], vals[2])

	case "none":
		return "", err
	}

	return "", errors.New("nuh uh")
}

func (usr *User) ColoredUsername() string {
	if usr.Color == "" {
		return usr.Name
	}
	return usr.Color + usr.Name + colors.ResetColorsTag
}
