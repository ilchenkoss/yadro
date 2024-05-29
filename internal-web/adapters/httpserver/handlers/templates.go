package handlers

import (
	"fmt"
	"html/template"
	"myapp/internal-web/core/domain"
	"net/http"
)

type TemplateExecutor struct {
}

func NewTemplateExecutor() TemplateExecutor {
	return TemplateExecutor{}
}

func (te *TemplateExecutor) Home(w *http.ResponseWriter, data domain.HomeTemplateData) error {
	files := []string{
		"./internal-web/storage/template/index.html",
		"./internal-web/storage/template/body/body.html",
		"./internal-web/storage/template/body/main/home.html",
		"./internal-web/storage/template/body/nav/nav.html",
	}

	ts, pfErr := template.ParseFiles(files...)
	if pfErr != nil {
		return fmt.Errorf("error parsing template Home: %w", pfErr)
	}

	exErr := ts.Execute(*w, data)
	if exErr != nil {
		return fmt.Errorf("error executing template Home: %w", pfErr)
	}
	return nil
}

func (te *TemplateExecutor) Login(w *http.ResponseWriter, data domain.LoginTemplateData) error {

	files := []string{
		"./internal-web/storage/template/index.html",
		"./internal-web/storage/template/body/body.html",
		"./internal-web/storage/template/body/main/login.html",
		"./internal-web/storage/template/body/nav/nav.html",
	}

	ts, pfErr := template.ParseFiles(files...)
	if pfErr != nil {
		return fmt.Errorf("error parsing template Login: %w", pfErr)
	}

	exErr := ts.Execute(*w, data)
	if exErr != nil {
		return fmt.Errorf("error executing template Login: %w", pfErr)
	}
	return nil
}

func (te *TemplateExecutor) Comics(w *http.ResponseWriter, data domain.ComicsTemplateData) error {

	files := []string{
		"./internal-web/storage/template/index.html",
		"./internal-web/storage/template/body/body.html",
		"./internal-web/storage/template/body/main/comics.html",
		"./internal-web/storage/template/body/nav/nav.html",
	}

	ts, pfErr := template.ParseFiles(files...)
	if pfErr != nil {
		return fmt.Errorf("error parsing template Comics: %w", pfErr)
	}

	exErr := ts.Execute(*w, data)
	if exErr != nil {
		return fmt.Errorf("error executing template Comics: %w", pfErr)
	}
	return nil
}
