# Удалите старые файлы
rm internal/templates/locations.html internal/templates/locations_list.html

# Создайте locations.html
cat > internal/templates/locations.html << 'EOF'
{{ define "locations.html" }}
{{ template "base.html" . }}
{{ end }}

{{ define "content" }}
<div class="header" style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem;">
    <h1>📍 Управление локациями</h1>
    <button class="btn btn-primary" 
            hx-get="/locations/form" 
            hx-target="#modal-body"
            onclick="showModal()">
        ➕ Добавить локацию
    </button>
</div>

<div class="card">
    <div id="locations-table">
        {{ template "locations_list.html" . }}
    </div>
</div>
{{ end }}
EOF

# Создайте locations_list.html
cat > internal/templates/locations_list.html << 'EOF'
{{ define "locations_list.html" }}
<table class="table">
    <thead>
        <tr>
            <th>ID</th>
            <th>Название</th>
            <th>Адрес</th>
            <th>Контактное лицо</th>
            <th>Телефон</th>
            <th>Аренда (₽)</th>
            <th>День оплаты</th>
            <th>Статус</th>
            <th>Действия</th>
        </tr>
    </thead>
    <tbody>
        {{range .Locations}}
        <tr>
            <td>{{.ID}}</td>
            <td>{{.Name}}</td>
            <td>{{.Address}}</td>
            <td>{{.ContactPerson}}</td>
            <td>{{.ContactPhone}}</td>
            <td>{{.MonthlyRent}} ₽</td>
            <td>{{.RentDueDay}}</td>
            <td>
                <span class="status-badge {{if .IsActive}}status-active{{else}}status-inactive{{end}}">
                    {{if .IsActive}}Активна{{else}}Неактивна{{end}}
                </span>
            </td>
            <td>
                <div style="display: flex; gap: 0.5rem;">
                    <button class="btn btn-primary" 
                            hx-get="/locations/form?id={{.ID}}"
                            hx-target="#modal-body"
                            _="on htmx:afterOnLoad call #modal.showModal()">
                        ✏️
                    </button>
                    <button class="btn btn-danger" 
                            hx-delete="/locations/delete?id={{.ID}}"
                            hx-target="#locations-table"
                            hx-confirm="Удалить локацию?">
                        🗑️
                    </button>
                </div>
            </td>
        </tr>
        {{else}}
        <tr>
            <td colspan="9" style="text-align: center; padding: 2rem; color: var(--secondary);">
                Нет локаций. <a href="#" hx-get="/locations/form" hx-target="#modal-body">Добавить первую локацию</a>
            </td>
        </tr>
        {{end}}
    </tbody>
</table>
{{ end }}
EOF
