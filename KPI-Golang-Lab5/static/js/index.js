document.addEventListener("DOMContentLoaded", function () {
    document.querySelectorAll(".calc-link-button").forEach(button => {
        button.addEventListener("click", function () {
            const targetUrl = this.getAttribute("data-href");
            if (targetUrl) {
                window.location.href = targetUrl;
            }
        });
    });
});