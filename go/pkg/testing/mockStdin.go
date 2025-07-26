package taurinetesting

import "os"

func MockStdin(filename string, stdinput string) (clean func(), err error) {
	oldOsStdin := os.Stdin

	tmpfile, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	content := []byte(stdinput)

	if _, err := tmpfile.Write(content); err != nil {
		return nil, err

	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		return nil, err
	}

	os.Stdin = tmpfile

	return func() {
		os.Stdin = oldOsStdin
		os.Remove(tmpfile.Name())
	}, nil
}