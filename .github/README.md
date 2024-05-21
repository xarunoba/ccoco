# @xarunoba/ccoco 🥥

![Static Badge](https://img.shields.io/badge/Made_with-%E2%9D%A4%EF%B8%8F-red?style=for-the-badge) ![NPM License](https://img.shields.io/npm/l/%40xarunoba%2Fccoco?style=for-the-badge)
![NPM Version](https://img.shields.io/npm/v/%40xarunoba%2Fccoco?style=for-the-badge&logo=npm)

**`ccoco`** — change config on checkout

Change your configurations on checkout for easy environment change. Integrate with git hooks.

## Why

My team didn't want to populate the `package.json` too much with scripts so I created a tool that can easily replace the main config files after every checkout. `ccoco` fixes this issue by automatically changing your main config files per branch name.

## Installation

### npm

```bash
npm install -D @xarunoba/ccoco
# run locally
npx ccoco
```

### pnpm

```bash
pnpm install -D @xarunoba/ccoco
# run locally
pnpm ccoco
```

### yarn

```bash
yarn add --dev @xarunoba/ccoco
# run locally
yarn run ccoco
```

### .gitignore

Please add `.ccoco` in your `.gitignore` file.

## Usage

## Run

```bash
# Runnign locally
npx ccoco
pnpm ccoco
yarn ccoco

# You can use the shorter alias `cc` when running locally
npx cc
pnpm cc
yarn cc

# Downloading from registry and executing it
npx @xarunoba/ccoco
pnpm dlx @xarunoba/ccoco
yarn dlx @xarunoba/ccoco
```

### Example usage

Create a directory that has the same file name as your config file in `.ccoco/configs`

```
.env -> .ccoco/configs/.env (as a directory not a file)
```

Create the file that is the same as your branch name and put your values inside.

```
# branch to file
main -> .ccoco/configs/.env/main

# .ccoco/configs/.env/main
value=1
```

Running `npx ccoco` or `npx cc` (for short) will copy the `.ccoco/configs/.env/main` file's value and overwrite the root `.env` file when the current branch is `main` IF and only if that `.ccoco/configs/.env/main` file exists otherwise it is skipped. This process can be automated using [git hooks](#with-simple-git-hooks).

```
# .env
<blank value>

# .ccoco/configs/.env/main
value=1

# execute ccoco
npx ccoco

# .env
value=1
```

### Using sub-branches

`ccoco` will recursively check if a sub-branch has a config file until it reaches the "root" of the sub-branch.

##### Example

1. Branch `nested/one/two` does not have a config file created.
2. `ccoco` will recursively check for the config file existing in `nested/one` up until the root `nested` and will fail if it cannot find one.

### Flags

You can add the following flags for `ccoco`:

- `-v, --version` — Show version number
- `-h, --help` — Show help
- `-q, --quiet` — Will not log anything to the console

## Configuring `ccoco`

You can configure `ccoco` by editing or creating the `ccoco.config.js` file. (This file will be automatically generated when running `ccoco`)

### Example

```javascript
export default {
  branchConfigDir: '.ccoco/env', // the directory where the config files per branch are located. defaults to ".ccoco/configs"
  files: ['.env', '.env.production', '.env.somemode'], // the main config files to be replaced, you can add more if you want. defaults to [".env"]
  mainConfigDir: '.', // the directory where the main config files are located. defaults to the root project directory, "."
}
```

## Integrations

### With [`simple-git-hooks`](https://github.com/toplenboren/simple-git-hooks)

Integrating with [`simple-git-hooks`](https://github.com/toplenboren/simple-git-hooks) is easy as a toasted bread:

```json
// package.json
{
  ...
  "simple-git-hooks": {
    // If you prefer on always using the latest version of ccoco
    // instead of using ccoco as a devDependency, use this:
    // "post-checkout": "npx @xarunoba/ccoco"
    "post-checkout": "npx ccoco"
  }
  ...
}
```
