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

// ============ DATA MODELS ============
function OS(name, version) {
  this.name = name;
  this.version = version || '';
}
OS.prototype.toString = function() {
  return this.version ? this.name + ' ' + this.version : this.name;
};

function Browser(name, version) {
  this.name = name;
  this.version = version || '';
}
Browser.prototype.toString = function() {
  return this.version ? this.name + ' ' + this.version : this.name;
};

function createMetaActive(namesArr, valuesArr) {
  return {
    names: namesArr.slice(0, 7),
    values: valuesArr.slice(0, 7),
    getOccupiedMask: function() {
      var mask = 0;
      for (var i = 0; i < 7; i++) {
        if (this.values[i] > 0) mask |= (1 << i);
      }
      return mask;
    }
  };
}

function createSiteData(siteName, baseRecords, seed) {
  function rand(min, max) {
    return Math.floor(min + (Math.sin(seed * 1000) * 0.5 + 0.5) * (max - min));
  }
  
  return {
    name: siteName,
    MetricActive: {
      TopIPs: createMetaActive(
        ['192.168.1.1', '10.0.0.5', '172.16.0.2', '8.8.8.8', '1.1.1.1', '203.0.113.5', '198.51.100.7'],
        [1450, 1230, 980, 870, 720, 650, 510].map(function(v) { return rand(v * 0.7, v * 1.3); })
      ),
      TopASN: createMetaActive(
        ['AS15169 Google', 'AS13335 Cloudflare', 'AS32934 Facebook', 'AS16509 Amazon', 'AS8075 Microsoft', 'AS24940 Hetzner', 'AS16276 OVH'],
        [2100, 1890, 1560, 1340, 1100, 890, 760].map(function(v) { return rand(v * 0.7, v * 1.3); })
      ),
      TopCountries: createMetaActive(
        ['US', 'DE', 'GB', 'JP', 'FR', 'CA', 'AU'],
        [3200, 2780, 2100, 1650, 1420, 1180, 970].map(function(v) { return rand(v * 0.7, v * 1.3); })
      ),
      TopCities: createMetaActive(
        ['New York', 'London', 'Tokyo', 'Berlin', 'Paris', 'Singapore', 'Sydney'],
        [1800, 1620, 1480, 1250, 1120, 980, 810].map(function(v) { return rand(v * 0.7, v * 1.3); })
      ),
      TopURLs: createMetaActive(
        ['/api/v1/data', '/home', '/login', '/dashboard', '/products', '/contact', '/about'],
        [2250, 1980, 1670, 1430, 1210, 950, 780].map(function(v) { return rand(v * 0.7, v * 1.3); })
      ),
      TopOperatingSystems: createMetaActive(
        [new OS('Windows', '10'), new OS('macOS', '14'), new OS('Linux', '6.5'), new OS('iOS', '17'), new OS('Android', '14'), new OS('Windows', '11'), new OS('Ubuntu', '22.04')],
        [1750, 1520, 1310, 1180, 990, 870, 690].map(function(v) { return rand(v * 0.7, v * 1.3); })
      ),
      TopBrowsers: createMetaActive(
        [new Browser('Chrome', '120'), new Browser('Firefox', '121'), new Browser('Safari', '17'), new Browser('Edge', '120'), new Browser('Opera', '105'), new Browser('Brave', '1.60'), new Browser('Vivaldi', '6.5')],
        [2100, 1680, 1450, 1220, 890, 760, 620].map(function(v) { return rand(v * 0.7, v * 1.3); })
      ),
      RecordsPerPeriod: { value: baseRecords }
    }
  };
}

// DataCenter
var dataCenter = {
  data: {
    site1: createSiteData('example.com', 28450, 1),
    site2: createSiteData('shop.example.com', 15200, 2),
    site3: createSiteData('analytics.example.com', 32100, 3),
    site4: createSiteData('corp.example.com', 8900, 4),
    site5: createSiteData('global.example.com', 45700, 5)
  }
};

// State
var currentSite = 'site1';
var currentPeriod = 'hour';
var currentView = 'grid';
window.chartInstance = null;

// Users data
var users = [
  { id: 1, name: 'John Doe', email: 'john@example.com', role: 'admin', sites: ['site1', 'site2', 'site3', 'site4', 'site5'], lastActive: '2026-08-06' },
  { id: 2, name: 'Jane Smith', email: 'jane@example.com', role: 'editor', sites: ['site1', 'site3'], lastActive: '2026-08-05' },
  { id: 3, name: 'Bob Wilson', email: 'bob@example.com', role: 'viewer', sites: ['site2'], lastActive: '2026-08-04' }
];

// ============ LOGIN ============
document.getElementById('loginButton').addEventListener('click', function() {
  var email = document.getElementById('loginEmail').value.trim();
  var password = document.getElementById('loginPassword').value.trim();
  
  if (email === 'admin@metricactive.com' && password === 'admin123') {
    // Hide login, show app
    document.getElementById('loginScreen').classList.add('hidden');
    document.getElementById('appScreen').classList.add('visible');
    document.getElementById('loginError').classList.remove('show');
    
    // Update user info
    document.getElementById('userAvatar').textContent = 'JD';
    document.getElementById('userNameDisplay').textContent = 'John Doe';
    
    // Initialize dashboard
    initApp();
  } else {
    document.getElementById('loginError').classList.add('show');
  }
});

// Allow Enter key to login
document.getElementById('loginPassword').addEventListener('keypress', function(e) {
  if (e.key === 'Enter') {
    document.getElementById('loginButton').click();
  }
});

document.getElementById('loginEmail').addEventListener('keypress', function(e) {
  if (e.key === 'Enter') {
    document.getElementById('loginButton').click();
  }
});

function logout() {
  document.getElementById('loginScreen').classList.remove('hidden');
  document.getElementById('appScreen').classList.remove('visible');
  document.getElementById('loginPassword').value = '';
  document.getElementById('loginError').classList.remove('show');
  
  // Close mobile sidebar
  document.getElementById('sidebar').classList.remove('mobile-open');
  document.getElementById('sidebarOverlay').classList.remove('active');
  
  // Destroy chart
  if (window.chartInstance) {
    window.chartInstance.destroy();
    window.chartInstance = null;
  }
}

document.getElementById('desktopLogoutBtn').addEventListener('click', logout);
document.getElementById('mobileLogoutBtn').addEventListener('click', logout);

// ============ NAVIGATION ============
function navigateToPage(page) {
  // Update sidebar
  document.querySelectorAll('.nav-item').forEach(function(item) {
    item.classList.toggle('active', item.dataset.page === page);
  });
  
  // Update mobile nav
  document.querySelectorAll('.mobile-nav-item').forEach(function(item) {
    item.classList.toggle('active', item.dataset.page === page);
  });
  
  // Show correct page
  document.querySelectorAll('.page').forEach(function(p) {
    p.classList.remove('active');
  });
  
  var pageEl = document.getElementById(page + '-page');
  if (pageEl) pageEl.classList.add('active');
  
  // Close mobile sidebar
  document.getElementById('sidebar').classList.remove('mobile-open');
  document.getElementById('sidebarOverlay').classList.remove('active');
  
  if (page === 'analytics') refreshAll();
}

// Desktop nav
document.querySelectorAll('.nav-item').forEach(function(item) {
  item.addEventListener('click', function(e) {
    e.preventDefault();
    navigateToPage(this.dataset.page);
  });
});

// Mobile nav
document.querySelectorAll('.mobile-nav-item').forEach(function(item) {
  item.addEventListener('click', function() {
    navigateToPage(this.dataset.page);
  });
});

// Mobile sidebar
document.getElementById('hamburgerBtn').addEventListener('click', function() {
  document.getElementById('sidebar').classList.add('mobile-open');
  document.getElementById('sidebarOverlay').classList.add('active');
});

document.getElementById('sidebarOverlay').addEventListener('click', function() {
  document.getElementById('sidebar').classList.remove('mobile-open');
  document.getElementById('sidebarOverlay').classList.remove('active');
});

// Close sidebar when nav item clicked on mobile
document.querySelectorAll('.sidebar .nav-item').forEach(function(item) {
  item.addEventListener('click', function() {
    document.getElementById('sidebar').classList.remove('mobile-open');
    document.getElementById('sidebarOverlay').classList.remove('active');
  });
});

// ============ VIEW TOGGLE ============
document.querySelectorAll('.view-btn').forEach(function(btn) {
  btn.addEventListener('click', function() {
    document.querySelectorAll('.view-btn').forEach(function(b) {
      b.classList.remove('active');
    });
    this.classList.add('active');
    currentView = this.dataset.view;
    renderMetrics();
  });
});

// ============ USERS ============
function renderUsers() {
  var siteNames = {
    site1: 'example.com',
    site2: 'shop.example.com',
    site3: 'analytics.example.com',
    site4: 'corp.example.com',
    site5: 'global.example.com'
  };
  
  var html = '';
  users.forEach(function(u) {
    html += '<tr>';
    html += '<td><strong>' + u.name + '</strong></td>';
    html += '<td>' + u.email + '</td>';
    html += '<td><span class="role-badge role-' + u.role + '">' + u.role + '</span></td>';
    html += '<td>' + u.sites.map(function(s) { return siteNames[s]; }).join(', ') + '</td>';
    html += '<td>' + u.lastActive + '</td>';
    html += '<td>';
    html += '<button class="action-btn edit" onclick="editUser(' + u.id + ')"><i class="fas fa-edit"></i></button>';
    html += '<button class="action-btn delete" onclick="deleteUser(' + u.id + ')"><i class="fas fa-trash"></i></button>';
    html += '</td>';
    html += '</tr>';
  });
  
  document.getElementById('usersTableBody').innerHTML = html;
}

window.openUserModal = function(userId) {
  var form = document.getElementById('userForm');
  form.reset();
  
  if (userId) {
    document.getElementById('modalTitle').textContent = 'Edit User';
    var u = users.find(function(x) { return x.id === userId; });
    if (u) {
      document.getElementById('userId').value = u.id;
      document.getElementById('userName').value = u.name;
      document.getElementById('userEmail').value = u.email;
      document.getElementById('userRole').value = u.role;
      
      var sitesSelect = document.getElementById('userSites');
      Array.from(sitesSelect.options).forEach(function(o) {
        o.selected = u.sites.includes(o.value);
      });
    }
  } else {
    document.getElementById('modalTitle').textContent = 'Add User';
    document.getElementById('userId').value = '';
  }
  
  document.getElementById('userModal').classList.add('active');
};

window.closeUserModal = function() {
  document.getElementById('userModal').classList.remove('active');
};

window.editUser = function(id) {
  window.openUserModal(id);
};

window.deleteUser = function(id) {
  if (confirm('Delete this user?')) {
    users = users.filter(function(u) { return u.id !== id; });
    renderUsers();
  }
};

document.getElementById('userForm').addEventListener('submit', function(e) {
  e.preventDefault();
  
  var id = document.getElementById('userId').value;
  var data = {
    name: document.getElementById('userName').value,
    email: document.getElementById('userEmail').value,
    role: document.getElementById('userRole').value,
    sites: Array.from(document.getElementById('userSites').selectedOptions).map(function(o) { return o.value; }),
    lastActive: new Date().toISOString().split('T')[0]
  };
  
  if (id) {
    var index = users.findIndex(function(u) { return u.id === parseInt(id); });
    if (index !== -1) {
      users[index] = Object.assign({}, users[index], data);
    }
  } else {
    data.id = Math.max.apply(Math, users.map(function(u) { return u.id; })) + 1;
    users.push(data);
  }
  
  renderUsers();
  window.closeUserModal();
});

document.getElementById('userModal').addEventListener('click', function(e) {
  if (e.target === this) {
    window.closeUserModal();
  }
});

// ============ METRICS ============
function getCurrentMetric() {
  var siteData = dataCenter.data[currentSite];
  return siteData ? siteData.MetricActive : null;
}

function refreshAll() {
  var metric = getCurrentMetric();
  document.getElementById('recordsValue').textContent = metric ? metric.RecordsPerPeriod.value.toLocaleString() : '0';
  renderMetrics();
  renderChart();
}

function escapeHtml(text) {
  return String(text).replace(/[&<>"]/g, function(m) {
    if (m === '&') return '&amp;';
    if (m === '<') return '&lt;';
    if (m === '>') return '&gt;';
    if (m === '"') return '&quot;';
    return m;
  });
}

function renderMetrics() {
  var metric = getCurrentMetric();
  if (!metric) return;
  
  var defs = [
    { key: 'TopIPs', title: 'Top IPs', icon: 'fa-network-wired' },
    { key: 'TopASN', title: 'Top ASN', icon: 'fa-cloud' },
    { key: 'TopCountries', title: 'Top Countries', icon: 'fa-globe' },
    { key: 'TopCities', title: 'Top Cities', icon: 'fa-city' },
    { key: 'TopURLs', title: 'Top URLs', icon: 'fa-link' },
    { key: 'TopOperatingSystems', title: 'Top OS', icon: 'fa-laptop' },
    { key: 'TopBrowsers', title: 'Top Browsers', icon: 'fa-chrome' }
  ];
  
  var container = document.getElementById('metricsContainer');
  var html = '';
  
  if (currentView === 'horizontal') {
    // HORIZONTAL ROW VIEW
    container.className = 'metrics-horizontal';
    
    // Find max value for bar scaling
    var maxVal = 0;
    defs.forEach(function(m) {
      var meta = metric[m.key];
      if (meta) {
        meta.values.forEach(function(v) {
          if (v > maxVal) maxVal = v;
        });
      }
    });
    
    defs.forEach(function(m) {
      var meta = metric[m.key];
      if (!meta) return;
      
      var mask = meta.getOccupiedMask();
      var items = '';
      
      for (var i = 0; i < 7; i++) {
        if (mask & (1 << i)) {
          var name = String(meta.names[i]);
          var pct = maxVal > 0 ? (meta.values[i] / maxVal * 100) : 0;
          
          items += '<div class="metric-row-item">';
          items += '<span class="rank">#' + (i + 1) + '</span>';
          items += '<span class="name">' + escapeHtml(name) + '</span>';
          items += '<span class="value">' + meta.values[i].toLocaleString() + '</span>';
          items += '<div class="metric-row-bar"><div class="metric-row-bar-fill" style="width:' + pct + '%"></div></div>';
          items += '</div>';
        }
      }
      
      if (!items) {
        items = '<span class="empty-metric">No data</span>';
      }
      
      html += '<div class="metric-row-card">';
      html += '<div class="metric-row-header"><i class="fas ' + m.icon + '"></i><span>' + m.title + '</span></div>';
      html += '<div class="metric-row-items">' + items + '</div>';
      html += '</div>';
    });
    
  } else {
    // GRID CARD VIEW (default)
    container.className = 'metrics-grid';
    
    defs.forEach(function(m) {
      var meta = metric[m.key];
      if (!meta) return;
      
      var mask = meta.getOccupiedMask();
      var items = '';
      
      for (var i = 0; i < 7; i++) {
        if (mask & (1 << i)) {
          var name = String(meta.names[i]);
          items += '<li class="top-item">';
          items += '<span class="item-name"><i class="fas fa-circle" style="font-size:0.5rem;color:var(--accent-primary);"></i> ' + escapeHtml(name) + '</span>';
          items += '<span class="item-value">' + meta.values[i].toLocaleString() + '</span>';
          items += '</li>';
        }
      }
      
      if (!items) {
        items = '<div class="empty-metric">No data</div>';
      }
      
      html += '<div class="metric-card">';
      html += '<div class="card-header"><span>' + m.title + '</span><i class="fas ' + m.icon + '"></i></div>';
      html += '<ul class="top-list">' + items + '</ul>';
      html += '</div>';
    });
  }
  
  container.innerHTML = html;
}

function renderChart() {
  var metric = getCurrentMetric();
  if (!metric) return;
  
  var months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
  var now = new Date();
  var labels = [];
  var values = [];
  
  for (var i = 5; i >= 0; i--) {
    var d = new Date(now.getFullYear(), now.getMonth() - i, 1);
    labels.push(months[d.getMonth()] + ' ' + d.getFullYear());
    values.push(Math.floor(metric.RecordsPerPeriod.value * (0.5 + Math.random() * 0.5)));
  }
  
  labels.push(months[now.getMonth()] + ' ' + now.getFullYear() + ' (active)');
  values.push(metric.RecordsPerPeriod.value);
  
  var isDark = getTheme() === 'dark';
  var ctx = document.getElementById('historyChart').getContext('2d');
  
  if (window.chartInstance) {
    window.chartInstance.destroy();
  }
  
  window.chartInstance = new Chart(ctx, {
    type: 'line',
    data: {
      labels: labels,
      datasets: [{
        data: values,
        borderColor: isDark ? '#60a5fa' : '#2563eb',
        backgroundColor: isDark ? 'rgba(96,165,250,0.1)' : 'rgba(37,99,235,0.08)',
        borderWidth: 3,
        tension: 0.2,
        fill: true,
        pointBackgroundColor: isDark ? '#93c5fd' : '#1e40af',
        pointRadius: 4
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: { display: false }
      },
      scales: {
        x: {
          grid: { color: isDark ? '#334155' : '#e2e8f0' },
          ticks: { color: isDark ? '#94a3b8' : '#64748b' }
        },
        y: {
          grid: { color: isDark ? '#334155' : '#e2e8f0' },
          ticks: {
            color: isDark ? '#94a3b8' : '#64748b',
            callback: function(v) { return v.toLocaleString(); }
          }
        }
      }
    }
  });
}

// ============ PERIOD TOGGLE ============
document.querySelectorAll('.period-btn').forEach(function(btn) {
  btn.addEventListener('click', function() {
    document.querySelectorAll('.period-btn').forEach(function(b) {
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

document.getElementById('quickDayBtn').addEventListener('click', function() {
  var dayBtn = document.querySelector('.period-btn[data-period="day"]');
  if (dayBtn) dayBtn.click();
});

// ============ SITE SELECTOR ============
document.getElementById('siteSelector').addEventListener('change', function() {
  currentSite = this.value;
  refreshAll();
});


// ============ INITIALIZATION ============
function getDaysInMonth(year, month) {
  // month is 0-indexed (0 = January, 11 = December)
  return new Date(year, month + 1, 0).getDate();
}

function populateDays(year, month, selectedDay) {
  var daysInMonth = getDaysInMonth(year, month);
  var dayHtml = '';
  for (var d = 1; d <= daysInMonth; d++) {
    dayHtml += '<option value="' + d + '"' + (d === selectedDay ? ' selected' : '') + '>' + d + '</option>';
  }
  document.getElementById('daySelect').innerHTML = dayHtml;
}

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
  
  renderUsers();
  refreshAll();
}