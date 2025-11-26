# Инструкция по обновлению роутера

Для завершения реструктуризации необходимо обновить файл роутера:

## 1. Обновите internal/handlers/router.go

Замените секцию загрузки шаблонов на:

```go
template.Must(template.ParseFS(templatesFS,
    "layouts/*.html",
    "pages/*.html", 
    "partials/*.html",
    "forms/*.html",
    "components/*.html",
))
