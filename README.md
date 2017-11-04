# Gitea Updater #

This is a simple tool which helps to upgrade your [Gitea](https://github.com/go-gitea/gitea) 
instance to the next version.

It's currently only working on my server, but I try to make it usable with lots of 
Gitea installations.

## Installation ##

Upload the executable to the directory where Gitea is installed and run it as root.

Currently it only works if you:
* Using Gitea in combination with [Systemd](https://en.wikipedia.org/wiki/Systemd)
* Have a symlink called `gitea` in your Gitea directory

### Why running as root? ###
This application was designed to run as root, because my Gitea instance
needs a special linux capability to claim the default ssh port (22). 
This capability can only be added when the command is run by root.
Maybe we can change this in the future, so that the application only
requires sudo access when it's needed.


## Building ##

Just get this repository with
```
go get github.com/LukWebsForge/GiteaUpdater
```