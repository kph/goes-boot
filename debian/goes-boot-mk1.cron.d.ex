#
# Regular cron jobs for the goes-boot-mk1 package
#
0 4	* * *	root	[ -x /usr/bin/goes-boot-mk1_maintenance ] && /usr/bin/goes-boot-mk1_maintenance
