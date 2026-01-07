// Bring in the elements
const methodSelect = document.getElementById('method');
const urlInput = document.getElementById('url');
const headersInput = document.getElementById('headers');
const bodyInput = document.getElementById('body');
const sendButton = document.getElementById('send');
const responseEl = document.getElementById('response');
const statusEl = document.getElementById('status');
const bodyGroup = document.getElementById('body-group');
const sidebar = document.querySelector('.sidebar');
const toggleButton = document.getElementById('toggleSidebar');
const toggleText = document.getElementById('toggleText');

// Toggle sidebar
toggleButton.addEventListener('click', () => {
    sidebar.classList.toggle('hidden');
    toggleText.textContent = sidebar.classList.contains('hidden') ? 'Show Sidebar' : 'Hide Sidebar';
});

// API endpoint templates
const endpoints = {
    health: {
        method: 'GET',
        url: 'http://localhost:8080/health',
        headers: {},
        body: null
    },
    register: {
        method: 'POST',
        url: 'http://localhost:8080/api/register',
        headers: { 'Content-Type': 'application/json' },
        body: {
            name: 'John Doe',
            email: 'john@example.com',
            phone: '+1234567890',
            password: 'Password123'
        }
    },
    login: {
        method: 'POST',
        url: 'http://localhost:8080/api/login',
        headers: { 'Content-Type': 'application/json' },
        body: {
            email: 'john@example.com',
            password: 'Password123'
        }
    },
    logout: {
        method: 'POST',
        url: 'http://localhost:8080/api/logout',
        headers: { 
            'Content-Type': 'application/json'
        },
        body: null,
        note: 'Token will be sent automatically via cookie'
    },
    'customers-list': {
        method: 'GET',
        url: 'http://localhost:8080/api/customers',
        headers: { 
            'Content-Type': 'application/json'
        },
        body: null,
        note: 'Token will be sent automatically via cookie'
    },
    'customers-create': {
        method: 'POST',
        url: 'http://localhost:8080/api/customers',
        headers: { 
            'Content-Type': 'application/json'
        },
        body: {
            name: 'Jane Doe',
            email: 'john@example.com',
            phone: '+1234567890',
            password: 'Password123'
        },
    },
};

// Handle endpoint selection
document.querySelectorAll('.endpoint-item').forEach(item => {
    item.addEventListener('click', function() {
        const endpointKey = this.getAttribute('data-endpoint');
        const endpoint = endpoints[endpointKey];
        
        // Update active state
        document.querySelectorAll('.endpoint-item').forEach(el => el.classList.remove('active'));
        this.classList.add('active');
        
        // Populate form
        methodSelect.value = endpoint.method;
        urlInput.value = endpoint.url;
        headersInput.value = JSON.stringify(endpoint.headers, null, 2);
        
        if (endpoint.body) {
            bodyInput.value = JSON.stringify(endpoint.body, null, 2);
            bodyGroup.style.display = 'block';  // Already there, but add !important via style
        } else {
            bodyInput.value = '';
            bodyGroup.style.display = endpoint.method === 'GET' || endpoint.method === 'DELETE' ? 'none' : 'block';
        }
        
        // Show note if present
        if (endpoint.note) {
            responseEl.textContent = `Note: ${endpoint.note}\n\nClick Send to test this endpoint.`;
        } else {
            responseEl.textContent = 'Click Send to test this endpoint.';
        }
        statusEl.classList.add('hidden');
    });
});

// Show/hide body input based on method
methodSelect.addEventListener('change', () => {
    const method = methodSelect.value;
    if (method === 'GET' || method === 'DELETE') {
        bodyGroup.style.display = 'none';
    } else {
        bodyGroup.style.display = 'block';
    }
});

// Handle Enter key in URL input
urlInput.addEventListener('keydown', (e) => {
    if (e.key === 'Enter') {
        sendRequest();
    }
});

sendButton.addEventListener('click', sendRequest);

async function sendRequest() {
    const method = methodSelect.value;
    const url = urlInput.value.trim();

    if (!url) {
        responseEl.textContent = 'Error: Please enter a URL or select an endpoint';
        statusEl.classList.add('hidden');
        return;
    }

    // Parse headers
    let headers = {};
    if (headersInput.value.trim()) {
        try {
            headers = JSON.parse(headersInput.value);
        } catch (e) {
            responseEl.textContent = 'Error: Invalid JSON in headers';
            statusEl.classList.add('hidden');
            return;
        }
    }

    // Build fetch options
    const options = {
        method: method,
        headers: headers,
        credentials: 'include' // Include cookies for auth
    };

    // Add body for POST, PUT, PATCH
    if (method !== 'GET' && method !== 'DELETE' && bodyInput.value.trim()) {
        try {
            JSON.parse(bodyInput.value); // Validate JSON
            options.body = bodyInput.value;
        } catch (e) {
            responseEl.textContent = 'Error: Invalid JSON in body';
            statusEl.classList.add('hidden');
            return;
        }
    }

    // Show loading state
    responseEl.textContent = 'Loading...';
    statusEl.classList.add('hidden');
    sendButton.disabled = true;

    try {
        const startTime = Date.now();
        const response = await fetch(url, options);
        const endTime = Date.now();
        const duration = endTime - startTime;

        // Update status
        statusEl.textContent = `${response.status} ${response.statusText} (${duration}ms)`;
        statusEl.classList.remove('hidden', 'success', 'error');
        statusEl.classList.add(response.ok ? 'success' : 'error');

        // Parse response
        const contentType = response.headers.get('content-type');
        let data;
        
        if (contentType && contentType.includes('application/json')) {
            data = await response.json();
            responseEl.textContent = JSON.stringify(data, null, 2);
        } else {
            data = await response.text();
            responseEl.textContent = data;
        }

    } catch (error) {
        responseEl.textContent = `Error: ${error.message}`;
        statusEl.textContent = 'Request Failed';
        statusEl.classList.remove('hidden', 'success');
        statusEl.classList.add('error');
    } finally {
        sendButton.disabled = false;
    }
}