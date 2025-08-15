# TGCli
Telegram client for interaction in cli

## start

    download tg-cli zip
    add to .env file your creds(or add ENV variables API_ID, API_HASH directly)
    run tg-cli

----

    you also can use tg-cli through params
    -cap string
        caption to photo or file
    -chat string
            chat id
    -f string
            send any file(includes photos)
    -p string
            send photo

## dev start
- Build `tdlib`. Example for debian base os

		sudo apt update
		sudo apt upgrade
		sudo apt install make git zlib1g-dev libssl-dev gperf php-cli cmake g++
		git clone https://github.com/tdlib/td.git
		cd td
		rm -rf build
		mkdir build
		cd build
		cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX:PATH=/usr/local ..
		cmake --build . --target install
		cd ..
		cd ..
		ls -l /usr/local
- Get creds

 	    receive telegram api
		add ENV API_ID and API_HASH via .env file or manual
- Start

  		go run .
