// Theme Toggle Functionality
const ThemeManager = {
    // Initialize theme
    init: function() {
        this.loadTheme();
        this.setupEventListeners();
    },

    // Load saved theme preference
    loadTheme: function() {
        const savedTheme = localStorage.getItem('theme') || 'dark';
        this.setTheme(savedTheme);
    },

    // Set theme
    setTheme: function(theme) {
        document.body.className = theme + '-theme';
        localStorage.setItem('theme', theme);
        this.updateThemeIcon(theme);
    },

    // Toggle between light and dark themes
    toggleTheme: function() {
        const currentTheme = document.body.classList.contains('dark-theme') ? 'dark' : 'light';
        const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
        this.setTheme(newTheme);
    },

    // Update theme toggle button icon
    updateThemeIcon: function(theme) {
        const themeToggle = document.getElementById('theme-toggle');
        if (themeToggle) {
            themeToggle.innerHTML = theme === 'dark' ? '☀️' : '🌙';
            themeToggle.title = theme === 'dark' ? 'Switch to Light Mode' : 'Switch to Dark Mode';
        }
    },

    // Setup event listeners
    setupEventListeners: function() {
        // Theme toggle button
        document.addEventListener('click', (e) => {
            if (e.target.id === 'theme-toggle' || e.target.closest('#theme-toggle')) {
                this.toggleTheme();
            }
        });

        // System theme preference
        if (window.matchMedia) {
            window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', e => {
                if (!localStorage.getItem('theme')) {
                    this.setTheme(e.matches ? 'dark' : 'light');
                }
            });
        }
    }
};

// Initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    ThemeManager.init();
});


