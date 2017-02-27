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

const TRG_FILE_EXTENSION string = ".html"
const RGX_MDWN_HYPERLINK string = `\[(?P<desc>.+)\]\((?P<link>(/?\w+)+(\.\w+)?)\)`

var url_base_prefix string = "w"

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

	// Check if file already exists and overwrite flag is not set
	trg_file_path := build_trg_filepath(fl.input_file, fl.output_dir)
	if _, err := os.Stat(trg_file_path); os.IsNotExist(err) && !fl.force {
		fmt.Println("Conversion of file %v aborted: File does exist and force flag is set to 0.", trg_file_path)
		os.Exit(0)
	}

	var md_in, md_processed string
	var md_in_raw []byte
	var renderer blackfriday.Renderer
	var doc_title string

	// Set document title
	doc_title = find_document_title()
	renderer = blackfriday.HtmlRenderer(0, doc_title, fl.css_file)

	// Read input file in markdown format
	md_in = ReadFile(fl.input_file)
	md_processed = prefixRelativeHyperlinks(md_in, RGX_MDWN_HYPERLINK, "w")
	md_in_raw = []byte(md_processed)

	// Convert markdown content to html
	html_out_raw := blackfriday.Markdown(md_in_raw, renderer, 0)
	html_out := string(html_out_raw)

	fmt.Println(string(html_out))
}

// TODO: complete function logic
func find_document_title() string {
	return "do title"
}

// Construct the target file path where the content will be saved later
func build_trg_filepath(src_filepath string, trg_dir string) string {
	var src_filename, file_basename, trg_filename string

	src_filename = filepath.Base(src_filepath)
	file_basename = strings.TrimSuffix(src_filename, filepath.Ext(src_filename))
	trg_filename = file_basename + TRG_FILE_EXTENSION

	return filepath.Join(trg_dir, trg_filename)
}

// Add a prefix to relative links so that they work
// with custom web server configurations
func prefixRelativeHyperlinks(mdwn_content string, hyperlink_mdwn_ident string, link_pref string) string {
	var patternPrefixedRelativeLink, ret string
	var patternNames []string

	// Find the solution here with regex stuff:
	// https://play.golang.org/p/IeAJmtkwB7 (OLD)
	// https://play.golang.org/p/NzQ3R8FHem (OLD)
	// https://play.golang.org/p/c0DwYWV-gl

	r, err := regexp.Compile(hyperlink_mdwn_ident)
	if err != nil {
		panic("Could not compile hyperlink markdown regex.")
	}

	patternNames = r.SubexpNames()
	patternPrefixedRelativeLink = fmt.Sprintf("%s/${%s}", path.Clean(link_pref), patternNames[2])
	//fmt.Println(patternPrefixedRelativeLink)
	pattern := fmt.Sprintf("[${%s}](%s)", patternNames[1], patternPrefixedRelativeLink)
	//fmt.Println(pattern)

	ret = r.ReplaceAllString(mdwn_content, pattern)

	//r.ReplaceAllStringFunc(mdwn_content, func(sstr string) string {
	//	var ret, formatted_url string

	//	sm := r.FindStringSubmatch(sstr)
	//	formatted_url = path.Join(link_pref, sm[2])
	//	ret = fmt.Sprintf("[%v](%v)", sm[1], formatted_url)

	//  return ret
	//})

	fmt.Println(ret)
	return ret
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
