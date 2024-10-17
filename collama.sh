#!/bin/bash

echo "Using sudo to switch to the user 'ollama' "
sudo -H -u ollama  OLLAMA_MODELS="/usr/share/ollama/.ollama/models/" ./build/collama $@