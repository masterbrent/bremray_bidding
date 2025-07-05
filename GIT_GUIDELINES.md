# Git Repository Guidelines

## Files that SHOULD be committed:
- Source code files (.go, .ts, .svelte, .js, .html, .css)
- Configuration examples (.env.example)
- Documentation (.md files, except sensitive ones)
- Package manifests (package.json, go.mod, go.sum)
- Docker configurations (docker-compose.yml, Dockerfile)
- Scripts (start-all.sh, stop-all.sh, etc.)
- Public assets (images, fonts in frontend/public)

## Files that should NEVER be committed:
- Environment files (.env, .env.local, .env.production)
- Build artifacts (binaries, .exe files, dist/, build/)
- Log files (*.log, server.log)
- Temporary files (.pids, *.tmp)
- IDE configurations (.vscode/, .idea/)
- OS files (.DS_Store, Thumbs.db)
- Dependencies (node_modules/, vendor/)
- Sensitive documentation (API keys, credentials)
- Database files (*.db, *.sqlite)
- Compiled binaries (backend/server, backend/bin/)

## Before committing:
1. Run `git status` to check what will be committed
2. Ensure no sensitive information is included
3. Check that .gitignore is properly configured
4. Use `git add .` carefully - prefer adding specific files

## If you accidentally commit sensitive files:
1. Remove from tracking: `git rm --cached <file>`
2. Add to .gitignore
3. Commit the changes
4. Consider rotating any exposed credentials
