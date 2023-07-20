@echo off
setlocal ENABLEDELAYEDEXPANSION

if [%1]==[] (
    pushd cmd
    for /d %%G in (*) do (
        go build -o ../builds/%%G/ ./%%G 
        set progs=!progs! split-pane builds\%%G\%%G.exe;
    )
)
popd
wt %progs:~12,-1%
