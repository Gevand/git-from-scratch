#!/bin/sh
alias inflate='ruby -r zlib -e "STDOUT.write Zlib::Inflate.inflate(STDIN.read)"'
#delete everything except the script
rm -rf .git
rm -rf d
rm -rf b
rm -rf empty_folder
find . ! -name '*.sh' -type f -exec rm -rf {} +
unset GIT_AUTHOR_NAME
unset GITH_AUTHOR_EMAIL


echo "Setting up environmental variables"
export GIT_AUTHOR_NAME="Geo Balayan"
export GIT_AUTHOR_EMAIL="geo@balayan.com"

echo "Buildling..."
rm -rf "git-from-sratch"
go build -o "geo-git" ../../bin/main.go

echo "Running innit"
rm -rf ".git"
./geo-git "init"

echo "Generate files to stage"
echo "a" > "a.txt"
mkdir b
echo "b" > "b/b.txt"

./geo-git add a.txt b/b.txt
./geo-git status
echo "My First Commit Message" | ./geo-git "commit"

echo "c" > "c.txt"
./geo-git add "c.txt"
echo "My Second Commit Message" | ./geo-git "commit"

mkdir "d"
mkdir "d/e"
echo  "f" > "d/e/f.txt"
echo "Modified c" > "c.txt"
chmod +x b/b.txt #making b executable, modifies it

./geo-git add "d/e/f.txt"
./geo-git "status" "--porcelain"
echo "My Third Commit Message" | ./geo-git "commit"

more .git/HEAD
git cat-file -p HEAD^{tree}

echo "g" > "g.txt"
echo "h" > "h.txt"
mkdir "empty_folder"
./geo-git "add" h.txt g.txt
./geo-git "status" "--porcelain"
./geo-git "status"
./geo-git "diff"

./geo-git "add" "file_that_does_not_exist.txt"

./geo-git "showhead"

#tree .git
#more .git/HEAD
find . ! -name '*.sh' -type f -exec rm -rf {} +
rm -rf .git
rm -rf b
rm -rf d
rm -rf empty_folder
unset GIT_AUTHOR_NAME
unset GITH_AUTHOR_EMAIL

