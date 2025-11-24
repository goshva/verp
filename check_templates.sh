#!/bin/bash

set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
log_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

check_file() {
    local file=$1
    if [ -f "$file" ]; then
        log_success "✓ $file"
        return 0
    else
        log_error "✗ $file - не найден"
        return 1
    fi
}

check_directory() {
    local dir=$1
    if [ -d "$dir" ]; then
        log_success "✓ $dir"
        return 0
    else
        log_error "✗ $dir - не найдена"
        return 1
    fi
}

validate_template() {
    local file=$1
    if [ -f "$file" ]; then
        # Проверяем базовый синтаксис шаблонов
        if grep -q "{{" "$file" && grep -q "}}" "$file"; then
            # Проверяем незакрытые шаблоны
            local open_braces=$(grep -o "{{" "$file" | wc -l)
            local close_braces=$(grep -o "}}" "$file" | wc -l)
            
            if [ "$open_braces" -eq "$close_braces" ]; then
                log_success "✓ $file - синтаксис OK"
            else
                log_error "✗ $file - незакрытые шаблоны: {{=$open_braces }}=$close_braces"
            fi
        fi
    fi
}

echo "Проверка структуры шаблонов..."
echo "==============================="

# Проверка директорий
log_info "Проверка директорий:"
check_directory "internal/templates"
check_directory "internal/templates/components"
check_directory "internal/templates/layouts"
check_directory "internal/templates/partials"
check_directory "internal/templates/static"

echo

# Проверка основных файлов
log_info "Проверка основных файлов:"
check_file "internal/templates/layouts/base.html"
check_file "internal/templates/static/styles.css"
check_file "internal/templates/static/app.js"
check_file "internal/templates/dashboard.html"

echo

# Проверка компонентов
log_info "Проверка компонентов:"
check_file "internal/templates/components/button.html"
check_file "internal/templates/components/card.html"
check_file "internal/templates/components/form.html"
check_file "internal/templates/components/form_field.html"

echo

# Валидация шаблонов
log_info "Валидация синтаксиса шаблонов:"
validate_template "internal/templates/layouts/base.html"
validate_template "internal/templates/dashboard.html"
validate_template "internal/templates/components/button.html"

echo

# Проверка CSS
log_info "Проверка CSS:"
if [ -f "internal/templates/static/styles.css" ]; then
    local css_size=$(wc -c < "internal/templates/static/styles.css")
    local css_lines=$(wc -l < "internal/templates/static/styles.css")
    log_success "✓ CSS файл: $css_size байт, $css_lines строк"
fi

# Проверка JavaScript
log_info "Проверка JavaScript:"
if [ -f "internal/templates/static/app.js" ]; then
    local js_size=$(wc -c < "internal/templates/static/app.js")
    local js_lines=$(wc -l < "internal/templates/static/app.js")
    log_success "✓ JavaScript файл: $js_size байт, $js_lines строк"
fi

echo
log_success "Проверка завершена!"
