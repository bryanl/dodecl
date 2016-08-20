# Design

What does a declarative syntax look like? Assume, we want to boot minio in
a highly available fashion. This configuration will require four droplets;
two droplets that will be set up as a highly available load balancer, and two
droplets that will be the backend. It will also require a floating IP, dns
for the floating IP, and a volume for data storage.

The ultimate goal is to have a tool that can detrmine state and figure out
exactly what to do. Since that's hard, I'm going to punt and have it generate
infrastructure in a one off mechanism.
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

## Config file

### Dockerfile style

Let's start with the description of our highly available load balancer.

```
FLOATINGIP nyc1
DROPLET lb 2 nyc1 ubuntu-16-04-x64 2gb
```

#### Notes

Using this format will be difficult to create the necessary relations. Also passing
arguments to the commands will be a problem

### YAML Manifest

```yaml
resources:
- type: compute.v2.droplet
  properties:
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

```toml
[dodecl]

project-name = "fancy project"

[[resource]]
type = "compute.v2.droplet
name = "lb"

[resource.properties]
region = "nyc1"
count = 2
image = "ubuntu-16-04-x64"
size = "2gb"
keys = ["104064"]

[[resource]]
type = "compute.v2.floating-ip"

[resource.properties]
region = "nyc1"
```

```jsonnet
{
  resources: [
    {
      type: "compute.v2.droplet",
      properties: {
        name: "lb",
        region: "nyc1",
        count: 2,
        image: "ubuntu-16-04-x64",
        size: "2gb",
        keys: [ 104064 ]
      }
    },
    {
      type: "compute.v2.floating-ip",
      properties: {
        region: "nyc1"
      }
    }
  ]
}
```

## App functions

The app should work in the same spirit as [kubectl](). We'll have commands like

```sh
dodecl create -f <file.yml>
dodecl get <droplet|fip>
dodecl delete <droplet|fip> <name>
```

## Keeping state

One of the issue I have with terrform is the opaque state file. At the current time, it looks like
`dodecl` will have the same thing. The future goal is to allow `dodecl` to bootstrap itself.
Once this happens, it can store the config else where.