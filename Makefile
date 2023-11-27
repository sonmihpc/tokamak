VERSION := $(shell git describe --tags --always --match='v*')
version := $(shell echo $(VERSION) |grep -Eo '[0-9]+\.[0-9]+\.[0-9]')

run:
	go run main.go -c configs/config.yaml
tokamak:
	rm -rf ~/rpmbuild/SOURCES/tokamakd-$(version)*
	rm -f ~/rpmbuild/SPECS/tokamakd.spec
	mkdir -p ~/rpmbuild/SOURCES/tokamakd-$(version)
	go build -o  ~/rpmbuild/SOURCES/tokamakd-$(version)/tokamakd main.go
	cp -r scripts/tokamakd.service ~/rpmbuild/SOURCES/tokamakd-$(version)/
	cp -r configs/config.yaml ~/rpmbuild/SOURCES/tokamakd-$(version)/
	cd ~/rpmbuild/SOURCES;tar -cvzf tokamakd-$(version).tar.gz tokamakd-$(version)/;rm -rf tokamakd-$(version)/
	cp -r scripts/tokamakd.spec ~/rpmbuild/SPECS/
	sed -i 's/1.0.0/$(version)/g' ~/rpmbuild/SPECS/tokamakd.spec
	rpmbuild -bb ~/rpmbuild/SPECS/tokamakd.spec
clean:
	rm -rf ~/rpmbuild/SOURCES/tokamakd-$(version)*
	rm -f ~/rpmbuild/SPECS/tokamakd.spec
	rm -f ~/rpmbuild/RPMS/x86_64/tokamakd-*