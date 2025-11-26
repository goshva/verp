package handlers

import (
    "fmt"
    "html/template"
    "net/http"
    "os"
    "path/filepath"
)

// TemplateRenderer handles template rendering for handlers
type TemplateRenderer struct {
    tmpl *template.Template
}

// NewTemplateRenderer creates a new template renderer
func NewTemplateRenderer() *TemplateRenderer {
    tr := &TemplateRenderer{}
    
    // Parse all HTML files recursively
    tmpl := template.New("")
    
    err := filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() && filepath.Ext(path) == ".html" {
            data, err := os.ReadFile(path)
            if err != nil {
                return err
            }
            _, err = tmpl.New(path).Parse(string(data))
            if err != nil {
                return err
            }
        }
        return nil
    })
    
    if err != nil {
        fmt.Printf("Warning: Could not load templates: %v\n", err)
        // Create a basic template as fallback
        tr.tmpl = template.Must(template.New("base").Parse("Template system not initialized"))
    } else {
        tr.tmpl = tmpl
    }
    return tr
}

// RenderTemplate renders a template with data
func (tr *TemplateRenderer) RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
    fmt.Printf("DEBUG: Rendering template: %s with data: %+v\n", name, data)
    
    if tr.tmpl == nil {
        fmt.Printf("ERROR: Templates not initialized!\n")
        http.Error(w, "Template system not initialized", http.StatusInternalServerError)
        return
    }
    
    // Try different template paths
    templatePaths := []string{
        name,
        "templates/" + name,
        "templates/pages/" + name,
        "templates/partials/" + name,
    }
    
    var foundTemplate *template.Template
    for _, path := range templatePaths {
        foundTemplate = tr.tmpl.Lookup(path)
        if foundTemplate != nil {
            break
        }
    }
    
    if foundTemplate == nil {
        fmt.Printf("ERROR: Template '%s' not found! Available templates:\n", name)
        http.Error(w, "Template not found: "+name, http.StatusInternalServerError)
        return
    }
    /*
    w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
    w.Header().Set("Pragma", "no-cache")
    w.Header().Set("Expires", "0")
    */
    err := foundTemplate.Execute(w, data)
    if err != nil {
        fmt.Printf("DEBUG: Template error: %v\n", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}