VERSION := 0.3.0

.PHONY: build push up

build:
	@sed -i'' '/revision:/s/: .*/: v$(VERSION)/' config/west.yml
	@sed -i'' 's/@.*/@v$(VERSION)/' .github/workflows/build.yml
	@make -s -C keymap-builder

push: build
	git commit -am auto
	git push

up:
	~/.local/libexec/update-ergonaut-firmware
