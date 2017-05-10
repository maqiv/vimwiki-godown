// TODO: Insert copyright crap

// TODO: Write a nice short description about stuff going on here

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/russross/blackfriday"
)

const RELATIVE_LINKS_PREFIX = "w"
const TRG_FILE_EXTENSION string = ".html"
const RGX_MDWN_HYPERLINK string = `\[(?P<desc>.+)\]\((?P<link>(/?\w+)+(\.\w+)?)\)`

var urlBasePrefix string = "w"

type Flags struct {
	Force           bool
	Syntax          string
	Extension       string
	OutputDirectory string
	InputFile       string
	CssFile         string
	TmplPath        string
	TmplDefault     string
	TmplExtension   string
	RootPath        string
}

func main() {

	var targetFilePath string
	var mdInput, mdInputProcessed string
	var mdInputRaw []byte
	var renderer blackfriday.Renderer
	var docTitle string

	fl := parseArguments(os.Args)

	// Check if file already exists and overwrite flag is not set
	targetFilePath = BuildTargetFilepath(fl.InputFile, fl.OutputDirectory)
	if _, err := os.Stat(targetFilePath); os.IsNotExist(err) && !fl.Force {
		fmt.Println("Conversion of file %v aborted: File does exist and force flag is set to 0.", targetFilePath)
		os.Exit(0)
	}

	docTitle = "insert some titlestuff here"

	// Set document title
	renderer = blackfriday.HtmlRenderer(0, docTitle, fl.CssFile)

	// Read input file in markdown format
	mdInput = ReadFile(fl.InputFile)
	mdInputProcessed = PrefixRelativeHyperlinks(mdInput, RGX_MDWN_HYPERLINK, RELATIVE_LINKS_PREFIX)
	mdInputRaw = []byte(mdInputProcessed)

	// Convert markdown content to html
	htmlOutputRaw := blackfriday.Markdown(mdInputRaw, renderer, 0)
	htmlOutput := string(htmlOutputRaw)

	fmt.Println(string(htmlOutput))

	// TODO
	// Write html content to html file
	tempFile, _ := ioutil.TempFile("/tmp", "vgdwn")
	tempFile.WriteString(htmlOutput)
	fmt.Println("filename -> ", tempFile.Name())
}

// Construct the target file path where the content will be saved later
func BuildTargetFilepath(sourceFilepath string, targetDirectory string) string {
	var sourceFilename, filenameBase, targetFilename string

	sourceFilename = filepath.Base(sourceFilepath)
	filenameBase = strings.TrimSuffix(sourceFilename, filepath.Ext(sourceFilename))
	targetFilename = filenameBase + TRG_FILE_EXTENSION

	return filepath.Join(targetDirectory, targetFilename)
}

// Add a prefix to relative links so that they work
// with custom web server configurations
func PrefixRelativeHyperlinks(mdContent string, hyperlinkMdIdentifier string, hyperlinkPrefix string) string {
	var patternPrefixedRelativeLink, returnVal string
	var patternNames []string
	var r *regexp.Regexp
	var err error

	// Find the solution here with regex stuff:
	// https://play.golang.org/p/IeAJmtkwB7 (OLD)
	// https://play.golang.org/p/NzQ3R8FHem (OLD)
	// https://play.golang.org/p/c0DwYWV-gl

	if r, err = regexp.Compile(hyperlinkMdIdentifier); err != nil {
		panic("Could not compile hyperlink markdown regex.")
	}

	patternNames = r.SubexpNames()
	patternPrefixedRelativeLink = fmt.Sprintf("%s/${%s}", path.Clean(hyperlinkPrefix), patternNames[2])
	//fmt.Println(patternPrefixedRelativeLink)
	pattern := fmt.Sprintf("[${%s}](%s)", patternNames[1], patternPrefixedRelativeLink)
	//fmt.Println(pattern)

	returnVal = r.ReplaceAllString(mdContent, pattern)

	//r.ReplaceAllStringFunc(mdwn_content, func(sstr string) string {
	//	var ret, formatted_url string

	//	sm := r.FindStringSubmatch(sstr)
	//	formatted_url = path.Join(link_pref, sm[2])
	//	ret = fmt.Sprintf("[%v](%v)", sm[1], formatted_url)

	//  return ret
	//})

	fmt.Println(returnVal)
	return returnVal
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
