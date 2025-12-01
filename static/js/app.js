// VendERP Application JavaScript
const VendERP = {
    // Charts module
    charts: {
        machines: {
            data: null,
            miniChart: null,
            fullChart: null,
            
            // Загрузка данных
            load: async function() {
                try {
                    console.log('Loading machines chart data...');
                    const res = await fetch('/api/charts/machines');
                    if (!res.ok) {
                        throw new Error(`HTTP error! status: ${res.status}`);
                    }
                    this.data = await res.json();
                    console.log('Machines chart data loaded:', this.data);
                    this.updateUI();
                } catch (err) {
                    console.error('Ошибка загрузки графика автоматов:', err);
                    this.showError();
                }
            },
            
            // Обновление интерфейса
            updateUI: function() {
                if (!this.data || !this.data.series || this.data.series.length === 0) {
                    console.warn('No data for machines chart');
                    this.showError();
                    return;
                }
                
                const series = this.data.series[0];
                const labels = this.data.labels || [];
                const counts = series.data ? series.data.map(d => d.count || d.value || 0) : [];
                
                console.log('Updating machines UI with:', { 
                    total: this.data.total, 
                    labels: labels.length,
                    data: counts.length 
                });
                
                // Мини-график
                const totalElement = document.getElementById('machines-mini-total');
                if (totalElement) {
                    totalElement.textContent = this.data.total || '-';
                }
                
                this.updateMiniChart(labels, counts, series.color);
                
                // Полноэкранный график
                const fullTotal = document.getElementById('machines-full-total');
                if (fullTotal) {
                    fullTotal.textContent = this.data.total || '-';
                }
                
                const periodElement = document.getElementById('machines-full-period');
                if (periodElement) {
                    periodElement.textContent = this.data.period || '30 дней';
                }
                
                const changeElement = document.getElementById('machines-full-change');
                if (changeElement) {
                    const change = this.data.change || 0;
                    const changeText = change >= 0 ? `+${change}` : change;
                    changeElement.textContent = changeText;
                    changeElement.className = `stat-value ${change >= 0 ? 'positive' : 'negative'}`;
                }
                
                const trendElement = document.getElementById('machines-full-trend');
                if (trendElement) {
                    const trend = this.data.trend || 0;
                    const trendInfo = this.getTrendInfo(trend);
                    trendElement.innerHTML = `${trendInfo.icon} ${trendInfo.text}`;
                    trendElement.className = `stat-trend ${trendInfo.class}`;
                }
                
                const infoElement = document.getElementById('machines-data-info');
                if (infoElement) {
                    infoElement.textContent = 
                        `Данные за ${this.data.period || '30 дней'} • Обновлено: ${new Date().toLocaleTimeString('ru-RU')}`;
                }
            },
            
            // Мини-график
            updateMiniChart: function(labels, data, color) {
                const canvas = document.getElementById('machines-mini-chart');
                if (!canvas) {
                    console.warn('Machines mini chart canvas not found');
                    return;
                }
                
                const ctx = canvas.getContext('2d');
                
                if (this.miniChart) this.miniChart.destroy();
                
                // Если данных нет, показываем placeholder
                if (data.length === 0) {
                    console.log('No data for mini chart, showing placeholder');
                    this.miniChart = new Chart(ctx, {
                        type: 'line',
                        data: {
                            labels: ['', '', '', '', '', ''],
                            datasets: [{
                                data: [1, 2, 1, 3, 2, 1],
                                borderColor: '#E5E7EB',
                                backgroundColor: 'rgba(229, 231, 235, 0.2)',
                                borderWidth: 1,
                                fill: true,
                                tension: 0.4,
                                pointRadius: 0
                            }]
                        },
                        options: {
                            responsive: true,
                            maintainAspectRatio: false,
                            plugins: { 
                                legend: { display: false },
                                tooltip: { enabled: false }
                            },
                            scales: { 
                                x: { display: false }, 
                                y: { display: false }
                            }
                        }
                    });
                    return;
                }
                
                // Фильтруем данные для мини-графика (каждую 3-ю точку)
                const filteredLabels = [];
                const filteredData = [];
                for (let i = 0; i < labels.length; i++) {
                    if (i % 3 === 0 || i === labels.length - 1) {
                        filteredLabels.push(labels[i]);
                        filteredData.push(data[i] || 0);
                    }
                }
                
                this.miniChart = new Chart(ctx, {
                    type: 'line',
                    data: {
                        labels: filteredLabels,
                        datasets: [{
                            data: filteredData,
                            borderColor: color || '#4F46E5',
                            backgroundColor: (color || '#4F46E5') + '20',
                            borderWidth: 1.5,
                            fill: false,
                            tension: 0.3,
                            pointRadius: 0
                        }]
                    },
                    options: {
                        responsive: true,
                        maintainAspectRatio: false,
                        plugins: { 
                            legend: { display: false },
                            tooltip: {
                                enabled: true,
                                callbacks: {
                                    label: function(context) {
                                        return `${context.raw} автоматов`;
                                    }
                                }
                            }
                        },
                        scales: { 
                            x: { 
                                display: false,
                                grid: { display: false }
                            }, 
                            y: { 
                                display: false,
                                grid: { display: false }
                            }
                        }
                    }
                });
            },
            
            // Полноэкранный график
            updateFullChart: function() {
                const canvas = document.getElementById('machines-full-chart');
                if (!canvas || !this.data || !this.data.series || this.data.series.length === 0) {
                    console.warn('Cannot update full chart - missing data or canvas');
                    return;
                }
                
                const ctx = canvas.getContext('2d');
                const series = this.data.series[0];
                const labels = this.data.labels || [];
                const counts = series.data ? series.data.map(d => d.count || d.value || 0) : [];
                const dates = series.data ? series.data.map(d => d.date || '') : [];
                
                if (this.fullChart) this.fullChart.destroy();
                
                const gradient = ctx.createLinearGradient(0, 0, 0, canvas.height);
                gradient.addColorStop(0, (series.color || '#4F46E5') + 'CC');
                gradient.addColorStop(1, (series.color || '#4F46E5') + '22');
                
                this.fullChart = new Chart(ctx, {
                    type: 'line',
                    data: {
                        labels: labels,
                        datasets: [{
                            label: 'Автоматы',
                            data: counts,
                            borderColor: series.color || '#4F46E5',
                            backgroundColor: gradient,
                            borderWidth: 3,
                            fill: true,
                            tension: 0.3,
                            pointBackgroundColor: series.color || '#4F46E5',
                            pointBorderColor: '#fff',
                            pointBorderWidth: 2,
                            pointRadius: 4,
                            pointHoverRadius: 6
                        }]
                    },
                    options: {
                        responsive: true,
                        maintainAspectRatio: false,
                        plugins: {
                            legend: { display: false },
                            tooltip: {
                                callbacks: {
                                    label: function(context) {
                                        const date = dates[context.dataIndex] || '';
                                        const label = date ? `${date}: ` : '';
                                        return `${label}${context.raw} автоматов`;
                                    }
                                }
                            }
                        },
                        scales: {
                            x: { 
                                grid: { display: false },
                                ticks: {
                                    maxTicksLimit: 10
                                }
                            },
                            y: { 
                                beginAtZero: true,
                                ticks: {
                                    callback: function(value) {
                                        return value;
                                    }
                                }
                            }
                        }
                    }
                });
            },
            
            // Разворачивание/сворачивание
            expand: function() {
                const fullscreen = document.getElementById('machines-chart-fullscreen');
                if (fullscreen) {
                    fullscreen.style.display = 'block';
                    document.body.style.overflow = 'hidden';
                    setTimeout(() => this.updateFullChart(), 100);
                }
            },
            
            collapse: function() {
                const fullscreen = document.getElementById('machines-chart-fullscreen');
                if (fullscreen) {
                    fullscreen.style.display = 'none';
                    document.body.style.overflow = 'auto';
                }
            },
            
            // Вспомогательные функции
            getTrendInfo: function(trend) {
                if (trend === 1) return { class: 'up', icon: '📈', text: 'Рост' };
                if (trend === -1) return { class: 'down', icon: '📉', text: 'Спад' };
                return { class: 'stable', icon: '➡️', text: 'Стабильно' };
            },
            
            refresh: function() {
                console.log('Refreshing machines chart...');
                if (this.miniChart) this.miniChart.destroy();
                if (this.fullChart) this.fullChart.destroy();
                this.load();
            },
            
            showError: function() {
                const elem = document.getElementById('machines-mini-total');
                if (elem) {
                    elem.textContent = 'Ошибка';
                    elem.style.color = 'var(--danger)';
                }
            }
        },
        
        operations: {
            data: null,
            miniChart: null,
            fullChart: null,
            
            // Загрузка данных
            load: async function() {
                try {
                    console.log('Loading operations chart data...');
                    const res = await fetch('/api/charts/operations');
                    if (!res.ok) {
                        throw new Error(`HTTP error! status: ${res.status}`);
                    }
                    this.data = await res.json();
                    console.log('Operations chart data loaded:', this.data);
                    this.updateUI();
                } catch (err) {
                    console.error('Ошибка загрузки графика операций:', err);
                    this.showError();
                }
            },
            
            // Обновление интерфейса
            updateUI: function() {
                if (!this.data || !this.data.series || this.data.series.length === 0) {
                    console.warn('No data for operations chart');
                    this.showError();
                    return;
                }
                
                const series = this.data.series[0];
                const labels = this.data.labels || [];
                const counts = series.data ? series.data.map(d => d.count || d.value || 0) : [];
                
                console.log('Updating operations UI with:', { 
                    total: this.data.total, 
                    labels: labels.length,
                    data: counts.length 
                });
                
                // Мини-график
                const totalElement = document.getElementById('operations-mini-total');
                if (totalElement) {
                    totalElement.textContent = this.data.total || '-';
                }
                
                this.updateMiniChart(labels, counts, series.color);
                
                // Полноэкранный график
                const fullTotal = document.getElementById('operations-full-total');
                if (fullTotal) {
                    fullTotal.textContent = this.data.total || '-';
                }
                
                const periodElement = document.getElementById('operations-full-period');
                if (periodElement) {
                    periodElement.textContent = this.data.period || '30 дней';
                }
                
                const changeElement = document.getElementById('operations-full-change');
                if (changeElement) {
                    const change = this.data.change || 0;
                    const changeText = change >= 0 ? `+${change}` : change;
                    changeElement.textContent = changeText;
                    changeElement.className = `stat-value ${change >= 0 ? 'positive' : 'negative'}`;
                }
                
                const trendElement = document.getElementById('operations-full-trend');
                if (trendElement) {
                    const trend = this.data.trend || 0;
                    const trendInfo = this.getTrendInfo(trend);
                    trendElement.innerHTML = `${trendInfo.icon} ${trendInfo.text}`;
                    trendElement.className = `stat-trend ${trendInfo.class}`;
                }
                
                const infoElement = document.getElementById('operations-data-info');
                if (infoElement) {
                    infoElement.textContent = 
                        `Данные за ${this.data.period || '30 дней'} • Обновлено: ${new Date().toLocaleTimeString('ru-RU')}`;
                }
            },
            
            // Мини-график
            updateMiniChart: function(labels, data, color) {
                const canvas = document.getElementById('operations-mini-chart');
                if (!canvas) {
                    console.warn('Operations mini chart canvas not found');
                    return;
                }
                
                const ctx = canvas.getContext('2d');
                
                if (this.miniChart) this.miniChart.destroy();
                
                // Если данных нет, показываем placeholder
                if (data.length === 0) {
                    console.log('No data for operations mini chart, showing placeholder');
                    this.miniChart = new Chart(ctx, {
                        type: 'line',
                        data: {
                            labels: ['', '', '', '', '', ''],
                            datasets: [{
                                data: [1, 2, 1, 3, 2, 1],
                                borderColor: '#E5E7EB',
                                backgroundColor: 'rgba(229, 231, 235, 0.2)',
                                borderWidth: 1,
                                fill: true,
                                tension: 0.4,
                                pointRadius: 0
                            }]
                        },
                        options: {
                            responsive: true,
                            maintainAspectRatio: false,
                            plugins: { 
                                legend: { display: false },
                                tooltip: { enabled: false }
                            },
                            scales: { 
                                x: { display: false }, 
                                y: { display: false }
                            }
                        }
                    });
                    return;
                }
                
                // Фильтруем данные для мини-графика (каждую 3-ю точку)
                const filteredLabels = [];
                const filteredData = [];
                for (let i = 0; i < labels.length; i++) {
                    if (i % 3 === 0 || i === labels.length - 1) {
                        filteredLabels.push(labels[i]);
                        filteredData.push(data[i] || 0);
                    }
                }
                
                this.miniChart = new Chart(ctx, {
                    type: 'line',
                    data: {
                        labels: filteredLabels,
                        datasets: [{
                            data: filteredData,
                            borderColor: color || '#10B981',
                            backgroundColor: (color || '#10B981') + '20',
                            borderWidth: 1.5,
                            fill: false,
                            tension: 0.3,
                            pointRadius: 0
                        }]
                    },
                    options: {
                        responsive: true,
                        maintainAspectRatio: false,
                        plugins: { 
                            legend: { display: false },
                            tooltip: {
                                enabled: true,
                                callbacks: {
                                    label: function(context) {
                                        return `${context.raw} операций`;
                                    }
                                }
                            }
                        },
                        scales: { 
                            x: { 
                                display: false,
                                grid: { display: false }
                            }, 
                            y: { 
                                display: false,
                                grid: { display: false }
                            }
                        }
                    }
                });
            },
            
            // Полноэкранный график
            updateFullChart: function() {
                const canvas = document.getElementById('operations-full-chart');
                if (!canvas || !this.data || !this.data.series || this.data.series.length === 0) {
                    console.warn('Cannot update operations full chart - missing data or canvas');
                    return;
                }
                
                const ctx = canvas.getContext('2d');
                const series = this.data.series[0];
                const labels = this.data.labels || [];
                const counts = series.data ? series.data.map(d => d.count || d.value || 0) : [];
                const dates = series.data ? series.data.map(d => d.date || '') : [];
                
                if (this.fullChart) this.fullChart.destroy();
                
                const gradient = ctx.createLinearGradient(0, 0, 0, canvas.height);
                gradient.addColorStop(0, (series.color || '#10B981') + 'CC');
                gradient.addColorStop(1, (series.color || '#10B981') + '22');
                
                this.fullChart = new Chart(ctx, {
                    type: 'line',
                    data: {
                        labels: labels,
                        datasets: [{
                            label: 'Операции',
                            data: counts,
                            borderColor: series.color || '#10B981',
                            backgroundColor: gradient,
                            borderWidth: 3,
                            fill: true,
                            tension: 0.3,
                            pointBackgroundColor: series.color || '#10B981',
                            pointBorderColor: '#fff',
                            pointBorderWidth: 2,
                            pointRadius: 4,
                            pointHoverRadius: 6
                        }]
                    },
                    options: {
                        responsive: true,
                        maintainAspectRatio: false,
                        plugins: {
                            legend: { display: false },
                            tooltip: {
                                callbacks: {
                                    label: function(context) {
                                        const date = dates[context.dataIndex] || '';
                                        const label = date ? `${date}: ` : '';
                                        return `${label}${context.raw} операций`;
                                    }
                                }
                            }
                        },
                        scales: {
                            x: { 
                                grid: { display: false },
                                ticks: {
                                    maxTicksLimit: 10
                                }
                            },
                            y: { 
                                beginAtZero: true,
                                ticks: {
                                    callback: function(value) {
                                        return value;
                                    }
                                }
                            }
                        }
                    }
                });
            },
            
            // Разворачивание/сворачивание
            expand: function() {
                const fullscreen = document.getElementById('operations-chart-fullscreen');
                if (fullscreen) {
                    fullscreen.style.display = 'block';
                    document.body.style.overflow = 'hidden';
                    setTimeout(() => this.updateFullChart(), 100);
                }
            },
            
            collapse: function() {
                const fullscreen = document.getElementById('operations-chart-fullscreen');
                if (fullscreen) {
                    fullscreen.style.display = 'none';
                    document.body.style.overflow = 'auto';
                }
            },
            
            // Вспомогательные функции
            getTrendInfo: function(trend) {
                if (trend === 1) return { class: 'up', icon: '📈', text: 'Рост' };
                if (trend === -1) return { class: 'down', icon: '📉', text: 'Спад' };
                return { class: 'stable', icon: '➡️', text: 'Стабильно' };
            },
            
            refresh: function() {
                console.log('Refreshing operations chart...');
                if (this.miniChart) this.miniChart.destroy();
                if (this.fullChart) this.fullChart.destroy();
                this.load();
            },
            
            showError: function() {
                const elem = document.getElementById('operations-mini-total');
                if (elem) {
                    elem.textContent = 'Ошибка';
                    elem.style.color = 'var(--danger)';
                }
            }
        },
        
        // Initialize all charts
        init: function() {
            console.log('Initializing charts...');
            
            // Проверяем, есть ли контейнеры для графиков на странице
            const hasMachinesChart = document.getElementById('machines-chart-mini');
            const hasOperationsChart = document.getElementById('operations-chart-mini');
            
            // Инициализируем Chart.js если не загружен
            if ((hasMachinesChart || hasOperationsChart) && typeof Chart === 'undefined') {
                console.log('Loading Chart.js...');
                const s = document.createElement('script');
                s.src = 'https://cdn.jsdelivr.net/npm/chart.js';
                s.onload = () => {
                    console.log('Chart.js loaded, initializing charts...');
                    if (hasMachinesChart) this.machines.load();
                    if (hasOperationsChart) this.operations.load();
                };
                s.onerror = (err) => {
                    console.error('Failed to load Chart.js:', err);
                };
                document.head.appendChild(s);
            } else {
                console.log('Chart.js already loaded or not needed');
                if (hasMachinesChart) {
                    console.log('Loading machines chart...');
                    this.machines.load();
                }
                if (hasOperationsChart) {
                    console.log('Loading operations chart...');
                    this.operations.load();
                }
            }
        }
    },

    // Show modal function
    showModal: function() {
        const modal = document.getElementById('modal');
        if (modal) {
            modal.style.display = 'flex';
            modal.classList.add('show');
            document.body.classList.add('modal-open');
            
            // Focus management for accessibility
            setTimeout(() => {
                const firstInput = modal.querySelector('input, select, textarea');
                if (firstInput) {
                    firstInput.focus();
                }
            }, 100);
        }
    },

    // Hide modal function
    hideModal: function() {
        const modal = document.getElementById('modal');
        if (modal) {
            modal.style.display = 'none';
            modal.classList.remove('show');
            document.body.classList.remove('modal-open');
        }
        // Clear modal content
        const modalBody = document.getElementById('modal-body');
        if (modalBody) {
            modalBody.innerHTML = '';
        }
    },

    // Setup event listeners - ДОБАВЛЕННЫЙ МЕТОД
    setupEventListeners: function() {
        console.log('Setting up event listeners...');

        // Close modal with Escape key
        document.addEventListener('keydown', function(e) {
            if (e.key === 'Escape') {
                VendERP.hideModal();
            }
        });

        // Show modal when form is loaded via HTMX
        document.addEventListener('htmx:afterSwap', function(evt) {
            if (evt.detail.target.id === 'modal-body' && evt.detail.xhr.response) {
                VendERP.showModal();
                
                // Remove any conflicting hideModal calls from loaded content
                const modalBody = document.getElementById('modal-body');
                if (modalBody) {
                    const scripts = modalBody.getElementsByTagName('script');
                    for (let script of scripts) {
                        if (script.textContent.includes('hideModal()')) {
                            script.textContent = script.textContent.replace(
                                /hideModal\(\)/g, 
                                'VendERP.hideModal()'
                            );
                        }
                    }
                    
                    // Also replace onclick attributes
                    const buttons = modalBody.querySelectorAll('[onclick*="hideModal()"]');
                    buttons.forEach(button => {
                        const onclick = button.getAttribute('onclick');
                        if (onclick) {
                            button.setAttribute('onclick', onclick.replace('hideModal()', 'VendERP.hideModal()'));
                        }
                    });

                    // Add large class for wider forms (like operations)
                    const form = modalBody.querySelector('form');
                    if (form && form.querySelectorAll('.form-group').length > 8) {
                        const modalContent = document.querySelector('.modal-content');
                        if (modalContent) {
                            modalContent.classList.add('large');
                        }
                    }
                }
            }
        });

        // Close modal after successful save for various tables
        document.addEventListener('htmx:beforeSwap', function(evt) {
            const targets = ['accounts-table', 'machines-table', 'locations-table', 'operations-table'];
            if (targets.includes(evt.detail.target.id) && evt.detail.shouldSwap) {
                VendERP.hideModal();
            }
        });

        // Handle HTMX errors
        document.addEventListener('htmx:responseError', function(evt) {
            console.error('HTMX Error:', evt.detail);
        });
    },

    // Utility function to format dates
    formatDate: function(dateString) {
        const date = new Date(dateString);
        return date.toLocaleDateString('ru-RU');
    },

    // Utility function to format currency
    formatCurrency: function(amount) {
        return new Intl.NumberFormat('ru-RU', {
            style: 'currency',
            currency: 'RUB'
        }).format(amount);
    },

    // Enhanced function to handle form loading with cleanup
    loadForm: function(url, params = '') {
        const fullUrl = params ? `${url}?${params}` : url;
        htmx.ajax('GET', fullUrl, {
            target: '#modal-body',
            swap: 'innerHTML'
        });
    },

    // Initialize application
    init: function() {
        this.setupEventListeners();
        this.charts.init();
        console.log('VendERP initialized with charts support');
    },
};

// Initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    VendERP.init();
});

// Global functions for backward compatibility
function showModal() {
    VendERP.showModal();
}

function hideModal() {
    VendERP.hideModal();
}

// Override any existing hideModal functions that might be loaded later
window.hideModal = VendERP.hideModal;
window.showModal = VendERP.showModal;
// Add chart functions to global scope if needed
window.refreshMachinesChart = () => VendERP.charts.machines.refresh();
window.refreshOperationsChart = () => VendERP.charts.operations.refresh();