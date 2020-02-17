goes-boot:
	goes-build coreboot-platina-mk1.rom coreboot-example-amd64.rom

install:
#	$(INSTALL) -m 0644 -d $(DESTDIR)/boot/goes
#	$(INSTALL) goes-boot $(DESTDIR)/boot/goes

clean:
	rm -f *.rom debian/debhelper-build-stamp debian/files debian/*.substvars *.vmlinuz *.xz
	rm -rf debian/.debhelper debian/goes-boot-mk1 debian/goes-boot-example-amd64

binpkg-deb:
	debuild -i -Iworktrees/*

.PHONY: goes-boot install binpkg-deb clean
