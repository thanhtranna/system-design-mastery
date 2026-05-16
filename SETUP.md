# Running this course with mdBook

> **Warning:** Before you start, make sure this repo is on `main` and your local `sd-mastery/` folder contains the correct `book.toml`, `src/`, and `.github/workflows/deploy.yml` files. If you change the published branch or remove the workflow, GitHub Pages deployment can fail.

Three ways to read this course:

1. **GitHub Pages** — push to GitHub, auto-deploy a public/private site (no install)
2. **Local** — `mdbook serve` on your machine
3. **Just files** — open `.md` files directly in VS Code or Obsidian

---

## Option 1: Deploy to GitHub Pages (recommended — no install)

The GitHub Actions workflow is already included. You only need to push and enable Pages.

### Step 1: Create a GitHub repo

On github.com, create a new repo (e.g. `system-design-mastery`). **Public or private — both work.**

### Step 2: Push the `sd-mastery/` folder

From your local `sd-mastery/` folder:

```bash
cd sd-mastery

git init
git add .
git commit -m "Initial commit: system design course mastery"
git branch -M main
git remote add origin https://github.com/YOUR_USERNAME/system-design-mastery.git
git push -u origin main
```

### Step 3: Enable Pages

On the GitHub repo page:

1. **Settings** → **Pages** (left sidebar)
2. Under "Build and deployment", **Source**: select **GitHub Actions**
3. Done.

### Step 4: Watch it deploy

1. Go to the **Actions** tab
2. You'll see "Deploy mdBook to GitHub Pages" running (~1-2 min)
3. When green, site is live at:

   ```bash
   https://YOUR_USERNAME.github.io/system-design-mastery/
   ```

### Update workflow

Every `git push` to `main` auto-rebuilds. No manual steps after setup.

```bash
vim src/modules/01-thinking-in-systems.md
git add . && git commit -m "Add notes"
git push
# ~2 min later: site updates
```

### Custom domain (optional)

Settings → Pages → Custom domain. Point a CNAME at `YOUR_USERNAME.github.io`.

---

## Option 2: Run Locally

### Install mdBook + Mermaid plugin

The course uses Mermaid diagrams (architecture, sequences, state machines), so you need both `mdbook` and `mdbook-mermaid`.

**macOS**:

```bash
brew install mdbook
cargo install mdbook-mermaid    # requires Rust; see below if not installed
```

**Cross-platform (Rust)**:

```bash
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
cargo install mdbook mdbook-mermaid
```

Verify:

```bash
mdbook --version
mdbook-mermaid --version
```

### One-time: install Mermaid assets

mdbook-mermaid needs to drop its JS files into your book folder once:

```bash
mdbook-mermaid install .
```

This creates `mermaid.min.js` and `mermaid-init.js` in the book root. Already wired into `book.toml`.

### Serve

From the `sd-mastery/` folder (where `book.toml` lives):

```bash
mdbook serve --open
```

Auto-opens at `http://localhost:3000`, reloads on file edits.

To just build static files:

```bash
mdbook build      # output in ./book/
```

---

## Option 3: Just Read the Files

Zero install:

- **VS Code**: open `sd-mastery/` folder, click any `.md`, `Cmd+Shift+V` for preview
- **Obsidian**: open `sd-mastery/src/` as a vault
- **GitHub web UI**: after pushing, you can read `.md` files directly on github.com

The GitHub web view has no sidebar nav or cross-file search — mdBook is much better. But this is zero-install.

---

## File Structure

```bash
sd-mastery/
├── book.toml                       # mdBook config
├── .gitignore
├── .github/workflows/deploy.yml    # GitHub Actions auto-deploy
├── SETUP.md
└── src/
    ├── SUMMARY.md                  # sidebar table of contents
    ├── introduction.md
    └── modules/
        ├── 01-thinking-in-systems.md
        ├── ...
        └── 08-architects-craft.md
```

---

## Customization

**Theme**: in `book.toml`, set `default-theme` to: `light`, `rust`, `coal`, `navy`, `ayu`.

**Add your notes**: create more `.md` files in `src/`, link them in `SUMMARY.md`. Great for personal annotations.

**Progress tracking**: create `src/progress.md`, add to `SUMMARY.md`. Use as a checklist.

**Custom CSS**: create `theme/css/general.css`. mdBook picks it up.

---

## Troubleshooting

### `command not found: mdbook` (after `cargo install`)

Add Cargo bin to PATH:

```bash
echo 'export PATH="$HOME/.cargo/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

### `Operation not permitted` on macOS

Folder is in a sandboxed location:

```bash
mv ~/Downloads/sd-mastery ~/
xattr -dr com.apple.quarantine ~/sd-mastery
cd ~/sd-mastery/sd-mastery
```

If still failing: System Settings → Privacy & Security → Full Disk Access → add Terminal/iTerm/VS Code, then restart it.

### `TOML parse error: unknown field 'multilingual'`

Old config. Delete the `multilingual = false` line from `book.toml`. Already removed in the included version.

### GitHub Actions failing

- Check **Actions** tab → click the failed run for logs
- Most common: forgot to set Pages source to "GitHub Actions" in Settings → Pages
- Or: pushed to a non-`main` branch (workflow only triggers on `main`)

### Port 3000 in use locally

```bash
mdbook serve --port 4000 --open
```

---

## Workflow Going Forward

Once setup:

```bash
# Edit + preview locally
mdbook serve --open

# Publish:
git add . && git commit -m "..." && git push
# ~2 min: live site updates
```

GitHub Pages is the easiest path — no install, anh đọc course từ điện thoại / iPad / bất kỳ máy nào.
