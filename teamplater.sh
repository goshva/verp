#!/bin/bash

set -e

echo "[INFO] Начало реорганизации шаблонов..."

# Создаем структуру директорий
mkdir -p ./internal/templates/{layouts,components,pages,partials}
mkdir -p ./internal/templates/static/{css,js}

# Функция для создания базового шаблона
create_base_template() {
    cat > ./internal/templates/layouts/base.html << 'EOF'
{{ define "base.html" }}
<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="Cache-Control" content="no-cache, no-store, must-revalidate">
    <meta http-equiv="Pragma" content="no-cache">
    <meta http-equiv="Expires" content="0">
    <title>{{.Title}} - VendERP</title>
    <script src="https://unpkg.com/htmx.org@1.9.6"></script>
    <link rel="stylesheet" href="/static/css/styles.css">
</head>
<body>
    {{ template "sidebar" . }}
    
    <div class="main-content">
        {{ template "content" . }}
    </div>

    {{ template "modal" }}
    {{ template "notifications" }}
    
    <script src="/static/js/app.js"></script>
</body>
</html>
{{ end }}
EOF
    echo "[OK] Базовый шаблон создан"
}

# Создаем компонент сайдбара
create_sidebar() {
    cat > ./internal/templates/components/sidebar.html << 'EOF'
{{ define "sidebar" }}
<div class="sidebar">
    <h2 style="margin-bottom: 2rem; color: var(--primary);">VendERP</h2>
    <nav>
        <a href="/dashboard" class="nav-link {{if eq .Active "dashboard"}}active{{end}}">📊 Дашборд</a>
        <a href="/users" class="nav-link {{if eq .Active "users"}}active{{end}}">👥 Пользователи</a>
        <a href="/machines" class="nav-link {{if eq .Active "machines"}}active{{end}}">🤖 Автоматы</a>
        <a href="/locations" class="nav-link {{if eq .Active "locations"}}active{{end}}">📍 Локации</a>
        <a href="/finance" class="nav-link {{if eq .Active "finance"}}active{{end}}">💰 Финансы</a>
        <a href="/partners" class="nav-link {{if eq .Active "partners"}}active{{end}}">🤝 Партнеры</a>
        <a href="/maintenance" class="nav-link {{if eq .Active "maintenance"}}active{{end}}">🔧 Обслуживание</a>
    </nav>
</div>
{{ end }}
EOF
    echo "[OK] Компонент сайдбара создан"
}

# Создаем компонент модального окна
create_modal() {
    cat > ./internal/templates/components/modal.html << 'EOF'
{{ define "modal" }}
<div id="modal" class="modal">
    <div class="modal-content">
        <button class="modal-close" onclick="VendERP.hideModal()">×</button>
        <div class="modal-header">
            <h3 class="modal-title" id="modal-title">Форма</h3>
        </div>
        <div class="modal-body">
            <div id="modal-body">
                <!-- Контент модального окна будет загружаться здесь -->
            </div>
        </div>
        <div class="modal-footer" id="modal-footer">
            <!-- Кнопки модального окна будут загружаться здесь -->
        </div>
    </div>
</div>
{{ end }}
EOF
    echo "[OK] Компонент модального окна создан"
}

# Создаем компонент уведомлений
create_notifications() {
    cat > ./internal/templates/components/notifications.html << 'EOF'
{{ define "notifications" }}
<div id="notifications" class="notifications-container"></div>
{{ end }}
EOF
    echo "[OK] Компонент уведомлений создан"
}

# Создаем компонент заголовка страницы
create_page_header() {
    cat > ./internal/templates/components/page_header.html << 'EOF'
{{ define "page_header" }}
<div class="header" style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem;">
    <h1>{{.Icon}} {{.Title}}</h1>
    {{if .ShowAddButton}}
    <button class="btn btn-primary"
            hx-get="{{.AddURL}}"
            hx-target="#modal-body"
            onclick="VendERP.showModal()">
        ➕ {{.AddButtonText}}
    </button>
    {{end}}
</div>
{{ end }}
EOF
    echo "[OK] Компонент заголовка страницы создан"
}

# Создаем переработанные страницы
create_users_page() {
    cat > ./internal/templates/pages/users.html << 'EOF'
{{ define "users.html" }}
{{ template "base.html" . }}
{{ end }}

{{ define "content" }}
    {{ template "page_header" dict 
        "Icon" "👥" 
        "Title" "Управление пользователями" 
        "ShowAddButton" true
        "AddURL" "/users/form"
        "AddButtonText" "Добавить пользователя"
    }}

    <div class="card">
        <div id="users-table">
            {{ template "users_list.html" . }}
        </div>
    </div>
{{ end }}
EOF
    echo "[OK] Страница пользователей создана"
}

create_machines_page() {
    cat > ./internal/templates/pages/machines.html << 'EOF'
{{ define "machines.html" }}
{{ template "base.html" . }}
{{ end }}

{{ define "content" }}
    {{ template "page_header" dict 
        "Icon" "🤖" 
        "Title" "Управление автоматами" 
        "ShowAddButton" true
        "AddURL" "/machines/form"
        "AddButtonText" "Добавить автомат"
    }}

    <div class="card">
        <div id="machines-table">
            {{ template "machines_list.html" . }}
        </div>
    </div>
{{ end }}
EOF
    echo "[OK] Страница автоматов создана"
}

create_locations_page() {
    cat > ./internal/templates/pages/locations.html << 'EOF'
{{ define "locations.html" }}
{{ template "base.html" . }}
{{ end }}

{{ define "content" }}
    {{ template "page_header" dict 
        "Icon" "📍" 
        "Title" "Управление локациями" 
        "ShowAddButton" true
        "AddURL" "/locations/form"
        "AddButtonText" "Добавить локацию"
    }}

    <div class="card">
        <div id="locations-table">
            {{ template "locations_list.html" . }}
        </div>
    </div>
{{ end }}
EOF
    echo "[OK] Страница локаций создана"
}

create_dashboard_page() {
    cat > ./internal/templates/pages/dashboard.html << 'EOF'
{{ define "dashboard.html" }}
{{ template "base.html" . }}
{{ end }}

{{ define "content" }}
<div class="header" style="margin-bottom: 2rem;">
    <h1>📊 Панель управления VendERP</h1>
    <p style="color: var(--secondary); margin-top: 0.5rem;">Система управления вендинговыми автоматами</p>
</div>

<div class="stats-grid" id="stats-grid" hx-get="/api/stats" hx-trigger="load">
    {{ template "stat_card" dict "Label" "Всего автоматов" "Value" "-" }}
    {{ template "stat_card" dict "Label" "Активные автоматы" "Value" "-" }}
    {{ template "stat_card" dict "Label" "Всего пользователей" "Value" "-" }}
    {{ template "stat_card" dict "Label" "Всего локаций" "Value" "-" }}
    {{ template "stat_card" dict "Label" "Общая выручка" "Value" "-" }}
    {{ template "stat_card" dict "Label" "Ожидающие задачи" "Value" "-" }}
</div>

<div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem;">
    {{ template "card" dict 
        "Title" "Последние пользователи"
        "Content" `<div id="recent-users" hx-get="/users" hx-trigger="load" hx-target="this">Загрузка...</div>`
    }}
    
    {{ template "card" dict 
        "Title" "Последние автоматы"
        "Content" `<div id="recent-machines" hx-get="/machines" hx-trigger="load" hx-target="this">Загрузка...</div>`
    }}
</div>

{{ template "card" dict 
    "Title" "📍 Активные локации"
    "Content" `
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem;">
        <h3>📍 Активные локации</h3>
        <button class="btn btn-primary"
                hx-get="/locations/form"
                hx-target="#modal-body"
                onclick="VendERP.showModal()">
            ➕ Добавить локацию
        </button>
    </div>
    <div id="locations-table" hx-get="/locations" hx-trigger="load">
        Загрузка локаций...
    </div>`
}}

<script>
    document.addEventListener('DOMContentLoaded', function() {
        setInterval(() => {
            htmx.ajax('GET', '/api/stats', { target: '#stats-grid' });
        }, 30000);
    });
</script>
{{ end }}
EOF
    echo "[OK] Страница дашборда создана"
}

# Создаем компонент карточки статистики
create_stat_card() {
    cat > ./internal/templates/components/stat_card.html << 'EOF'
{{ define "stat_card" }}
<div class="stat-card">
    <div class="stat-label">{{.Label}}</div>
    <div class="stat-number">{{.Value}}</div>
</div>
{{ end }}
EOF
    echo "[OK] Компонент карточки статистики создан"
}

# Создаем улучшенные CSS стили
create_styles() {
    cat > ./internal/templates/static/css/styles.css << 'EOF'
:root {
    --primary: #2563eb;
    --secondary: #64748b;
    --success: #10b981;
    --warning: #f59e0b;
    --danger: #ef4444;
    --light: #f8fafc;
    --dark: #1e293b;
    --border: #e2e8f0;
    --shadow: 0 1px 3px rgba(0,0,0,0.1);
    --radius: 8px;
}

* { 
    margin: 0; 
    padding: 0; 
    box-sizing: border-box; 
}

body { 
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; 
    background: #f1f5f9; 
    line-height: 1.6;
}

/* Layout */
.sidebar { 
    width: 250px; 
    background: white; 
    height: 100vh; 
    position: fixed; 
    padding: 1rem; 
    box-shadow: 2px 0 10px rgba(0,0,0,0.1); 
}

.main-content { 
    margin-left: 250px; 
    padding: 2rem; 
    min-height: 100vh;
}

/* Navigation */
.nav-link { 
    display: block; 
    padding: 0.75rem 1rem; 
    color: var(--dark); 
    text-decoration: none; 
    border-radius: 6px; 
    margin-bottom: 0.5rem; 
    transition: all 0.2s ease;
}

.nav-link:hover { 
    background: var(--light); 
    transform: translateX(4px);
}

.nav-link.active { 
    background: var(--primary); 
    color: white; 
}

/* Cards */
.card { 
    background: white; 
    border-radius: var(--radius); 
    padding: 1.5rem; 
    box-shadow: var(--shadow); 
    margin-bottom: 1.5rem; 
    border: 1px solid var(--border);
}

.card-header {
    border-bottom: 1px solid var(--border);
    padding-bottom: 1rem;
    margin-bottom: 1rem;
}

.card-title {
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--dark);
}

/* Buttons */
.btn { 
    padding: 0.5rem 1rem; 
    border: none; 
    border-radius: 6px; 
    cursor: pointer; 
    text-decoration: none; 
    display: inline-block; 
    font-size: 0.875rem;
    font-weight: 500;
    transition: all 0.2s ease;
    border: 1px solid transparent;
}

.btn:hover {
    transform: translateY(-1px);
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.btn-primary { 
    background: var(--primary); 
    color: white; 
}

.btn-primary:hover {
    background: #1d4ed8;
}

.btn-success { 
    background: var(--success); 
    color: white; 
}

.btn-danger { 
    background: var(--danger); 
    color: white; 
}

.btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
    transform: none;
}

/* Tables */
.table { 
    width: 100%; 
    border-collapse: collapse; 
    font-size: 0.875rem;
}

.table th, .table td { 
    padding: 0.75rem; 
    text-align: left; 
    border-bottom: 1px solid var(--border); 
}

.table th { 
    background: var(--light); 
    font-weight: 600; 
    color: var(--dark);
    font-size: 0.75rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
}

.table tr:hover {
    background: #f8fafc;
}

/* Forms */
.form-group { 
    margin-bottom: 1rem; 
}

.form-label { 
    display: block; 
    margin-bottom: 0.5rem; 
    font-weight: 500; 
    color: var(--dark);
    font-size: 0.875rem;
}

.form-input, 
.form-select, 
.form-textarea { 
    width: 100%; 
    padding: 0.5rem; 
    border: 1px solid #d1d5db; 
    border-radius: 4px; 
    font-size: 0.875rem;
    transition: border-color 0.2s ease;
}

.form-input:focus,
.form-select:focus,
.form-textarea:focus {
    outline: none;
    border-color: var(--primary);
    box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
}

.form-select { 
    background: white; 
}

.form-textarea { 
    min-height: 80px; 
    resize: vertical; 
}

.form-help {
    font-size: 0.75rem;
    color: var(--secondary);
    margin-top: 0.25rem;
}

/* Modal */
.modal {
    display: none;
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0,0,0,0.5);
    z-index: 1000;
    backdrop-filter: blur(4px);
}

.modal.show {
    display: flex;
    align-items: center;
    justify-content: center;
}

.modal-content {
    background: white;
    padding: 2rem;
    border-radius: var(--radius);
    width: 90%;
    max-width: 600px;
    max-height: 90vh;
    overflow-y: auto;
    position: relative;
    box-shadow: 0 20px 25px -5px rgba(0,0,0,0.1);
}

.modal-close {
    position: absolute;
    top: 1rem;
    right: 1rem;
    background: none;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    color: var(--secondary);
    width: 32px;
    height: 32px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
}

.modal-close:hover {
    background: var(--light);
}

.modal-header {
    margin-bottom: 1.5rem;
    padding-right: 2rem;
}

.modal-title {
    font-size: 1.5rem;
    font-weight: 600;
    color: var(--dark);
}

.modal-footer {
    margin-top: 2rem;
    display: flex;
    gap: 1rem;
    justify-content: flex-end;
}

/* Status badges */
.status-badge { 
    padding: 0.25rem 0.75rem; 
    border-radius: 9999px; 
    font-size: 0.75rem; 
    font-weight: 500; 
    text-transform: uppercase;
    letter-spacing: 0.05em;
}

.status-active { 
    background: #dcfce7; 
    color: #166534; 
}

.status-inactive { 
    background: #f3f4f6; 
    color: #374151; 
}

.status-pending { 
    background: #fef3c7; 
    color: #92400e; 
}

/* Statistics */
.stats-grid { 
    display: grid; 
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); 
    gap: 1.5rem; 
    margin-bottom: 2rem; 
}

.stat-card { 
    background: white; 
    padding: 1.5rem; 
    border-radius: var(--radius); 
    box-shadow: var(--shadow); 
    border: 1px solid var(--border);
    transition: transform 0.2s ease;
}

.stat-card:hover {
    transform: translateY(-2px);
}

.stat-number { 
    font-size: 2rem; 
    font-weight: bold; 
    margin: 0.5rem 0; 
    color: var(--primary);
}

.stat-label { 
    color: var(--secondary); 
    font-size: 0.875rem; 
    font-weight: 500;
}

/* Notifications */
.notifications-container {
    position: fixed;
    top: 1rem;
    right: 1rem;
    z-index: 1100;
}

.notification {
    background: white;
    padding: 1rem;
    border-radius: var(--radius);
    box-shadow: 0 10px 15px -3px rgba(0,0,0,0.1);
    border-left: 4px solid var(--primary);
    margin-bottom: 0.5rem;
    min-width: 300px;
    animation: slideIn 0.3s ease;
}

.notification.success {
    border-left-color: var(--success);
}

.notification.error {
    border-left-color: var(--danger);
}

.notification.warning {
    border-left-color: var(--warning);
}

@keyframes slideIn {
    from {
        transform: translateX(100%);
        opacity: 0;
    }
    to {
        transform: translateX(0);
        opacity: 1;
    }
}

/* Responsive */
@media (max-width: 768px) {
    .sidebar {
        width: 100%;
        height: auto;
        position: relative;
    }
    
    .main-content {
        margin-left: 0;
        padding: 1rem;
    }
    
    .stats-grid {
        grid-template-columns: 1fr;
    }
    
    .modal-content {
        margin: 1rem;
        width: calc(100% - 2rem);
    }
}

/* Loading states */
.htmx-request {
    opacity: 0.7;
    pointer-events: none;
}

.loading {
    display: inline-block;
    width: 20px;
    height: 20px;
    border: 2px solid #f3f3f3;
    border-top: 2px solid var(--primary);
    border-radius: 50%;
    animation: spin 1s linear infinite;
}

@keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
}
EOF
    echo "[OK] CSS стили созданы"
}

# Создаем улучшенный JavaScript
create_javascript() {
    cat > ./internal/templates/static/js/app.js << 'EOF'
// VendERP Global Namespace
const VendERP = {
    // Modal functions
    showModal: function(title = 'Форма') {
        const modal = document.getElementById('modal');
        const modalTitle = document.getElementById('modal-title');
        
        if (modalTitle) {
            modalTitle.textContent = title;
        }
        
        modal.classList.add('show');
        document.body.style.overflow = 'hidden';
    },

    hideModal: function() {
        const modal = document.getElementById('modal');
        const modalBody = document.getElementById('modal-body');
        const modalFooter = document.getElementById('modal-footer');
        
        modal.classList.remove('show');
        document.body.style.overflow = '';
        
        // Очищаем содержимое модального окна
        if (modalBody) modalBody.innerHTML = '';
        if (modalFooter) modalFooter.innerHTML = '';
    },

    // Notification functions
    showNotification: function(message, type = 'info', duration = 5000) {
        const container = document.getElementById('notifications');
        if (!container) return;

        const notification = document.createElement('div');
        notification.className = `notification ${type}`;
        notification.innerHTML = `
            <div style="display: flex; justify-content: between; align-items: start;">
                <div style="flex: 1;">${message}</div>
                <button onclick="this.parentElement.parentElement.remove()" 
                        style="background: none; border: none; font-size: 1.25rem; cursor: pointer; color: var(--secondary); margin-left: 1rem;">
                    ×
                </button>
            </div>
        `;

        container.appendChild(notification);

        // Auto remove after duration
        if (duration > 0) {
            setTimeout(() => {
                if (notification.parentElement) {
                    notification.remove();
                }
            }, duration);
        }
    },

    // Form handling
    handleFormResponse: function(evt) {
        const targetId = evt.detail.target.id;
        
        // Если форма успешно отправлена и target - это таблица, закрываем модальное окно
        if (targetId && targetId.includes('-table') && !evt.detail.xhr.response) {
            VendERP.hideModal();
            VendERP.showNotification('Операция выполнена успешно', 'success');
        }
    },

    // Utility functions
    formatCurrency: function(amount) {
        return new Intl.NumberFormat('ru-RU', {
            style: 'currency',
            currency: 'RUB'
        }).format(amount);
    },

    formatDate: function(dateString) {
        if (!dateString) return '-';
        return new Date(dateString).toLocaleDateString('ru-RU');
    },

    debounce: function(func, wait) {
        let timeout;
        return function executedFunction(...args) {
            const later = () => {
                clearTimeout(timeout);
                func(...args);
            };
            clearTimeout(timeout);
            timeout = setTimeout(later, wait);
        };
    }
};

// Event Listeners
document.addEventListener('DOMContentLoaded', function() {
    // Закрытие модального окна при клике вне его
    document.addEventListener('click', function(e) {
        const modal = document.getElementById('modal');
        if (e.target === modal) {
            VendERP.hideModal();
        }
    });

    // Закрытие модального окна при нажатии Escape
    document.addEventListener('keydown', function(e) {
        if (e.key === 'Escape') {
            VendERP.hideModal();
        }
    });

    // Показываем модальное окно при загрузке формы через htmx
    document.addEventListener('htmx:afterSwap', function(evt) {
        if (evt.detail.target.id === 'modal-body' && evt.detail.xhr.response) {
            VendERP.showModal();
        }
        
        // Обработка успешных ответов форм
        VendERP.handleFormResponse(evt);
    });

    // Обработка ошибок htmx
    document.addEventListener('htmx:responseError', function(evt) {
        VendERP.showNotification('Произошла ошибка при выполнении запроса', 'error');
    });

    // Подтверждение удаления
    document.addEventListener('click', function(e) {
        if (e.target.hasAttribute('hx-delete') && !e.target.hasAttribute('hx-confirm')) {
            e.preventDefault();
            const message = e.target.getAttribute('data-confirm') || 'Вы уверены, что хотите удалить этот элемент?';
            if (confirm(message)) {
                htmx.trigger(e.target, 'htmx:confirm');
            }
        }
    });
});

// HTMX конфигурация
htmx.defineExtension('debug', {
    onEvent: function (name, evt) {
        if (console.debug) {
            console.debug(name, evt);
        }
    }
});
EOF
    echo "[OK] JavaScript создан"
}

# Создаем утилитарные компоненты
create_utility_components() {
    # Компонент карточки
    cat > ./internal/templates/components/card.html << 'EOF'
{{ define "card" }}
<div class="card {{.Class}}">
    {{if .Title}}
    <div class="card-header">
        <h3 class="card-title">{{.Title}}</h3>
    </div>
    {{end}}
    <div class="card-body">
        {{.Content}}
    </div>
</div>
{{ end }}
EOF

    # Компонент кнопки
    cat > ./internal/templates/components/button.html << 'EOF'
{{ define "button" }}
<button class="btn btn-{{.Variant}} {{.Class}}"
        {{if .ID}}id="{{.ID}}"{{end}}
        {{if .HXGet}}hx-get="{{.HXGet}}"{{end}}
        {{if .HXPost}}hx-post="{{.HXPost}}"{{end}}
        {{if .HXTarget}}hx-target="{{.HXTarget}}"{{end}}
        {{if .HXTrigger}}hx-trigger="{{.HXTrigger}}"{{end}}
        {{if .OnClick}}onclick="{{.OnClick}}"{{end}}
        {{if .Disabled}}disabled{{end}}>
    {{if .Icon}}{{.Icon}} {{end}}{{.Text}}
</button>
{{ end }}
EOF

    # Компонент поля формы
    cat > ./internal/templates/components/form_field.html << 'EOF'
{{ define "form_field" }}
<div class="form-group">
    <label class="form-label" for="{{.ID}}">{{.Label}}{{if .Required}} *{{end}}</label>
    {{if eq .Type "select"}}
    <select class="form-select" id="{{.ID}}" name="{{.Name}}" {{if .Required}}required{{end}}>
        <option value="">Выберите...</option>
        {{range .Options}}
        <option value="{{.Value}}" {{if .Selected}}selected{{end}}>{{.Text}}</option>
        {{end}}
    </select>
    {{else if eq .Type "textarea"}}
    <textarea class="form-textarea" id="{{.ID}}" name="{{.Name}}"
              {{if .Required}}required{{end}}
              {{if .Placeholder}}placeholder="{{.Placeholder}}"{{end}}
              {{if .Rows}}rows="{{.Rows}}"{{end}}>{{.Value}}</textarea>
    {{else}}
    <input type="{{.Type}}" class="form-input" id="{{.ID}}" name="{{.Name}}"
           value="{{.Value}}"
           {{if .Required}}required{{end}}
           {{if .Placeholder}}placeholder="{{.Placeholder}}"{{end}}>
    {{end}}
    {{if .HelpText}}
    <div class="form-help">{{.HelpText}}</div>
    {{end}}
</div>
{{ end }}
EOF

    echo "[OK] Утилитарные компоненты созданы"
}

# Основная функция миграции
migrate_templates() {
    echo "[INFO] Начало миграции шаблонов..."
    
    # Создаем базовые компоненты
    create_base_template
    create_sidebar
    create_modal
    create_notifications
    create_page_header
    create_stat_card
    
    # Создаем страницы
    create_users_page
    create_machines_page
    create_locations_page
    create_dashboard_page
    
    # Создаем статические файлы
    create_styles
    create_javascript
    create_utility_components
    
    # Копируем существующие partials (списки и формы)
    echo "[INFO] Копирование существующих partials..."
    cp ./internal/templates/users_list.html ./internal/templates/partials/
    cp ./internal/templates/machines_list.html ./internal/templates/partials/
    cp ./internal/templates/locations_list.html ./internal/templates/partials/
    cp ./internal/templates/user_form.html ./internal/templates/partials/
    cp ./internal/templates/machine_form.html ./internal/templates/partials/
    cp ./internal/templates/location_form.html ./internal/templates/partials/
    
    echo "[SUCCESS] Миграция завершена успешно!"
    echo ""
    echo "Новая структура шаблонов:"
    echo "├── layouts/"
    echo "│   └── base.html          (базовый шаблон)"
    echo "├── components/"
    echo "│   ├── sidebar.html       (навигация)"
    echo "│   ├── modal.html         (модальное окно)"
    echo "│   ├── notifications.html (уведомления)"
    echo "│   ├── page_header.html   (заголовок страницы)"
    echo "│   ├── stat_card.html     (карточка статистики)"
    echo "│   ├── card.html          (универсальная карточка)"
    echo "│   ├── button.html        (компонент кнопки)"
    echo "│   └── form_field.html    (поле формы)"
    echo "├── pages/"
    echo "│   ├── dashboard.html     (дашборд)"
    echo "│   ├── users.html         (пользователи)"
    echo "│   ├── machines.html      (автоматы)"
    echo "│   └── locations.html     (локации)"
    echo "├── partials/"
    echo "│   ├── *_list.html        (списки элементов)"
    echo "│   └── *_form.html        (формы)"
    echo "└── static/"
    echo "    ├── css/"
    echo "    │   └── styles.css     (стили)"
    echo "    └── js/"
    echo "        └── app.js         (JavaScript)"
    echo ""
    echo "Преимущества новой структуры:"
    echo "✅ Переиспользование компонентов"
    echo "✅ Единый источник истины для стилей и скриптов"
    echo "✅ Упрощенное обслуживание"
    echo "✅ Лучшая производительность"
    echo "✅ Легкое расширение"
}

# Запуск миграции
migrate_templates
