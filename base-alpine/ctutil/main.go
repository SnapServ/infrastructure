package main

import (
	"log"

	"github.com/alecthomas/kong"
)

var cli struct {
	Account   AccountCmd   `kong:"cmd,help='Create a service account with sane defaults'"`
	Directory DirectoryCmd `kong:"cmd,help='Ensure a directory exists with given attributes'"`
	Log       LogCmd       `kong:"cmd,help='Print a log message with optional format arguments'"`
	Relocate  RelocateCmd  `kong:"cmd,help='Relocate a directory from an old to a new path and add a symlink'"`
	Run       RunCmd       `kong:"cmd,help='Run another process using execve(2), optionally changing privileges'"`
	Secret    SecretCmd    `kong:"cmd,help='Read secret from environment, either directly or from file'"`
	Template  TemplateCmd  `kong:"cmd,help='Generate configuration file based on Go template'"`
}

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds)

	ctx := kong.Parse(&cli)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
