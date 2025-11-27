package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

type TemplateRenderer struct {
	templates map[string]*template.Template
	funcMap   template.FuncMap
}

func NewTemplateRenderer() *TemplateRenderer {
	renderer := &TemplateRenderer{
		templates: make(map[string]*template.Template),
	}

	renderer.addCustomFuncs()
	renderer.loadTemplates()

	fmt.Println("DEBUG: Loaded templates:")
	for name := range renderer.templates {
		fmt.Printf("  - %s\n", name)
	}

	return renderer
}

func (tr *TemplateRenderer) addCustomFuncs() {
	tr.funcMap = template.FuncMap{
		"mult": func(a int, b float64) float64 {
			return float64(a) * b
		},
		"percent": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return int(float64(a) / float64(b) * 100)
		},
		"subtract": func(a, b int) int {
			return a - b
		},
	}
}

func (tr *TemplateRenderer) loadTemplates() {
	// Загружаем базовые шаблоны с кастомными функциями
	baseTmpl := template.New("").Funcs(tr.funcMap)
	baseTmpl = template.Must(baseTmpl.ParseFiles(
		"templates/layouts/base.html",
		"templates/partials/sidebar.html",
	))

	// Загружаем основные страницы и их partials
	mainPages := []struct {
		page    string
		partial string
	}{
		{"accounts_page.html", "accounts_list.html"},
		{"locations_page.html", "locations_list.html"},
		{"machines_page.html", "machines_list.html"},
		{"operations_page.html", "operations_list.html"},
		{"warehouses_page.html", "warehouses_list.html"},
		{"dashboard_page.html", ""}, // Добавлена дашборд страница
		{"auth.html", ""},
	}

	for _, mp := range mainPages {
		// Создаем клон базового шаблона для каждой страницы
		pageTmpl := template.Must(baseTmpl.Clone())

		// Парсим основную страницу
		pagePath := "templates/" + mp.page
		pageTmpl = template.Must(pageTmpl.ParseFiles(pagePath))

		// Парсим partial если он есть
		if mp.partial != "" {
			partialPath := "templates/partials/" + mp.partial
			pageTmpl = template.Must(pageTmpl.ParseFiles(partialPath))
		}

		// Сохраняем под именем основной страницы
		tr.templates[mp.page] = pageTmpl
	}

	// Загружаем формы отдельно
	forms := []string{
		"account_form.html",
		"location_form.html",
		"machine_form.html",
		"operation_form.html",
		"warehouse_form.html",
		"inventory_form.html",
		"quick_action_form.html",
	}

	for _, form := range forms {
		formPath := "templates/partials/" + form
		formTmpl := template.New("").Funcs(tr.funcMap)
		formTmpl = template.Must(formTmpl.ParseFiles(formPath))
		tr.templates[form] = formTmpl
	}

	// Загружаем partials отдельно для HTMX
	partials := []string{
		"accounts_list.html",
		"locations_list.html",
		"machines_list.html",
		"operations_list.html",
		"warehouses_list.html",
		"sidebar.html",
	}

	for _, partial := range partials {
		partialPath := "templates/partials/" + partial
		partialTmpl := template.New("").Funcs(tr.funcMap)
		partialTmpl = template.Must(partialTmpl.ParseFiles(partialPath))
		tr.templates[partial] = partialTmpl
	}
}

func (tr *TemplateRenderer) Render(w http.ResponseWriter, name string, data interface{}) {
	fmt.Printf("DEBUG: Attempting to render template: %s\n", name)

	// Заголовки против кэширования
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	tmpl, exists := tr.templates[name]
	if !exists {
		fmt.Printf("ERROR: Template %s not found in registry\n", name)
		fmt.Printf("DEBUG: Available templates: %v\n", tr.getTemplateNames())
		http.Error(w, "Template not found: "+name, http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		fmt.Printf("ERROR: Template %s execution failed: %v\n", name, err)
		fmt.Printf("DEBUG: Template details: %+v\n", tmpl.DefinedTemplates())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("DEBUG: Successfully rendered template: %s\n", name)
}

func (tr *TemplateRenderer) getTemplateNames() []string {
	names := make([]string, 0, len(tr.templates))
	for name := range tr.templates {
		names = append(names, name)
	}
	return names
}
