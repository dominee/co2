# CO2 - Collaborator Colorizer

**CO2** - Collaborator Colorizer is a simple tool to enhance Burpsuite collaborator logs. Reads input from stdin, processes it and prints it out. 

Current features:

* Amends the source IP address with the name of its resource holder (IP block owner) obtained from RIPE DB.
* Color highlighting to look more 1337 and to replace [ccze](https://github.com/cornet/ccze).


Usage

```
$ co2 [-i]
  -i : Only show lines containing interactions
```

BEFORE: Example pipeline with public collaborator and `ccze`

```bash
curl -sA "${USER}-curl" "http://polling.burpcollaborator.net/burpresults?biid=${BIID}" | tail -n 32 | grep 'IDs:' | ccze -m ansi
```

NOW: Example pipeline with private collaborator and `co2`

```bash
curl -sA "${USER}-curl" -u login:pass "https://yourcollab.pwn:1337/burp.txt?${RANDOM}" | grep "yourcollab" | tail -n 32 | co2 -i
```

Some sample fictional output.

![example-output](example.png)

Note : *Yes, this is just a golang playground project to play with different tooling*