package vcs

import (
	"testing"
)

func TestVCSLookup(t *testing.T) {
	// TODO: Expand to make sure it detected the right vcs.
	urlList := map[string]struct {
		work bool
		t    Type
	}{
		"https://github.com/masterminds":                                   {work: false, t: Git},
		"https://github.com/Masterminds/VCSTestRepo":                       {work: true, t: Git},
		"https://bitbucket.org/mattfarina/testhgrepo":                      {work: true, t: Hg},
		"https://launchpad.net/govcstestbzrrepo/trunk":                     {work: true, t: Bzr},
		"https://launchpad.net/~mattfarina/+junk/mygovcstestbzrrepo":       {work: true, t: Bzr},
		"https://launchpad.net/~mattfarina/+junk/mygovcstestbzrrepo/trunk": {work: true, t: Bzr},
		"https://git.launchpad.net/govcstestgitrepo":                       {work: true, t: Git},
		"https://git.launchpad.net/~mattfarina/+git/mygovcstestgitrepo":    {work: true, t: Git},
		"http://farbtastic.googlecode.com/svn/":                            {work: true, t: Svn},
		"http://farbtastic.googlecode.com/svn/trunk":                       {work: true, t: Svn},
		"https://code.google.com/p/farbtastic":                             {work: false, t: Svn},
		"https://code.google.com/p/plotinum":                               {work: true, t: Hg},
		"https://example.com/foo/bar.git":                                  {work: true, t: Git},
		"https://example.com/foo/bar.svn":                                  {work: true, t: Svn},
		"https://example.com/foo/bar/baz.bzr":                              {work: true, t: Bzr},
		"https://example.com/foo/bar/baz.hg":                               {work: true, t: Hg},
		"https://gopkg.in/tomb.v1":                                         {work: true, t: Git},
		"https://golang.org/x/net":                                         {work: true, t: Git},
		"https://speter.net/go/exp/math/dec/inf":                           {work: true, t: Git},
	}

	for u, c := range urlList {
		ty, _, err := detectVcsFromRemote(u)
		if err == nil && c.work == false {
			t.Errorf("Error detecting VCS from URL(%s)", u)
		}

		if err == ErrCannotDetectVCS && c.work == true {
			t.Errorf("Error detecting VCS from URL(%s)", u)
		}

		if err != nil && c.work == true {
			t.Errorf("Error detecting VCS from URL(%s): %s", u, err)
		}

		if c.work == true && ty != c.t {
			t.Errorf("Incorrect VCS type returned(%s)", u)
		}
	}
}
