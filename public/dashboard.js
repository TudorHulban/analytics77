function refreshAll() {
  var metric = getCurrentMetric();
  document.getElementById('recordsValue').textContent = metric ? metric.RecordsPerPeriod.value.toLocaleString() : '0';

  console.log("refresh was invoked");
}

// ============ THEME MANAGEMENT ============
function getTheme() {
  return localStorage.getItem('metricactive-theme') || 'light';
}

function setTheme(theme) {
  document.documentElement.setAttribute('data-theme', theme);
  localStorage.setItem('metricactive-theme', theme);

  var isDark = theme === 'dark';
  var icons = document.querySelectorAll('#desktopThemeToggle i, #mobileThemeToggle i');
  icons.forEach(function (icon) {
    icon.className = isDark ? 'fas fa-sun' : 'fas fa-moon';
  });

  var themeLabel = document.querySelector('#desktopThemeToggle span');
  if (themeLabel) {
    themeLabel.textContent = isDark ? 'Light Mode' : 'Dark Mode';
  }

  if (window.chartInstance) {
    updateChartTheme();
  }
}

function toggleTheme() {
  var current = getTheme();
  setTheme(current === 'dark' ? 'light' : 'dark');
}

function updateChartTheme() {
  if (!window.chartInstance) return;
  var isDark = getTheme() === 'dark';
  var chart = window.chartInstance;

  chart.options.scales.x.grid.color = isDark ? '#334155' : '#e2e8f0';
  chart.options.scales.y.grid.color = isDark ? '#334155' : '#e2e8f0';
  chart.options.scales.x.ticks.color = isDark ? '#94a3b8' : '#64748b';
  chart.options.scales.y.ticks.color = isDark ? '#94a3b8' : '#64748b';
  chart.data.datasets[0].borderColor = isDark ? '#60a5fa' : '#2563eb';
  chart.data.datasets[0].backgroundColor = isDark ? 'rgba(96,165,250,0.1)' : 'rgba(37,99,235,0.08)';
  chart.data.datasets[0].pointBackgroundColor = isDark ? '#93c5fd' : '#1e40af';
  chart.update();
}

// Initialize theme
setTheme(getTheme());

// Theme toggle event listeners
document.getElementById('desktopThemeToggle').addEventListener('click', toggleTheme);
document.getElementById('mobileThemeToggle').addEventListener('click', toggleTheme);


// ============ METRICS ============
function getCurrentMetric() {
  var siteData = dataCenter.data[currentSite];
  return siteData ? siteData.MetricActive : null;
}

// ============ PERIOD TOGGLE ============
document.querySelectorAll('.period-btn').forEach(function (btn) {
  btn.addEventListener('click', function () {
    document.querySelectorAll('.period-btn').forEach(function (b) {
      b.classList.remove('active');
    });
    this.classList.add('active');
    currentPeriod = this.dataset.period;

    document.getElementById('hourGroup').classList.toggle('disabled', currentPeriod !== 'hour');
    document.getElementById('dayGroup').classList.toggle('disabled', currentPeriod === 'month');

    var labels = { hour: 'Hourly', day: 'Daily', month: 'Monthly' };
    document.getElementById('chartContextLabel').textContent = labels[currentPeriod] + ' data view';

    refreshAll();
  });
});

// ============ SITE SELECTOR ============
document.getElementById('siteSelector').addEventListener('change', function () {
  currentSite = this.value;
  refreshAll();
});
