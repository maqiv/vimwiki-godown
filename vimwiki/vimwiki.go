package vimwiki

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	// regexp to recognize relative Markdown links to other Vimwiki pages
	RGX_MDWN_HYPERLINK string = `\[(?P<desc>.+)\]\((?P<link>(?:/?\w+)+)(?P<extension>\.\w+)?\)`

	// regexp to check if the relative link is a link to an image
	RGX_MDWN_IMAGE_EXT string = `\.(?i:gif|jpe?g|bmp|png|webp)`

	// regexp to find checkboxes in Markdown syntax
	RGX_MDWN_CHECKBOX string = `\[(\W|\.|o|O|X){1}\]\W{1}`

	// regexp to figure out if we can guess a page title
	RGX_MDWN_TITLE string = `(?m:^\s*\#(?P<title>(\s\w+)+)$)`

	HTML_CKB_UNCHECKED string = `<input type="checkbox" disabled>`
	HTML_CKB_CHECKED   string = `<input type="checkbox" disabled checked>`

	TRG_FILE_EXTENSION string = ".html"
)

// Flags represent the parsed parameters coming from Vimwiki.
// For further details about the parameters that are handed over see
// https://github.com/vimwiki/vimwiki/blob/dev/doc/vimwiki.txt
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
	UrlBasePrefix   string
}

// Try to guess the page title by crawling the page content with specific rules
func FindPageTitle(mdContent string) string {
	var title, firstLine string
	var regTitleLine, regTitle *regexp.Regexp

	regTitleLine = regexp.MustCompile(RGX_MDWN_TITLE)

	for _, s := range strings.Split(mdContent, "\n") {
		ts := strings.TrimSpace(s)
		if len(ts) > 0 {
			firstLine = ts
			break
		}
	}

	if firstMatch := regTitleLine.FindString(mdContent); strings.TrimSpace(firstMatch) == firstLine {
		regTitle = regexp.MustCompile(RGX_MDWN_TITLE)

		title = regTitle.ReplaceAllString(firstLine, "${title}")
		title = strings.TrimSpace(title)
	}

	return title
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
func ProcessRelativeLinks(mdContent string, relLinkPrefix string) string {
	var fileExtension, returnVal string
	var regRelLink, regFileExt *regexp.Regexp
	var regexpGroups = map[string]int{}

	regRelLink = regexp.MustCompile(RGX_MDWN_HYPERLINK)
	regFileExt = regexp.MustCompile(RGX_MDWN_IMAGE_EXT)

	// Cleanup and replace relative links to be valid with custom url prefix
	if len(relLinkPrefix) > 0 {
		relLinkPrefix = fmt.Sprintf("%s/", path.Clean(relLinkPrefix))
	}

	for i, g := range regRelLink.SubexpNames() {
		if len(g) > 0 {
			regexpGroups[g] = i
		}
	}

	returnVal = regRelLink.ReplaceAllStringFunc(mdContent, func(s string) string {
		submatch := regRelLink.FindStringSubmatch(s)

		// Check if the link points to an image
		fileExtension = submatch[regexpGroups["extension"]]
		if !regFileExt.MatchString(submatch[regexpGroups["extension"]]) {
			fileExtension = TRG_FILE_EXTENSION
		}

		return fmt.Sprintf(
			"[%s](%s%s%s)",
			submatch[regexpGroups["desc"]],
			relLinkPrefix,
			submatch[regexpGroups["link"]],
			fileExtension)
	})

	return returnVal
}

// Convert markdown styled checkboxes to HTML encoded checkboxes like Github
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
		var encodedHtml string

		encodedHtml = HTML_CKB_UNCHECKED
		if strings.Contains(s, "X") {
			encodedHtml = HTML_CKB_CHECKED
		}

		return encodedHtml
	})

	return returnVal
}
