package patcher

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type PatchSet struct {
	path string
}

func NewPatchSet(path string) PatchSet {
	return PatchSet{path}
}

type Version struct {
	Major          int
	Minor          int
	Patch          int
	Ref            string
	Patches        []string
	SubmoduleBumps map[string]string
}

func (ps PatchSet) VersionsToApplyFor(version string) ([]Version, error) {
	majorVersion, minorVersion, patchVersion, err := ps.parseVersion(version)
	if err != nil {
		return nil, err
	}

	startingVersions, err := ps.parseStartingVersionsFile(majorVersion, minorVersion)
	if err != nil {
		return nil, err
	}

	var versions []Version
	var currentVersion Version
	for _, v := range startingVersions.Versions {
		vers := Version{
			Major:          majorVersion,
			Minor:          minorVersion,
			Patch:          v.Version,
			Ref:            v.Ref,
			SubmoduleBumps: map[string]string{},
		}

		releaseDirName := fmt.Sprintf("%d.%d", majorVersion, minorVersion)

		for _, patch := range v.Patches {
			vers.Patches = append(vers.Patches, filepath.Join(ps.path, releaseDirName, patch))
		}

		for path, submodule := range v.Submodules {
			vers.SubmoduleBumps[path] = submodule.Ref
		}

		if v.Version == patchVersion {
			currentVersion = vers
		}

		versions = append(versions, vers)
	}

	var versionsToApply []Version
	for _, v := range versions {
		if v.Ref == currentVersion.Ref && v.Patch <= currentVersion.Patch {
			versionsToApply = append(versionsToApply, v)
		}
	}

	return versionsToApply, nil
}

func (ps PatchSet) SubmodulePatchesFor(version Version) (map[string][]string, error) {
	patchVersionDir := fmt.Sprintf("%d.%d/%d", version.Major, version.Minor, version.Patch)
	patchDir := filepath.Join(ps.path, patchVersionDir)

	submodulePatches := make(map[string][]string)

	err := filepath.Walk(patchDir, filepath.WalkFunc(func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".patch" && filepath.Dir(path) != patchDir {
			rel, err := filepath.Rel(patchDir, path)
			if err != nil {
				return err
			}

			relativePatchDir, err := filepath.Rel(patchDir, filepath.Dir(path))
			if err != nil {
				return err
			}

			absolutePathToPatch := filepath.Join(patchDir, rel)

			submodulePatches[relativePatchDir] = append(submodulePatches[relativePatchDir], absolutePathToPatch)
		}
		return nil
	}))
	if err != nil {
		return nil, err
	}

	return submodulePatches, nil
}

func (ps PatchSet) parseVersion(version string) (int, int, int, error) {
	versionParts := strings.Split(version, ".")
	majorVersion, err := strconv.Atoi(versionParts[0])
	if err != nil {
		return -1, -1, -1, err
	}

	minorVersion, err := strconv.Atoi(versionParts[1])
	if err != nil {
		return -1, -1, -1, err
	}

	patchVersion, err := strconv.Atoi(versionParts[2])
	if err != nil {
		return -1, -1, -1, err
	}

	return majorVersion, minorVersion, patchVersion, nil
}

func (ps PatchSet) parseStartingVersionsFile(majorVersion, minorVersion int) (StartingVersions, error) {
	startingVersionsYAML, err := ioutil.ReadFile(filepath.Join(ps.path, fmt.Sprintf("%d.%d", majorVersion, minorVersion), "starting-versions.yml"))
	if err != nil {
		return StartingVersions{}, err
	}

	var startingVersions StartingVersions
	err = yaml.Unmarshal(startingVersionsYAML, &startingVersions)
	if err != nil {
		return StartingVersions{}, err
	}

	return startingVersions, nil
}
