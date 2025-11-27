#!/bin/bash

set -e  # Выход при ошибке

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Конфигурация
TEMPLATE_DIR="./templates"
STATIC_DIR="./static"
CSS_FILES=("./static/css/styles.css" "./static/css/dark-theme.css")
BACKUP_DIR="./backups"
REPORT_DIR="./reports"

# Создание директорий
mkdir -p "$BACKUP_DIR" "$REPORT_DIR"

log() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Функция для создания бэкапа
create_backup() {
    local file="$1"
    local filename=$(basename "$file")
    local backup_path="$BACKUP_DIR/${filename}.backup_$(date +%Y%m%d_%H%M%S)"
    
    cp "$file" "$backup_path"
    log "Создан бэкап: $backup_path"
}

# Функция для извлечения CSS классов из файла
extract_css_classes() {
    local css_file="$1"
    local classes_file="$REPORT_DIR/css_classes.txt"
    
    if [[ ! -f "$css_file" ]]; then
        warn "CSS файл не найден: $css_file"
        return
    fi
    
    # Извлекаем классы с помощью grep и sed
    grep -Eo '\.[a-zA-Z0-9_-]+\s*[^{]*\{' "$css_file" | \
    sed 's/\.\([a-zA-Z0-9_-]*\)\s*[^{]*{/\1/' | \
    grep -v ':' >> "$classes_file"
}

# Функция для поиска HTML файлов
find_html_files() {
    local dir="$1"
    find "$dir" -name "*.html" -type f
}

# Функция для добавления классов к таблицам
add_table_classes() {
    local file="$1"
    local temp_file="${file}.tmp"
    
    # Добавляем table-container к таблицам
    sed -E 's/<table([^>]*)>/<table\1 class="table-container">/g' "$file" > "$temp_file"
    
    # Проверяем изменения
    if ! cmp -s "$file" "$temp_file"; then
        mv "$temp_file" "$file"
        echo "table"
        return
    fi
    
    rm -f "$temp_file"
    echo ""
}

# Функция для добавления классов к input элементам
add_input_classes() {
    local file="$1"
    local temp_file="${file}.tmp"
    local changes=0
    
    # Создаем временный файл
    cp "$file" "$temp_file"
    
    # Поисковые поля
    sed -i -E 's/<input([^>]*placeholder[^>]*Поиск[^>]*)>/<input\1 class="search-input enhanced-search">/g' "$temp_file"
    
    # Текстовые поля
    sed -i -E 's/<input([^>]*type="text"[^>]*)>/<input\1 class="form-input">/g' "$temp_file"
    
    # Email поля
    sed -i -E 's/<input([^>]*type="email"[^>]*)>/<input\1 class="form-input">/g' "$temp_file"
    
    # Password поля
    sed -i -E 's/<input([^>]*type="password"[^>]*)>/<input\1 class="form-input">/g' "$temp_file"
    
    # Number поля
    sed -i -E 's/<input([^>]*type="number"[^>]*)>/<input\1 class="form-input">/g' "$temp_file"
    
    # Проверяем изменения
    if ! cmp -s "$file" "$temp_file"; then
        changes=1
        mv "$temp_file" "$file"
    else
        rm -f "$temp_file"
    fi
    
    [[ $changes -eq 1 ]] && echo "input" || echo ""
}

# Функция для добавления классов к кнопкам
add_button_classes() {
    local file="$1"
    local temp_file="${file}.tmp"
    local changes=0
    
    cp "$file" "$temp_file"
    
    # Кнопки "Добавить"
    sed -i -E 's/<button([^>]*>[^<]*Добавить[^<]*<\/button>)/<button\1 class="btn btn-primary">/g' "$temp_file"
    sed -i -E 's/<button([^>]*>[^<]*➕[^<]*<\/button>)/<button\1 class="btn btn-primary">/g' "$temp_file"
    
    # Кнопки "Удалить"
    sed -i -E 's/<button([^>]*>[^<]*Удалить[^<]*<\/button>)/<button\1 class="btn btn-danger">/g' "$temp_file"
    sed -i -E 's/<button([^>]*>[^<]*🗑️[^<]*<\/button>)/<button\1 class="btn btn-danger">/g' "$temp_file"
    
    # Кнопки "Редактировать"
    sed -i -E 's/<button([^>]*>[^<]*Редактировать[^<]*<\/button>)/<button\1 class="btn btn-primary">/g' "$temp_file"
    sed -i -E 's/<button([^>]*>[^<]*✏️[^<]*<\/button>)/<button\1 class="btn btn-primary">/g' "$temp_file"
    
    # Остальные кнопки
    sed -i -E 's/<button([^>]*)>/<button\1 class="btn">/g' "$temp_file"
    
    if ! cmp -s "$file" "$temp_file"; then
        changes=1
        mv "$temp_file" "$file"
    else
        rm -f "$temp_file"
    fi
    
    [[ $changes -eq 1 ]] && echo "button" || echo ""
}

# Функция для добавления классов к статусным элементам
add_status_classes() {
    local file="$1"
    local temp_file="${file}.tmp"
    local changes=0
    
    cp "$file" "$temp_file"
    
    # Активный статус
    sed -i -E 's/<span([^>]*)>Активен<\/span>/<span\1 class="status-badge status-active">Активен<\/span>/g' "$temp_file"
    
    # Неактивный статус
    sed -i -E 's/<span([^>]*)>Неактивен<\/span>/<span\1 class="status-badge status-inactive">Неактивен<\/span>/g' "$temp_file"
    
    # Ожидание
    sed -i -E 's/<span([^>]*)>Ожидание<\/span>/<span\1 class="status-badge status-pending">Ожидание<\/span>/g' "$temp_file"
    
    if ! cmp -s "$file" "$temp_file"; then
        changes=1
        mv "$temp_file" "$file"
    else
        rm -f "$temp_file"
    fi
    
    [[ $changes -eq 1 ]] && echo "status" || echo ""
}

# Функция для добавления классов к select элементам
add_select_classes() {
    local file="$1"
    local temp_file="${file}.tmp"
    
    sed -E 's/<select([^>]*)>/<select\1 class="form-select">/g' "$file" > "$temp_file"
    
    if ! cmp -s "$file" "$temp_file"; then
        mv "$temp_file" "$file"
        echo "select"
        return
    fi
    
    rm -f "$temp_file"
    echo ""
}

# Функция для добавления контейнера мобильных карточек
add_mobile_cards() {
    local file="$1"
    
    # Проверяем, есть ли уже мобильные карточки
    if grep -q "mobile-cards" "$file"; then
        echo ""
        return
    fi
    
    # Проверяем, есть ли accounts-table
    if grep -q "id=\"accounts-table\"" "$file"; then
        local mobile_section="\n<!-- Mobile Cards View -->\n<div class=\"mobile-cards\" id=\"mobile-accounts\">\n    <!-- Cards will be generated by JavaScript -->\n</div>"
        
        # Вставляем после accounts-table
        sed -i "/id=\"accounts-table\"/a\\$mobile_section" "$file"
        echo "mobile_cards"
        return
    fi
    
    echo ""
}

# Функция для анализа использования CSS классов
analyze_css_usage() {
    local file="$1"
    local classes_file="$2"
    local analysis_file="$REPORT_DIR/$(basename "$file").analysis"
    
    echo "File: $file" > "$analysis_file"
    echo "=====================" >> "$analysis_file"
    
    local total_classes=0
    local used_classes=0
    
    while IFS= read -r class; do
        [[ -z "$class" ]] && continue
        ((total_classes++))
        
        if grep -q "class=\".*$class" "$file" || grep -q "\"$class\"" "$file"; then
            ((used_classes++))
            echo "✅ $class" >> "$analysis_file"
        else
            echo "❌ $class" >> "$analysis_file"
        fi
    done < "$classes_file"
    
    local coverage=0
    if [[ $total_classes -gt 0 ]]; then
        coverage=$(echo "scale=2; $used_classes * 100 / $total_classes" | bc)
    fi
    
    echo "Coverage: $coverage% ($used_classes/$total_classes)" >> "$analysis_file"
    echo "$file:$coverage:$used_classes:$total_classes" >> "$REPORT_DIR/coverage_summary.txt"
}

# Основная функция
main() {
    log "🚀 Запуск миграции CSS классов..."
    
    # Извлекаем CSS классы
    log "📊 Извлекаем CSS классы..."
    > "$REPORT_DIR/css_classes.txt"
    for css_file in "${CSS_FILES[@]}"; do
        if [[ -f "$css_file" ]]; then
            extract_css_classes "$css_file"
        fi
    done
    
    local total_classes=$(wc -l < "$REPORT_DIR/css_classes.txt" | tr -d ' ')
    log "Найдено CSS классов: $total_classes"
    
    # Ищем HTML файлы
    log "📁 Поиск HTML файлов..."
    local html_files=()
    while IFS= read -r file; do
        html_files+=("$file")
    done < <(find_html_files "$TEMPLATE_DIR")
    
    local total_files=${#html_files[@]}
    log "Найдено HTML файлов: $total_files"
    
    # Обрабатываем файлы
    local processed_files=0
    local total_changes=0
    
    > "$REPORT_DIR/coverage_summary.txt"
    
    for file in "${html_files[@]}"; do
        log "🔧 Обрабатываем: $(basename "$file")"
        
        # Создаем бэкап
        create_backup "$file"
        
        local changes=()
        
        # Применяем преобразования
        changes+=($(add_table_classes "$file"))
        changes+=($(add_input_classes "$file"))
        changes+=($(add_button_classes "$file"))
        changes+=($(add_status_classes "$file"))
        changes+=($(add_select_classes "$file"))
        changes+=($(add_mobile_cards "$file"))
        
        # Фильтруем пустые значения
        local non_empty_changes=()
        for change in "${changes[@]}"; do
            [[ -n "$change" ]] && non_empty_changes+=("$change")
        done
        
        if [[ ${#non_empty_changes[@]} -gt 0 ]]; then
            ((processed_files++))
            total_changes=$((total_changes + ${#non_empty_changes[@]}))
            log "Добавлено изменений: ${#non_empty_changes[@]} (${non_empty_changes[*]})"
        fi
        
        # Анализируем покрытие
        analyze_css_usage "$file" "$REPORT_DIR/css_classes.txt"
    done
    
    # Генерируем отчет
    generate_report "$total_files" "$processed_files" "$total_changes" "$total_classes"
}

# Функция для генерации отчета
generate_report() {
    local total_files=$1
    local processed_files=$2
    local total_changes=$3
    local total_classes=$4
    
    log "📈 Генерация отчета..."
    
    # Анализируем покрытие
    local total_used_classes=0
    local total_coverage=0
    local file_count=0
    
    while IFS=':' read -r file coverage used classes; do
        total_used_classes=$((total_used_classes + used))
        total_coverage=$(echo "scale=2; $total_coverage + $coverage" | bc)
        ((file_count++))
    done < "$REPORT_DIR/coverage_summary.txt"
    
    local avg_coverage=0
    if [[ $file_count -gt 0 ]]; then
        avg_coverage=$(echo "scale=2; $total_coverage / $file_count" | bc)
    fi
    
    local overall_coverage=0
    if [[ $total_classes -gt 0 ]]; then
        overall_coverage=$(echo "scale=2; $total_used_classes * 100 / ($total_classes * $file_count)" | bc)
    fi
    
    # Создаем итоговый отчет
    cat > "$REPORT_DIR/final_report.txt" << EOF
ОТЧЕТ О МИГРАЦИИ CSS КЛАССОВ
============================
Дата: $(date)
Обработано файлов: $processed_files/$total_files
Всего изменений: $total_changes
Всего CSS классов: $total_classes
Среднее покрытие: $avg_coverage%
Общее покрытие: $overall_coverage%

ДЕТАЛИ:
EOF

    # Добавляем статистику по файлам
    echo -e "\nСТАТИСТИКА ПО ФАЙЛАМ:" >> "$REPORT_DIR/final_report.txt"
    while IFS=':' read -r file coverage used classes; do
        echo "  $(basename "$file"): $coverage% ($used/$classes)" >> "$REPORT_DIR/final_report.txt"
    done < "$REPORT_DIR/coverage_summary.txt"
    
    # Показываем итоги
    echo -e "\n${GREEN}✅ МИГРАЦИЯ ЗАВЕРШЕНА${NC}"
    echo -e "${BLUE}📊 ИТОГИ:${NC}"
    echo -e "   Обработано файлов: ${GREEN}$processed_files/${total_files}${NC}"
    echo -e "   Всего изменений: ${GREEN}$total_changes${NC}"
    echo -e "   Общее покрытие CSS: ${GREEN}$overall_coverage%${NC}"
    echo -e "   Отчеты сохранены в: ${YELLOW}$REPORT_DIR/${NC}"
}

# Запуск скрипта
main "$@"
