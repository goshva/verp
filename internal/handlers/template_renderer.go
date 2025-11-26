package handlers

import (
    "fmt"
    "html/template"
    "net/http"
)

type TemplateRenderer struct {
    templates map[string]*template.Template
}

func NewTemplateRenderer() *TemplateRenderer {
    renderer := &TemplateRenderer{
        templates: make(map[string]*template.Template),
    }
    
    // Загружаем шаблоны явно по категориям
    renderer.loadTemplates()
    
    // Выводим список загруженных шаблонов
    fmt.Println("DEBUG: Loaded templates:")
    for name := range renderer.templates {
        fmt.Printf("  - %s\n", name)
    }
    
    return renderer
}

func (tr *TemplateRenderer) loadTemplates() {
    // Загружаем базовые шаблоны
    baseTmpl := template.Must(template.ParseFiles(
        "templates/layouts/base.html",
        "templates/components/theme_toggle.html",
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
    }
    
    for _, form := range forms {
        formPath := "templates/partials/" + form
        formTmpl := template.Must(template.ParseFiles(formPath))
        tr.templates[form] = formTmpl
    }
    
    // Загружаем partials отдельно для HTMX
    partials := []string{
        "accounts_list.html",
        "locations_list.html",
        "machines_list.html", 
        "operations_list.html",
        "sidebar.html",
    }
    
    for _, partial := range partials {
        partialPath := "templates/partials/" + partial
        partialTmpl := template.Must(template.ParseFiles(partialPath))
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