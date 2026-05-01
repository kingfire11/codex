// designapi.ink — instructions page
const BASE_URL = "https://api.designapi.ink";
const PLACEHOLDER_TOKEN = "YOUR_API_KEY";
const INSTALLER_REPO = "kingfire11/codex";
const INSTALLER_BASE = `https://github.com/${INSTALLER_REPO}/releases/latest/download`;

const ASSETS = [
  "codex-cli.macos.sh", "codex-cli.linux.sh", "codex-cli.windows.ps1",
  "codex-vscode.macos.sh", "codex-vscode.linux.sh", "codex-vscode.windows.ps1",
  "codex-app.macos.sh", "codex-app.windows.ps1",
  "opencode.macos.sh", "opencode.linux.sh", "opencode.windows.ps1",
  "opencode.config.json",
];

const templates = {};
const tokenInput = document.getElementById("token");
const modelSelect = document.getElementById("model");

function currentToken() {
  const v = tokenInput.value.trim();
  return v || PLACEHOLDER_TOKEN;
}
function currentModel() {
  return modelSelect.value || "gpt-5.5";
}

function fill(text) {
  return text
    .replaceAll("__API_KEY__", currentToken())
    .replaceAll("__BASE_URL__", BASE_URL)
    .replaceAll("__MODEL__", currentModel());
}

function oneliner(name) {
  const isPs = name.endsWith(".ps1");
  const url = `${location.origin}${location.pathname.replace(/\/[^/]*$/, "")}/scripts/${name}`;
  if (isPs) {
    return `$env:OPENAI_API_KEY="${currentToken()}"; iwr -useb ${url} | iex`;
  }
  return `OPENAI_API_KEY="${currentToken()}" bash <(curl -fsSL ${url})`;
}

function render() {
  // полные скрипты
  document.querySelectorAll("pre[data-script]").forEach(pre => {
    const tpl = templates[pre.dataset.script];
    if (tpl != null) pre.textContent = fill(tpl);
  });
  // однострочники
  document.querySelectorAll("pre[data-oneliner]").forEach(pre => {
    pre.textContent = oneliner(pre.dataset.oneliner);
  });
  // installer run-команды (macOS / Linux / Windows)
  const tok = currentToken();
  document.querySelectorAll("pre[data-oneliner-mac]").forEach(pre => {
    const arch = (navigator.userAgent || "").includes("ARM") || /arm/i.test(navigator.platform || "")
      ? "arm64" : "amd64";
    pre.textContent =
`cd ~/Downloads
xattr -d com.apple.quarantine ./designapi-installer-darwin-${arch} 2>/dev/null || true
chmod +x ./designapi-installer-darwin-${arch}
./designapi-installer-darwin-${arch} doctor --token=${tok}`;
  });
  document.querySelectorAll("pre[data-oneliner-linux]").forEach(pre => {
    pre.textContent =
`cd ~/Downloads
chmod +x ./designapi-installer-linux-amd64
./designapi-installer-linux-amd64 doctor --token=${tok}`;
  });
  document.querySelectorAll("pre[data-oneliner-win]").forEach(pre => {
    pre.textContent =
`cd $HOME\\Downloads
Unblock-File .\\designapi-installer-windows-amd64.exe
.\\designapi-installer-windows-amd64.exe doctor --token=${tok}`;
  });
  // статичные блоки — оставляем как есть
  document.querySelectorAll("pre[data-static]").forEach(pre => {
    if (!pre.textContent.trim()) pre.textContent = pre.dataset.static;
  });
}

async function loadTemplates() {
  await Promise.all(ASSETS.map(async name => {
    try {
      const res = await fetch(`scripts/${name}`);
      if (res.ok) templates[name] = await res.text();
    } catch {}
  }));
  render();
}

// Tabs
function initTabs() {
  document.querySelectorAll("nav.tabs").forEach(nav => {
    const buttons = nav.querySelectorAll("button");
    const container = nav.parentElement;
    function activate(name) {
      buttons.forEach(b => b.classList.toggle("active", b.dataset.tab === name));
      // Только прямые потомки контейнера, а не вложенные подтабы
      Array.from(container.children).forEach(child => {
        if (child.classList && child.classList.contains("tab-panel")) {
          child.style.display = (child.dataset.panel === name) ? "" : "none";
        }
      });
    }
    buttons.forEach(b => b.addEventListener("click", () => activate(b.dataset.tab)));
    activate(buttons[0].dataset.tab);
  });
}

// Copy / Download buttons (event delegation)
document.addEventListener("click", async e => {
  const t = e.target;
  if (t.classList.contains("copy")) {
    const pre = t.parentElement.querySelector("pre");
    try {
      await navigator.clipboard.writeText(pre.textContent);
      t.textContent = "✓"; setTimeout(() => t.textContent = "copy", 1200);
    } catch { t.textContent = "ошибка"; }
  }
  if (t.classList.contains("dl")) {
    const pre = t.parentElement.querySelector("pre");
    const blob = new Blob([pre.textContent], { type: "text/plain;charset=utf-8" });
    const a = document.createElement("a");
    a.href = URL.createObjectURL(blob);
    a.download = t.dataset.name || "designapi.txt";
    document.body.appendChild(a); a.click(); a.remove();
    setTimeout(() => URL.revokeObjectURL(a.href), 2000);
  }
});

// Show / hide token
document.getElementById("toggle-token").addEventListener("click", () => {
  tokenInput.type = tokenInput.type === "password" ? "text" : "password";
});

// Live re-render
tokenInput.addEventListener("input", render);
modelSelect.addEventListener("change", render);

// Installer download links + autodetect
function setInstallerLinks() {
  const map = {
    "dl-mac-arm":   "designapi-installer-darwin-arm64",
    "dl-mac-x64":   "designapi-installer-darwin-amd64",
    "dl-linux-x64": "designapi-installer-linux-amd64",
    "dl-linux-arm": "designapi-installer-linux-arm64",
    "dl-win-x64":   "designapi-installer-windows-amd64.exe",
  };
  Object.entries(map).forEach(([id, file]) => {
    const a = document.getElementById(id);
    if (a) a.href = `${INSTALLER_BASE}/${file}`;
  });

  // Подсветим кнопку под текущую систему
  const ua = (navigator.userAgent || "").toLowerCase();
  let id = null;
  if (ua.includes("mac"))     id = navigator.userAgent.includes("ARM") ? "dl-mac-arm" : "dl-mac-arm";
  if (ua.includes("windows")) id = "dl-win-x64";
  if (ua.includes("linux") && !ua.includes("android")) id = "dl-linux-x64";
  if (id) document.getElementById(id)?.classList.add("primary");
}

initTabs();
setInstallerLinks();
loadTemplates();
