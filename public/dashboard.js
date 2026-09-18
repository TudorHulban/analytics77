function initApp() {
  var months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
  var now = new Date();
  var currentYear = now.getFullYear();
  var currentMonth = now.getMonth();
  var currentDay = now.getDate();
  
  // Month select
  var monthHtml = '';
  for (var i = 5; i >= 0; i--) {
    var d = new Date(currentYear, currentMonth - i, 1);
    monthHtml += '<option value="' + d.getFullYear() + '-' + d.getMonth() + '">' + months[d.getMonth()] + ' ' + d.getFullYear() + '</option>';
  }
  monthHtml += '<option value="' + currentYear + '-' + currentMonth + '" selected>' + months[currentMonth] + ' ' + currentYear + ' (active)</option>';
  document.getElementById('monthSelect').innerHTML = monthHtml;
  
  // Day select - dynamic based on current month
  populateDays(currentYear, currentMonth, currentDay);
  
  // Hour select
  var hourHtml = '';
  for (var h = 0; h < 24; h++) {
    hourHtml += '<option value="' + h + '"' + (h === now.getHours() ? ' selected' : '') + '>' + String(h).padStart(2, '0') + ':00</option>';
  }
  document.getElementById('hourSelect').innerHTML = hourHtml;
  
  // Add event listener to update days when month changes
  document.getElementById('monthSelect').addEventListener('change', function() {
    var parts = this.value.split('-');
    var year = parseInt(parts[0]);
    var month = parseInt(parts[1]);
    
    // Get currently selected day, or default to 1
    var currentSelectedDay = parseInt(document.getElementById('daySelect').value) || 1;
    
    // Cap the day to the max days in the new month
    var maxDays = getDaysInMonth(year, month);
    var newDay = Math.min(currentSelectedDay, maxDays);
    
    populateDays(year, month, newDay);
  });
  

// add render chart here
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
  icons.forEach(function(icon) {
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