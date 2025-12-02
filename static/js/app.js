// VendERP Application JavaScript
const VendERP = {
    // Charts module
    charts: {
        cash: {
            data: null,
            miniChart: null,
            fullChart: null,

            // Загрузка данных
            load: async function () {
                try {
                    console.log('DEBUG: Loading cash chart data...');
                    const res = await fetch('/api/charts/cash');
                    if (!res.ok) {
                        throw new Error(`HTTP error! status: ${res.status}`);
                    }
                    this.data = await res.json();
                    console.log('DEBUG: Cash chart data loaded:', this.data);
                    this.updateUI();
                } catch (err) {
                    console.error('Ошибка загрузки графика денег:', err);
                    this.showError();
                }
            },

            // Обновление интерфейса
            updateUI: function () {
                if (!this.data || !this.data.series || this.data.series.length === 0) {
                    console.warn('No data for cash chart');
                    this.showError();
                    return;
                }

                const series = this.data.series[0];
                const labels = this.data.labels || [];
                const amounts = series.data ? series.data.map(d => d.value || 0) : [];

                console.log('DEBUG: Updating cash UI with:', {
                    total: this.data.total,
                    labels: labels.length,
                    data: amounts.length
                });

                // Мини-график
                const totalElement = document.getElementById('cash-mini-total');
                if (totalElement) {
                    totalElement.textContent = this.formatCurrency(this.data.total || 0);
                }

                this.updateMiniChart(labels, amounts, series.color);

                // Полноэкранный график
                const fullTotal = document.getElementById('cash-full-total');
                if (fullTotal) {
                    fullTotal.textContent = this.formatCurrency(this.data.total || 0);
                }

                const periodElement = document.getElementById('cash-full-period');
                if (periodElement) {
                    periodElement.textContent = this.data.period || '30 дней';
                }

                const changeElement = document.getElementById('cash-full-change');
                if (changeElement) {
                    const change = this.data.change || 0;
                    const changeText = change >= 0 ? `+${this.formatCurrency(change)}` : this.formatCurrency(change);
                    changeElement.textContent = changeText;
                    changeElement.className = `stat-value ${change >= 0 ? 'positive' : 'negative'}`;
                }

                const trendElement = document.getElementById('cash-full-trend');
                if (trendElement) {
                    const trend = this.data.trend || 0;
                    const trendInfo = this.getTrendInfo(trend);
                    trendElement.innerHTML = `${trendInfo.icon} ${trendInfo.text}`;
                    trendElement.className = `stat-trend ${trendInfo.class}`;
                }

                const infoElement = document.getElementById('cash-data-info');
                if (infoElement) {
                    infoElement.textContent =
                        `Данные за ${this.data.period || '30 дней'} • Обновлено: ${new Date().toLocaleTimeString('ru-RU')}`;
                }
            },

            // Форматирование валюты
            formatCurrency: function (amount) {
                return new Intl.NumberFormat('ru-RU', {
                    style: 'currency',
                    currency: 'RUB',
                    minimumFractionDigits: 0,
                    maximumFractionDigits: 0
                }).format(amount);
            },

            // Мини-график
            updateMiniChart: function (labels, data, color) {
                const canvas = document.getElementById('cash-mini-chart');
                if (!canvas) {
                    console.warn('Cash mini chart canvas not found');
                    return;
                }

                const ctx = canvas.getContext('2d');

                if (this.miniChart) this.miniChart.destroy();

                // Если данных нет, показываем placeholder
                if (data.length === 0) {
                    console.log('No data for cash mini chart, showing placeholder');
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
                                    label: function (context) {
                                        return new Intl.NumberFormat('ru-RU', {
                                            style: 'currency',
                                            currency: 'RUB'
                                        }).format(context.raw);
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
            updateFullChart: function () {
                const canvas = document.getElementById('cash-full-chart');
                if (!canvas || !this.data || !this.data.series || this.data.series.length === 0) {
                    console.warn('Cannot update cash full chart - missing data or canvas');
                    return;
                }

                const ctx = canvas.getContext('2d');
                const series = this.data.series[0];
                const labels = this.data.labels || [];
                const amounts = series.data ? series.data.map(d => d.value || 0) : [];
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
                            label: 'Деньги',
                            data: amounts,
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
                                    label: function (context) {
                                        const date = dates[context.dataIndex] || '';
                                        const label = date ? `${date}: ` : '';
                                        return `${label}${new Intl.NumberFormat('ru-RU', {
                                            style: 'currency',
                                            currency: 'RUB'
                                        }).format(context.raw)}`;
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
                                    callback: function (value) {
                                        return new Intl.NumberFormat('ru-RU', {
                                            style: 'currency',
                                            currency: 'RUB',
                                            minimumFractionDigits: 0,
                                            maximumFractionDigits: 0
                                        }).format(value);
                                    }
                                }
                            }
                        }
                    }
                });
            },

            // Разворачивание/сворачивание
            expand: function () {
                const fullscreen = document.getElementById('cash-chart-fullscreen');
                if (fullscreen) {
                    fullscreen.style.display = 'block';
                    document.body.style.overflow = 'hidden';
                    setTimeout(() => this.updateFullChart(), 100);
                }
            },

            collapse: function () {
                const fullscreen = document.getElementById('cash-chart-fullscreen');
                if (fullscreen) {
                    fullscreen.style.display = 'none';
                    document.body.style.overflow = 'auto';
                }
            },

            // Вспомогательные функции
            getTrendInfo: function (trend) {
                if (trend === 1) return { class: 'up', icon: '📈', text: 'Рост' };
                if (trend === -1) return { class: 'down', icon: '📉', text: 'Спад' };
                return { class: 'stable', icon: '➡️', text: 'Стабильно' };
            },

            refresh: function () {
                console.log('Refreshing cash chart...');
                if (this.miniChart) this.miniChart.destroy();
                if (this.fullChart) this.fullChart.destroy();
                this.load();
            },

            showError: function () {
                const elem = document.getElementById('cash-mini-total');
                if (elem) {
                    elem.textContent = 'Ошибка';
                    elem.style.color = 'var(--danger)';
                }
            }
        },
        toys: {
            data: null,
            miniChart: null,
            fullChart: null,

            // Загрузка данных
            load: async function () {
                try {
                    console.log('DEBUG: Loading toys chart data...');
                    const res = await fetch('/api/charts/toys');
                    if (!res.ok) {
                        throw new Error(`HTTP error! status: ${res.status}`);
                    }
                    this.data = await res.json();
                    console.log('DEBUG: Toys chart data loaded:', this.data);
                    this.updateUI();
                } catch (err) {
                    console.error('Ошибка загрузки графика игрушек:', err);
                    this.showError();
                }
            },

            // Обновление интерфейса
            updateUI: function () {
                if (!this.data || !this.data.series || this.data.series.length === 0) {
                    console.warn('No data for toys chart');
                    this.showError();
                    return;
                }

                const series = this.data.series[0];
                const labels = this.data.labels || [];
                const counts = series.data ? series.data.map(d => d.count || d.value || 0) : [];

                console.log('DEBUG: Updating toys UI with:', {
                    total: this.data.total,
                    labels: labels.length,
                    data: counts.length
                });

                // Мини-график
                const totalElement = document.getElementById('toys-mini-total');
                if (totalElement) {
                    totalElement.textContent = this.data.total || '-';
                }

                this.updateMiniChart(labels, counts, series.color);

                // Полноэкранный график
                const fullTotal = document.getElementById('toys-full-total');
                if (fullTotal) {
                    fullTotal.textContent = this.data.total || '-';
                }

                const periodElement = document.getElementById('toys-full-period');
                if (periodElement) {
                    periodElement.textContent = this.data.period || '30 дней';
                }

                const changeElement = document.getElementById('toys-full-change');
                if (changeElement) {
                    const change = this.data.change || 0;
                    const changeText = change >= 0 ? `+${change}` : change;
                    changeElement.textContent = changeText;
                    changeElement.className = `stat-value ${change >= 0 ? 'positive' : 'negative'}`;
                }

                const trendElement = document.getElementById('toys-full-trend');
                if (trendElement) {
                    const trend = this.data.trend || 0;
                    const trendInfo = this.getTrendInfo(trend);
                    trendElement.innerHTML = `${trendInfo.icon} ${trendInfo.text}`;
                    trendElement.className = `stat-trend ${trendInfo.class}`;
                }

                const infoElement = document.getElementById('toys-data-info');
                if (infoElement) {
                    infoElement.textContent =
                        `Данные за ${this.data.period || '30 дней'} • Обновлено: ${new Date().toLocaleTimeString('ru-RU')}`;
                }
            },

            // Мини-график
            updateMiniChart: function (labels, data, color) {
                const canvas = document.getElementById('toys-mini-chart');
                if (!canvas) {
                    console.warn('Toys mini chart canvas not found');
                    return;
                }

                const ctx = canvas.getContext('2d');

                if (this.miniChart) this.miniChart.destroy();

                // Если данных нет, показываем placeholder
                if (data.length === 0) {
                    console.log('No data for toys mini chart, showing placeholder');
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
                            borderColor: color || '#8B5CF6',
                            backgroundColor: (color || '#8B5CF6') + '20',
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
                                    label: function (context) {
                                        return `${context.raw} игрушек`;
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
            updateFullChart: function () {
                const canvas = document.getElementById('toys-full-chart');
                if (!canvas || !this.data || !this.data.series || this.data.series.length === 0) {
                    console.warn('Cannot update toys full chart - missing data or canvas');
                    return;
                }

                const ctx = canvas.getContext('2d');
                const series = this.data.series[0];
                const labels = this.data.labels || [];
                const counts = series.data ? series.data.map(d => d.count || d.value || 0) : [];
                const dates = series.data ? series.data.map(d => d.date || '') : [];

                if (this.fullChart) this.fullChart.destroy();

                const gradient = ctx.createLinearGradient(0, 0, 0, canvas.height);
                gradient.addColorStop(0, (series.color || '#8B5CF6') + 'CC');
                gradient.addColorStop(1, (series.color || '#8B5CF6') + '22');

                this.fullChart = new Chart(ctx, {
                    type: 'line',
                    data: {
                        labels: labels,
                        datasets: [{
                            label: 'Игрушки',
                            data: counts,
                            borderColor: series.color || '#8B5CF6',
                            backgroundColor: gradient,
                            borderWidth: 3,
                            fill: true,
                            tension: 0.3,
                            pointBackgroundColor: series.color || '#8B5CF6',
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
                                    label: function (context) {
                                        const date = dates[context.dataIndex] || '';
                                        const label = date ? `${date}: ` : '';
                                        return `${label}${context.raw} игрушек`;
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
                                    callback: function (value) {
                                        return value;
                                    }
                                }
                            }
                        }
                    }
                });
            },

            // Разворачивание/сворачивание
            expand: function () {
                const fullscreen = document.getElementById('toys-chart-fullscreen');
                if (fullscreen) {
                    fullscreen.style.display = 'block';
                    document.body.style.overflow = 'hidden';
                    setTimeout(() => this.updateFullChart(), 100);
                }
            },

            collapse: function () {
                const fullscreen = document.getElementById('toys-chart-fullscreen');
                if (fullscreen) {
                    fullscreen.style.display = 'none';
                    document.body.style.overflow = 'auto';
                }
            },

            // Вспомогательные функции
            getTrendInfo: function (trend) {
                if (trend === 1) return { class: 'up', icon: '📈', text: 'Рост' };
                if (trend === -1) return { class: 'down', icon: '📉', text: 'Спад' };
                return { class: 'stable', icon: '➡️', text: 'Стабильно' };
            },

            refresh: function () {
                console.log('Refreshing toys chart...');
                if (this.miniChart) this.miniChart.destroy();
                if (this.fullChart) this.fullChart.destroy();
                this.load();
            },

            showError: function () {
                const elem = document.getElementById('toys-mini-total');
                if (elem) {
                    elem.textContent = 'Ошибка';
                    elem.style.color = 'var(--danger)';
                }
            }
        },
        machines: {
            data: null,
            activeData: null,
            dualMiniChart: null,
            dualFullChart: null,

            // Загрузка данных
            load: async function () {
                try {
                    console.log('DEBUG: Loading machines chart data...');

                    // Загружаем данные для всех автоматов
                    const totalRes = await fetch('/api/charts/machines');
                    if (!totalRes.ok) {
                        throw new Error(`HTTP error! status: ${totalRes.status}`);
                    }
                    this.data = await totalRes.json();
                    console.log('DEBUG: Total machines chart data loaded:', this.data);

                    // Загружаем данные для активных автоматов
                    const activeRes = await fetch('/api/charts/machines/active');
                    if (!activeRes.ok) {
                        throw new Error(`HTTP error! status: ${activeRes.status}`);
                    }
                    this.activeData = await activeRes.json();
                    console.log('DEBUG: Active machines chart data loaded:', this.activeData);

                    this.updateUI();
                } catch (err) {
                    console.error('Ошибка загрузки графика автоматов:', err);
                    this.showError();
                }
            },

            // Обновление интерфейса
            updateUI: function () {
                // Обновляем данные из data-атрибутов
                const chartContainer = document.getElementById('machines-chart-mini');
                if (chartContainer) {
                    //const totalCount = document.getElementById('machines-total-count');
                    //const activeCount = document.getElementById('machines-active-count');
                    //const activePercent = document.getElementById('machines-active-percent');

                    //if (totalCount) totalCount.textContent = chartContainer.dataset.totalMachines || '0';
                    //if (activeCount) activeCount.textContent = chartContainer.dataset.activeMachines || '0';
                    //if (activePercent) activePercent.textContent = `${chartContainer.dataset.activePercent || '0'}%`;
                }

                // Обновляем один график с двумя линиями
                if (this.data && this.data.series && this.data.series.length > 0 &&
                    this.activeData && this.activeData.series && this.activeData.series.length > 0) {

                    const totalSeries = this.data.series[0];
                    const activeSeries = this.activeData.series[0];
                    const labels = this.data.labels || [];

                    const totalCounts = totalSeries.data ? totalSeries.data.map(d => d.count || d.value || 0) : [];
                    const activeCounts = activeSeries.data ? activeSeries.data.map(d => d.count || d.value || 0) : [];

                    this.updateDualMiniChart(labels, totalCounts, activeCounts);
                }

                // Обновляем полноэкранные данные
                // this.updateFullscreenUI();
            },

            // Мини-график с двумя линиями
            updateDualMiniChart: function (labels, totalData, activeData) {
                const canvas = document.getElementById('machines-dual-chart');
                if (!canvas) {
                    console.warn('Dual machines mini chart canvas not found');
                    return;
                }

                const ctx = canvas.getContext('2d');

                if (this.dualMiniChart) this.dualMiniChart.destroy();

                if (totalData.length === 0 || activeData.length === 0) {
                    this.dualMiniChart = new Chart(ctx, {
                        type: 'line',
                        data: {
                            labels: ['', '', '', '', '', ''],
                            datasets: [
                                {
                                    label: 'Все автоматы',
                                    data: [1, 2, 1, 3, 2, 1],
                                    borderColor: '#4F46E5',
                                    backgroundColor: 'rgba(79, 70, 229, 0.2)',
                                    borderWidth: 1,
                                    fill: false,
                                    tension: 0.4,
                                    pointRadius: 0
                                },
                                {
                                    label: 'Активные',
                                    data: [1, 1, 2, 2, 1, 1],
                                    borderColor: '#10B981',
                                    backgroundColor: 'rgba(16, 185, 129, 0.2)',
                                    borderWidth: 1,
                                    fill: false,
                                    tension: 0.4,
                                    pointRadius: 0
                                }
                            ]
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

                // Фильтруем данные для мини-графика
                const filteredLabels = [];
                const filteredTotalData = [];
                const filteredActiveData = [];

                for (let i = 0; i < labels.length; i++) {
                    if (i % 3 === 0 || i === labels.length - 1) {
                        filteredLabels.push(labels[i]);
                        filteredTotalData.push(totalData[i] || 0);
                        filteredActiveData.push(activeData[i] || 0);
                    }
                }

                this.dualMiniChart = new Chart(ctx, {
                    type: 'line',
                    data: {
                        labels: filteredLabels,
                        datasets: [
                            {
                                label: 'Все автоматы',
                                data: filteredTotalData,
                                borderColor: '#4F46E5',
                                backgroundColor: 'rgba(79, 70, 229, 0.1)',
                                borderWidth: 1.5,
                                fill: false,
                                tension: 0.3,
                                pointRadius: 0,
                                borderDash: [0, 0]
                            },
                            {
                                label: 'Активные',
                                data: filteredActiveData,
                                borderColor: '#10B981',
                                backgroundColor: 'rgba(16, 185, 129, 0.1)',
                                borderWidth: 1.5,
                                fill: false,
                                tension: 0.3,
                                pointRadius: 0,
                                borderDash: [5, 5]
                            }
                        ]
                    },
                    options: {
                        responsive: true,
                        maintainAspectRatio: false,
                        plugins: {
                            legend: { display: false },
                            tooltip: {
                                enabled: true,
                                mode: 'index',
                                intersect: false,
                                callbacks: {
                                    label: function (context) {
                                        const label = context.dataset.label || '';
                                        return `${label}: ${context.parsed.y} автоматов`;
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
                        },
                        interaction: {
                            mode: 'index',
                            intersect: false
                        },
                        elements: {
                            line: {
                                tension: 0.3
                            }
                        }
                    }
                });
            },

            // Обновление полноэкранного интерфейса
            updateFullscreenUI: function () {
                // Общая статистика
                const fullTotal = document.getElementById('machines-full-total');
                if (fullTotal && this.data) {
                    fullTotal.textContent = this.data.total || '-';
                }

                const fullActive = document.getElementById('machines-full-active');
                if (fullActive && this.activeData) {
                    fullActive.textContent = this.activeData.total || '-';
                }

                const fullPercent = document.getElementById('machines-full-percent');
                if (fullPercent && this.data && this.activeData && this.data.total > 0) {
                    const percent = Math.round((this.activeData.total / this.data.total) * 100);
                    fullPercent.textContent = `${percent}%`;
                }

                const periodElement = document.getElementById('machines-full-period');
                if (periodElement && this.data) {
                    periodElement.textContent = this.data.period || '30 дней';
                }

                const infoElement = document.getElementById('machines-data-info');
                if (infoElement) {
                    infoElement.textContent =
                        `Данные за 30 дней • Обновлено: ${new Date().toLocaleTimeString('ru-RU')}`;
                }
            },

            // Полноэкранный график с двумя линиями
            updateFullChart: function () {
                const canvas = document.getElementById('machines-full-dual-chart');
                if (!canvas || !this.data || !this.data.series || this.data.series.length === 0 ||
                    !this.activeData || !this.activeData.series || this.activeData.series.length === 0) {
                    console.warn('Cannot update full dual chart - missing data');
                    return;
                }

                const ctx = canvas.getContext('2d');
                const totalSeries = this.data.series[0];
                const activeSeries = this.activeData.series[0];
                const labels = this.data.labels || [];

                const totalCounts = totalSeries.data ? totalSeries.data.map(d => d.count || d.value || 0) : [];
                const activeCounts = activeSeries.data ? activeSeries.data.map(d => d.count || d.value || 0) : [];

                const totalDates = totalSeries.data ? totalSeries.data.map(d => d.date || '') : [];
                const activeDates = activeSeries.data ? activeSeries.data.map(d => d.date || '') : [];

                if (this.dualFullChart) this.dualFullChart.destroy();

                // Создаем градиенты для заливки
                const totalGradient = ctx.createLinearGradient(0, 0, 0, canvas.height);
                totalGradient.addColorStop(0, '#4F46E5' + '99');
                totalGradient.addColorStop(1, '#4F46E5' + '22');

                const activeGradient = ctx.createLinearGradient(0, 0, 0, canvas.height);
                activeGradient.addColorStop(0, '#10B981' + '99');
                activeGradient.addColorStop(1, '#10B981' + '22');

                this.dualFullChart = new Chart(ctx, {
                    type: 'line',
                    data: {
                        labels: labels,
                        datasets: [
                            {
                                label: 'Все автоматы',
                                data: totalCounts,
                                borderColor: '#4F46E5',
                                backgroundColor: totalGradient,
                                borderWidth: 3,
                                fill: true,
                                tension: 0.3,
                                pointBackgroundColor: '#4F46E5',
                                pointBorderColor: '#fff',
                                pointBorderWidth: 2,
                                pointRadius: 4,
                                pointHoverRadius: 6
                            },
                            {
                                label: 'Активные автоматы',
                                data: activeCounts,
                                borderColor: '#10B981',
                                backgroundColor: activeGradient,
                                borderWidth: 3,
                                fill: true,
                                tension: 0.3,
                                pointBackgroundColor: '#10B981',
                                pointBorderColor: '#fff',
                                pointBorderWidth: 2,
                                pointRadius: 4,
                                pointHoverRadius: 6
                            }
                        ]
                    },
                    options: {
                        responsive: true,
                        maintainAspectRatio: false,
                        plugins: {
                            legend: {
                                display: true,
                                position: 'top',
                                labels: {
                                    padding: 20,
                                    usePointStyle: true,
                                    pointStyle: 'circle',
                                    font: {
                                        size: 14
                                    }
                                }
                            },
                            tooltip: {
                                mode: 'index',
                                intersect: false,
                                callbacks: {
                                    label: function (context) {
                                        const date = totalDates[context.dataIndex] || activeDates[context.dataIndex] || '';
                                        const label = date ? `${date}: ` : '';
                                        return `${label}${context.dataset.label}: ${context.parsed.y} автоматов`;
                                    }
                                }
                            }
                        },
                        scales: {
                            x: {
                                grid: {
                                    color: 'rgba(0,0,0,0.05)'
                                },
                                ticks: {
                                    maxTicksLimit: 10,
                                    font: {
                                        size: 12
                                    }
                                }
                            },
                            y: {
                                beginAtZero: true,
                                grid: {
                                    color: 'rgba(0,0,0,0.05)'
                                },
                                ticks: {
                                    font: {
                                        size: 12
                                    },
                                    callback: function (value) {
                                        return value;
                                    }
                                }
                            }
                        },
                        interaction: {
                            mode: 'index',
                            intersect: false
                        },
                        elements: {
                            line: {
                                tension: 0.3
                            }
                        }
                    }
                });
            },

            // Разворачивание/сворачивание
            expand: function () {
                const fullscreen = document.getElementById('machines-chart-fullscreen');
                if (fullscreen) {
                    fullscreen.style.display = 'block';
                    document.body.style.overflow = 'hidden';
                    setTimeout(() => this.updateFullChart(), 100);
                }
            },

            collapse: function () {
                const fullscreen = document.getElementById('machines-chart-fullscreen');
                if (fullscreen) {
                    fullscreen.style.display = 'none';
                    document.body.style.overflow = 'auto';
                }
            },

            // Вспомогательные функции
            getTrendInfo: function (trend) {
                if (trend === 1) return { class: 'up', icon: '📈', text: 'Рост' };
                if (trend === -1) return { class: 'down', icon: '📉', text: 'Спад' };
                return { class: 'stable', icon: '➡️', text: 'Стабильно' };
            },

            refresh: function () {
                console.log('Refreshing machines charts...');
                if (this.dualMiniChart) this.dualMiniChart.destroy();
                if (this.dualFullChart) this.dualFullChart.destroy();
                this.load();
            },

            showError: function () {
                const elem = document.getElementById('machines-total-count');
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
            load: async function () {
                try {
                    console.log('DEBUG: Loading operations chart data...');
                    const res = await fetch('/api/charts/operations');
                    if (!res.ok) {
                        throw new Error(`HTTP error! status: ${res.status}`);
                    }
                    this.data = await res.json();
                    console.log('DEBUG: Operations chart data loaded:', this.data);
                    this.updateUI();
                } catch (err) {
                    console.error('Ошибка загрузки графика операций:', err);
                    this.showError();
                }
            },

            // Обновление интерфейса
            updateUI: function () {
                if (!this.data || !this.data.series || this.data.series.length === 0) {
                    console.warn('No data for operations chart');
                    this.showError();
                    return;
                }

                console.log('DEBUG: Updating operations UI with:', {
                    seriesCount: this.data.series.length,
                    labels: this.data.labels?.length || 0,
                    total: this.data.total
                });

                // Считаем общие суммы для каждой категории
                const totals = {};
                this.data.series.forEach(series => {
                    totals[series.name] = series.data.reduce((sum, d) => sum + (d.count || 0), 0);
                });

                // Обновляем мини-карточки
                this.updateMiniCards(totals);

                // Обновляем мини-график (группированные столбцы)
                this.updateMiniChart();

                // Обновляем полноэкранные данные
                this.updateFullscreenUI(totals);
            },

            // Обновление мини-карточек с иконками
            updateMiniCards: function (totals) {
                const cards = [
                    { id: 'operations-replenish-card', type: 'Пополнение', icon: '🔄', color: '#10B981' },
                    { id: 'operations-collection-card', type: 'Инкассация', icon: '💵', color: '#F59E0B' },
                    { id: 'operations-service-card', type: 'Обслуживание', icon: '🔧', color: '#EF4444' }
                ];

                cards.forEach(card => {
                    const cardElement = document.getElementById(card.id);
                    if (cardElement) {
                        const countElement = cardElement.querySelector('.stat-number');
                        const iconElement = cardElement.querySelector('.stat-icon');
                        const textElement = cardElement.querySelector('.stat-text');

                        if (countElement) {
                            countElement.textContent = totals[card.type] || '0';
                        }
                        if (iconElement) {
                            iconElement.textContent = card.icon;
                            iconElement.style.color = card.color;
                        }
                        if (textElement) {
                            textElement.textContent = card.type;
                        }
                    }
                });

                // Обновляем общий тотал
                const totalElement = document.getElementById('operations-mini-total');
                if (totalElement) {
                    totalElement.textContent = this.data.total || '-';
                }
            },

            // Мини-график с группированными столбцами
            updateMiniChart: function () {
                const canvas = document.getElementById('operations-mini-chart');
                if (!canvas) {
                    console.warn('Operations mini chart canvas not found');
                    return;
                }

                const ctx = canvas.getContext('2d');

                if (this.miniChart) this.miniChart.destroy();

                if (!this.data.series || this.data.series.length === 0) {
                    this.miniChart = new Chart(ctx, {
                        type: 'bar',
                        data: {
                            labels: ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб'],
                            datasets: [
                                {
                                    label: 'Пополнения',
                                    data: [3, 2, 4, 5, 3, 2],
                                    backgroundColor: '#10B981',
                                    borderColor: '#10B981',
                                    borderWidth: 1,
                                    barPercentage: 0.6,
                                    categoryPercentage: 0.8
                                },
                                {
                                    label: 'Инкассации',
                                    data: [2, 3, 1, 4, 2, 1],
                                    backgroundColor: '#F59E0B',
                                    borderColor: '#F59E0B',
                                    borderWidth: 1,
                                    barPercentage: 0.6,
                                    categoryPercentage: 0.8
                                },
                                {
                                    label: 'Обслуживания',
                                    data: [1, 2, 1, 3, 1, 2],
                                    backgroundColor: '#EF4444',
                                    borderColor: '#EF4444',
                                    borderWidth: 1,
                                    barPercentage: 0.6,
                                    categoryPercentage: 0.8
                                }
                            ]
                        },
                        options: {
                            responsive: true,
                            maintainAspectRatio: false,
                            plugins: {
                                legend: { display: false },
                                tooltip: {
                                    enabled: false,
                                    mode: 'index',
                                    intersect: false
                                }
                            },
                            scales: {
                                x: {
                                    display: false,
                                    grid: { display: false }
                                },
                                y: {
                                    display: false,
                                    grid: { display: false },
                                    beginAtZero: true
                                }
                            }
                        }
                    });
                    return;
                }

                // Для мини-графика берем только последние 7 дней
                const labels = this.data.labels || [];
                const series = this.data.series;

                // Определяем индексы последних 7 точек
                const startIndex = Math.max(0, labels.length - 7);
                const miniLabels = labels.slice(startIndex);

                const datasets = series.map(s => {
                    const data = s.data || [];
                    const miniData = data.slice(startIndex).map(d => d.count || 0);

                    return {
                        label: s.name,
                        data: miniData,
                        backgroundColor: s.color + 'CC',
                        borderColor: s.color,
                        borderWidth: 1,
                        barPercentage: 0.6,
                        categoryPercentage: 0.8
                    };
                });

                this.miniChart = new Chart(ctx, {
                    type: 'bar',
                    data: {
                        labels: miniLabels,
                        datasets: datasets
                    },
                    options: {
                        responsive: true,
                        maintainAspectRatio: false,
                        plugins: {
                            legend: { display: false },
                            tooltip: {
                                enabled: true,
                                mode: 'index',
                                intersect: false,
                                callbacks: {
                                    label: function (context) {
                                        return `${context.dataset.label}: ${context.parsed.y} операций`;
                                    }
                                }
                            }
                        },
                        scales: {
                            x: {
                                display: false,
                                grid: { display: false },
                                stacked: false
                            },
                            y: {
                                display: false,
                                grid: { display: false },
                                beginAtZero: true,
                                stacked: false
                            }
                        },
                        interaction: {
                            mode: 'index',
                            intersect: false
                        }
                    }
                });
            },

            // Обновление полноэкранного интерфейса
            updateFullscreenUI: function (totals) {
                // Обновляем карточки в полноэкранном режиме
                const categories = [
                    { id: 'operations-full-replenish', type: 'Пополнение', color: '#10B981' },
                    { id: 'operations-full-collection', type: 'Инкассация', color: '#F59E0B' },
                    { id: 'operations-full-service', type: 'Обслуживание', color: '#EF4444' }
                ];

                categories.forEach(cat => {
                    const element = document.getElementById(cat.id);
                    if (element) {
                        element.textContent = totals[cat.type] || '0';
                        element.style.color = cat.color;
                    }
                });

                // Общий тотал
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

            // Полноэкранный график с группированными столбцами
            updateFullChart: function () {
                const canvas = document.getElementById('operations-full-chart');
                if (!canvas || !this.data || !this.data.series || this.data.series.length === 0) {
                    console.warn('Cannot update operations full chart - missing data or canvas');
                    return;
                }

                const ctx = canvas.getContext('2d');
                const labels = this.data.labels || [];

                if (this.fullChart) this.fullChart.destroy();

                // Подготавливаем данные для группированных столбцов
                const datasets = this.data.series.map(series => {
                    const data = series.data || [];
                    const counts = data.map(d => d.count || 0);
                    const dates = data.map(d => d.date || '');

                    return {
                        label: series.name,
                        data: counts,
                        backgroundColor: series.color + 'CC',
                        borderColor: series.color,
                        borderWidth: 1,
                        dates: dates
                    };
                });

                this.fullChart = new Chart(ctx, {
                    type: 'bar',
                    data: {
                        labels: labels,
                        datasets: datasets
                    },
                    options: {
                        responsive: true,
                        maintainAspectRatio: false,
                        plugins: {
                            legend: {
                                display: true,
                                position: 'top',
                                labels: {
                                    padding: 20,
                                    usePointStyle: true,
                                    pointStyle: 'rect',
                                    font: {
                                        size: 14
                                    }
                                }
                            },
                            tooltip: {
                                mode: 'index',
                                intersect: false,
                                callbacks: {
                                    title: function (context) {
                                        const date = context[0].dataset.dates?.[context[0].dataIndex] || '';
                                        return date || `День ${context[0].label}`;
                                    },
                                    label: function (context) {
                                        return `${context.dataset.label}: ${context.parsed.y} операций`;
                                    }
                                }
                            }
                        },
                        scales: {
                            x: {
                                grid: {
                                    display: false
                                },
                                ticks: {
                                    maxTicksLimit: 15,
                                    font: {
                                        size: 11
                                    }
                                },
                                stacked: false
                            },
                            y: {
                                beginAtZero: true,
                                grid: {
                                    color: 'rgba(0,0,0,0.05)'
                                },
                                ticks: {
                                    font: {
                                        size: 12
                                    },
                                    callback: function (value) {
                                        return value;
                                    },
                                    stepSize: 1
                                },
                                stacked: false
                            }
                        },
                        interaction: {
                            mode: 'index',
                            intersect: false
                        },
                        barPercentage: 0.8,
                        categoryPercentage: 0.9
                    }
                });
            },

            // Альтернатива: сложенный (stacked) bar chart
            updateStackedFullChart: function () {
                const canvas = document.getElementById('operations-full-chart');
                if (!canvas || !this.data || !this.data.series || this.data.series.length === 0) {
                    console.warn('Cannot update operations full chart - missing data or canvas');
                    return;
                }

                const ctx = canvas.getContext('2d');
                const labels = this.data.labels || [];

                if (this.fullChart) this.fullChart.destroy();

                // Для stacked chart используем полупрозрачные цвета
                const datasets = this.data.series.map(series => {
                    const data = series.data || [];
                    const counts = data.map(d => d.count || 0);
                    const dates = data.map(d => d.date || '');

                    return {
                        label: series.name,
                        data: counts,
                        backgroundColor: series.color + '99', // Более прозрачный для stacked
                        borderColor: series.color,
                        borderWidth: 0,
                        dates: dates
                    };
                });

                this.fullChart = new Chart(ctx, {
                    type: 'bar',
                    data: {
                        labels: labels,
                        datasets: datasets
                    },
                    options: {
                        responsive: true,
                        maintainAspectRatio: false,
                        plugins: {
                            legend: {
                                display: true,
                                position: 'top',
                                labels: {
                                    padding: 20,
                                    usePointStyle: true,
                                    pointStyle: 'rect',
                                    font: {
                                        size: 14
                                    }
                                }
                            },
                            tooltip: {
                                mode: 'index',
                                intersect: false,
                                callbacks: {
                                    title: function (context) {
                                        const date = context[0].dataset.dates?.[context[0].dataIndex] || '';
                                        return date || `День ${context[0].label}`;
                                    },
                                    label: function (context) {
                                        return `${context.dataset.label}: ${context.parsed.y} операций`;
                                    }
                                }
                            }
                        },
                        scales: {
                            x: {
                                grid: {
                                    display: false
                                },
                                ticks: {
                                    maxTicksLimit: 15,
                                    font: {
                                        size: 11
                                    }
                                },
                                stacked: true
                            },
                            y: {
                                beginAtZero: true,
                                grid: {
                                    color: 'rgba(0,0,0,0.05)'
                                },
                                ticks: {
                                    font: {
                                        size: 12
                                    },
                                    callback: function (value) {
                                        return value;
                                    },
                                    stepSize: 1
                                },
                                stacked: true
                            }
                        },
                        interaction: {
                            mode: 'index',
                            intersect: false
                        }
                    }
                });
            },

            // Разворачивание/сворачивание с возможностью выбора типа графика
            expand: function (chartType = 'grouped') {
                const fullscreen = document.getElementById('operations-chart-fullscreen');
                if (fullscreen) {
                    fullscreen.style.display = 'block';
                    document.body.style.overflow = 'hidden';

                    setTimeout(() => {
                        if (chartType === 'stacked') {
                            this.updateStackedFullChart();
                        } else {
                            this.updateFullChart();
                        }
                    }, 100);
                }
            },

            collapse: function () {
                const fullscreen = document.getElementById('operations-chart-fullscreen');
                if (fullscreen) {
                    fullscreen.style.display = 'none';
                    document.body.style.overflow = 'auto';
                }
            },

            // Переключение между типами графиков в полноэкранном режиме
            toggleChartType: function () {
                const toggleBtn = document.getElementById('operations-chart-toggle');
                if (!toggleBtn) return;

                const currentType = toggleBtn.dataset.type || 'grouped';
                const newType = currentType === 'grouped' ? 'stacked' : 'grouped';

                toggleBtn.dataset.type = newType;
                toggleBtn.innerHTML = newType === 'stacked' ?
                    '<i class="icon-layers"></i> Сгруппировать' :
                    '<i class="icon-stack"></i> Сложить';

                if (this.fullChart) {
                    this.fullChart.destroy();
                    if (newType === 'stacked') {
                        this.updateStackedFullChart();
                    } else {
                        this.updateFullChart();
                    }
                }
            },

            // Вспомогательные функции
            getTrendInfo: function (trend) {
                if (trend === 1) return { class: 'up', icon: '📈', text: 'Рост' };
                if (trend === -1) return { class: 'down', icon: '📉', text: 'Спад' };
                return { class: 'stable', icon: '➡️', text: 'Стабильно' };
            },

            refresh: function () {
                console.log('Refreshing operations chart...');
                if (this.miniChart) this.miniChart.destroy();
                if (this.fullChart) this.fullChart.destroy();
                this.load();
            },

            showError: function () {
                const elem = document.getElementById('operations-mini-total');
                if (elem) {
                    elem.textContent = 'Ошибка';
                    elem.style.color = 'var(--danger)';
                }
            }
        },

        // Initialize all charts
        init: function () {
            console.log('DEBUG: charts.init() called');

            // Проверяем, есть ли контейнеры для графиков на странице
            const hasMachinesChart = document.getElementById('machines-chart-mini');
            const hasOperationsChart = document.getElementById('operations-chart-mini');
            const hasCashChart = document.getElementById('cash-chart-mini');
            const hasToysChart = document.getElementById('toys-chart-mini'); // ДОБАВЛЕНО

            console.log('DEBUG: Chart elements check:', {
                machines: hasMachinesChart ? 'FOUND' : 'NOT FOUND',
                operations: hasOperationsChart ? 'FOUND' : 'NOT FOUND',
                cash: hasCashChart ? 'FOUND' : 'NOT FOUND'
            });

            // Если нет ни одного графика, выходим
            if (!hasMachinesChart && !hasOperationsChart && !hasCashChart) {
                console.log('DEBUG: No chart containers found, skipping initialization');
                return;
            }

            // Инициализируем Chart.js если не загружен
            if ((hasMachinesChart || hasOperationsChart || hasCashChart) && typeof Chart === 'undefined') {
                console.log('DEBUG: Loading Chart.js...');
                const s = document.createElement('script');
                s.src = 'https://cdn.jsdelivr.net/npm/chart.js';
                s.onload = () => {
                    console.log('DEBUG: Chart.js loaded');
                    this.loadAllCharts();
                };
                s.onerror = (err) => {
                    console.error('DEBUG: Failed to load Chart.js:', err);
                };
                document.head.appendChild(s);
            } else {
                console.log('DEBUG: Chart.js already loaded');
                this.loadAllCharts();
            }
        },

        // Новый метод для загрузки всех графиков
        loadAllCharts: function () {
            console.log('DEBUG: loadAllCharts() called');

            const charts = [
                { id: 'machines-chart-mini', module: 'machines', name: 'machines' },
                { id: 'operations-chart-mini', module: 'operations', name: 'operations' },
                { id: 'cash-chart-mini', module: 'cash', name: 'cash' },
                { id: 'toys-chart-mini', module: 'toys', name: 'toys' } // ДОБАВЛЕНО

            ];

            charts.forEach(chart => {
                const element = document.getElementById(chart.id);
                if (element && this[chart.module] && typeof this[chart.module].load === 'function') {
                    console.log(`DEBUG: Loading ${chart.name} chart...`);
                    this[chart.module].load();
                } else {
                    console.log(`DEBUG: Skipping ${chart.name} chart:`, {
                        element: !!element,
                        module: !!this[chart.module],
                        loadMethod: this[chart.module] ? typeof this[chart.module].load : 'no module'
                    });
                }
            });
        }
    },

    // Show modal function
    showModal: function () {
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
    hideModal: function () {
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
    setupEventListeners: function () {
        console.log('DEBUG: Setting up event listeners...');

        // Close modal with Escape key
        document.addEventListener('keydown', function (e) {
            if (e.key === 'Escape') {
                VendERP.hideModal();
            }
        });

        // Show modal when form is loaded via HTMX
        document.addEventListener('htmx:afterSwap', function (evt) {
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
        document.addEventListener('htmx:beforeSwap', function (evt) {
            const targets = ['accounts-table', 'machines-table', 'locations-table', 'operations-table'];
            if (targets.includes(evt.detail.target.id) && evt.detail.shouldSwap) {
                VendERP.hideModal();
            }
        });

        // Handle HTMX errors
        document.addEventListener('htmx:responseError', function (evt) {
            console.error('HTMX Error:', evt.detail);
        });
    },

    // Utility function to format dates
    formatDate: function (dateString) {
        const date = new Date(dateString);
        return date.toLocaleDateString('ru-RU');
    },

    // Utility function to format currency
    formatCurrency: function (amount) {
        return new Intl.NumberFormat('ru-RU', {
            style: 'currency',
            currency: 'RUB'
        }).format(amount);
    },

    // Enhanced function to handle form loading with cleanup
    loadForm: function (url, params = '') {
        const fullUrl = params ? `${url}?${params}` : url;
        htmx.ajax('GET', fullUrl, {
            target: '#modal-body',
            swap: 'innerHTML'
        });
    },

    // Initialize application
    init: function () {
        this.setupEventListeners();
        console.log('DEBUG: VendERP.init() called');
        this.charts.init();
        console.log('DEBUG: VendERP initialized with charts support');
    },
};

// Initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', function () {
    console.log('DEBUG: DOMContentLoaded fired');
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
window.refreshCashChart = () => VendERP.charts.cash.refresh();
window.refreshToysChart = () => VendERP.charts.toys.refresh(); // ДОБАВЛЕНО

// ДОБАВЛЕНО: Автоматическая проверка и загрузка графика денег
setTimeout(function () {
    console.log('DEBUG: Delayed check for cash chart');

    // Проверяем наличие элемента
    const cashElement = document.getElementById('cash-chart-mini');
    console.log('DEBUG: cash-chart-mini element:', cashElement ? 'FOUND' : 'NOT FOUND');

    // Проверяем наличие модуля
    console.log('DEBUG: VendERP.charts.cash:', VendERP?.charts?.cash ? 'EXISTS' : 'MISSING');

    // Если элемент есть, но данные не загружены
    if (cashElement && VendERP?.charts?.cash && !VendERP.charts.cash.data) {
        console.log('DEBUG: Cash chart element found but data not loaded, loading now...');
        VendERP.charts.cash.load();
    }
}, 3000);

// ДОБАВЛЕНО: Упрощенная альтернативная инициализация
function initializeCashChart() {
    console.log('DEBUG: initializeCashChart() called');

    // Проверяем наличие элемента
    const cashElement = document.getElementById('cash-chart-mini');
    if (!cashElement) {
        console.log('DEBUG: cash-chart-mini not found on page');
        return;
    }

    // Проверяем наличие модуля
    if (!window.VendERP || !window.VendERP.charts || !window.VendERP.charts.cash) {
        console.log('DEBUG: VendERP.charts.cash not available');
        return;
    }

    // Проверяем, загружены ли уже данные
    if (window.VendERP.charts.cash.data) {
        console.log('DEBUG: Cash chart data already loaded');
        return;
    }

    // Загружаем данные
    console.log('DEBUG: Loading cash chart data...');
    window.VendERP.charts.cash.load();
}

// Пытаемся инициализировать несколько раз с разными задержками
document.addEventListener('DOMContentLoaded', function () {
    // Первая попытка через 500мс
    setTimeout(initializeCashChart, 500);

    // Вторая попытка через 2 секунды
    setTimeout(initializeCashChart, 2000);

    // Третья попытка через 5 секунд
    setTimeout(initializeCashChart, 5000);
});

// Также при полной загрузке страницы
window.addEventListener('load', function () {
    setTimeout(initializeCashChart, 1000);
});

// ДОБАВЛЕНО: Экстренная загрузка при клике на элемент графика
document.addEventListener('click', function (event) {
    if (event.target.closest('#cash-chart-mini')) {
        console.log('DEBUG: Cash chart clicked, checking data...');
        setTimeout(() => {
            if (VendERP?.charts?.cash && !VendERP.charts.cash.data) {
                console.log('DEBUG: Data missing, loading cash chart...');
                VendERP.charts.cash.load();
            }
        }, 100);
    }
});