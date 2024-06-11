package handlers

import (
	"fmt"
	"html/template"
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

	exErr := ts.Execute(w, data)
	if exErr != nil {
		return fmt.Errorf("error executing template Home: %w", pfErr)
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

	exErr := ts.Execute(w, data)
	if exErr != nil {
		return fmt.Errorf("error executing template Login: %w", pfErr)
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

	exErr := ts.Execute(w, data)
	if exErr != nil {
		return fmt.Errorf("error executing template Comics: %w", pfErr)
	}
	return nil
}
