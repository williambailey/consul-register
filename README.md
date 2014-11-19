consul-register
===============

[![Build Status](https://travis-ci.org/williambailey/consul-register.svg)](https://travis-ci.org/williambailey/consul-register)

Consul-register allows you to define items in configuration and then register those items on a running consul cluster.

This is useful in situations where you have your infrastructure defined as code and want ensure that important parts of your consul configuration are also under version control.

Consul-register currently supports the registration and exporting of:
* ACLs
* External services
* Key/Value pairs

Usage
-----

```
$ consul-register
consul-register is a tool for managing consul key value storage.
Usage:
consul-register command [arguments]

The commands are:

apply    Apply a list of actions to the consul server.
export   Export consul configuration.

Use "consul-register help [command]" for more information about a command.
```

Please see [example.json](example.json) for the JSON structure that consul-register uses.

Contributing
------------

1. Fork the repository on GitHub
2. Create a named feature branch (i.e. `add-new-feature`)
3. Write your change
4. Submit a Pull Request

Authors
-------

- William Bailey - [@cowboysfromhell](https://twitter.com/cowboysfromhell) - ([mail@williambailey.org.uk](mailto:mail@williambailey.org.uk))

License
-------

Licensed under a [MIT license](LICENSE.txt).
