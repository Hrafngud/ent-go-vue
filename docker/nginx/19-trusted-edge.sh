#!/bin/sh
set -eu

# Only numeric CIDRs may enter generated configuration. Nginx validates addresses
# at startup; reject catch-all trust and configuration injection before writing.
awk 'BEGIN {
    count = split(ENVIRON["TRUSTED_EDGE_CIDRS"], entries, ",")
    for (i = 1; i <= count; i++) {
        cidr = entries[i]
        gsub(/^[[:space:]]+|[[:space:]]+$/, "", cidr)
        if (cidr == "") continue
        if (cidr !~ /^[0-9a-fA-F:.]+\/[0-9]+$/) {
            print "TRUSTED_EDGE_CIDRS requires numeric CIDRs" > "/dev/stderr"
            exit 1
        }
        split(cidr, parts, "/")
        if (parts[2] + 0 == 0) {
            print "TRUSTED_EDGE_CIDRS must not trust all addresses" > "/dev/stderr"
            exit 1
        }
        prefixes[++valid] = cidr
    }
    print "geo $realip_remote_addr $trusted_edge {"
    print "    default 0;"
    for (i = 1; i <= valid; i++) print "    " prefixes[i] " 1;"
    print "}"
    for (i = 1; i <= valid; i++) print "set_real_ip_from " prefixes[i] ";"
    print "real_ip_header X-Forwarded-For;"
    print "real_ip_recursive on;"
}' > /etc/nginx/conf.d/trusted-edge.conf
