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
			`Suspendisse enim, vivamus [a little link](3abc/blamuh.JPG) non nec libero nam, magna suspendisse www.exmpl.org ac etiam et
eget enim, congue pede lacus fringilla tempus non at, magna erat vel.`,
			"",
			`Suspendisse enim, vivamus [a little link](3abc/blamuh.JPG) non nec libero nam, magna suspendisse www.exmpl.org ac etiam et
eget enim, congue pede lacus fringilla tempus non at, magna erat vel.`,
		},
		{
			`Harenae pedes ducibusque stantem subiectis gerens molire, bene ad ascendunt rubescere signa pressaeque ipse vidit *ille
dixit* venti comitum? [Ille](blah/muh.jpg) torvis teneret oravere bis, sive. Fato ferrum dissuaserat erit coniungere curvos
exsereret gener, sequar inpensaque soporem cumque, quid Phoebum sine.`,
			"nuanu/",
			`Harenae pedes ducibusque stantem subiectis gerens molire, bene ad ascendunt rubescere signa pressaeque ipse vidit *ille
dixit* venti comitum? [Ille](nuanu/blah/muh.jpg) torvis teneret oravere bis, sive. Fato ferrum dissuaserat erit coniungere curvos
exsereret gener, sequar inpensaque soporem cumque, quid Phoebum sine.`,
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

func TestFindPageTitle(t *testing.T) {

	tests := []struct {
		in, want string
	}{
		{
			`# Fertque caput in enixa

## Timoris doceo

Lorem markdownum armis; nam *orandus* undis torsi, forte. Auris sed glandes
**postera natasque** manerent reposco animam quater fuit.

## Noluit capillis rude modo viderat malum

Hoc ea factis pariter **hostibus huic** est sparsus quibus, ventorumque vocat
*quoque* femina sorsque, non. En ambo etiamnum ardua in hunc Acmon nigraque
Hyperionis et maiores venturi *vero oras mensis*!`,
			"Fertque caput in enixa",
		},
		{
			`Postquam celare Prytaninque patria sunt opposui positum *nais* ineunt flammas
usus. Cum poterit, est ut aestu sumptis incenduntque, luna mihi.

# Cum caelo vaticinatus ignis

Lorem markdownum colle responsaque ultro Epimethida Iolen. **Resque in quoque**
obitumque; angulus *nisi vota tua* aegra.

## Repetita instructa sensere flumina

Optantemque dictis esses, fera prohibebere tardius est *spatioque ira* oculi
vestigia tacuit captatam, in, ut Lampetie. Feritate adulter pia, *placet et
eodem*: dumque!`,
			"",
		},
	}

	for _, tt := range tests {
		got := FindPageTitle(tt.in)
		if tt.want != got {
			t.Errorf("FindPageTitle(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
