package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
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
	// Загружаем все шаблоны из templates/ папки
	templatePatterns := []string{
		"templates/*.html",
		"templates/layouts/*.html",
		"templates/partials/*.html",
		"templates/components/*.html",
	}

	// Собираем все файлы шаблонов
	var templateFiles []string
	for _, pattern := range templatePatterns {
		files, err := filepath.Glob(pattern)
		if err != nil {
			fmt.Printf("WARN: Error globbing pattern %s: %v\n", pattern, err)
			continue
		}
		templateFiles = append(templateFiles, files...)
	}

	if len(templateFiles) == 0 {
		panic("No template files found")
	}

	fmt.Printf("DEBUG: Found %d template files:\n", len(templateFiles))
	for _, file := range templateFiles {
		fmt.Printf("  - %s\n", file)
	}

	// Собираем ВСЕ файлы, которые нужны для базового шаблона
	baseFiles := []string{
		"templates/layouts/base.html",
		"templates/partials/sidebar.html",
		// Добавляем ВСЕ partials списков
		"templates/partials/accounts_list.html",
		"templates/partials/locations_list.html",
		"templates/partials/machines_list.html",
		"templates/partials/operations_list.html",
		"templates/partials/warehouses_list.html",
		// Добавляем ВСЕ формы
		"templates/partials/account_form.html",
		"templates/partials/location_form.html",
		"templates/partials/machine_form.html",
		"templates/partials/operation_form.html",
		"templates/partials/warehouse_form.html",
		"templates/partials/inventory_form.html",
		"templates/partials/quick_action_form.html",
		"templates/components/machines_chart.html",
		"templates/components/operations_chart.html",
		"templates/components/cash_chart.html",
	}

	// Проверяем существование файлов перед добавлением
	var existingBaseFiles []string
	for _, file := range baseFiles {
		if tr.fileExists(file) {
			existingBaseFiles = append(existingBaseFiles, file)
			fmt.Printf("DEBUG: Adding to base template: %s\n", file)
		} else {
			fmt.Printf("WARN: Base template file not found: %s\n", file)
		}
	}

	// Создаем базовый шаблон с функциями и ВСЕМИ partials
	baseTmpl := template.New("").Funcs(tr.funcMap)
	baseTmpl = template.Must(baseTmpl.ParseFiles(existingBaseFiles...))

	// Парсим все основные страницы
	mainPages := []string{
		"templates/accounts_page.html",
		"templates/locations_page.html",
		"templates/machines_page.html",
		"templates/operations_page.html",
		"templates/warehouses_page.html",
		"templates/dashboard_page.html",
		"templates/auth.html",
	}

	for _, pagePath := range mainPages {
		if !tr.fileExists(pagePath) {
			fmt.Printf("WARN: Main page not found: %s\n", pagePath)
			continue
		}

		// Создаем клон базового шаблона для каждой страницы
		pageTmpl := template.Must(baseTmpl.Clone())
		pageTmpl = template.Must(pageTmpl.ParseFiles(pagePath))

		// Извлекаем имя файла без пути
		name := filepath.Base(pagePath)
		tr.templates[name] = pageTmpl
		fmt.Printf("DEBUG: Loaded main page: %s\n", name)
	}

	// Также загружаем формы отдельно для HTMX запросов
	forms := []string{
		"templates/partials/account_form.html",
		"templates/partials/location_form.html",
		"templates/partials/machine_form.html",
		"templates/partials/operation_form.html",
		"templates/partials/warehouse_form.html",
		"templates/partials/inventory_form.html",
		"templates/partials/quick_action_form.html",
	}

	for _, formPath := range forms {
		if !tr.fileExists(formPath) {
			fmt.Printf("WARN: Form not found: %s\n", formPath)
			continue
		}

		formTmpl := template.New("").Funcs(tr.funcMap)
		formTmpl = template.Must(formTmpl.ParseFiles(formPath))
		
		name := filepath.Base(formPath)
		tr.templates[name] = formTmpl
		fmt.Printf("DEBUG: Loaded form separately: %s\n", name)
	}

	// Также загружаем partials списков отдельно для HTMX
	partials := []string{
		"templates/partials/accounts_list.html",
		"templates/partials/locations_list.html",
		"templates/partials/machines_list.html",
		"templates/partials/operations_list.html",
		"templates/partials/warehouses_list.html",
	}

	for _, partialPath := range partials {
		if !tr.fileExists(partialPath) {
			fmt.Printf("WARN: Partial not found: %s\n", partialPath)
			continue
		}

		partialTmpl := template.New("").Funcs(tr.funcMap)
		partialTmpl = template.Must(partialTmpl.ParseFiles(partialPath))
		
		name := filepath.Base(partialPath)
		tr.templates[name] = partialTmpl
		fmt.Printf("DEBUG: Loaded partial separately: %s\n", name)
	}
}

func (tr *TemplateRenderer) fileExists(path string) bool {
	files, err := filepath.Glob(path)
	return err == nil && len(files) > 0
}

func (tr *TemplateRenderer) Render(w http.ResponseWriter, name string, data interface{}) {
	fmt.Printf("DEBUG: Attempting to render template: %s\n", name)

	// Проверяем, есть ли шаблон
	tmpl, exists := tr.templates[name]
	if !exists {
		// Пробуем найти с расширением .html
		if !strings.HasSuffix(name, ".html") {
			nameWithExt := name + ".html"
			tmpl, exists = tr.templates[nameWithExt]
		}
		
		if !exists {
			fmt.Printf("ERROR: Template %s not found in registry\n", name)
			fmt.Printf("DEBUG: Available templates:\n")
			for tname := range tr.templates {
				fmt.Printf("  - %s\n", tname)
			}
			http.Error(w, "Template not found: "+name, http.StatusInternalServerError)
			return
		}
	}

	// Заголовки против кэширования
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// Определяем, какой шаблон выполнять
	tmplName := name
	
	// Для основных страниц используем имя страницы
	if strings.HasSuffix(name, "_page.html") || name == "auth.html" || name == "dashboard_page.html" {
		// Оставляем как есть
	} else if strings.HasSuffix(name, "_form.html") || strings.HasSuffix(name, "_list.html") {
		// Для partials используем базовое имя
		tmplName = filepath.Base(name)
	}

	// Проверяем, определен ли шаблон в наборе
	if tmpl.Lookup(tmplName) == nil {
		// Если нет, пробуем базовое имя
		baseName := filepath.Base(tmplName)
		if tmpl.Lookup(baseName) != nil {
			tmplName = baseName
		} else {
			// Исполняем первый найденный шаблон
			fmt.Printf("WARN: Template %s not found, using default execution\n", tmplName)
			fmt.Printf("DEBUG: Available templates in this tmpl: %s\n", tmpl.DefinedTemplates())
			if err := tmpl.Execute(w, data); err != nil {
				fmt.Printf("ERROR: Template %s execution failed: %v\n", name, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}

	if err := tmpl.ExecuteTemplate(w, tmplName, data); err != nil {
		fmt.Printf("ERROR: Template %s execution failed: %v\n", name, err)
		fmt.Printf("DEBUG: Available templates in this tmpl: %s\n", tmpl.DefinedTemplates())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("DEBUG: Successfully rendered template: %s as %s\n", name, tmplName)
}

func (tr *TemplateRenderer) getTemplateNames() []string {
	names := make([]string, 0, len(tr.templates))
	for name := range tr.templates {
		names = append(names, name)
	}
	return names
}