# Design

What does a declarative syntax look like? Assume, we want to boot minio in
a highly available fashion. This configuration will require four droplets;
two droplets that will be set up as a highly available load balancer, and two
droplets that will be the backend. It will also require a floating IP, dns
for the floating IP, and a volume for data storage.

```
 ┌───────────────────────────────────────────────┐
 │                  floating ip                  │
 └───────────────────────────────────────────────┘
               │                   │
               │
         ┌─────┘                   └ ─ ─   ┌────────────────┐
         │                              │  │  hot standby   │
         │                                 └────────────────┘
         ▼                              ▼
┌─────────────────┐            ┌─────────────────┐
│                 │            │                 │
│                 │            │                 │
│       lb1       │            │       lb2       │
│                 │            │                 │
│                 │            │                 │
└─────────────────┘            └─────────────────┘
         │                              │
         │                              │
         ├──────────────────────────────┤
         │                              │
         │                              │
         ▼                              ▼
┌─────────────────┐            ┌─────────────────┐
│                 │            │                 │
│                 │            │                 │
│       os1       │            │       os2       │
│                 │            │                 │
│                 │            │                 │
└─────────────────┘            └─────────────────┘
         │                              │
         │                              │
         └────────────────┬─────────────┘
                          │
                          │
                          │
                 ┌────────────────┐
                 │                │
                 │                │
                 │     volume     │
                 │                │
                 │                │
                 │                │
                 └────────────────┘
```

## Dockerfile style

Let's start with the description of our highly available load balancer.

```
FLOATINGIP nyc1
DROPLET lb 2 nyc1 ubuntu-16-04-x64 2gb
```

### Notes

Using this format will be difficult to create the necessary relations. Also passing
arguments to the commands will be a problem

## YAML Manifest

```yaml
resources:
- type: compute.v2.droplet
  proerties:
    name: lb
    region: nyc1
    count: 2
    image: ubuntu-16-04-x64
    size: 2gb
    keys:
      - 104064
- type: compute.v2.floating-ip
  properties:
    region: nyc1
```

