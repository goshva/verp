package handlers

import (
    "fmt"
    "html/template"
    "net/http"
)

var globalTemplates *template.Template

func InitTemplates() error {
    var err error
    globalTemplates, err = template.ParseGlob("internal/templates/*.html")
    if err != nil {
        return err
    }
    
    fmt.Println("DEBUG: Global templates loaded:")
    for _, t := range globalTemplates.Templates() {
        fmt.Printf("  - %s\n", t.Name())
    }
    
    return nil
}

func RenderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
    if globalTemplates == nil {
        http.Error(w, "Templates not initialized", http.StatusInternalServerError)
        return
    }
    
    t := globalTemplates.Lookup(tmplName)
    if t == nil {
        fmt.Printf("DEBUG: Template %s not found\n", tmplName)
        http.Error(w, "Template not found", http.StatusInternalServerError)
        return
    }
    
    err := t.Execute(w, data)
    if err != nil {
        fmt.Printf("DEBUG: Template execution error: %v\n", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}