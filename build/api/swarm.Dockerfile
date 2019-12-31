FROM xanderflood/fruit-pi-server:local

COPY ./start.sh ./start.sh

CMD ["./start.sh"]
