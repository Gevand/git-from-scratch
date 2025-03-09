#!/bin/sh

#delete everything except the script
rm -rf .git
rm -rf d
find . ! -name 'test_script.sh' -type f -exec rm -rf {} +
unset GIT_AUTHOR_NAME
unset GITH_AUTHOR_EMAIL


echo "Setting up environmental variables"
export GIT_AUTHOR_NAME="Geo Balayan"
export GIT_AUTHOR_EMAIL="geo@balayan.com"

echo "Buildling..."
rm -rf "git-from-sratch"
go build -o "geo-git" ../bin/main.go

echo "Running innit"
rm -rf ".git"
./geo-git "init"

echo "Generate files to stage"
echo "a" > "a.txt"
echo "b" > "b.txt"
chmod +x b.txt
echo "My First Commit Message" | ./geo-git "commit"

echo "c" > "c.txt"
echo "My Second Commit Message" | ./geo-git "commit"

mkdir "d"
mkdir "d/e"
echo  "f" > "d/e/f.txt"
echo "My Third Commit Message" | ./geo-git "commit"

tree .git

more .git/HEAD

git cat-file -p HEAD^{tree}

find . ! -name 'test_script.sh' -type f -exec rm -rf {} +
rm -rf .git
rm -rf d
unset GIT_AUTHOR_NAME
unset GITH_AUTHOR_EMAIL

