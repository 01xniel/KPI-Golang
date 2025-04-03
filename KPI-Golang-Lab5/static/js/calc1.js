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

        const response = await fetch("/evaluate1", {
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
            <p class="result-paragraph">Частота відмов одноколової системи: 
                <b>${data.failureRateSCS.toFixed(3)} рік⁻¹</b></p>
            <p class="result-paragraph">Середня тривалість відновлення одноколової<br>системи: 
                <b>${data.averageRecoveryTime.toFixed(3)} год.</b></p>
            <p class="result-paragraph">Коефіцієнт аварійного простою одноколової<br>системи: 
                <b>${data.coefEmergencyDowntimeSCS.toExponential(4)}</b></p>
            <p class="result-paragraph">Коефіцієнт планового простою одноколової<br>системи: 
                <b>${data.coefScheduledDowntimeSCS.toExponential(3)}</b></p>
            <p class="result-paragraph">Частота відмов двоколової системи (не враховуючи секційний вимикач): 
                <b>${data.failureRateTCS.toExponential(3)} рік⁻¹</b></p>
            <p class="result-paragraph">Частота відмов двоколової системи (враховуючи секційний вимикач): 
                <b>${data.failureRateWithSectionalizerTCS.toExponential(2)} рік⁻¹</b></p>
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