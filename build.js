const esbuild = require("esbuild");

esbuild.build({
  entryPoints: ["static/common.js"], // Входной файл
  bundle: true, // Собираем в один файл
  minify: true, // Минифицируем
  outfile: "static/bundle.js", // Выходной файл
	platform: "browser",
}).catch(() => process.exit(1));
