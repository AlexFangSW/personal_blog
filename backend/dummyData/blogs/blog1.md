---
title: just a dummy blog 
description: Use some key words here to help users find this post
pined: false
visible: true
tags:
- tag-1
- tag-2
topics:
- topic-1
---

## Some background
We ran a k3s cluster with 38 nodes in one of our production enviroment. 

k3s is a light weight k8s alternative, instead of running k8s component such as **kube-controller-manager**, **kub-proxy**, **kubelete** as seperate services, it is packed into a single binary and runs as seperate go routines.

## The issue
One day, random deployments started having network issues, unable to interact with other services. All pods seems to run just fine, service settings binds to the correct pods as well.

So I dug deeper, k8s/k3s by default uses iptables to bind services to pods, a typical network routing goes like this:

Assuming we are using SERVICE_NAME.NAMESPACE as hostnames. ex: http://serviceA.test-namespace:8080/api/v1/xxx

- First, we need to do a dns lookup with **kube-dns**, it is the dns of the cluster, responsible for resolving service name to ip

- After we get the ip, we now enter iptables's NAT table, following the PREROUTING entrypoint. 
  ```bash
  -A PREROUTING -m comment --comment "kubernetes service portals" -j KUBE-SERVICES
  ```

- By following the chains of rules, we end up at a rule which tells us what the target pod ip is
  ```bash
  -A KUBE-SEP-PTGQTQ6HWPYJOHOM -p tcp -m comment --comment "<namespace>/<service name>:http" -m tcp -j DNAT --to-destination 10.42.4.191:8000
  ```

- Now that we know the target pod ip, we follow the routing table of the current node. Bellow is an example of what a typical routing table from a node looks like
  ```bash
  # there is a total of 5 node in this cluster
  # each having pods of different ip range (10.42.xx.0/24)
  # pods on this node have ip range of 10.42.4.0/24

  # if the pod is on other nodes, it will need to go through flannel tunnel
  10.42.0.0/24 via 10.42.0.0 dev flannel.1 onlink 
  10.42.1.0/24 via 10.42.1.0 dev flannel.1 onlink 
  10.42.2.0/24 via 10.42.2.0 dev flannel.1 onlink 

  # you can see that it points to it self when going to pods within this node
  10.42.4.0/24 dev cni0 proto kernel scope link src 10.42.4.1 
  ```

- If the destination is indeed on the current node, we stop here, the packet successfully arived at the target pod.

- If the destination is NOT on this node, the packet will be sent to the target node throught flannel tunnel.

After some digging, I found that on some nodes, iptables are not updated correctly, the current pod ip doesn't exist in them.

And because of this, packets can't be sent to their correct destination.

Going through k3s logs, I couldn't find anything related. And ended up just retarting the k3s service with `systemctl restart k3s.service`

After restart, everything whent back to normal... for a while.

With the root cause unknown, this issue kept comming backup. And I kept restarting k3s when ever this happend.

One day we found a bugfix in k8s 1.28 changelog.
> https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.28.md#bug-or-regression-4

Which stated that after 1.27 there was a bug related to iptables that caused updates getting lost.

Although we were at 1.25, we gave it a shoot. k3s mostly follows k8s, even the release versions match. 

While updateing the cluster, most of the nodes whent smothly, only one node got stuck.

After checking the logs, we found that is was due to networking issues, not just within the cluster, it can't event connect to google.

Anything that run on that node will have the same problem.

A thought hit me, in order for iptables to be up to date, all nodes must continuously communicate with each other.

With this node haveing networking issues, it might be the root cause of our fragmented iptable.

I isolated the node, preventing any pod from being schedualed. So that this node never needs to send a iptable update.

And the iptable issue never happend again. : )

## Afterword
About the node having network issue, my college that was responsible for hardware told me that is was related to it's RAM.
# AAA  
