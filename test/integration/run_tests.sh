echo "Setting up environmental variables"
export GIT_AUTHOR_NAME="Geo Balayan"
export GIT_AUTHOR_EMAIL="geo@balayan.com"

echo "Buildling..."
rm -rf "git-from-sratch"
go build -o "geo-git" ../../bin/main.go
chmod 777 geo-git
export PATH=$PATH:$(pwd)

go test &
program_pid=$!
wait $program_pid
rm -rf "geo-git"