package vcs

import (
	"io/ioutil"
	//"log"
	"os"
	"testing"
)

// To verify svn is working we perform intergration testing
// with a known svn service.

// Canary test to ensure SvnRepo implements the Repo interface.
var _ Repo = &SvnRepo{}

func TestSvn(t *testing.T) {

	tempDir, err := ioutil.TempDir("", "go-vcs-svn-tests")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err = os.RemoveAll(tempDir)
		if err != nil {
			t.Error(err)
		}
	}()

	repo, err := NewSvnRepo("https://github.com/Masterminds/VCSTestRepo/trunk", tempDir+"/VCSTestRepo")
	if err != nil {
		t.Error(err)
	}

	if repo.Vcs() != Svn {
		t.Error("Svn is detecting the wrong type")
	}

	// Check the basic getters.
	if repo.Remote() != "https://github.com/Masterminds/VCSTestRepo/trunk" {
		t.Error("Remote not set properly")
	}
	if repo.LocalPath() != tempDir+"/VCSTestRepo" {
		t.Error("Local disk location not set properly")
	}

	//Logger = log.New(os.Stdout, "", log.LstdFlags)

	// Do an initial checkout.
	err = repo.Get()
	if err != nil {
		t.Errorf("Unable to checkout SVN repo. Err was %s", err)
	}

	// Verify SVN repo is a SVN repo
	if repo.CheckLocal() == false {
		t.Error("Problem checking out repo or SVN CheckLocal is not working")
	}

	// Verify an incorrect remote is caught when NewSvnRepo is used on an existing location
	_, nrerr := NewSvnRepo("https://github.com/Masterminds/VCSTestRepo/unknownbranch", tempDir+"/VCSTestRepo")
	if nrerr != ErrWrongRemote {
		t.Error("ErrWrongRemote was not triggered for SVN")
	}

	// Test internal lookup mechanism used outside of Hg specific functionality.
	ltype, err := DetectVcsFromFS(tempDir + "/VCSTestRepo")
	if err != nil {
		t.Error("detectVcsFromFS unable to Svn repo")
	}
	if ltype != Svn {
		t.Errorf("detectVcsFromFS detected %s instead of Svn type", ltype)
	}

	// Commenting out auto-detection tests for SVN. NewRepo automatically detects
	// GitHub to be a Git repo and that's an issue for this test. Need an
	// SVN host that can autodetect from before using this test again.
	//
	// Test NewRepo on existing checkout. This should simply provide a working
	// instance without error based on looking at the local directory.
	// nrepo, nrerr := NewRepo("https://github.com/Masterminds/VCSTestRepo/trunk", tempDir+"/VCSTestRepo")
	// if nrerr != nil {
	// 	t.Error(nrerr)
	// }
	// // Verify the right oject is returned. It will check the local repo type.
	// if nrepo.CheckLocal() == false {
	// 	t.Error("Wrong version returned from NewRepo")
	// }

	// Update the version to a previous version.
	err = repo.UpdateVersion("r2")
	if err != nil {
		t.Errorf("Unable to update SVN repo version. Err was %s", err)
	}

	// Use Version to verify we are on the right version.
	v, err := repo.Version()
	if v != "2" {
		t.Error("Error checking checked SVN out version")
	}
	if err != nil {
		t.Error(err)
	}

	// Perform an update which should take up back to the latest version.
	err = repo.Update()
	if err != nil {
		t.Error(err)
	}

	// Make sure we are on a newer version because of the update.
	v, err = repo.Version()
	if v == "2" {
		t.Error("Error with version. Still on old version. Update failed")
	}
	if err != nil {
		t.Error(err)
	}

	tags, err := repo.Tags()
	if err != nil {
		t.Error(err)
	}
	if len(tags) != 0 {
		t.Error("Svn is incorrectly returning tags")
	}

	branches, err := repo.Branches()
	if err != nil {
		t.Error(err)
	}
	if len(branches) != 0 {
		t.Error("Svn is incorrectly returning branches")
	}

	if repo.IsReference("r4") != true {
		t.Error("Svn is reporting a reference is not one")
	}

	if repo.IsReference("55") == true {
		t.Error("Svn is reporting a non-existant reference is one")
	}

	if repo.IsDirty() == true {
		t.Error("Svn incorrectly reporting dirty")
	}
}

func TestSvnCheckLocal(t *testing.T) {
	// Verify repo.CheckLocal fails for non-SVN directories.
	// TestSvn is already checking on a valid repo
	tempDir, err := ioutil.TempDir("", "go-vcs-svn-tests")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err = os.RemoveAll(tempDir)
		if err != nil {
			t.Error(err)
		}
	}()

	repo, _ := NewSvnRepo("", tempDir)
	if repo.CheckLocal() == true {
		t.Error("SVN CheckLocal does not identify non-SVN location")
	}

	// Test NewRepo when there's no local. This should simply provide a working
	// instance without error based on looking at the remote localtion.
	_, nrerr := NewRepo("https://github.com/Masterminds/VCSTestRepo/trunk", tempDir+"/VCSTestRepo")
	if nrerr != nil {
		t.Error(nrerr)
	}
}
