"use strict";

// Fixed synthetic illustration. No API calls, persistence, or evaluation output.
const records = {
  en: [
    ["Define the work", "A fictional literature review begins with a goal, a language constraint, and a deadline.", "November 12 (initial instruction).", "Collect primary sources.", "Expected behavior: preserve the goal and constraint at the first handoff."],
    ["Update the context", "During the local model B stage, the fictional user corrects the deadline. The goal and language remain unchanged.", "November 15 (corrected from November 12).", "Compare the collected sources.", "Expected behavior: retain the correction and its origin, rather than treating both dates as current."],
    ["Resume with the correction", "The task returns to local model A. This record illustrates the expected context, not an observed model response.", "November 15 (latest instruction retained).", "Draft the comparison in Japanese.", "Evaluation question: does the model use the corrected deadline and continue the task without inventing missing facts?"]
  ],
  ja: [
    ["作業を定義する", "架空の文献調査を、目的・使用言語・締切とともに開始します。", "11月12日（最初の指示）。", "一次資料を集める。", "期待する動作：最初の引き継ぎで目的と制約を保持する。"],
    ["文脈を更新する", "ローカルモデルBの段階で、架空のユーザーが締切を訂正します。目的と使用言語は変わりません。", "11月15日（12日から訂正）。", "収集した資料を比較する。", "期待する動作：訂正とその出典を残し、二つの日付を同時に有効な締切として扱わない。"],
    ["訂正を反映して再開する", "作業をローカルモデルAに戻します。ここに示すのは期待する文脈であり、実際のモデルの出力ではありません。", "11月15日（最新の指示を保持）。", "日本語で比較の草稿を作る。", "評価の問い：訂正後の締切を使い、足りない情報を捏造せずに作業を続けられるか。"]
  ]
};
const language = document.documentElement.lang === "ja" ? "ja" : "en";
const fields = ["step-title", "step-description", "demo-deadline", "demo-next", "demo-note"];
const buttons = document.querySelectorAll("[data-step]");
for (const button of buttons) {
  button.addEventListener("click", () => {
    const step = Number(button.dataset.step);
    fields.forEach((id, index) => {
      document.getElementById(id).textContent = records[language][step][index];
    });
    document.getElementById("step-count").textContent = `0${step + 1} / 03`;
    for (const candidate of buttons) {
      candidate.setAttribute("aria-pressed", String(candidate === button));
    }
  });
}
