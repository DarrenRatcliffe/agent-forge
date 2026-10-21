# Git Workflow Rules

## Branches
- Branch from the main branch, or from the parent branch when stacking dependent PRs. Name them `type/short-description`: `feat/`, `fix/`, `chore/`, `docs/`, `refactor/`.
- Don't commit directly to `main`.

## Commits
- Short imperative subject line, 72 characters or fewer: "Add order pagination", not "added pagination stuff".
- Body explains why, not what, when the reason isn't obvious from the diff.
- One logical change per commit. Don't mix a refactor with a behaviour change.
- Never commit secrets, `.env` files, or build output.

## Pull requests
- Keep them small enough to review in one sitting. If a PR needs a walkthrough, split it.
- Description states what changed, why, and how it was tested.
- In a stack, each PR targets its parent branch and its description states the stack order.
- Call out breaking changes, migrations, and new dependencies explicitly.
- CI must pass before merge. Don't merge with failing or skipped checks.

## History
- Don't rewrite shared history (force-push to `main` or a shared branch).
- Rebase or squash your own branch before merge to keep history readable.
