#!/bin/sh
alias inflate='ruby -r zlib -e "STDOUT.write Zlib::Inflate.inflate(STDIN.read)"'
#delete everything except the script
rm -rf .git
rm -rf d
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

echo "a" > "a.txt"
echo "tbd" > "tbd.txt"
./geo-git add "a.txt" "tbd.txt"
echo "a and tbd added" | ./geo-git "commit"

rm "tbd.txt"
echo "a" >> "a.txt"
echo "b" > "b.txt"
echo "untracked" > "untracked.txt"
./geo-git add "b.txt"
./geo-git status 
./geo-git status --porcelain
./geo-git diff 


tree .git
more .git/HEAD
find . ! -name '*.sh' -type f -exec rm -rf {} +
rm -rf .git
rm -rf d
unset GIT_AUTHOR_NAME
unset GITH_AUTHOR_EMAIL
