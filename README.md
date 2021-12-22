# pidproxy

tableflip is a good way graceful process restarts in go, but not provided how to integration with supervisor,
pidproxy is the way you can do it.


https://github.com/cloudflare/tableflip


## Integration with supervisor

[program:Service]
command=/path/to/pidproxy /path/to/binary -some-flag
startsec=3
startretries=3
user=root
directory=/path
redirect_stderr=true
autorestart=true
environment=
