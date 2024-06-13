package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"myapp/internal-web/core/domain"
	"net/http"
)

type TemplateExecutor struct {
	templateDir string
}

func NewTemplateExecutor(td string) TemplateExecutor {
	return TemplateExecutor{
		templateDir: td,
	}
}

func (te *TemplateExecutor) Home(w http.ResponseWriter, data domain.HomeTemplateData) error {
	files := []string{
		fmt.Sprintf("%s/index.html", te.templateDir),
		fmt.Sprintf("%s/body/body.html", te.templateDir),
		fmt.Sprintf("%s/body/main/home.html", te.templateDir),
		fmt.Sprintf("%s/body/nav/nav.html", te.templateDir),
	}

	ts, pfErr := template.ParseFiles(files...)
	if pfErr != nil {
		return fmt.Errorf("error parsing template Home: %w", pfErr)
	}

	var buffer bytes.Buffer

	exErr := ts.Execute(&buffer, data)
	if exErr != nil {
		return fmt.Errorf("error executing template Home to buffer: %w", pfErr)
	}

	_, err := io.Copy(w, &buffer)
	if err != nil {
		return fmt.Errorf("error writing buffer to ResponseWriter: %w", err)
	}

	return nil
}

func (te *TemplateExecutor) Login(w http.ResponseWriter, data domain.LoginTemplateData) error {

	files := []string{
		fmt.Sprintf("%s/index.html", te.templateDir),
		fmt.Sprintf("%s/body/body.html", te.templateDir),
		fmt.Sprintf("%s/body/main/login.html", te.templateDir),
		fmt.Sprintf("%s/body/nav/nav.html", te.templateDir),
	}

	ts, pfErr := template.ParseFiles(files...)
	if pfErr != nil {
		return fmt.Errorf("error parsing template Login: %w", pfErr)
	}

	var buffer bytes.Buffer

	exErr := ts.Execute(&buffer, data)
	if exErr != nil {
		return fmt.Errorf("error executing template Login to buffer: %w", pfErr)
	}

	_, err := io.Copy(w, &buffer)
	if err != nil {
		return fmt.Errorf("error writing buffer to ResponseWriter: %w", err)
	}

	return nil
}

func (te *TemplateExecutor) Comics(w http.ResponseWriter, data domain.ComicsTemplateData) error {

	files := []string{
		fmt.Sprintf("%s/index.html", te.templateDir),
		fmt.Sprintf("%s/body/body.html", te.templateDir),
		fmt.Sprintf("%s/body/main/comics.html", te.templateDir),
		fmt.Sprintf("%s/body/nav/nav.html", te.templateDir),
	}

	ts, pfErr := template.ParseFiles(files...)
	if pfErr != nil {
		return fmt.Errorf("error parsing template Comics: %w", pfErr)
	}

	var buffer bytes.Buffer

	exErr := ts.Execute(&buffer, data)
	if exErr != nil {
		return fmt.Errorf("error executing template Comics to buffer: %w", pfErr)
	}

	_, err := io.Copy(w, &buffer)
	if err != nil {
		return fmt.Errorf("error writing buffer to ResponseWriter: %w", err)
	}
	return nil
}
