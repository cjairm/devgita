#!/bin/bash

# Get the list of existing tmux sessions
sessions=$(tmux list-sessions -F "#{session_name}" 2>/dev/null)

if [ -n "$sessions" ]; then
	# If there are existing sessions, attach to the first one
	tmux attach-session
else
	# If no sessions exist, start a new session
	tmux new-session -s "misc"
fi
