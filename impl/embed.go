package impl

import (
	"fmt"
	"os"

	"github.com/maja42/ember/embedding"
)

func Embed(base string, destination string, attachments map[string]string) error {
	// Open executable
	exe, err := os.Open(base)
	if err != nil {
		return fmt.Errorf("Failed to open executable %q: %s", base, err)
	}
	defer exe.Close()

	// Open output
	out, err := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0755)
	if err != nil {
		return fmt.Errorf("Failed to open output file %q: %s", destination, err)
	}
	defer func() {
		_ = out.Close()
		if err := recover(); err != nil { // execution failed; delete created output file
			_ = os.Remove(destination)
		}
	}()

	logger := func(format string, args ...interface{}) {
		fmt.Printf("\t"+format+"\n", args...)
	}

	embedding.SkipCompatibilityCheck = true
	return embedding.EmbedFiles(out, exe, attachments, logger)
}

func RemoveEmbed(base string, destination string) error {
	// Open executable
	exe, err := os.Open(base)
	if err != nil {
		return fmt.Errorf("Failed to open executable %q: %s", base, err)
	}
	defer exe.Close()

	// Open output
	out, err := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0755)
	if err != nil {
		return fmt.Errorf("Failed to open output file %q: %s", destination, err)
	}
	defer func() {
		_ = out.Close()
		if err := recover(); err != nil { // execution failed; delete created output file
			_ = os.Remove(destination)
		}
	}()

	logger := func(format string, args ...interface{}) {
		// fmt.Printf("\t"+format+"\n", args...)
	}

	embedding.SkipCompatibilityCheck = true
	return embedding.RemoveEmbedding(out, exe, logger)
}
