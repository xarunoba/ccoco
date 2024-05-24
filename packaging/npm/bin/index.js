#!/usr/bin/env node

// Heavily adapted from https://github.com/evilmartians/lefthook/tree/master/packaging/npm-installer

import { spawn } from "child_process";
import path from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const extension = ["win32", "cygwin"].includes(process.platform) ? ".exe" : "";
const exePath = path.join(__dirname, `ccoco${extension}`);

const command_args = process.argv.slice(2);
const child = spawn(exePath, command_args, { stdio: "inherit" });

child.on("close", (code) => {
  if (code !== 0) {
    process.exit(1);
  }
});
