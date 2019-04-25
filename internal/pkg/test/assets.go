package test

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func GetRepoRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "test: could not get current working directory")
	}

	dir := cwd
	for dir != filepath.VolumeName(dir) {
		if dir == "" {
			break
		}

		stat, err := os.Stat(filepath.Join(dir, ".git"))
		if err != nil && !os.IsNotExist(err) {
			return "", errors.Wrap(err, "test: failed to stat directory")
		}

		if (err != nil && os.IsNotExist(err)) || !stat.IsDir() {
			dir = filepath.Dir(dir)
			continue
		}

		return dir, nil
	}

	return "", errors.New("test: could not find the repo root")
}

func GetAssetsPath() (string, error) {
	repo, err := GetRepoRoot()
	if err != nil {
		return "", err
	}

	return filepath.Join(repo, "assets"), nil
}

func GetAssetPath(path ...string) (string, error) {
	assets, err := GetAssetsPath()
	if err != nil {
		return "", err
	}

	return filepath.Join(append([]string{assets}, path...)...), nil
}

func GetTestAssetsPath() (string, error) {
	repo, err := GetRepoRoot()
	if err != nil {
		return "", err
	}

	return filepath.Join(repo, "test"), nil
}

func GetTestAssetPath(path ...string) (string, error) {
	assets, err := GetTestAssetsPath()
	if err != nil {
		return "", err
	}

	return filepath.Join(append([]string{assets}, path...)...), nil
}
