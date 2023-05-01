default: install

install:
	# install task
	mkdir -p build
	sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b ~/.local/bin