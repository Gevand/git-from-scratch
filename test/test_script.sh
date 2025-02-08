#!/bin/sh

#delete everything except the script
find . ! -name 'test_script.sh' -type f -exec rm -f {} +

echo "Buildling..."
rm -rf "git-from-sratch"
go build -o "git-from-scratch" ../bin/main.go
echo "Running innit"
rm -rf ".git"
./git-from-scratch "init"
echo "Generate files to stage"
echo "a" > "a.txt"
echo "b" > "b.txt"
./git-from-scratch "commit"

tree .git

find . ! -name 'test_script.sh' -type f -exec rm -f {} +
