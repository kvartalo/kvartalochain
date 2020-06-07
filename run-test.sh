#!/bin/sh

SESSION='kvartalochain'

tmux new-session -d -s $SESSION
tmux split-window -d -t 0 -h

# init tendermint
# todo

rm -r tmp
rm -r data
go run main.go initNode
go run main.go initGenesis
# cd test && CLIENT=test go test -initBalance && cd ..

tmux send-keys -t 0 'go run main.go start' enter
sleep 2
tmux send-keys -t 1 'cd test && CLIENT=test go test' enter

tmux attach
