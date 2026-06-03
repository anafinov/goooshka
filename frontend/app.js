const API_URL = 'http://localhost:8080';

// DOM Elements
const authSection = document.getElementById('auth-section');
const dashboardSection = document.getElementById('dashboard-section');
const btnLogin = document.getElementById('btn-login');
const btnRegister = document.getElementById('btn-register');
const btnLogout = document.getElementById('btn-logout');
const btnRequest = document.getElementById('btn-request');
const usernameInput = document.getElementById('username');
const passwordInput = document.getElementById('password');
const certTypeInput = document.getElementById('cert-type');
const authError = document.getElementById('auth-error');
const requestError = document.getElementById('request-error');
const requestsBody = document.getElementById('requests-body');

// State
let pollingInterval = null;

// Initialize
function init() {
    const token = localStorage.getItem('token');
    if (token) {
        showDashboard();
    } else {
        showAuth();
    }
}

// UI Switching
function showAuth() {
    authSection.classList.add('active');
    authSection.classList.remove('hidden');
    dashboardSection.classList.remove('active');
    setTimeout(() => dashboardSection.style.display = 'none', 400);
    authSection.style.display = 'block';
    
    if (pollingInterval) clearInterval(pollingInterval);
}

function showDashboard() {
    dashboardSection.classList.add('active');
    dashboardSection.classList.remove('hidden');
    authSection.classList.remove('active');
    setTimeout(() => authSection.style.display = 'none', 400);
    dashboardSection.style.display = 'block';
    
    fetchRequests();
    // Poll every 3 seconds
    pollingInterval = setInterval(fetchRequests, 3000);
}

// Helpers
function showError(element, msg) {
    element.textContent = msg;
    setTimeout(() => element.textContent = '', 5000);
}

function getAuthHeaders() {
    return {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
    };
}

// Event Listeners
btnLogin.addEventListener('click', async () => {
    const username = usernameInput.value;
    const password = passwordInput.value;
    
    if (!username || !password) return showError(authError, 'Заполните все поля');

    try {
        const res = await fetch(`${API_URL}/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        });
        
        if (!res.ok) throw new Error('Неверные учетные данные');
        
        const data = await res.json();
        localStorage.setItem('token', data.token);
        showDashboard();
    } catch (err) {
        showError(authError, err.message);
    }
});

btnRegister.addEventListener('click', async () => {
    const username = usernameInput.value;
    const password = passwordInput.value;
    
    if (!username || !password) return showError(authError, 'Заполните все поля');

    try {
        const res = await fetch(`${API_URL}/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        });
        
        if (!res.ok) throw new Error('Пользователь уже существует');
        
        showError(authError, 'Успешная регистрация! Теперь войдите.');
        authError.style.color = 'var(--success)';
        setTimeout(() => authError.style.color = '', 3000);
    } catch (err) {
        showError(authError, err.message);
    }
});

btnLogout.addEventListener('click', () => {
    localStorage.removeItem('token');
    showAuth();
});

btnRequest.addEventListener('click', async () => {
    const type = certTypeInput.value;
    btnRequest.disabled = true;
    
    try {
        const res = await fetch(`${API_URL}/request`, {
            method: 'POST',
            headers: getAuthHeaders(),
            body: JSON.stringify({ type })
        });
        
        if (!res.ok) throw new Error('Ошибка создания заявки');
        
        fetchRequests();
    } catch (err) {
        showError(requestError, err.message);
    } finally {
        btnRequest.disabled = false;
    }
});

async function fetchRequests() {
    try {
        const res = await fetch(`${API_URL}/requests`, {
            headers: getAuthHeaders()
        });
        
        if (res.status === 401) {
            localStorage.removeItem('token');
            showAuth();
            return;
        }
        
        if (!res.ok) throw new Error('Ошибка получения данных');
        
        const data = await res.json();
        renderTable(data);
    } catch (err) {
        console.error(err);
    }
}

function renderTable(requests) {
    requestsBody.innerHTML = '';
    
    if (!requests || requests.length === 0) {
        requestsBody.innerHTML = `<tr><td colspan="4" style="text-align: center; color: var(--text-secondary)">У вас пока нет заявок</td></tr>`;
        return;
    }

    requests.forEach(req => {
        const date = new Date(req.created_at).toLocaleString('ru-RU');
        const statusClass = req.status === 'pending' ? 'status-pending' : 'status-ready';
        const statusText = req.status === 'pending' ? 'В обработке' : 'Готова';

        const tr = document.createElement('tr');
        tr.innerHTML = `
            <td>#${req.id}</td>
            <td>${req.type}</td>
            <td>${date}</td>
            <td><span class="status-badge ${statusClass}">${statusText}</span></td>
        `;
        requestsBody.appendChild(tr);
    });
}

// Start
init();
