document.addEventListener("DOMContentLoaded", function () {
    const signUpForm = document.getElementById("SignUpForm");
    if (signUpForm) {
        signUpForm.addEventListener("submit", async function (event) {
            event.preventDefault();

            const email = document.getElementById("SignUpEmail").value;
            const phone = document.getElementById("SignUpPhone").value;
            const firstName = document.getElementById("SignUpFirstName").value;
            const lastName = document.getElementById("SignUpLastName").value;
            const password = document.getElementById("SignUpPassword").value;
            const confirmPassword = document.getElementById("SignUpConfirmPassword").value
            const message = document.getElementById("SignUpMsg");

            if (password !== confirmPassword) {
                message.style.color = "red";
                message.textContent = "Password don't match!"
                return;
            }

            const response = await fetch("api/auth/signup", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ email, phone, firstName, lastName, password }),
            });

            const data = await response.json();

            if (response.ok) {
                message.style.color = "var(--success)";
                message.textContent = "Registration successful! Redirecting...";
                setTimeout(() => {
                    window.location.href = "/signin";
                }, 3000);

            } else {
                message.style.color = "var(--error)";
                message.textContent = data.error || "Registration failed!";
            }
        });
    }
});