#!/bin/bash

set -euo pipefail

. "${XDG_LIB_HOME:-$HOME/.local/lib}/sh/util"

# ~/p/k/ergonaut-one-firmware > udiskie --event-hook 'echo {event} {id_label}'                                                                                                                                                         (master) 21:19
# device_added None
# device_added XIAO-SENSE
# device_mounted XIAO-SENSE
# mounted /org/freedesktop/UDisks2/block_devices/sdb on /run/media/z/XIAO-SENSE

# cd "$HOME/project/ergonaut-one-firmware"

tmpdir

query='.[0].status'
while :; do
	status="$(gh run list --json status,name,databaseId,createdAt --limit 1 --jq "$query")"
	case $status in
		completed)
			break
			;;
		queued)
			sleep 2
			continue
			;;
		*)
			error unknown status $status
			;;
	esac
done

extract='if .[0].status != "completed" then halt_error(1) else .[0].databaseId end'
gh run list --json status,name,databaseId,createdAt --limit 1 --jq "$extract" |
	xargs -I% gh run download % -D "$tmpdir"

echo activate flashing

d="/run/media/$USER/XIAO-SENSE"
if on macos; then
	d="/Volumes/XIAO-SENSE"
fi
while sleep 1; do
	[ -d "$d" ] || echo wait for "$d"
	[ -d "$d" ] || continue

	echo found dir
	cp "$tmpdir/firmware/ergonaut_one_left-seeeduino_xiao_ble-zmk.uf2" "$d" || true
	break
done

echo done
