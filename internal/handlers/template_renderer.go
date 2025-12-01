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
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
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

	// Создаем базовый шаблон с функциями
	baseTmpl := template.New("").Funcs(tr.funcMap)

	// Сначала парсим базовые шаблоны (layouts)
	baseTmpl = template.Must(baseTmpl.ParseFiles(
		"templates/layouts/base.html",
		"templates/partials/sidebar.html",
		"templates/components/machines_chart.html",
        "templates/components/theme_toggle.html",
	))

	// Парсим все основные страницы
	mainPages := []string{
		"templates/accounts_page.html",
		"templates/locations_page.html",
		"templates/machines_page.html",
		"templates/operations_page.html",
		"templates/warehouses_page.html",
		"templates/dashboard_page.html", // Добавляем дашборд
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

	// Загружаем формы
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
		fmt.Printf("DEBUG: Loaded form: %s\n", name)
	}

	// Загружаем partials для HTMX
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
		fmt.Printf("DEBUG: Loaded partial: %s\n", name)
	}

	// Загружаем компоненты
components, _ := filepath.Glob("templates/components/*.html")
for _, componentPath := range components {
    componentTmpl := template.New("").Funcs(tr.funcMap)
    componentTmpl = template.Must(componentTmpl.ParseFiles(componentPath))
    
    // Используем имя файла как имя шаблона
    name := filepath.Base(componentPath)
    tr.templates[name] = componentTmpl
    fmt.Printf("DEBUG: Loaded component: %s as %s\n", componentPath, name)
    
    // Также регистрируем под полным путем для удобства
    relPath := strings.TrimPrefix(componentPath, "templates/")
    tr.templates[relPath] = componentTmpl
    fmt.Printf("DEBUG: Also registered as: %s\n", relPath)
}
	// Создаем псевдонимы для компонентов для удобства использования
	tr.createComponentAliases()
}

func (tr *TemplateRenderer) createComponentAliases() {
	// Алиасы для компонентов
	aliases := map[string]string{
		"components_machines_chart.html": "machines_chart",
		"components_chart_component.html": "chart_component",
		"components_theme_toggle.html": "theme_toggle",
	}

	for fullName, alias := range aliases {
		if tmpl, exists := tr.templates[fullName]; exists {
			tr.templates[alias] = tmpl
			fmt.Printf("DEBUG: Created alias: %s -> %s\n", alias, fullName)
		}
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
	if strings.Contains(name, "components/") || strings.Contains(name, "partials/") {
		// Для компонентов используем базовое имя файла
		tmplName = filepath.Base(name)
		if !strings.HasSuffix(tmplName, ".html") {
			tmplName += ".html"
		}
	}

	// Проверяем, определен ли шаблон в наборе
	if tmpl.Lookup(tmplName) == nil {
		// Если нет, ищем с базовым именем
		baseName := filepath.Base(tmplName)
		if tmpl.Lookup(baseName) == nil {
			// Исполняем первый найденный шаблон
			fmt.Printf("WARN: Template %s not found, using default execution\n", tmplName)
			if err := tmpl.Execute(w, data); err != nil {
				fmt.Printf("ERROR: Template %s execution failed: %v\n", name, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		tmplName = baseName
	}

	if err := tmpl.ExecuteTemplate(w, tmplName, data); err != nil {
		fmt.Printf("ERROR: Template %s execution failed: %v\n", name, err)
		fmt.Printf("DEBUG: Available templates in this tmpl: %s\n", tmpl.DefinedTemplates())
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