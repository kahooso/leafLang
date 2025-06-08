document.addEventListener('DOMContentLoaded', function () {
    const statusFilter = document.getElementById('status-filter');
    const searchInput = document.getElementById('search-input');
    const searchBtn = document.getElementById('search-btn');
    const prevPageBtn = document.getElementById('prev-page');
    const nextPageBtn = document.getElementById('next-page');
    const pageInfo = document.getElementById('page-info');

    let currentPage = 1;
    const limit = 20;
    let totalPages = 1;

    function loadWords(page = 1) {
        const status = statusFilter.value;
        const search = searchInput.value.trim();

        const url = new URL('/words', window.location.origin);
        url.searchParams.append('page', page);
        url.searchParams.append('limit', limit);
        if (status) url.searchParams.append('status', status);
        if (search) url.searchParams.append('search', search);

        window.location.href = url.toString();
    }

    statusFilter.addEventListener('change', () => loadWords(1));
    searchBtn.addEventListener('click', () => loadWords(1));
    searchInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') loadWords(1);
    });

    prevPageBtn.addEventListener('click', () => {
        if (currentPage > 1) loadWords(currentPage - 1);
    });

    nextPageBtn.addEventListener('click', () => {
        if (currentPage < totalPages) loadWords(currentPage + 1);
    });

    const editModal = document.getElementById('edit-word-modal');
    const editForm = document.getElementById('edit-word-form');

    document.querySelectorAll('.edit-btn').forEach(btn => {
        btn.addEventListener('click', function () {
            const row = this.closest('tr');
            const wordId = this.getAttribute('data-id');

            document.getElementById('edit-word-id').value = wordId;
            document.getElementById('edit-original-word').value = row.cells[0].textContent.trim();
            document.getElementById('edit-translation').value = row.cells[1].textContent.trim();
            document.getElementById('edit-example').value = row.cells[2].textContent.trim() === '-' ? '' : row.cells[2].textContent.trim();

            const statusName = row.cells[3].querySelector('.status-badge').classList[1];
            let statusValue = '1';
            if (statusName === 'Known')
                statusValue = '2';
            if (statusName === 'Learned')
                statusValue = '3';
            document.getElementById('edit-status').value = statusValue;

            editModal.style.display = 'block';
        })
    })

    document.querySelector('.close-modal').addEventListener('click', () => {
        editModal.style.display = 'none';
    });

    editForm.addEventListener('submit', function (e) {
        e.preventDefault();
        console.log('SADASD')

        const wordId = document.getElementById('edit-word-id').value;
        const originalWord = document.getElementById('edit-original-word').value;
        const translation = document.getElementById('edit-translation').value;
        const example = document.getElementById('edit-example').value;
        const statusId = document.getElementById('edit-status').value;
        console.log(statusId)

        fetch(`/words/${wordId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                original_word: originalWord,
                translation: translation,
                example: example,
                status_id: parseInt(statusId)
            })
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Update failed');
                }
                return response.json();
            })
            .then(data => {
                alert("Update successfully!");
                editModal.style.display = 'none';
                loadWords(currentPage);
            })
            .catch(error => {
                console.error('Error:', error);
                alert('Failed to update word');
            });
    });

    document.querySelectorAll('.delete-btn').forEach(btn => {
        btn.addEventListener('click', function (e) {
            if (!confirm('Are you sure you want to delete this word?')) return;

            const wordId = this.getAttribute('data-id');
            fetch(`/words/${wordId}`, { method: 'DELETE' })
                .then(() => loadWords(currentPage));
        });
    });

    currentPage = "{{ .Pagination.Page }}";
    totalPages = "{{ .Pagination.Pages }}";
    pageInfo.textContent = `Page ${currentPage} of ${totalPages}`;
    prevPageBtn.disabled = currentPage <= 1;
    nextPageBtn.disabled = currentPage >= totalPages;

    statusFilter.value = "{{ .StatusFilter }}";
    searchInput.value = "{{ .SearchQuery }}";

    function escapeHtml(unsafe) {
        return unsafe
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;");
    }
});