#!/bin/bash

# Script to deploy VendERP to GitHub

set -e

echo "🚀 Starting deployment process..."

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ Not a git repository. Initializing..."
    git init
fi

# Add all files
echo "📦 Adding files to git..."
git add .

# Check if there are changes to commit
if git diff-index --quiet HEAD --; then
    echo "✅ No changes to commit"
else
    # Commit changes
    echo "💾 Committing changes..."
    git commit -m "Deploy VERP v1.0 - $(date '+%Y-%m-%d %H:%M:%S')"
fi

# Check if remote origin existsu
if ! git remote | grep -q "origin"; then
    echo "❌ No remote origin found."
    echo "Please add your GitHub repository as origin:"
    echo "git remote add origin git@github.com:goshva/verp.git"
    exit 1
fi

# Push to GitHub
echo "📤 Pushing to GitHub..."
git push -u origin main

echo "✅ Successfully deployed to GitHub!"
echo "🌐 Your repository is available at: https://github.com/$(git config --get remote.origin.url | sed 's/.*github.com[:/]//' | sed 's/\.git$//')"

