package vimwiki

import "testing"

func TestBuildTargetFilepath(t *testing.T) {

	tests := []struct {
		in   []string
		want string
	}{
		{
			in:   []string{"path0/subfile", "path1/subpath/subsub/"},
			want: "path1/subpath/subsub/subfile.html",
		},
		{
			in:   []string{"path4/somefile", "anotherpath/like/this"},
			want: "anotherpath/like/this/somefile.html",
		},
	}

	for _, tt := range tests {
		if got := BuildTargetFilepath(tt.in[0], tt.in[1]); got != tt.want {
			t.Errorf("Building target filepath for '%v' and '%v' failed, got: '%v', want: '%v'.", tt.in[0], tt.in[1], got, tt.want)
		}
	}
}

func TestProcessRelativeLinks(t *testing.T) {

	tests := []struct {
		in, prefix, want string
	}{
		{
			`Pellentesque habitant [this is a link here](this/is/a) morbi tristique [an external link]
(https://www.example.org/?pbla=muh) senectus et netus et malesuada fames ac turpis egestas.`,
			"pref",
			`Pellentesque habitant [this is a link here](pref/this/is/a.html) morbi tristique [an external link]
(https://www.example.org/?pbla=muh) senectus et netus et malesuada fames ac turpis egestas.`,
		},
		{
			`Suspendisse enim, vivamus [a little link](3abc/blamuh.jpg) non nec libero nam, magna suspendisse www.exmpl.org ac etiam et
eget enim, congue pede lacus fringilla tempus non at, magna erat vel.`,
			"",
			`Suspendisse enim, vivamus [a little link](3abc/blamuh.jpg) non nec libero nam, magna suspendisse www.exmpl.org ac etiam et
eget enim, congue pede lacus fringilla tempus non at, magna erat vel.`,
		},
		{
			`Suspendisse enim, vivamus [a little link](3abc/blamuh.jpg) non nec libero nam, magna suspendisse www.exmpl.org ac
etiam et eget enim, congue pede lacus fringilla tempus non at, magna erat vel.`,
			"hello/",
			`Suspendisse enim, vivamus [a little link](hello/3abc/blamuh.jpg) non nec libero nam, magna suspendisse www.exmpl.org ac
etiam et eget enim, congue pede lacus fringilla tempus non at, magna erat vel.`,
		},
	}

	for _, tt := range tests {
		got := ProcessRelativeLinks(tt.in, tt.prefix)
		if tt.want != got {
			t.Errorf("ProcessRelativeLinks(%q, %q) = %q, want %q", tt.in, tt.prefix, got, tt.want)
		}
	}
}

func TestProcessHtmlCheckboxes(t *testing.T) {

	tests := []struct {
		in, want string
	}{
		{` * [ ] this is an open task`, ` * <input type="checkbox" disabled>this is an open task`},
		{` * [o] this is a started task`, ` * <input type="checkbox" disabled>this is a started task`},
		{` * [O] this is a progressing task`, ` * <input type="checkbox" disabled>this is a progressing task`},
		{` * [X] this is a closed task`, ` * <input type="checkbox" disabled checked>this is a closed task`},
	}

	for _, tt := range tests {
		got := ProcessHtmlCheckboxes(tt.in)
		if tt.want != got {
			t.Errorf("ProcessHtmlCheckboxes(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
