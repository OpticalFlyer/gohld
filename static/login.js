var loggedin = false;

        document.addEventListener('DOMContentLoaded', (event) => {
    var loginForm = document.getElementById('loginForm');

    loginForm.addEventListener('submit', function(event) {
        event.preventDefault(); // Prevent the form from submitting the default way

        var email = document.getElementById('emailInput').value;
        var password = document.getElementById('passwordInput').value;

        // Create an object with the email and password
        var loginCredentials = {
            username: email, // Assuming the backend expects 'username'
            password: password
        };

        // Use the fetch API to send a POST request to /login
        fetch('/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json', // Specify the content type as JSON
            },
            body: JSON.stringify(loginCredentials) // Convert the credentials object to a JSON string
        })
        .then(response => {
            if (response.ok) {
                return response.json();
            } else {
                throw new Error('Failed to log in');
            }
        })
        .then(data => {
            loggedin = true; // Set the logged in state to true
            toggleLoginDialog(); // Update the UI to reflect the logged in state
            // Handle any additional login success actions here
        })
        .catch(error => {
            console.error('Error during login:', error);
            // Handle login errors (show error message to the user, etc.)
        });
    });
});

function logout() {
    fetch('/logout', {
        method: 'POST',
        credentials: 'include' // Necessary to include cookies for the session
    })
    .then(response => {
        if (response.ok) {
            return response.json();
        } else {
            throw new Error('Logout failed');
        }
    })
    .then(data => {
        console.log(data.message); // Or handle the logout success message in the UI
        checkLoginStatus(); // Check login status to update UI
    })
    .catch(error => {
        console.error('Error during logout:', error);
        // Optionally handle logout errors in the UI
    });
}

document.addEventListener('DOMContentLoaded', (event) => {
    // ... (other event listeners)

    var logoutButton = document.getElementById('logout-button');
    logoutButton.addEventListener('click', function(event) {
        logout();
    });
});

        // Function to check login status
        function checkLoginStatus() {
            fetch('/check-auth')
                .then(response => {
                    // If the response is successful, assume the user is logged in
                    if (response.ok) {
                        return response.json();  // Parse the JSON response
                    } else {
                        // If the response is not successful, assume the user is not logged in
                        throw new Error('Not logged in');
                    }
                })
                .then(data => {
                    // Update the login state based on the response
                    loggedin = data.authenticated;
                    toggleLoginDialog(); // Update the login dialog visibility
                })
                .catch(error => {
                    console.error('Error checking login status:', error);
                    loggedin = false;
                    toggleLoginDialog(); // Ensure the login dialog is shown if not logged in
                });
        }

        // Function to toggle the login dialog visibility
        function toggleLoginDialog() {
            var loginDialog = document.getElementById('login-dialog-container');
            if (loggedin) {
                loginDialog.style.display = 'none'; // Hide it
            } else {
                loginDialog.style.display = 'flex'; // Show it, since the container is flex
            }
        }

        // Call the function to check the initial login status
        checkLoginStatus();