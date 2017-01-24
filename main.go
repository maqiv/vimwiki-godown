// TODO: Insert copyright crap

// TODO: Write a nice short description about stuff going on here

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type Flags struct {
	force            bool
	syntax           string
	extension        string
	output_dir       string
	input_file       string
	css_file         string
	template_path    string
	template_default string
	template_ext     string
	root_path        string
}

func main() {

	fl := parseArguments(os.Args)

	c := ReadFile(fl.input_file)
	fmt.Println(c)
}

// Commandline arguments are parsed like defined by vimwiki
// documentation at:
// https://github.com/vimwiki/vimwiki/blob/dev/doc/vimwiki.txt#L2091
func parseArguments(args []string) *Flags {

	f := new(Flags)

	frc, err := strconv.ParseBool(args[1])
	if err != nil {
		panic(err)
	}
	f.force = frc
	f.syntax = args[2]
	f.extension = args[3]
	f.output_dir = args[4]
	f.input_file = args[5]
	f.css_file = args[6]
	f.template_path = args[7]
	f.template_default = args[8]
	f.template_ext = args[9]
	f.root_path = args[10]

	return f
}

func ReadFile(filename string) string {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	content := string(data)

	return content
}
