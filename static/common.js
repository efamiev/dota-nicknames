const htmx = require("htmx.org"); // Принудительно импортируем CommonJS
window.htmx = htmx.default; // Делаем htmx глобальным

require("htmx-ext-sse"); // Подключаем расширение
