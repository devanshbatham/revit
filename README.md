<h1 align="center">
    revit
  <br>
</h1>

<h4 align="center">A command-line utility for performing reverse DNS lookups</h4>


<p align="center">
  <a href="#install">🏗️ Install</a>
  <a href="#usage">⛏️ Usage</a>
  <br>
</p>


![revit](https://github.com/devanshbatham/revit/blob/main/static/revit.png?raw=true)

# Install
To install revit, follow these steps:

```sh
go install github.com/devanshbatham/revit/cmd/revit@v0.0.1
```


# Usage
This utility allows you to perform reverse DNS lookups on IP addresses. Here are some examples of how to use the tool:

- Look up a single IP address:
```sh
revit -i "8.8.8.8"
```

- Look up a list of IP addresses from a file:
```sh
revit -l ip_list.txt
```

- Pipe input from another command:
```sh
echo "8.8.8.8" | revit
```

```sh
cat ip_list.txt | revit
```

Here are the available command-line flags:

| Flag        | Description                                                        | Example                    |
|-------------|--------------------------------------------------------------------|----------------------------|
| `-i`        | Specify a single target IP address for reverse DNS lookup.        | `revit -i "8.8.8.8"`    |
| `-l`        | Provide the path to a file containing a list of IP addresses.     | `revit -l ip_list.txt`     |
| `-c`        | Set the level of concurrency for concurrent DNS lookups (default: 10). | `revit -c 20`              |
| `-r`        | Specify resolvers for reverse DNS lookup.                         | `revit -r "8.8.8.8"`       |
|             | You can provide a single IP address or a path to a file.          | `revit -r resolvers.txt`   |



# Inspiration

**revit** was born out of curiosity and a desire to explore Golang. While there are existing tools like [hakrevdns](https://github.com/hakluke/hakrevdns) that perform similar tasks (and I have immense resprect for them), I decided to create this utility as a personal project to further my understanding of Go and enhance my programming skills.

The development of **revit** started as an exploration into concurrent programming and networking in Go. As I tinkered with the language's features and learned more about its capabilities, the utility began to take shape. 
