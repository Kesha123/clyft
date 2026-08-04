# CLyft

Tool to manage container deployments. This tool fils gap between emrged gitops concepts and container orchestration layers, such as kubernetes, and environments where such orchestration is excesive. This tool provides an ability to manage container deploymets, compose, queadlets, pods, via versioned OCI packages.

```
clyft stack ls
clyft stack sync [--stack, --path, --git]
clyft stack sync [--stack, --path, --git] --preview (diff)
clyft stack rollback [--stack, --path, --git]
clyft stack rollback [--stack, --path, --git] --preview (diff)

clyft init

clyft push <tag>

clyft login <domain>
clyft logout <domian>

clyft upgrade (latest)
clyft upgrade --version

clyft rollback --version

clyft artifact --tag --path

clyft prune (removes packages that a currently not used, i.e. containers by those packages definition are not running)
```
