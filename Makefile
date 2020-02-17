goes-boot:
	goes-build coreboot-platina-mk1.rom coreboot-example-amd64.rom

install:
#	$(INSTALL) -m 0644 -d $(DESTDIR)/boot/goes
#	$(INSTALL) goes-boot $(DESTDIR)/boot/goes

binpkg-deb:
	debuild -i -Iworktrees/*

.PHONY: goes-boot install binpkg-deb
