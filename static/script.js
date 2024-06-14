//  static/script.js

document.addEventListener("DOMContentLoaded", () => {
    // Check if signup form exists before adding event listener
    const signupForm = document.getElementById("signup-form");
    if (signupForm) {
        signupForm.addEventListener("submit", signupTeacher);
    }

    // Check if login form exists before adding event listener
    const loginForm = document.getElementById("login-form");
    if (loginForm) {
        loginForm.addEventListener("submit", loginTeacher);
    }
});

async function loginTeacher(event) {
    event.preventDefault();
    
    const email = document.getElementById("login-email").value;
    const password = document.getElementById("login-password").value;

    try {
        const response = await fetch("/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ email, password })
        });

        if (response.ok) {
            const data = await response.json();
            alert("Login successful!");
            // Save token and teacher data to localStorage or cookies for future API requests
            localStorage.setItem("token", data.token);
            localStorage.setItem("teacher_id", data.teacher_id);
            localStorage.setItem("teacher_name", data.teacher.name);
            // Redirect to main.html
            window.location.href = "main.html";
        } else {
            const errorData = await response.json();
            alert(`Login failed: ${errorData.message}`);
        }
    } catch (error) {
        console.error("Error logging in:", error);
        alert("An error occurred during login. Please try again later.");
    }
}


async function signupTeacher(event) {
    event.preventDefault();
    
    const name = document.getElementById("signup-name").value;
    const email = document.getElementById("signup-email").value;
    const password = document.getElementById("signup-password").value;

    try {
        const response = await fetch("/signup", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ name, email, password })
        });

        if (response.ok) {
            alert("Signup successful!");
        } else {
            const errorData = await response.json();
            alert(`Signup failed: ${errorData.message}`);
        }
    } catch (error) {
        console.error("Error signing up:", error);
        alert("An error occurred during signup. Please try again later.");
    }
}
