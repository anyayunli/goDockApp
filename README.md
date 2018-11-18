
## compile
    docker build -t go-dock-app .
    docker run -it --rm --name my-running-app go-dock-app
## run
    docker run -it -p 3344:3344 --rm --name my-running-app go-dock-app
    docker-compose up
