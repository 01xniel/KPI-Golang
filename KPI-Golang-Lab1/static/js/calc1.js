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
            <p class="result-paragraph">Коефіцієнт переходу від робочої до:</p>
            <p class="result-paragraph">  - сухої маси: <b>${data.coefs["dry"].toFixed(2)}</b></p>
            <p class="result-paragraph">  - горючої маси: <b>${data.coefs["combustible"].toFixed(2)}</b></p>
            
            <p class="result-paragraph">Склад сухої маси палива:</p>
            <p class="result-paragraph">  - Водень (H): 
                <b>${data.compositionDry["hydrogen"].toFixed(2)} %</b></p>
            <p class="result-paragraph">  - Вуглець (C): 
                <b>${data.compositionDry["carbon"].toFixed(2)} %</b></p>
            <p class="result-paragraph">  - Сірка (S): 
                <b>${data.compositionDry["sulfur"].toFixed(2)} %</b></p>
            <p class="result-paragraph">  - Азот (N): 
                <b>${data.compositionDry["nitrogen"].toFixed(2)} %</b></p>
            <p class="result-paragraph">  - Кисень (O): 
                <b>${data.compositionDry["oxygen"].toFixed(2)} %</b></p>
            <p class="result-paragraph">  - Зола (A): 
                <b>${data.compositionDry["ash"].toFixed(2)} %</b></p>
            
            <p class="result-paragraph">Склад горючої маси палива:</p>
            <p class="result-paragraph">  - Водень (H): 
                <b>${data.compositionCombustible["hydrogen"].toFixed(2)} %</b></p>
            <p class="result-paragraph">  - Вуглець (C): 
                <b>${data.compositionCombustible["carbon"].toFixed(2)} %</b></p>
            <p class="result-paragraph">  - Сірка (S): 
                <b>${data.compositionCombustible["sulfur"].toFixed(2)} %</b></p>
            <p class="result-paragraph">  - Азот (N): 
                <b>${data.compositionCombustible["nitrogen"].toFixed(2)} %</b></p>
            <p class="result-paragraph">  - Кисень (O): 
                <b>${data.compositionCombustible["oxygen"].toFixed(2)} %</b></p>
                
            <p class="result-paragraph">Нижча теплота згоряння для:</p>
            <p class="result-paragraph">  - робочої маси: 
                <b>${data.lowHeatingValues["raw"].toFixed(2)} МДж/кг.</b></p>
            <p class="result-paragraph">  - сухої маси: 
                <b>${data.lowHeatingValues["dry"].toFixed(2)} МДж/кг.</b></p>
            <p class="result-paragraph">  - горючої маси: 
                <b>${data.lowHeatingValues["combustible"].toFixed(2)} МДж/кг.</b></p>
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