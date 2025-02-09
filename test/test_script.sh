#!/bin/sh

#delete everything except the script
rm -rf .git
find . ! -name 'test_script.sh' -type f -exec rm -rf {} +
unset GIT_AUTHOR_NAME
unset GITH_AUTHOR_EMAIL


echo "Setting up environmental variables"
export GIT_AUTHOR_NAME="Geo Balayan"
export GIT_AUTHOR_EMAIL="geo@balayan.com"

echo "Buildling..."
rm -rf "git-from-sratch"
go build -o "git-from-scratch" ../bin/main.go

echo "Running innit"
rm -rf ".git"
./git-from-scratch "init"

echo "Generate files to stage"
echo "a" > "a.txt"
echo "b" > "b.txt"
echo "My first Commit Message" | ./git-from-scratch "commit"

tree .git

find . ! -name 'test_script.sh' -type f -exec rm -rf {} +
rm -rf .git
unset GIT_AUTHOR_NAME
unset GITH_AUTHOR_EMAIL

