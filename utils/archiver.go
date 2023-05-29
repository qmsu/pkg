package utils

import "github.com/mholt/archiver/v3"

func UnArchive(src, dest string) error {
	r := archiver.NewRar()
	err := r.Unarchive(src, dest)
	if err != nil {
		return err
	}
	return nil
}
