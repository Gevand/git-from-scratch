#!/bin/sh
alias inflate='ruby -r zlib -e "STDOUT.write Zlib::Inflate.inflate(STDIN.read)"'
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
go build -o "geo-git" ../../bin/main.go

echo "Running innit"
rm -rf ".git"
./geo-git "init"

echo "Generate files to stage"
echo "a" > "a.txt"
echo "b" > "b.txt"
chmod +x b.txt

./geo-git add a.txt b.txt
echo "My First Commit Message" | ./geo-git "commit"

echo "c" > "c.txt"
./geo-git add .
echo "My Second Commit Message" | ./geo-git "commit"

mkdir "d"
mkdir "d/e"
echo  "f" > "d/e/f.txt"
./geo-git add .
echo "My Third Commit Message" | ./geo-git "commit"

more .git/HEAD
git cat-file -p HEAD^{tree}
more .git/index

echo "g" > "g.txt"
echo "h" > "h.txt"
./geo-git "add" h.txt d g.txt

tree .git
find . ! -name 'test_script.sh' -type f -exec rm -rf {} +
rm -rf .git
rm -rf d
unset GIT_AUTHOR_NAME
unset GITH_AUTHOR_EMAIL

