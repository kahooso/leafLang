document.addEventListener("DOMContentLoaded", function () {
    const profileUpdateForm = document.getElementById("ProfileUpdateForm");
    const userId = "{{ .id }}";
    const fileInput = document.getElementById('ProfileImage');
    const fileLabel = document.querySelector('.file-label');

    fileInput.addEventListener('change', function () {
        if (this.files && this.files.length > 0) {
            fileLabel.textContent = this.files[0].name;
        } else {
            fileLabel.textContent = 'Choose profile image';
        }
    });

    profileUpdateForm.addEventListener("submit", async function (event) {
        event.preventDefault();

        const phone = document.getElementById("Phone").value.trim();
        const learningGoal = document.getElementById("LearningGoal").value;
        const profileImage = document.getElementById("ProfileImage").files[0];

        const formData = new FormData();
        formData.append("phone", phone);
        formData.append("learningGoal", learningGoal);
        if (profileImage) {
            formData.append("image", profileImage);
        }

        try {
            const token = localStorage.getItem("token");
            if (!token) {
                throw new Error("Authentication required");
            }

            const response = await fetch(`/profile/update/${userId}`, {
                method: "PUT",
                headers: {
                    "Authorization": `Bearer ${token}`
                },
                body: formData,
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || "Profile update failed");
            }

            alert("Profile updated successfully!");
            window.location.reload();
        } catch (error) {
            console.error("Update error:", error);
            alert(error.message);
        }
    });
});