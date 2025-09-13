document.addEventListener('DOMContentLoaded', () => {
    const registerForm = document.getElementById('registerForm');
    const loginForm = document.getElementById('loginForm');
    const messageElement = document.getElementById('message');

    if (registerForm) {
        registerForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const formData = new FormData(registerForm);
            const data = Object.fromEntries(formData.entries());

            try {
                const response = await fetch('/api/auth/register', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(data),
                });

                const result = await response.json();
                if (response.ok) {
                    showMessage(result.message, 'success');
                    registerForm.reset();
                    setTimeout(() => {
                        window.location.href = '/login.html';
                    }, 2000);
                } else {
                    showMessage(result.error || 'Registration failed', 'error');
                }
            } catch (error) {
                console.error('Error:', error);
                showMessage('An unexpected error occurred', 'error');
            }
        });
    }

    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const formData = new FormData(loginForm);
            const data = Object.fromEntries(formData.entries());

            try {
                const response = await fetch('/api/auth/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(data),
                });

                const result = await response.json();
                if (response.ok) {
                    showMessage(result.message, 'success');
                    loginForm.reset();
                    // Redirect to a protected page or dashboard
                    setTimeout(() => {
                        window.location.href = '/dashboard.html'; // Create a dashboard.html later
                    }, 2000);
                } else {
                    showMessage(result.error || 'Login failed', 'error');
                }
            } catch (error) {
                console.error('Error:', error);
                showMessage('An unexpected error occurred', 'error');
            }
        });
    }

    function showMessage(text, type) {
        messageElement.textContent = text;
        messageElement.className = `message ${type}`;
    }

    const logoutButton = document.getElementById('logoutButton');
    if (logoutButton) {
        logoutButton.addEventListener('click', async () => {
            try {
                const response = await fetch('/api/auth/logout', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                });

                if (response.ok) {
                    window.location.href = '/login.html';
                } else {
                    alert('Logout failed');
                }
            } catch (error) {
                console.error('Error:', error);
                alert('An unexpected error occurred during logout');
            }
        });
    }
});
