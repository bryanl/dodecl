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
