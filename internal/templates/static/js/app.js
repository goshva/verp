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
