#!/usr/bin/env sh

# 创建 COCOON 链
iptables -t nat -N COCOON

# 所有容器流量都转发给 1820 端口
iptables -t nat -A COCOON -p tcp -s 172.17.0.1/16 -j REDIRECT --to-ports 7820

# 非容器流量放行
iptables -t nat -A COCOON -j RETURN

# 在 PREROUTING 链前插入 SHADOWSOCKS 链,使其生效
iptables -t nat -I PREROUTING -p tcp -j COCOON