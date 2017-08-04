# Vimwiki-GoDown

Vimwiki-GoDown is a [Markdown][0] to HTML converter for `.md` files created with [Vimwiki][1]. The converter is implemented in [Golang][2]. It's idea is based on the preceding implementation [mikasa_md2html][3], which is written in [Python][4].

To convert Markdown syntax to HTML Vimwiki-GoDown uses the [Blackfriday Markdown processor][5]. The only conversion that is done manually are checkboxes.

One of the main reasons behind creating another Markdown to HTML converter for Vimwiki is, that compared to other converters, Vimwiki-GoDown has the additional ability to prefix relative links to other Vimwiki pages. This means that the converted HTML files can be hosted in a subdirectory or rather a URL subpath (of course it can also be used without prefixing links between Vimwiki pages).

### Installation

To run Vimwiki-GoDown an installation of [Golang][2] (Go 1 or above) is required. The following command will install Vimwiki-GoDown into the `$GOPATH/bin` directory:

	go get "github.com/maqiv/vimwiki-godown"

Please make sure that the environment variable `$GOPATH` is set correctly and `$GOPATH/bin` is included in the `$PATH` environment variable. For more information please consult the [Golang $GOPATH environment variable][6].

### Configuration

To get Vimwiki-GoDown running with Vimwiki, only a few settings need to be configured. Parameters that need to be adjusted can be found in the `~/.vim/vimrc` configuration file or rather in the variable `g:vimwiki_list`. The two parameters that are needed for Vimwiki-GoDown are:

* `custom_wiki2html`: Path to the Vimwiki-GoDown binary
* `custom_wiki2html_args`:  URL Subpath that is prefixed to every relative link between Vimwiki pages

Example:
```
let g:vimwiki_list = [{
...
  \'custom_wiki2html': '$GOPATH/bin/vimwiki-godown',
  \'custom_wiki2html_args': 'xyz/',
...
\}]
```

More information about the parameters can be found in the [Vimwiki help][7].

### License

This software is distributed under the [MIT license][8].

### ToDo

* Possibility to use templates
* Ensure compatility to percent codes used in misaka_md2html.py
* Customize title guessing so that TOC heading is not used as title
* _.... any additional ideas or contributions are very welcome :-)_


[0]: http://daringfireball.net/projects/markdown/ "Markdown"
[1]: https://vimwiki.github.io "Vimwiki"
[2]: http://golang.org/ "Go Language"
[3]: https://github.com/jason6/vimwiki_md2html "mikasa_md2html.py"
[4]: https://www.python.org/ "Python"
[5]: https://github.com/russross/blackfriday "Blackfriday Markdown processor"
[6]: https://golang.org/doc/code.html#GOPATH "Golang $GOPATH environment variable"
[7]: https://github.com/vimwiki/vimwiki/blob/dev/doc/vimwiki.txt "Vimwiki help"
[8]: LICENSE "MIT License"
