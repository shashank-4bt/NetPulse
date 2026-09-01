const { copyFileSync, mkdirSync } = require("node:fs");
const path = require("node:path");

const dist = path.join(path.dirname(require.resolve("maplibre-gl/package.json")), "dist");
const dest = path.join(__dirname, "..", "public", "maplibre");

mkdirSync(dest, { recursive: true });
for (const file of ["maplibre-gl-worker.mjs", "maplibre-gl-shared.mjs"]) {
  copyFileSync(path.join(dist, file), path.join(dest, file));
}
