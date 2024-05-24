// Heavily adapted from https://github.com/evilmartians/lefthook/tree/master/packaging/npm-installer

import path from "path";
import fs from "fs";
import { DownloaderHelper } from "node-downloader-helper";
import { fileURLToPath } from "url";
import { createRequire } from "module";
import { createGunzip } from "zlib";

const require = createRequire(import.meta.url);
const packageJson = require("./package.json");

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const iswin = ["win32", "cygwin"].includes(process.platform);

async function install() {
  if (process.env.CI) {
    return;
  }
  const exePath = await downloadBinary();
  if (!iswin) {
    fs.chmodSync(exePath, "755");
  }
}

function getDownloadURL() {
  // Detect OS
  // https://nodejs.org/api/process.html#process_process_platform
  let goOS = process.platform;
  if (iswin) {
    goOS = "windows";
  }

  // Convert the goOS to the os name in the download URL
  let downloadOS = goOS === "darwin" ? "macOS" : goOS;
  downloadOS = `${downloadOS.charAt(0).toUpperCase()}${downloadOS.slice(1)}`;

  // Detect architecture
  // https://nodejs.org/api/process.html#process_process_arch
  let arch = process.arch;
  switch (process.arch) {
    case "x64": {
      arch = "x86_64";
      break;
    }
    case "x32":
    case "ia32": {
      arch = "i386";
      break;
    }
  }
  const version = packageJson.version;

  return `https://github.com/xarunoba/ccoco/releases/download/v${version}/ccoco_${version}_${downloadOS}_${arch}.gz`;
}

async function downloadBinary() {
  const downloadURL = getDownloadURL();
  const fileName = `ccoco${iswin ? ".exe.gz" : ".gz"}`;
  const binDir = path.join(__dirname, "bin");
  const compressedFilePath = path.join(binDir, fileName);
  const decompressedFileName = `ccoco${iswin ? ".exe" : ""}`;
  const decompressedFilePath = path.join(binDir, decompressedFileName);

  const dl = new DownloaderHelper(downloadURL, binDir, {
    fileName,
    retry: { maxRetries: 5, delay: 50 },
  });

  dl.on("end", async () => {
    console.log("ccoco gzip was downloaded");
    try {
      await decompressFile(compressedFilePath, decompressedFilePath);
      fs.unlinkSync(compressedFilePath); // remove compressed file
      console.log("ccoco binary was decompressed");
    } catch (e) {
      console.error(`Failed to decompress ${fileName}: ${e.message}`);
      throw new Error(`Failed to decompress ${fileName}: ${e.message}`);
    }
  });

  try {
    await dl.start();
  } catch (e) {
    const message = `Failed to download ${fileName}: ${e.message} while fetching ${downloadURL}`;
    console.error(message);
    throw new Error(message);
  }
  return decompressedFilePath;
}

async function decompressFile(sourcePath, destinationPath) {
  return new Promise((resolve, reject) => {
    const gunzip = createGunzip();
    const source = fs.createReadStream(sourcePath);
    const destination = fs.createWriteStream(destinationPath);

    source
      .pipe(gunzip)
      .pipe(destination)
      .on("finish", resolve)
      .on("error", reject);
  });
}

// start:
try {
  await install();
} catch (e) {
  throw e;
}
