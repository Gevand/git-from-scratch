#!/bin/sh
alias inflate='ruby -r zlib -e "STDOUT.write Zlib::Inflate.inflate(STDIN.read)"'
#delete everything except the script
rm -rf .git
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
echo -e "Line1\nLine2\nLine3\n" > "a.txt"

./geo-git add a.txt
echo "My First Commit Message" | ./geo-git "commit"


echo -e "Line4\nLine5\n" >> "a.txt"
sed -i 's/Line2/NewLine2/g' a.txt
./geo-git add a.txt
./geo-git diff --cached

#delete everything except the script
rm -rf .git
find . ! -name '*.sh' -type f -exec rm -rf {} +
