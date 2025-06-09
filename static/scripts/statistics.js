document.addEventListener('DOMContentLoaded', function () {
    initWordStatusChart();
    initActivityChart();

    loadStatisticsData();

    setupEventListeners();
});

function initWordStatusChart() {
    const ctx = document.getElementById('wordStatusChart').getContext('2d');
    const chartData = {
        labels: ['Learned', 'Known', 'To Learn'],
        datasets: [{
            data: [0, 0, 0],
            backgroundColor: [
                '#4caf50',
                '#8bc34a',
                '#ffb74d'
            ],
            borderWidth: 0,
            hoverOffset: 10
        }]
    };

    window.wordStatusChart = new Chart(ctx, {
        type: 'doughnut',
        data: chartData,
        options: {
            responsive: true,
            maintainAspectRatio: false,
            cutout: '70%',
            plugins: {
                legend: {
                    position: 'bottom',
                    labels: {
                        boxWidth: 12,
                        padding: 20,
                        font: {
                            family: "'Segoe UI', system-ui, -apple-system, sans-serif",
                            size: 12
                        }
                    }
                },
                tooltip: {
                    backgroundColor: 'rgba(0, 0, 0, 0.7)',
                    titleFont: {
                        size: 14,
                        weight: 'bold'
                    },
                    bodyFont: {
                        size: 12
                    },
                    padding: 12,
                    cornerRadius: 4
                }
            }
        }
    });
}

function initActivityChart() {
    const ctx = document.getElementById('activityChart').getContext('2d');
    const chartData = {
        labels: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
        datasets: [{
            label: 'Words Practiced',
            data: [0, 0, 0, 0, 0, 0, 0],
            backgroundColor: '#81c784',
            borderRadius: 6,
            borderWidth: 0,
            barPercentage: 0.7
        }]
    };

    window.activityChart = new Chart(ctx, {
        type: 'bar',
        data: chartData,
        options: {
            responsive: true,
            maintainAspectRatio: false,
            scales: {
                y: {
                    beginAtZero: true,
                    grid: {
                        drawBorder: false,
                        color: 'rgba(0, 0, 0, 0.05)'
                    },
                    ticks: {
                        stepSize: 1,
                        font: {
                            family: "'Segoe UI', system-ui, -apple-system, sans-serif"
                        }
                    }
                },
                x: {
                    grid: {
                        display: false,
                        drawBorder: false
                    },
                    ticks: {
                        font: {
                            family: "'Segoe UI', system-ui, -apple-system, sans-serif"
                        }
                    }
                }
            },
            plugins: {
                legend: {
                    display: false
                },
                tooltip: {
                    backgroundColor: 'rgba(0, 0, 0, 0.7)',
                    titleFont: {
                        size: 14,
                        weight: 'bold'
                    },
                    bodyFont: {
                        size: 12
                    },
                    padding: 12,
                    cornerRadius: 4
                }
            }
        }
    });
}

async function loadStatisticsData() {
    console.log("Loading statistics...");
    try {
        showLoadingIndicator();

        const token = localStorage.getItem('token');
        if (!token) {
            throw new Error("User not authenticated");
        }

        const response = await fetch('/api/statistics', {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        console.log("Data received:", data);
        updateStatisticsUI(data);
    } catch (error) {
        console.error("Error loading statistics:", error);
        showNotification(`Error: ${error.message}`, 'error');
    } finally {
        hideLoadingIndicator();
    }
}

function updateStatisticsUI(data) {
    if (!data) {
        console.error("No data provided!");
        return;
    }

    const elements = {
        'total-words': data.total_words ?? 0,
        'learned-words': data.learned_words ?? 0,
        'known-words': data.known_words ?? 0,
        'to-learn-words': data.to_learn_words ?? 0,
        'success-rate': `${data.success_rate ?? 0}%`,
        'sessions-week': data.sessions_week ?? 0,
        'words-week': data.words_week ?? 0,
        'last-practice': formatDate(data.last_practice) || 'No data',
        'registration-date': formatDate(data.registration_date) || 'No data'
    };

    for (const [id, value] of Object.entries(elements)) {
        const element = document.getElementById(id);
        if (element) {
            element.textContent = value;
        } else {
            console.error(`Element with ID "${id}" not found!`);
        }
    }

    updateCharts(data);
    updateRecentWords(data.recent_words ?? []);
}

function updateStatisticsUI(data) {
    document.getElementById('total-words').textContent = data.total_words || 0;
    document.getElementById('learned-words').textContent = data.learned_words || 0;
    document.getElementById('known-words').textContent = data.known_words || 0;
    document.getElementById('to-learn-words').textContent = data.to_learn_words || 0;
    document.getElementById('success-rate').textContent = `${data.success_rate || 0}%`;
    document.getElementById('sessions-week').textContent = data.sessions_week || 0;
    document.getElementById('words-week').textContent = data.words_week || 0;

    document.getElementById('last-practice').textContent = formatDate(data.last_practice) || 'No data';
    document.getElementById('registration-date').textContent = formatDate(data.registration_date) || 'No data';

    updateCharts(data);

    updateRecentWords(data.recent_words || []);
}

function updateCharts(data) {
    if (window.wordStatusChart) {
        window.wordStatusChart.data.datasets[0].data = [
            data.learned_words || 0,
            data.known_words || 0,
            data.to_learn_words || 0
        ];
        window.wordStatusChart.update();
    }

    if (window.activityChart) {
        const activityData = data.daily_activity && data.daily_activity.length === 7 ?
            data.daily_activity :
            [0, 0, 0, 0, 0, 0, 0];

        window.activityChart.data.datasets[0].data = activityData;
        window.activityChart.update();
    }
}

function updateRecentWords(words) {
    const wordsList = document.querySelector('.words-list');

    if (!words || words.length === 0) {
        wordsList.innerHTML = `
            <div class="empty-message">
                <i class="fas fa-inbox"></i>
                <p>No recently practiced words</p>
            </div>
        `;
        return;
    }

    wordsList.innerHTML = words.map(word => `
        <div class="word-item ${word.status.replace(' ', '-')}">
            <div class="word-text">
                <strong>${word.word}</strong>
                <span>${word.translation}</span>
            </div>
            <div class="word-status">
                ${getStatusIcon(word.status)}
            </div>
        </div>
    `).join('');
}

function getStatusIcon(status) {
    switch (status.toLowerCase()) {
        case 'learned': return '<i class="fas fa-check-circle"></i>';
        case 'known': return '<i class="fas fa-leaf"></i>';
        case 'to learn': return '<i class="fas fa-seedling"></i>';
        default: return '<i class="fas fa-question"></i>';
    }
}

function formatDate(dateString) {
    if (!dateString) return null;

    const date = new Date(dateString);
    if (isNaN(date.getTime())) return null;

    const now = new Date();
    const diffDays = Math.floor((now - date) / (1000 * 60 * 60 * 24));

    if (diffDays === 0) return 'Today';
    if (diffDays === 1) return 'Yesterday';
    if (diffDays < 7) return `${diffDays} days ago`;

    return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
    });
}

function showLoadingIndicator() {
    const loadingIndicator = document.createElement('div');
    loadingIndicator.id = 'loading-indicator';
    loadingIndicator.innerHTML = `
        <div class="loading-spinner">
            <i class="fas fa-spinner fa-spin"></i>
        </div>
    `;
    document.body.appendChild(loadingIndicator);
}

function hideLoadingIndicator() {
    const indicator = document.getElementById('loading-indicator');
    if (indicator) {
        indicator.remove();
    }
}

function showNotification(message, type) {
    const notification = document.createElement('div');
    notification.className = `notification ${type}`;
    notification.innerHTML = `
        <i class="fas fa-${type === 'error' ? 'exclamation-circle' : 'check-circle'}"></i>
        <span>${message}</span>
    `;

    document.body.appendChild(notification);

    setTimeout(() => {
        notification.classList.add('fade-out');
        setTimeout(() => notification.remove(), 500);
    }, 3000);
}

function setupEventListeners() {
    document.addEventListener('visibilitychange', () => {
        if (!document.hidden) {
            loadStatisticsData();
        }
    });

    const refreshBtn = document.getElementById('refresh-btn');
    if (refreshBtn) {
        refreshBtn.addEventListener('click', loadStatisticsData);
    }
}

window.refreshStatistics = loadStatisticsData;