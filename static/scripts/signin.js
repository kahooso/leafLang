document.addEventListener("DOMContentLoaded", function () {
    const signInForm = document.getElementById("SignInForm");
    if (signInForm) {
        signInForm.addEventListener("submit", async function (event) {
            event.preventDefault();

            const email = document.getElementById("SignInEmail").value;
            const password = document.getElementById("SignInPassword").value;
            const message = document.getElementById("SignInMsg");

            const response = await fetch("api/auth/signin", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ email, password }),
            });

            const data = await response.json();

            if (response.ok) {
                localStorage.setItem("token", data.token);
                localStorage.setItem("role", data.role);

                message.style.color = "var(--success)";
                message.textContent = "Login successful! Redirecting...";
                setTimeout(() => {
                    window.location.href = "/main";
                }, 1500);
            } else {
                message.style.color = "var(--error)";
                message.textContent = data.error || "Invalid email or password!";
            }
        });
    }
});