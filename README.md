consul-register
===============

**Currently pre-alpha. Please don't expect it to work just yet.**

Consul-register is a utility app that allows you to define items in configuration and then register those items in a running consul cluster.

This is useful in situations where you have your infrastructure defined as code and want ensure that important parts of your consul configuration are also under version control.

Consul-register allows for the registration of:
* Key/Value pairs
* ACLs
* External services
