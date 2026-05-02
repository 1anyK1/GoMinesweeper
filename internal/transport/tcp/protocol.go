package tcp

import (
	"strconv"
	"strings"
)

type Command struct {
	Name string
	Args []string
}

func ParseCommand(input string) Command {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return Command{}
	}

	return Command{
		Name: parts[0],
		Args: parts[1:],
	}
}

func ParseXY(args []string) (int, int, bool) {
	if len(args) != 2 {
		return 0, 0, false
	}

	x, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, 0, false
	}

	y, err := strconv.Atoi(args[1])
	if err != nil {
		return 0, 0, false
	}

	return x, y, true
}
