# Git Commit Policy

This repository uses a single configured Git identity. Follow these rules for every commit.

## 1. Configured identity

Only `shashank-4bt` may be used as this repository's configured Git identity.

Set identity at the **repository-local** level (preferred over global):

```bash
git config user.name "shashank-4bt"
git config user.email "<the email configured in this repository>"
```

Verify before committing:

```bash
git config --local user.name
git config --local user.email
```

## 2. Configured email

Use only the email configured in this repository's local Git config.

Do not guess, substitute, or replace that email. Do not use a different personal, work, or tool-generated address.

## 3. No Cursor or AI attribution

Cursor may be used as a development tool. That does **not** make Cursor a Git author, committer, or GitHub co-author.

Do not intentionally add Cursor, Cursor Agent, cursoragent, or any equivalent AI attribution to:

- commit authors
- committers
- commit messages
- commit trailers
- Git configuration
- Git hooks
- repository metadata

## 4. No AI co-author trailers

Do not add `Co-Authored-By:` trailers unless a human co-author is explicitly requested.

Never add:

```text
Co-Authored-By: Cursor
Co-Authored-By: cursoragent
```

or any equivalent AI co-author trailer.

## 5. No Cursor contributor metadata

Do not add Cursor contributor metadata, Generated-by-Cursor markers, or GitHub AI co-author metadata.

## 6. No automated AI attribution

Do not add automated labels such as:

- AI Author
- AI Contributor
- AI Attribution
- Generated-by-Cursor

## 7. Verify commit metadata before committing

Before every commit:

1. Confirm local identity:

   ```bash
   git config --local user.name
   git config --local user.email
   ```

2. Review staged files (`git status`, `git diff --cached`).
3. Inspect the commit message for forbidden attribution before finalizing.
4. After committing, inspect:

   ```bash
   git log -1 --format=fuller
   ```

Expected author and committer:

```text
shashank-4bt <the email configured in this repository>
```

If forbidden metadata is present, stop the commit and remove it before continuing.

## 8. Do not rewrite history unless explicitly requested

Never rewrite, amend (except when required to fix a just-created unpushed commit that accidentally received forbidden attribution), rebase, or otherwise alter published history unless the repository owner explicitly requests it.
