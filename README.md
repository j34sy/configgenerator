# Small command-line tool to generate CISCO Router and Switch configs from a network YAML file

## DUE TO OTHER PROJECT CURRENTLY BEING ON BREAK, FOR IPV6 ONLY /64 SUBNETS ARE ALLOWED

## Idea
- Command-line tool to generate configs
    - Hostname
    - Interfaces (IP addresses v4 & v6)
    - VLANs
    - Routing (static & OSPF)
    - basic setup (logging, ip domain-lookup, users)
    - ssh (v2)
    - descriptions
- Own YAML format as input
    - Reduce repetitive input

## To-Dos
- Command-line
    - Config Writer for the data models
    - Syntax checker for the YAML
    - Error catching 
    - Allow every field to be NULL without any errors...

- Loading Network plan
    - Adding fields for more features
    - Writing guide for creating network YAML


- Missing Features
    - IPv6 support
    - Enable secret

