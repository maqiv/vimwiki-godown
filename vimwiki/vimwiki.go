package vimwiki

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	RGX_MDWN_HYPERLINK string = `\[(?P<desc>.+)\]\((?P<link>(/?\w+)+(\.\w+)?)\)`
	RGX_MDWN_CHECKBOX  string = `\[(\W|\.|o|O|X){1}\]\W{1}`

	HTML_CKB_UNCHECKED string = `<input type="checkbox" disabled>`
	HTML_CKB_CHECKED   string = `<input type="checkbox" disabled checked>`

	TRG_FILE_EXTENSION string = ".html"
)

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
func ProcessRelativeLinks(mdContent string, relLinkPrefix string) string {
	var err error
	var ptnPrefixedRelLink, ptnRelLink, returnVal string
	var ptnRelLinkGroups []string
	var regRelLink *regexp.Regexp

	// Find the solution here with regex stuff:
	// https://play.golang.org/p/IeAJmtkwB7 (OLD)
	// https://play.golang.org/p/NzQ3R8FHem (OLD)
	// https://play.golang.org/p/c0DwYWV-gl

	if regRelLink, err = regexp.Compile(RGX_MDWN_HYPERLINK); err != nil {
		panic("Could not compile relative link markdown regex.")
	}

	// Cleanup and replace relative links to be valid with custom url prefix
	ptnRelLinkGroups = regRelLink.SubexpNames()
	if len(relLinkPrefix) > 0 {
		relLinkPrefix = fmt.Sprintf("%s/", path.Clean(relLinkPrefix))
	}

	ptnPrefixedRelLink = fmt.Sprintf("%s${%s}%s", relLinkPrefix, ptnRelLinkGroups[2], TRG_FILE_EXTENSION)
	ptnRelLink = fmt.Sprintf("[${%s}](%s)", ptnRelLinkGroups[1], ptnPrefixedRelLink)

	returnVal = regRelLink.ReplaceAllString(mdContent, ptnRelLink)

	return returnVal
}

// Convert markdown styled checkboxes to HTML coded checkboxes like Github
// styled markdown. Decide to set them checked or unchecked based on whether
// the checkbox is set with "X" or not.
func ProcessHtmlCheckboxes(mdContent string) string {
	var err error
	var returnVal string
	var regCheckbox *regexp.Regexp

	if regCheckbox, err = regexp.Compile(RGX_MDWN_CHECKBOX); err != nil {
		panic("Could not compile checkbox markdown regex.")
	}

	returnVal = regCheckbox.ReplaceAllStringFunc(mdContent, func(s string) string {
		var html string

		html = HTML_CKB_UNCHECKED
		if strings.Contains(s, "X") {
			html = HTML_CKB_CHECKED
		}

		return html
	})

	return returnVal
}
