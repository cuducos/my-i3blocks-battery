# My [`i3blocks`](https://vivien.github.io/i3blocks/) battery block

A simple program to generate a block using `acpi` and [Font Awesome 5](https://fontawesome.com/).

How I use it:

1. Build it with `go build main.go`
2. Change its permissions with `chmod a+x my-i3blocks-battery`
3. Move it to `/usr/local/bin`
4. Adds it to `~/.config/i3/i3blocks.conf`:

```
[battery]
interval=60
command=my-i3blocks-battery
```