all:
	@make -s -C keymap-builder

up:
	~/.local/libexec/update-ergonaut-firmware
