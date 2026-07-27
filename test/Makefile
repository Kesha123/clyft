ZOT_POD     := e2e/zot/zot-pod.yaml
ZOT_NETWORK := zot-net
ZOT_SUBNET  := 172.29.0.0/24

.PHONY: start-zot stop-zot

start-zot:
	@podman network create --driver bridge --subnet $(ZOT_SUBNET) $(ZOT_NETWORK) 2>/dev/null || true
	@podman kube play $(ZOT_POD) --network $(ZOT_NETWORK)

stop-zot:
	@-podman kube down $(ZOT_POD)
	@-podman network rm $(ZOT_NETWORK)
