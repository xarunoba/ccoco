# @xarunoba/ccoco 🥥

![Static Badge](https://img.shields.io/badge/Made_with-%E2%9D%A4%EF%B8%8F-red?style=for-the-badge) ![NPM License](https://img.shields.io/npm/l/%40xarunoba%2Fccoco?style=for-the-badge)
![GitHub package.json version](https://img.shields.io/github/package-json/v/xarunoba/ccoco?style=for-the-badge&logo=npm)

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

## Usage

### npm

```bash
npx @xarunoba/ccoco
```

### pnpm

```bash
pnpm dlx @xarunoba/ccoco
```

### yarn

```bash
yarn dlx @xarunoba/ccoco
```

### flags

You can add the following flags for `ccoco`:

- `-v, --version` — Show version number
- `-h, --help` — Show help
- `-q, --quiet` — Will not log anything to the console

### sub-branches

`ccoco` will recursively check if a sub-branch has a config file until it reaches the "root" of the sub-branch.

### Example

1. Branch `nested/one/two` does not have a config file created.
2. `ccoco` will recursively check for the config file existing in `nested/one` up until the root `nested` and will fail if it cannot find one.

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
