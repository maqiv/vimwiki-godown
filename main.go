// TODO: Write a nice short description about stuff going on here

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	v "g.maqiv.com/vimwiki_godown/vimwiki"

	"github.com/russross/blackfriday"
)

func main() {
	var targetFilePath string
	var mdRaw, mdOutput string
	var renderer blackfriday.Renderer
	var docTitle string
	var htmlFlags int
	var mdExtensions int
	var targetFile *os.File
	var err error

	fl := parseArguments(os.Args)

	// Check if file already exists and overwrite flag is not set
	targetFilePath = v.BuildTargetFilepath(fl.InputFile, fl.OutputDirectory)
	if _, err = os.Stat(targetFilePath); os.IsNotExist(err) && !fl.Force {
		fmt.Println("Conversion of file %v aborted: File does exist and force flag is set to 0.", targetFilePath)
		// Exit with error code different from 0
		os.Exit(1)
	}

	htmlFlags = 0 |
		blackfriday.HTML_COMPLETE_PAGE |
		blackfriday.HTML_HREF_TARGET_BLANK

	mdExtensions = 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS

	// Read input file in markdown format
	mdRaw = readFile(fl.InputFile)

	// Process markdown content
	if len(fl.UrlBasePrefix) > 0 && fl.UrlBasePrefix != "-" {
		// Prefix sub-url path to each relative url
		mdOutput = v.ProcessRelativeLinks(mdRaw, fl.UrlBasePrefix)
	}
	mdOutput = v.ProcessHtmlCheckboxes(mdOutput)

	// Set document title
	docTitle = v.FindPageTitle(mdOutput)
	if len(docTitle) == 0 {
		docTitle = fl.InputFile
	}

	renderer = blackfriday.HtmlRenderer(htmlFlags, docTitle, fl.CssFile)

	// Convert markdown content to html
	htmlOutputRaw := blackfriday.Markdown([]byte(mdOutput), renderer, mdExtensions)
	fmt.Println(string(htmlOutputRaw))

	// Write processed and converted html content to html file
	if targetFile, err = os.OpenFile(targetFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666); err != nil {
		panic(err)
	}
	if _, err = targetFile.Write(htmlOutputRaw); err != nil {
		panic(err)
	}

	os.Exit(0)
}

// Commandline arguments are parsed like defined by vimwiki
// documentation at:
// https://github.com/vimwiki/vimwiki/blob/dev/doc/vimwiki.txt#L2091
//
// Additionally a url prefix is parsed (UrlBasePrefix) that is prefixed to
// all relative urls later. This functionality is currently not officially
// implemented in vimwiki, but at the time writing this I sent a pull request
// that implements this feature:
// https://github.com/vimwiki/vimwiki/pull/348
func parseArguments(args []string) *v.Flags {

	f := new(v.Flags)

	frc, err := strconv.ParseBool(args[1])
	if err != nil {
		panic(err)
	}
	f.Force = frc
	f.Syntax = args[2]
	f.Extension = args[3]
	f.OutputDirectory = args[4]
	f.InputFile = args[5]
	f.CssFile = args[6]
	f.TmplPath = args[7]
	f.TmplDefault = args[8]
	f.TmplExtension = args[9]
	f.RootPath = args[10]
	f.UrlBasePrefix = args[11]

	return f
}

func readFile(filename string) string {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	content := string(data)

	return content
}
