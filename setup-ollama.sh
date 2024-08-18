#!/usr/bin/env bash

ollama serve &
sleep 10
ollama pull orca-mini