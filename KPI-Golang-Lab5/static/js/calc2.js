document.addEventListener("DOMContentLoaded", function () {
    const backButton = document.getElementById("back-button");
    if (backButton) {
        backButton.addEventListener("click", function () {
            window.history.back();
        });
    }

    const form = document.getElementById("calculator-form");
    form.addEventListener("submit", async function (event) {
        event.preventDefault();

        const formData = new FormData(form);

        const response = await fetch("/evaluate2", {
            method: "POST",
            body: formData
        });

        if (!response.ok) {
            if (response.status === 400) {
                alert("Помилка у введених даних...");
            } else if (response.status === 500) {
                alert("Помилка сервера...");
            } else {
                alert("Невідома помилка...");
            }
            return;
        }

        const data = await response.json();

        showResultSection(`
            <h2>Результат</h2>
            <p class="result-paragraph">Очікується, що планове недовідпущення<br>електроенергії складе: 
                <b>${data.expectedOutagesScheduled} кВт⋅год.</b></p>
            <p class="result-paragraph">Очікується, що аварійне недовідпущення<br>електроенергії складе: 
                <b>${data.expectedOutagesEmergency} кВт⋅год.</b></p>
            <p class="result-paragraph">Коефіцієнт аварійного простою одноколової<br>системи: 
                <b>${data.expectedLosses} грн.</b></p>
        `)
    });
});

function showResultSection(htmlContent) {
    let resultSection = document.getElementById("result-section");

    if (!resultSection) {
        resultSection = document.createElement("div");
        resultSection.id = "result-section";
        document.getElementById("app-container").appendChild(resultSection);
    }

    resultSection.innerHTML = htmlContent;
}