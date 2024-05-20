#!/usr/bin/env node

import { resolve } from "node:path";
import fs from "node:fs/promises";
import { pathToFileURL } from "url";
import { argv, exit } from "node:process";
import { promisify } from "node:util";
import { exec } from "node:child_process";
import yargs from "yargs/yargs";
import { hideBin } from "yargs/helpers";

const execAsync = promisify(exec);

const flagsDescription = {
  quiet: "Will not log anything to the console",
};

const flagsAlias = {
  v: "version",
  h: "help",
  q: "quiet",
};

const flags = yargs(hideBin(argv))
  .strict()
  .alias(flagsAlias)
  .options({
    quiet: {
      describe: flagsDescription.quiet,
      type: "boolean",
    },
  })
  .help().argv;

/**
 * Check the enabled flags and log them to the console.
 */
const checkFlags = () => {
  let enabledFlags = Object.keys(flagsDescription).filter(
    (flag) => flags[flag],
  );
  if (enabledFlags.length) {
    if (!flags.q) console.log("◼ Enabled flag/s:");
    enabledFlags.forEach((flag) => {
      if (!flags.q) console.log(`    ❯ ${flag}`);
    });
  }
};

/**
 * Check if git is installed.
 */
const checkGitVersion = async () => {
  try {
    await execAsync("git --version");
    if (!flags.q) console.log("✔ Git is installed!");
  } catch (err) {
    console.error("✘ Git is not installed!");
    exit(1);
  }
};

/**
 * Check if the current directory is a git repository and if HEAD exists.
 */
const checkGitRepo = async () => {
  try {
    await execAsync("git rev-parse --is-inside-work-tree");
    if (!flags.q) console.log("✔ Current directory is a git repository!");
  } catch (err) {
    console.error("✘ Not a git repository!");
    exit(1);
  }
};

/**
 * Set config default values.
 */
const setConfigDefaults = (config) => {
  if (!config.branchConfigDir) config.branchConfigDir = ".ccoco";
  if (!config.files || !config.files.length) config.files = [".env"];
  if (!config.mainConfigDir) config.mainConfigDir = ".";
  return config;
};

/**
 * The generated default config file for ccoco.
 */
const defaultConfigFile = `export default {
  // branchConfigDir: ".ccoco", // the directory where the config files per branch are located
  // files: [".env"], // the main config files to be replaced
  // mainConfigDir: ".", // the directory where the main config files are located
};`;

/**
 *  Load the config file.
 */
const loadConfig = async () => {
  try {
    const configPath = resolve(process.cwd(), "ccoco.config.js");
    const configUrl = pathToFileURL(configPath).href;
    try {
      const { default: config } = await import(configUrl);
      return setConfigDefaults(config);
    } catch (error) {
      console.warn(
        `⚠ Config file not found. Generating a default config file.`,
      );
      try {
        await fs.writeFile(configPath, defaultConfigFile);
        const { default: config } = await import(configUrl);
        return setConfigDefaults(config);
      } catch (error) {
        console.error(`Error creating config file: ${error}`);
        exit(1);
      }
    }
  } catch (error) {
    console.error(`Error loading config: ${error}`);
    exit(1);
  }
};

/**
 * Replace the config files.
 */
const replaceConfigFiles = async (config) => {
  await Promise.all(
    config.files.map(async (file) => {
      try {
        // Check if file is string
        if (typeof file !== "string") {
          console.warn(`⚠ ${file} is not a file string. Skipped.`);
          return;
        }

        // Resolve main config file path
        let mainConfigFile = `./${file}`;
        if (typeof config.mainConfigDir === "string") {
          mainConfigFile = `${config.mainConfigDir}/${file}`;
        }

        // Get branch name
        const { stdout } = await execAsync("git branch --show-current");
        const branch = stdout.trim();

        // Resolve branch config file path
        let branchFile = "";
        if (config.branchConfigDir) {
          branchFile = `./${config.branchConfigDir}/${file}.${branch}`;
        } else branchFile = `${file}.${branch}`;

        // Check if branch config file exists
        try {
          await fs.access(branchFile);
          if (!flags.q)
            console.log(`✔ ${file} ❯ Found branch config file: ${branchFile}`);
        } catch (error) {
          if (!flags.q)
            console.log(
              `✘ ${file} ❯ Could not find config file: ${branchFile}. Skipped.`,
            );
          return;
        }

        // Check if main config file exists and delete it
        try {
          await fs.access(mainConfigFile);
          await fs.rm(mainConfigFile);
        } catch (error) {
          if (!flags.q)
            console.warn(
              `⚠ ${file} ❯ Could not find config file: ${mainConfigFile}. Skipped deletion.`,
            );
        }

        // Check if mainConfigDir exists and create it if not
        try {
          await fs.access(config.mainConfigDir);
        } catch (error) {
          if (!flags.q)
            console.warn(
              `⚠ ${file} ❯ Cannot find main config directory: ${config.mainConfigDir}`,
            );
          try {
            await fs.mkdir(config.mainConfigDir, { recursive: true });
            if (!flags.q)
              console.log(
                `✔ ${file} ❯ Created main config directory: ${config.mainConfigDir}`,
              );
          } catch (error) {
            console.error(`✘ ${file} ❯ ${error}`);
            exit(1);
          }
        }

        // Copy branch config file to main config file
        await fs.copyFile(branchFile, mainConfigFile);
        if (!flags.q)
          console.log(
            `✔ ${file} ❯ Successfully replaced config file: ${branchFile} -> ${mainConfigFile}`,
          );
      } catch (error) {
        console.error(`An error occurred when handling config file: ${error}`);
        exit(1);
      }
    }),
  );
};

/**
 * Main function to execute the ccoco.
 */
const main = async () => {
  const startTime = Date.now();
  if (!flags.q) console.log("❯ Preparing ccoco...");

  checkFlags();

  // Check if git is installed
  await checkGitVersion();

  // Check if the current directory is a git repository and if HEAD exists
  await checkGitRepo();

  // Load config
  const config = await loadConfig();

  await replaceConfigFiles(config);

  if (!flags.q) console.log("✔ Finished ccoco!");
  const endTime = Date.now();
  const elapsedTime = (endTime - startTime) / 1000;
  if (!flags.q) console.log(`Total time elapsed: ${elapsedTime} seconds`);
  exit(0);
};

main();
