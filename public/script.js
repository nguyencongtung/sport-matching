document.addEventListener('DOMContentLoaded', () => {
    const registerForm = document.getElementById('registerForm');
    const loginForm = document.getElementById('loginForm');
    const profileForm = document.getElementById('profileForm');
    const messageElement = document.getElementById('message'); // Assuming a message element exists for feedback

    // Helper function to get JWT token from cookies
    function getCookie(name) {
        const value = `; ${document.cookie}`;
        const parts = value.split(`; ${name}=`);
        if (parts.length === 2) return parts.pop().split(';').shift();
        return null;
    }

    // Helper function to decode JWT token and get user ID
    function getUserIdFromToken() {
        const token = getCookie('token'); // Assuming the JWT token is stored in a cookie named 'token'
        if (!token) {
            console.error('No token found');
            return null;
        }
        try {
            const base64Url = token.split('.')[1];
            const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
            const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
                return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
            }).join(''));
            const payload = JSON.parse(jsonPayload);
            return payload.user_id; // Assuming the user ID is stored as 'user_id' in the token payload
        } catch (e) {
            console.error('Error decoding token:', e);
            return null;
        }
    }

    // Function to show messages
    function showMessage(text, type) {
        if (messageElement) {
            messageElement.textContent = text;
            messageElement.className = `message ${type}`;
        } else {
            console.log(`Message (${type}): ${text}`);
        }
    }

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
                        window.location.href = '/dashboard.html';
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

    if (profileForm) {
        const userId = getUserIdFromToken();

        // Fetch user profile data
        const fetchUserProfile = async () => {
            if (!userId) {
                showMessage('User not logged in or ID not found.', 'error');
                return;
            }
            try {
                const response = await fetch(`/api/user/${userId}`, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${getCookie('token')}` // Assuming token is needed for auth
                    },
                });
                const result = await response.json();
                if (response.ok && result.user) {
                    const user = result.user;
                    document.getElementById('names').value = user.names || '';
                    document.getElementById('gender').value = user.gender || '';
                    document.getElementById('date_of_birth').value = user.date_of_birth ? user.date_of_birth.split('T')[0] : ''; // Format date for input type="date"
                    document.getElementById('bio').value = user.bio || '';
                    document.getElementById('interests').value = user.interests || '';
                    document.getElementById('looking_for').value = user.looking_for || '';
                    // Handle profile pictures
                    const pictureFramesContainer = document.querySelector('.picture-frames');
                    pictureFramesContainer.innerHTML = ''; // Clear existing frames

                    const profilePictureURLs = user.profile_picture_urls || [];
                    for (let i = 0; i < 6; i++) {
                        const frameDiv = document.createElement('div');
                        frameDiv.className = 'picture-frame';

                        const img = document.createElement('img');
                        img.src = profilePictureURLs[i] || 'https://via.placeholder.com/100x100?text=Add+Image'; // Placeholder image
                        img.alt = `Profile Picture ${i + 1}`;
                        img.style.width = '100px';
                        img.style.height = '100px';
                        img.style.border = '1px solid #ccc';
                        img.style.objectFit = 'cover';

                        const input = document.createElement('input');
                        input.type = 'url';
                        input.className = 'profile-picture-input';
                        input.placeholder = `Picture ${i + 1} URL`;
                        input.value = profilePictureURLs[i] || '';
                        input.dataset.index = i; // Store index for easy access

                        frameDiv.appendChild(img);
                        frameDiv.appendChild(input);
                        pictureFramesContainer.appendChild(frameDiv);

                        // Update image preview on input change
                        input.addEventListener('input', (event) => {
                            img.src = event.target.value || 'https://via.placeholder.com/100x100?text=Add+Image';
                        });
                    }
                    document.getElementById('location').value = user.location || '';
                    document.getElementById('distance_preference').value = user.distance_preference || 0;
                } else {
                    showMessage(result.error || 'Failed to fetch profile data', 'error');
                }
            } catch (error) {
                console.error('Error fetching profile:', error);
                showMessage('An unexpected error occurred while fetching profile', 'error');
            }
        };

        // Update user profile data
        profileForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            if (!userId) {
                showMessage('User not logged in or ID not found.', 'error');
                return;
            }

            const data = {
                names: document.getElementById('names').value,
                gender: document.getElementById('gender').value,
                date_of_birth: document.getElementById('date_of_birth').value,
                bio: document.getElementById('bio').value,
                interests: document.getElementById('interests').value,
                looking_for: document.getElementById('looking_for').value,
                location: document.getElementById('location').value,
                distance_preference: parseInt(document.getElementById('distance_preference').value, 10) || 0,
                profile_picture_urls: Array.from(document.querySelectorAll('.profile-picture-input')).map(input => input.value).filter(url => url !== '')
            };

            try {
                const response = await fetch(`/api/user/profile/${userId}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${getCookie('token')}` // Assuming token is needed for auth
                    },
                    body: JSON.stringify(data),
                });

                const result = await response.json();
                if (response.ok) {
                    showMessage(result.message, 'success');
                } else {
                    showMessage(result.error || 'Failed to update profile', 'error');
                }
            } catch (error) {
                console.error('Error updating profile:', error);
                showMessage('An unexpected error occurred while updating profile', 'error');
            }
        });

        // Fetch profile data when the page loads
        fetchUserProfile();
    }
});
