document.addEventListener("DOMContentLoaded", function () {
    const form = document.getElementById("calculator-form");

    form.addEventListener("submit", async function (event) {
        event.preventDefault();

        const formData = new FormData(form);

        const response = await fetch("/evaluate", {
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
            <p class="result-paragraph"><b>${data.Profit.toFixed(1)} тис. грн.</b> - дохід від генерації 
                енергії без небалансів <b>(${data.ElectricityNoImbalance.toFixed(1)} МВт⋅год)</b></p>
            <p class="result-paragraph"><b>${data.Penalty.toFixed(1)} тис. грн.</b> - штраф за генерацію 
                енергії з небалансами <b>(${data.ElectricityImbalance.toFixed(1)} МВт⋅год)</b></p>
            <p class="result-paragraph">${data.NetProfit >= 0 ? "Прибуток" : "Збитки"}: 
                <b>${Math.abs(data.NetProfit.toFixed(1))} тис. грн.</b></p>
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
