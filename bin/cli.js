#!/usr/bin/env node
import { resolve, posix } from "node:path";
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
 * Check cache to skip git checks
 */
const checkCache = async () => {
  try {
    await fs.access(process.cwd() + "/.ccoco/.cache");
    console.log("✔ .cache file exists. Skipping git checks.");
    return true;
  } catch (error) {
    console.log("✘ .cache file does not exist. Git checks will be performed.");
    return false;
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
  if (!config.branchConfigDir) config.branchConfigDir = ".ccoco/configs";
  if (!config.files || !config.files.length) config.files = [".env"];
  if (!config.mainConfigDir) config.mainConfigDir = ".";
  return config;
};

/**
 * The generated default config file for ccoco.
 */
const defaultConfigFile = `export default {
  // branchConfigDir: ".ccoco/configs", // the directory where the config files per branch are located
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
      if (!flags.q)
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
          if (!flags.q)
            console.warn(`⚠ ${file} is not a file string. Skipped.`);
          return;
        }

        // Get main config directory and file
        const mainFile = posix.resolve(file);

        // Get branch name
        const { stdout } = await execAsync("git branch --show-current");
        const branch = stdout.trim();

        // Get branch config directory and file
        let branchFileDir = posix.resolve(config.branchConfigDir, file);
        let branchFile = posix.resolve(branchFileDir, branch);

        // Check if branch is a subranch
        const isSubBranch = branch.split("/").length > 1;
        if (isSubBranch && !flags.q)
          console.log(`✔ ${file} ❯ ${branch} ❯ Sub-branch detected!`);
        let subBranchFallbackFile = branchFile
          .split("/")
          .slice(0, -1)
          .join("/");

        // Check if file directory exists
        try {
          await fs.access(branchFileDir);
        } catch (error) {
          if (!flags.q)
            console.warn(
              `⚠ ${file} ❯ Cannot find branch config directory. Attempting to create.`,
            );
          try {
            await fs.mkdir(branchFileDir, { recursive: true });
            if (!flags.q) console.log(`✔ ${file} ❯ Branch config created!`);
          } catch (error) {
            console.error(`✘ ${file} ❯ ${error}`);
            exit(1);
          }
        }

        // Check if config file for branch exists
        try {
          await fs.access(branchFile);
        } catch (error) {
          if (!flags.q && error.code === "ENOENT") {
            if (isSubBranch) {
              if (!flags.q)
                console.warn(
                  `⚠ ${file} ❯ ${branch} ❯ Cannot find sub-branch config file.\n⚠ ${file} ❯ ${branch} ❯ Checking for fallback file...`,
                );
              let HasFallbackFile = false;
              while (subBranchFallbackFile !== branchFileDir) {
                try {
                  await fs.access(subBranchFallbackFile);
                  branchFile = subBranchFallbackFile;
                  HasFallbackFile = true;
                  break;
                } catch (error) {
                  if (error.code === "ENOENT") {
                    subBranchFallbackFile = subBranchFallbackFile
                      .split("/")
                      .slice(0, -1)
                      .join("/");
                  } else {
                    console.error(`✘ ${file} ❯ ${branch} ❯ ${error}`);
                    exit(1);
                  }
                }
              }

              if (HasFallbackFile) {
                if (!flags.q)
                  console.log(
                    `✔ ${file} ❯ ${branch} ❯ Found fallback config!`,
                  );
                return;
              } else {
                if (!flags.q)
                  console.warn(
                    `⚠ ${file} ❯ ${branch} ❯ Cannot find fallback config file. Skipped.\n⚠ Please create this file \"${branchFile}\"`,
                  );
                return;
              }
            } else {
              if (!flags.q)
                console.warn(
                  `⚠ ${file} ❯ ${branch} ❯ Cannot find branch config file. Skipped.\n⚠ Please create this file \"${branchFile}\"`,
                );
              return;
            }
          } else console.error(`✘ ${file} ❯ ${branch} ❯ ${error}`);
        }

        // Remove existing main config file
        try {
          await fs.rm(mainFile);
        } catch (error) {
          if (!flags.q && error.code === "ENOENT")
            console.log(
              `❯ ${file} ❯ Main config file not found. Skipped removal.`,
            );
          else console.error(`✘ ${file} ❯ ${error}`);
        }

        // Copy config file
        console.log(mainFile);
        await fs.copyFile(branchFile, mainFile);
        if (!flags.q) console.log(`✔ ${file} ❯ ${branch} ❯ Replaced config!`);
      } catch (error) {
        console.error(`✘ ${file} ❯ ${error}`);
        exit(1);
      }
    }),
  );
};

/**
 * The main function to execute the ccoco.
 */
const main = async () => {
  const startTime = Date.now();
  if (!flags.q) console.log("❯ Preparing ccoco...");

  checkFlags();

  // Check if .cache file exists
  const hasCache = await checkCache();
  if (!hasCache) {
    // Check if git is installed
    await checkGitVersion();

    // Check if the current directory is a git repository and if HEAD exists
    await checkGitRepo();

    // Add .cache file
    await fs.writeFile(process.cwd() + "/.ccoco/.cache", `keep me`);
  }

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
