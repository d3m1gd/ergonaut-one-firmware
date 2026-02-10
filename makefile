VERSION := 0.3.0

.PHONY: build push up

build:
	@sed -i.bak -e '/revision:/s/: .*/: v$(VERSION)/' config/west.yml
	@sed -i.bak -e 's/@.*/@v$(VERSION)/' .github/workflows/build.yml
	@make -s -C keymap-builder

push: build
	git commit -am auto
	git push

up:
	./update.sh
