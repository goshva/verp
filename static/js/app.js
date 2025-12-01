// VendERP Application JavaScript
const VendERP = {
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

    // Initialize application
    init: function() {
        this.setupEventListeners();
        console.log('VendERP initialized');
    },

    // Setup event listeners
    setupEventListeners: function() {


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
    }
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
    document.getElementById('theme-toggle').addEventListener('click', function() {
  document.body.classList.toggle('dark-theme');})
// Override any existing hideModal functions that might be loaded later
window.hideModal = VendERP.hideModal;
window.showModal = VendERP.showModal;