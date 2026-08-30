package common

import (
	"html/template"
	"os"
	"path/filepath"
)

type Generator struct {
	outputDir    string
	templatesDir string
}

func NewGenerator(templatesDir string, outputDir string) *Generator {
	return &Generator{
		outputDir:    outputDir,
		templatesDir: templatesDir,
	}
}

func (this *Generator) Generate() {
	os.Mkdir(this.outputDir, 0755)

	partials, err := template.ParseGlob(this.templatesDir +  "/partials/*.html")
	if err != nil {
		panic(err)
	}

	entries, err := os.ReadDir(this.templatesDir)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		f, err := os.Create(filepath.Join(this.outputDir, entry.Name()))
		if err != nil {
			panic(err)
		}

		defer f.Close()

		templ, err := partials.Clone()
		if err != nil {
			panic(err)
		}

	   _, err = templ.ParseFiles(filepath.Join(this.templatesDir, entry.Name()))
		if err != nil {
			panic(err)
		}

		err = templ.ExecuteTemplate(f, entry.Name(), nil)
		if err != nil {
			panic(err)
		}
	}
}
