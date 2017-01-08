package main

import (
	"flag"
	"fmt"
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

	f := new(Flags)

	flag.BoolVar(&f.force, "force", true, "Overwrite an existing file")
	flag.StringVar(&f.syntax, "syntax", "markdown", "The syntax chosen for this wiki")
	flag.StringVar(&f.extension, "extension", ".md", "The file extension for this wiki")
	flag.StringVar(&f.output_dir, "output_dir", "~/vimwiki/public_html/", "The full path of the output directory")
	flag.StringVar(&f.input_file, "input_file", "", "The full path of the wiki page")
	flag.StringVar(&f.css_file, "css_file", "~/vimwiki/style.css", "The full path of the css file for this wiki")
	flag.StringVar(&f.template_path, "template_path", "~/vimwiki/templates/", "The full path to the wiki's templates")
	flag.StringVar(&f.template_default, "template_default", "default", "The default template name")
	flag.StringVar(&f.template_ext, "template_ext", ".tpl", "The extension of template files")
	flag.StringVar(&f.root_path, "root_path", "./", "A count of ../ for pages buried in subdirs")

	flag.Parse()

	fmt.Print(*f)
}
