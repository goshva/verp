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
    git commit -m "Deploy VendERP v1.0 - $(date '+%Y-%m-%d %H:%M:%S')"
fi

# Check if remote origin exists, if not add it
if ! git remote | grep -q "origin"; then
    echo "🌐 Adding GitHub remote origin..."
    git remote add origin git@github.com:goshva/verp.git
fi

# Set branch to main
echo "🌿 Setting branch to main..."
git branch -M main

# Push to GitHub
echo "📤 Pushing to GitHub..."
git push -u origin main

echo "✅ Successfully deployed to GitHub!"
echo "🌐 Your repository is available at: https://github.com/goshva/verp"
