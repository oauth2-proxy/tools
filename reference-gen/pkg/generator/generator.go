package generator

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Generator interface {
	Run() error
}

func NewGenerator(packageName string, requestedTypesList []string, headerTextFile string, outputFileName string, templateDirectory string) (Generator, error) {
	if packageName == "" {
		return nil, errors.New("a package name must be specified")
	}

	headerText, err := loadHeaderText(headerTextFile)
	if err != nil {
		return nil, fmt.Errorf("error loading header text: %v", err)
	}

	if err := checkTemplateDir(templateDirectory); err != nil {
		return nil, fmt.Errorf("invalid template directory: %v", err)
	}

	return &generator{
		packageName:       packageName,
		requestedTypes:    newStringSet(requestedTypesList),
		headerText:        headerText,
		outputFileName:    outputFileName,
		templateDirectory: templateDirectory,
	}, nil
}

// checkTemplateDir checks whether the template directory given exists and can be read
func checkTemplateDir(dir string) error {
	if dir == "" {
		return nil
	}
	path, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	if fi, err := os.Stat(path); err != nil {
		return fmt.Errorf("cannot read directory %q: %v", path, err)
	} else if !fi.IsDir() {
		return fmt.Errorf("path %q is not a directory", path)
	}
	return nil
}

// loadHeaderText loads the header text from the file if a filename was given
func loadHeaderText(fileName string) ([]byte, error) {
	if fileName == "" {
		return []byte{}, nil
	}

	headerText, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}
	return headerText, nil
}

type generator struct {
	packageName       string
	requestedTypes    stringSet
	headerText        []byte
	outputFileName    string
	templateDirectory string
}

// Run runs the generation logic for the generator
func (g *generator) Run() error {
	return nil
}
