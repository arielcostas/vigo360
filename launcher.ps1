# SPDX-FileCopyrightText: 2022-2025 Ariel Costas <ariel@costas.dev>
#
# SPDX-License-Identifier: BSD-3-Clause

prebuild() {
	mkdir -p assets/{extra,images,papers,profile,thumb}
	cp -r static/* assets/
	sass --no-source-map -s compressed styles/:assets/
	chmod -R a+r assets/
}

if [ "$1" == "run" ];
then
	prebuild
	export $(cat .env | grep -v '^#' | xargs)
	go run -ldflags "-X main.version=$(git rev-parse --short HEAD)" .
elif  [ "$1" == "build" ];
then
	prebuild
	go build -o vigo360 -ldflags "-X main.version=$(git rev-parse --short HEAD)" .
else
	printf "Vigo360 launcher script\nusage: ${0} [command]\n\n	build: compiles vigo360 to a binary to be ran in production\n	  run: runs vigo360 via in memory.\n"
fi