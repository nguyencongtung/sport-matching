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
                    // Redirect to the welcome page for profile setup
                    setTimeout(() => {
                        window.location.href = '/welcome.html';
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

    // Profile setup navigation
    const agreeButton = document.getElementById('agreeButton');
    if (agreeButton) {
        agreeButton.addEventListener('click', () => {
            window.location.href = '/profile-setup-name.html';
        });
    }

    const firstNameInput = document.getElementById('firstName');
    if (firstNameInput) {
        const nextButton = document.getElementById('nextButton');
        nextButton.addEventListener('click', () => {
            const firstName = firstNameInput.value.trim();
            if (firstName) {
                // In a real application, you'd send this to the server
                console.log('First Name:', firstName);
                window.location.href = '/profile-setup-details.html';
            } else {
                alert('Please enter your first name.');
            }
        });
    }

    const dobInput = document.getElementById('dob');
    const genderSelect = document.getElementById('gender');
    const showGenderCheckbox = document.getElementById('showGender');
    if (dobInput && genderSelect && showGenderCheckbox) {
        const nextButton = document.getElementById('nextButton');
        nextButton.addEventListener('click', () => {
            const dob = dobInput.value;
            const gender = genderSelect.value;
            const showGender = showGenderCheckbox.checked;

            if (dob && gender) {
                console.log('DOB:', dob, 'Gender:', gender, 'Show Gender:', showGender);
                window.location.href = '/profile-setup-distance.html';
            } else {
                alert('Please fill in your birthday and select your gender.');
            }
        });
    }

    const distanceRange = document.getElementById('distanceRange');
    const distanceValue = document.getElementById('distanceValue');
    if (distanceRange && distanceValue) {
        distanceRange.addEventListener('input', () => {
            distanceValue.textContent = distanceRange.value;
        });

        const nextButton = document.getElementById('nextButton');
        nextButton.addEventListener('click', () => {
            const distance = distanceRange.value;
            console.log('Distance Preference:', distance);
            window.location.href = '/profile-setup-pictures.html';
        });
    }

    const pictureInput = document.getElementById('pictureInput');
    const imagePreview = document.getElementById('imagePreview');
    if (pictureInput && imagePreview) {
        let selectedFiles = [];

        pictureInput.addEventListener('change', (event) => {
            const files = Array.from(event.target.files);
            selectedFiles = [...selectedFiles, ...files].slice(0, 6); // Max 6 pictures
            renderImagePreviews();
        });

        function renderImagePreviews() {
            imagePreview.innerHTML = '';
            selectedFiles.forEach((file, index) => {
                const reader = new FileReader();
                reader.onload = (e) => {
                    const imgContainer = document.createElement('div');
                    imgContainer.className = 'image-preview-item';
                    const img = document.createElement('img');
                    img.src = e.target.result;
                    imgContainer.appendChild(img);

                    const removeButton = document.createElement('button');
                    removeButton.textContent = 'X';
                    removeButton.className = 'remove-image';
                    removeButton.addEventListener('click', () => {
                        selectedFiles.splice(index, 1);
                        renderImagePreviews();
                    });
                    imgContainer.appendChild(removeButton);
                    imagePreview.appendChild(imgContainer);
                };
                reader.readAsDataURL(file);
            });
        }

        const finishButton = document.getElementById('nextButton');
        if (finishButton) {
            finishButton.addEventListener('click', () => {
                if (selectedFiles.length >= 1) {
                    console.log('Selected Pictures:', selectedFiles);
                    alert('Profile setup complete! (Pictures would be uploaded to server)');
                    // In a real application, you'd upload these to the server
                    window.location.href = '/dashboard.html'; // Redirect to dashboard after setup
                } else {
                    alert('Please add at least 1 picture.');
                }
            });
        }
    }
});
