#!/bin/bash
SESSION=go_cqrs

tmux -2 new-session -d -s $SESSION

# Neues Fenster mit VIM.
tmux new-window -t $SESSION:1 -n 'VIM' 
tmux split-window -v
tmux resize-pane 10
tmux send-keys "nocorrect watch go test ./..." C-m
tmux select-pane -t 0
tmux send-keys "vim ." C-m

# Neues Fenster für Terminal.
tmux new-window -t $SESSION:2 -n 'Terminal'
tmux select-window -t $SESION:2
tmux split-window -h

# Vim-Fenster anzeigen.
tmux select-window -t $SESION:1

# Mit Tmux verbinden.
tmux -2 attach-session -t $SESSION
