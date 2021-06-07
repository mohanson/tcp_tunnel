set -ex

if [ ! -d ./bin ]; then
    mkdir bin
fi

go build -o bin github.com/mohanson/tcp_tunnel/cmd/tcp_tunnel
