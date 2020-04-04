package linx

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLinkParser(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("failed to get wd : %s\n", err.Error())
	}
	pagesDir := filepath.Dir(wd)
	for _, test := range []struct {
		name, path string
		want       []Link
	}{
		{
			name: "example 1",
			path: filepath.Join(pagesDir, "ex1.html"),
			want: []Link{
				Link{
					Href: "/other-page",
					Text: "A link to another page",
				},
			},
		},
		{
			name: "example 2",
			path: filepath.Join(pagesDir, "ex2.html"),
			want: []Link{
				Link{
					Href: "https://www.twitter.com/joncalhoun",
					Text: "Check me out on twitter",
				},
				Link{
					Href: "https://github.com/gophercises",
					Text: "Gophercises is on Github!",
				},
			},
		},
		{
			name: "example 3",
			path: filepath.Join(pagesDir, "ex3.html"),
			want: []Link{
				Link{
					Href: "#",
					Text: "Login",
				},
				Link{
					Href: "/lost",
					Text: "Lost? Need help?",
				},
				Link{
					Href: "https://twitter.com/marcusolsson",
					Text: "@marcusolsson",
				},
			},
		},
		{
			name: "example 4",
			path: filepath.Join(pagesDir, "ex4.html"),
			want: []Link{
				Link{
					Href: "/dog-cat",
					Text: "dog cat",
				},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			data, err := ioutil.ReadFile(test.path)
			if err != nil {
				t.Errorf("failed to read %s\nerror : %s\n", test.path, err)
			}

			links, err := GetLinks(data)
			if err != nil {
				t.Errorf("failed to get links for %s\nerror : %s\n", test.path, err.Error())
			}

			for i, link := range links {
				wantLink := test.want[i].Href
				wantText := test.want[i].Text
				linkTextMatches := func() bool { return strings.Compare(link.Text, wantText) == 0 }()
				linkMatches := func() bool { return strings.Compare(link.Href, wantLink) == 0 }()

				if !linkMatches {
					t.Errorf("\nwanted link : '%s'\ngot : '%s'\n", wantLink, link.Href)
				}

				if !linkTextMatches {
					t.Errorf("\nwanted text : '%s'\ngot : '%s'\n", wantText, link.Text)
				}
			}
		})
	}
}
