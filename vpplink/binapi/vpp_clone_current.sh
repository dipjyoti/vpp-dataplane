#!/bin/bash
VPP_COMMIT=884058096

if [ ! -d $1/.git ]; then
	rm -rf $1
	git clone "https://gerrit.fd.io/r/vpp" $1
	cd $1
	git reset --hard ${VPP_COMMIT}
else
	cd $1
	git fetch "https://gerrit.fd.io/r/vpp" && git reset --hard ${VPP_COMMIT}
fi

# ------------- 10us interrupt patches -------------
# This should be first to avoid hiding failures in the patches
# echo "diff --git a/src/vlib/unix/input.c b/src/vlib/unix/input.c
# index 7531dd197..94a2bfb12 100644
# --- a/src/vlib/unix/input.c
# +++ b/src/vlib/unix/input.c
# @@ -245,7 +245,7 @@ linux_epoll_input_inline (vlib_main_t * vm, vlib_node_runtime_t * node,
#  	      {
#  		/* Sleep for 100us at a time */
#  		ts.tv_sec = 0;
# -		ts.tv_nsec = 1000 * 100;
# +		ts.tv_nsec = 1000 * 10;
 
#  		while (nanosleep (&ts, &tsrem) < 0)
#  		  ts = tsrem;
# " | git apply -- && git add -A &&  git commit --author "Calico-vpp builder <>" -m "Use 10us interrupt sleep"
# ------------- 10us interrupt patches -------------

git fetch "https://gerrit.fd.io/r/vpp" refs/changes/67/30467/1 && git cherry-pick FETCH_HEAD # 30467: tap: fix the buffering index for gro | https://gerrit.fd.io/r/c/vpp/+/30467
git fetch "https://gerrit.fd.io/r/vpp" refs/changes/86/29386/9 && git cherry-pick FETCH_HEAD # 29386: virtio: DRAFT: multi tx support | https://gerrit.fd.io/r/c/vpp/+/29386

# ------------- interrupt patches -------------
git fetch "https://gerrit.fd.io/r/vpp" refs/changes/08/29808/25 && git cherry-pick FETCH_HEAD # 29808: interface: rx queue infra rework, part one | https://gerrit.fd.io/r/c/vpp/+/29808
git fetch "https://gerrit.fd.io/r/vpp" refs/changes/28/30128/5 && git cherry-pick FETCH_HEAD # 30128: virtio: update interrupt mode to new infra | https://gerrit.fd.io/r/c/vpp/+/30128
git fetch "https://gerrit.fd.io/r/vpp" refs/changes/91/30391/5 && git cherry-pick FETCH_HEAD # 30391: interface: fix rx-placement api/cli for new infra | https://gerrit.fd.io/r/c/vpp/+/30391
git fetch "https://gerrit.fd.io/r/vpp" refs/changes/27/30527/1 && git cherry-pick FETCH_HEAD # 30527: interface: let drivers control polling when down | https://gerrit.fd.io/r/c/vpp/+/30527
git fetch "https://gerrit.fd.io/r/vpp" refs/changes/30/30530/1 && git cherry-pick FETCH_HEAD # 30530: interfaces: fix vnet_hw_if_update_runtime_data | https://gerrit.fd.io/r/c/vpp/+/30530
git fetch "https://gerrit.fd.io/r/vpp" refs/changes/85/30485/6 && git cherry-pick FETCH_HEAD # 30485: devices: adapt to new vnet rxq framework | https://gerrit.fd.io/r/c/vpp/+/30485
# ------------- interrupt patches -------------

# ------------- Cnat patches -------------
git fetch "https://gerrit.fd.io/r/vpp" refs/changes/55/29955/4 && git cherry-pick FETCH_HEAD # 29955: cnat: Fix throttle hash & cleanup | https://gerrit.fd.io/r/c/vpp/+/29955
git fetch "https://gerrit.fd.io/r/vpp" refs/changes/73/30273/2 && git cherry-pick FETCH_HEAD # 30273: cnat: Fix session with deleted tr | https://gerrit.fd.io/r/c/vpp/+/30273
git fetch "https://gerrit.fd.io/r/vpp" refs/changes/75/30275/8 && git cherry-pick FETCH_HEAD # 30275: cnat: add input feature node | https://gerrit.fd.io/r/c/vpp/+/30275
git fetch "https://gerrit.fd.io/r/vpp" refs/changes/87/28587/28 && git cherry-pick FETCH_HEAD # 28587: cnat: k8s extensions
# ------------- Cnat patches -------------

# ------------- Policies patches -------------
git fetch "https://gerrit.fd.io/r/vpp" refs/changes/83/28083/14 && git cherry-pick FETCH_HEAD # 28083: acl: acl-plugin custom policies |  https://gerrit.fd.io/r/c/vpp/+/28083
git fetch "https://gerrit.fd.io/r/vpp" refs/changes/13/28513/15 && git cherry-pick FETCH_HEAD # 25813: capo: Calico Policies plugin | https://gerrit.fd.io/r/c/vpp/+/28513
# ------------- Policies patches -------------

git fetch "https://gerrit.fd.io/r/vpp" refs/changes/95/30695/2 && git cherry-pick FETCH_HEAD # 30695: wireguard: testing alternative timer dispatch | https://gerrit.fd.io/r/c/vpp/+/30695
git fetch "https://gerrit.fd.io/r/vpp" refs/changes/26/30726/3 && git cherry-pick FETCH_HEAD # 30726: avf: fix l2_len for csum offload | https://gerrit.fd.io/r/c/vpp/+/30726
git fetch "https://gerrit.fd.io/r/vpp" refs/changes/62/30762/1 && git cherry-pick FETCH_HEAD # 30762: avf: keep irq on main if PF cannot do WB_ON_ITR | https://gerrit.fd.io/r/c/vpp/+/30762

