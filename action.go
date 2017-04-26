package ec2go

import (
	"os"

	"github.com/urfave/cli"
	"strings"
)

func ListAction(c *cli.Context) {

	sess, err := newSession()
	if err != nil {
		os.Exit(1)
	}
	instances := getInstances(sess, false)

	writer := newTableWriter()
	header := []string{"InstanceId", "Name", "Type", "PrivateIP", "State"}
	setTable(writer, header, instances)

	displayTable(writer)
}

func ListAllSubAction(c *cli.Context) {
	sess, err := newSession()
	if err != nil {
		os.Exit(1)
	}

	instances := getInstances(sess, true)

	writer := newTableWriter()
	header := []string{"InstanceId", "Name", "Type", "PrivateIP", "State"}
	setTable(writer, header, instances)

	displayTable(writer)
}


func SshAction(c *cli.Context) {

	sess, err := newSession()
	if err != nil {
		os.Exit(1)
	}
	instances := getInstances(sess, false)

	instanceList := []string{}
	for _, v := range instances {
		instanceList = append(instanceList, strings.Join(v, " "))
	}

	s := &Screen{
		prompt:         QueryPrompt,
		cursorIdx:      InitialCursorIdx,
		selectedLine:   1,
		input:          []rune{},
		candidates:     instanceList,
		originContents: instanceList,
	}

	s.Run()
}
