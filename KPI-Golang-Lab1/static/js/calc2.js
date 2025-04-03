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
            <p class="result-paragraph">Склад сухої маси палива:</p>
            <p class="result-paragraph">  - Вуглець (C): 
                <b>${data.rawComposition["carbon"].toFixed(2)} %</b></p>
            <p class="result-paragraph">  - Сірка (S): 
                <b>${data.rawComposition["sulfur"].toFixed(2)} %</b></p>
            <p class="result-paragraph">  - Водень (H): 
                <b>${data.rawComposition["hydrogen"].toFixed(2)} %</b></p>
            <p class="result-paragraph">  - Кисень (O): 
                <b>${data.rawComposition["oxygen"].toFixed(2)} %</b></p>
            <p class="result-paragraph">  - Волога (W): 
                <b>${data.rawComposition["moisture"].toFixed(2)} %</b></p>
            <p class="result-paragraph">  - Зола (A): 
                <b>${data.rawComposition["ash"].toFixed(2)} %</b></p>
            <p class="result-paragraph">  - Ванадій (V): 
                <b>${data.rawComposition["vanadium"].toFixed(2)} мг/кг.</b></p>
            <p class="result-paragraph">Нижча теплота згоряння робочої маси<br>мазуту: 
                <b>${data.rawLHV.toFixed(2)} МДж/кг.</b></p>
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