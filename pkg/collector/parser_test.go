package collector

import (
	"testing"
	"path/filepath"
	"bufio"
	"log"
	"os"

	"github.com/EMnify/spu-exporter/pkg/transport"
)

var parsetests = []struct{
	inFile        string
	expectFailure bool
	expectedTrans *[]transport.Transport
}{
	{"../../test/example",false,&[]transport.Transport{}},
}

func TestParseLines(t *testing.T) {
	for _, tt := range parsetests {
		t.Run(tt.inFile, func(t *testing.T) {
			absPath,_ :=  filepath.Abs(tt.inFile)
			_, err := parseLines(readFromFile(absPath))
			if tt.expectFailure {
				if err == nil {
					t.Errorf("expected failure, but passed")
				}
			} else {
				if err != nil {
					t.Errorf("got parse error when it should not fail")
				}
			}
		})
	}
}

func readFromFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	return lines
}