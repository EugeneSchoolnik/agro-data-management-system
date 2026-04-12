import type { FieldReport } from "../types/models";
import html2pdf from "html2pdf.js";

const escapeHtml = (v: string) =>
  v.replace(
    /[&<>"']/g,
    (m) =>
      ({
        "&": "&amp;",
        "<": "&lt;",
        ">": "&gt;",
        '"': "&quot;",
        "'": "&#39;",
      })[m] || m,
  );

export async function downloadFieldReportPDF(
  report: FieldReport,
  fieldName: string,
  fromDate: string,
  toDate: string,
): Promise<void> {
  // 1. Створюємо контейнер БЕЗ fixed/absolute
  const container = document.createElement("div");
  container.id = "final-pdf-report";

  // Робимо його білим, видимим, але дуже тонким або прихованим через overflow
  Object.assign(container.style, {
    width: "794px",
    background: "white",
    color: "black",
    margin: "0",
    padding: "0",
    // Замість ховання координатами, просто робимо його "невидимим" у потоці
    height: "auto",
    overflow: "hidden",
  });

  container.innerHTML = `
    <div style="padding: 48px; font-family: Arial, Helvetica, sans-serif; color: #111111; background: white;">
      <div style="border-bottom: 1px solid #222222; padding-bottom: 10px; margin-bottom: 24px;">
        <h1 style="font-size: 24px; font-weight: 700; margin: 0;">Звіт по полю</h1>
        <p style="font-size: 14px; margin: 8px 0 0 0;">${escapeHtml(fieldName)}</p>
      </div>
      <div style="display: flex; justify-content: space-between; gap: 24px; margin-bottom: 24px;">
        <div style="min-width: 260px;">
          <p style="margin: 4px 0; font-size: 13px;"><strong>Період:</strong> ${escapeHtml(fromDate)} — ${escapeHtml(toDate)}</p>
          <p style="margin: 4px 0; font-size: 13px;"><strong>Згенеровано:</strong> ${escapeHtml(new Date().toLocaleString("uk-UA"))}</p>
        </div>
        <div style="min-width: 200px; text-align: right;">
          <p style="margin: 4px 0; font-size: 12px; color: #555555;">Прогноз шкідників</p>
          <p style="margin: 4px 0; font-size: 28px; font-weight: 700;">${Math.round(report.forecast_average_probability * 100)}%</p>
        </div>
      </div>
      <div style="margin-bottom: 16px;">
        <h2 style="font-size: 16px; font-weight: 700; letter-spacing: 0.02em; margin: 0 0 12px 0;">Оцінка параметрів</h2>
        <table style="width: 100%; border-collapse: collapse; font-size: 13px;">
          <thead>
            <tr>
              <th style="text-align: left; padding: 10px 8px; border-bottom: 1px solid #cccccc;">Показник</th>
              <th style="text-align: right; padding: 10px 8px; border-bottom: 1px solid #cccccc;">Середнє</th>
              <th style="text-align: right; padding: 10px 8px; border-bottom: 1px solid #cccccc;">Мінімум</th>
              <th style="text-align: right; padding: 10px 8px; border-bottom: 1px solid #cccccc;">Максимум</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td style="padding: 10px 8px; border-bottom: 1px solid #ededed;">Температура, °C</td>
              <td style="padding: 10px 8px; text-align: right; border-bottom: 1px solid #ededed;">${report.temperature.avg.toFixed(1)}</td>
              <td style="padding: 10px 8px; text-align: right; border-bottom: 1px solid #ededed;">${report.temperature.min.toFixed(1)}</td>
              <td style="padding: 10px 8px; text-align: right; border-bottom: 1px solid #ededed;">${report.temperature.max.toFixed(1)}</td>
            </tr>
            <tr>
              <td style="padding: 10px 8px; border-bottom: 1px solid #ededed;">Вологість повітря, %</td>
              <td style="padding: 10px 8px; text-align: right; border-bottom: 1px solid #ededed;">${report.air_humidity.avg.toFixed(1)}</td>
              <td style="padding: 10px 8px; text-align: right; border-bottom: 1px solid #ededed;">${report.air_humidity.min.toFixed(1)}</td>
              <td style="padding: 10px 8px; text-align: right; border-bottom: 1px solid #ededed;">${report.air_humidity.max.toFixed(1)}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div style="margin-top: 32px; display: flex; justify-content: space-between; gap: 24px;">
        <div style="min-width: 200px;">
          <p style="margin: 4px 0; font-size: 12px; color: #555555;"><strong>Підсумок</strong></p>
          <p style="margin: 4px 0; font-size: 14px;">Прогноз шкідників: <strong>${Math.round(report.forecast_average_probability * 100)}%</strong></p>
        </div>
      </div>
    </div>
  `;

  // 2. Додаємо в body (він з'явиться в самому низу сторінки на мілісекунду)
  document.body.appendChild(container);

  const safeFieldName =
    fieldName.replace(/[^a-zA-Z0-9_-]+/g, "_").slice(0, 40) || "report";
  const opt = {
    margin: 0,
    filename: `report_${safeFieldName}_${fromDate}_${toDate}.pdf`,
    image: { type: "jpeg", quality: 1 },
    html2canvas: {
      scale: 2,
      useCORS: true,
      // ВАЖЛИВО: Примусово кажемо canvas малювати саме цей елемент
      onclone: (clonedDoc: Document) => {
        // Можна переконатися, що клонований елемент видимий
        const el = clonedDoc.getElementById("final-pdf-report");
        if (el) el.style.opacity = "1";
      },
    },
    jsPDF: { unit: "mm", format: "a4", orientation: "portrait" },
  };

  try {
    // 3. Використовуємо worker-інтерфейс для кращого контролю
    const worker = html2pdf()
      .set(opt as any)
      .from(container);

    // Запускаємо рендер
    await worker.save();
  } catch (error) {
    console.error("PDF Fatal Error:", error);
  } finally {
    // 4. Видаляємо елемент
    if (document.body.contains(container)) {
      document.body.removeChild(container);
    }
  }
}
