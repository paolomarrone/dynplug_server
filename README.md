# Dynplug Server

A web server listening for audio plugins to be loaded into dynplug (https://github.com/paolomarrone/dynplug).

- On startup, it creates, if not existing, a named pipe "dynplug_magicpipe" in the default OS tempdir (e.g. /tmp/ in most unix systems)
- It listens on port 10001
- it saves the file in a temporary directory
- It communicates the path of the file to dynplug via magicpipe

The file needs to be removed by dynplug

# Execution
```
go run main.go
```
or build and run the compiled file.