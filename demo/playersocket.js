var playerSocket = null;

function setupPlayerSocket() {
  const url = 'http://0.0.0.0:8000/player';
  playerSocket = io(url, {
    transports: ['websocket'],
  });

  playerSocket.on('connect', () => {
    console.log('connected player socket...');
  });

  playerSocket.on('disconnect', () => {
    console.log('disconnected player socket...');

    reducer.player = null;
  });
}
