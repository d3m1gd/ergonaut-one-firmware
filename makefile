VERSION := 0.3.0

all:
	@sed -i'' '/revision:/s/: .*/: v$(VERSION)/' config/west.yml
	@sed -i'' 's/@.*/@v$(VERSION)/' .github/workflows/build.yml
	@make -s -C keymap-builder

up:
	~/.local/libexec/update-ergonaut-firmware
