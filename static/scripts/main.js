document.addEventListener('DOMContentLoaded', function () {
    const modal = document.getElementById('add-word-modal');
    const addWordBtn = document.querySelector('.add-words');
    const closeBtn = document.querySelector('.close-modal');
    const addWordForm = document.getElementById('add-word-form');

    addWordBtn.addEventListener('click', function (e) {
        e.preventDefault();
        modal.style.display = 'block';
        document.body.style.overflow = 'hidden';
    });

    closeBtn.addEventListener('click', function () {
        modal.style.display = 'none';
        document.body.style.overflow = 'auto';
    });

    window.addEventListener('click', function (e) {
        if (e.target === modal) {
            modal.style.display = 'none';
            document.body.style.overflow = 'auto';
        }
    });

    addWordForm.addEventListener('submit', async function (e) {
        e.preventDefault();

        const originalWord = document.getElementById('original-word').value;
        const translation = document.getElementById('translation').value;
        const example = document.getElementById('example').value;

        try {
            const response = await fetch('/word/add', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    original_word: originalWord,
                    translation: translation,
                    example: example
                })
            });

            const result = await response.json();

            if (result.success) {
                modal.style.display = 'none';
                document.body.style.overflow = 'auto';
                addWordForm.reset();

                showNotification('Word has successfully added!', 'success');
            } else {
                showNotification(result.error || 'Error while adding a word', 'error');
            }
        } catch (error) {
            showNotification('Server error', 'error');
            console.error('Error:', error);
        }
    });

    function showNotification(message, type) {
        const notification = document.createElement('div');
        notification.className = `notification ${type}`;
        notification.textContent = message;
        document.body.appendChild(notification);

        setTimeout(() => {
            notification.classList.add('fade-out');
            setTimeout(() => notification.remove(), 500);
        }, 3000);
    }
});